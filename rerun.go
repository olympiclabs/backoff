// Copyright © 2024 Timothy E. Peoples

// Package rerun implements logic to rerun a given function up to a set number
// of times with configurable wait periods interleaved between each attempt.
// XXX MORE XXX
package rerun

import (
	"context"
	"fmt"
	"time"
)

type delayUnits time.Duration

const (
	// The following values are available for use by types implementing
	// the Algorithm interface where none of their underlying paramaters
	// provide a time.Duration anchor value. They use an unexported type
	// (wrapping time.Duration) to ensure that only basic unit values are
	// used.  See LogarithmicDelay for an example of how this may be used.
	Nanosecond  = delayUnits(time.Nanosecond)
	Microsecond = delayUnits(time.Microsecond)
	Millisecond = delayUnits(time.Millisecond)
	Second      = delayUnits(time.Second)
	Minute      = delayUnits(time.Minute)
	Hour        = delayUnits(time.Hour)
)

// The Algorithm interface is implemented by types defining the waiting
// periods Rerun.Execute will interleave between each call to a provided
// Func.
type Algorithm interface {
	// OK is called to ensure the underlying parameters defining this Algorithm
	// implementation are valid and should only return nil if all calls to its
	// Wait method (up to the provided uint value) will successfully calculate
	// a valid wait time.
	OK(uint) error

	// Warmup should return a time.Duration for the waiting period to be
	// imposed by Rerun.Execute prior to the first execution of its Func.
	// If no warmup time is desired, this method should return zero.
	Warmup() time.Duration

	// Wait returns the time.Duration for the waiting period to be imposed
	// by Rerun.Execute before each retry of its associated Func.  The
	// provided uint value is the retry iteration number about to be attempted;
	// the first call to Wait will be retry #1, the second is #2, and so on.
	//
	// Note that since wait times are interleaved between each retry attempt,
	// the number of calls to Wait will always be 1 less than the number of
	// configured iterations.
	Wait(uint) time.Duration
}

// Rerun defines the behavior for running a given function up to a set number
// of times with configurable waiting periods interleaved between each attempt.
// The zero-value is unusable.
type Rerun struct {
	iterations uint
	algorithm  Algorithm
	function   Func
	err        error
}

// DefaultAlgorithm is the default Algorithm used by Rerun.Execute if no other
// Algorithm is specified. This default is a 1s FixedDelay algorithm with no
// warmup time and a fixed, 1s wait between each retry attempt.
const DefaultAlgorithm = Fixed1s

// New returns a new Rerun object configured for the given number of
// iterations using the DefaultAlgorithm. To employ a different Algorithm,
// use the WithAlgorithm option method.
func New(i uint) *Rerun {
	return &Rerun{iterations: i, algorithm: DefaultAlgorithm}
}

// WithAlgorithm returns a pointer to its receiver after updating its attached
// Algorithm to the given value. If algo is nil or its OK method returns an
// error subsequent calls to the receiver's Err method will return a non-nil
// error. Note that since this method does not employ a pointer receiver,
// only the return value will be updated (but not the caller's receiver value).
func (r Rerun) WithAlgorithm(algo Algorithm) *Rerun {
	if algo == nil {
		r.err = ErrNilAlgorithm
	} else {
		r.err = algo.OK(r.iterations)
	}
	r.algorithm = algo

	return &r
}

// WithFunction returns a pointer to its receiver after updating its associated
// Fun to the given value. If the receiver already has an associated Func value
// it will be silently overwritten and passing a nil Func here will clear the
// receiver's Func value (if any). Note that calling the Execute method with
// a Rerun having a nil Func associated will always results in an error.
func (r Rerun) WithFunction(function Func) *Rerun {
	r.function = function
	return &r
}

// Err returns any non-nil error that occurred during construction of its
// receiver or if the OK method for the receiver's Algorithm returns an
// error.
func (r Rerun) Err() error {
	if r.err == nil {
		r.err = r.algorithm.OK(r.iterations)
	}
	return r.err
}

// Func defines the signature for functions called by Rerun.Execute.
type Func func(uint) error

// Execute is used to repeatedly execute the reciever's configured Func while
// interleaving wait periods as defined by the Algorithm attached to the
// receiver. Execute's behavior is goverened by the following rules:
//
//   - If the receiver no associated Func configured, ErrNoFunction
//     is returned.
//
//   - If the receiver configured with fewer than 2 iterations,
//     ErrTooFewIterations is returned.
//
//   - If r.Err() returns a non-nil error, that error will be returned
//     immediately.
//
//   - If the receiver's configure Func  returns a nil error, Execute
//     returns immediately.  If the provided Context has not yet become
//     done, then Execute returns a nil error. Otherwise, Execute will
//     return ctx.Err().
//
//   - If the receiver's Func returns ErrDoRetry -- and Execute has not
//     yet exhausted all of the receiver's configured iterations -- then
//     Execute will pause for the Duration returned by Algorithm.Wait.
//     If the given Context becomes done during this wait period, Execute
//     will immediately return ctx.Err(). Otherwise, the receiver's Func
//     will be rerun after the alotted wait time.
//
//   - If the receiver's Func returns ErrDoRetry -- but all of the receiver's
//     configured iterations, have been exhausted -- then no pause will be
//     introduced and Execute instead ErrAttemptsExhausted immediately.
//
//   - If the receiver's Func causes a panic, it will be recovered and
//     returned as an error.
//
//   - Otherwise, Execute returns the error returned by the receiver's Func.
//
// Prior to executing the receiver's Func for the first time, Execute calls
// Algorithm.Warmup to determine whether it should pause for a warmup period
// and behaves accordingly based on what's returned:
//
//   - If Warmup returns a positive value, Execute will pause for that
//     Duration before its first attempt.  However, if the given Context
//     becomes done during this period, Execute immediately returns
//     ctx.Err().
//
//   - If Warmup returns 0, no delay will be imposed before the first call
//     to the Func.
//
//   - If Warmup returns a negative value, Execute returns ErrNegativeDuration
//
// Generally, regardless of the error returned by the receiver's Func, if ctx
// becomes done, Execute will err towards returning ctx.Err() as soon as that
// can be detected -- even during waiting periods (albeit, no effort is made
// to cover any race conditions so this is not guaranteed).
func (r Rerun) Execute(ctx context.Context) (err error) {
	defer func() {
		select {
		default:
		case <-ctx.Done():
			err = context.Cause(ctx)
		}
	}()

	if r.function == nil {
		return ErrNoFunction
	}

	if r.iterations < 2 {
		return ErrTooFewIterations
	}

	if err = r.Err(); err != nil {
		return err
	}

	// n.b. If Warmup returns 0, sleep will immediately return a nil error.
	if err = sleep(ctx, r.algorithm.Warmup()); err != nil {
		return err
	}

	for i := uint(0); i < r.iterations; i++ {
		if i > 0 {
			if err = sleep(ctx, r.algorithm.Wait(i)); err != nil {
				return err
			}
		}

		switch err = r.runFunction(i); err {
		case nil:
			return nil

		case ErrDoRetry:
			continue

		default:
			return err
		}
	}

	return ErrAttemptsExhausted
}

// runFunction executes the Func associated with the receiver. Any panic
// caused by doing so will be recovered and returned as an error.
func (r Rerun) runFunction(i uint) (err error) {
	defer func() {
		if perr := recover(); perr != nil {
			err = fmt.Errorf("recovered from panic: %v", perr)
		}
	}()

	return r.function(i)
}
