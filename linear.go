// Copyright © 2024 Timothy E. Peoples

package rerun

import "time"

// LinearDelay defines a delay Algorithm imposing wait periods along a
// straight line using that old familiar formula you learned in high
// school, "y = mx + b" -- such that:
//
//   - y is the wait period duration for a given iteration -- or
//     the return value from the Wait method
//
//   - m is the slope of the line defined by the Slope field
//
//   - x is 1 less than the current iteration number passed to Wait
//
//   - b is the formula's "Y Intercept" (or rather, the first rerun
//     delay) as defined by the Base field.
//
// Take note that, while both Base and Slope can be zero (which results
// in no waiting periods whatsoever), and Slope *may* be negative, Base
// cannot. A negative Base value will always cause the OK method to return
// ErrNegativeDuration while a negative Slope will only force OK to return
// this error if our line ever crosses the X axis for any x value from
// 1 to Retry's number of iterations.
//
// For example, A LinearDelay with a Base of 100ms and a Slope of -25
// used with a Rerun of 5 iterations would be fine -- since the final
// waiting period would be 25ms. But 7 iterations would cause a problem
// since the final waiting period then be -25ms and, last I checked,
// time travel is not possible (yet).
type LinearDelay struct {
	// Start defines the warmup time Rerun uses before its first call to a Func.
	// This value may be zero or positive but a negative value will cause the
	// OK method to return ErrNegativeDuration.
	Start time.Duration

	// Base defines the waiting period that Rerun.Execute interleaves between its
	// first call to Func and the first rerun attempt. A negative Base value will
	// cause the OK method to return ErrNegativeDuration.
	// i.e. It's the "b" in "y = mx + b".
	Base time.Duration

	// Slope is the slope of the line defined by this type. While Slope may be
	// negative the OK method will return an error if any possible call to
	// Wait would result in a negative Duration value.
	// i.e. It's the "m" in ""y = mx + b".
	Slope float64
}

// OK returns an error if its receiver is il-defined or it defines a line that
// cannot be used for the given number of iterations.
// This method contributes to implementing the Algorithm interface.
func (ld LinearDelay) OK(n uint) error {
	if ld.Start < 0 {
		return ErrNegativeDuration
	}

	if ld.Wait(n) >= 0 {
		return nil
	}

	for i := uint(1); i < n; i++ {
		if ld.Wait(i) < 0 {
			return ErrNegativeDuration
		}
	}
	return nil
}

// Warmup returns the  value of the receiver's Warm field in order to satisfy
// thre Algorithm interface.
func (ld LinearDelay) Warmup() time.Duration {
	return ld.Start
}

// Wait calculates the waiting period, along the receiver's defined line, for
// the given iteration number.
// Wait is part of the Algorithm interface..
func (ld LinearDelay) Wait(n uint) time.Duration {
	if n == 0 {
		return 0
	}

	return time.Duration(ld.Slope*(float64(n)-1)) + ld.Base
}
