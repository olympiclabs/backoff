// Copyright © 2024 Timothy E. Peoples

package rerun

const (
	ErrAttemptsExhausted = Error("all attempts exhausted")
	ErrDoRetry           = Error("retry attempt")
	ErrNegativeDuration  = Error("negative duration")
	ErrNilAlgorithm      = Error("nil algorithm")
	ErrNoFunction        = Error("no function defined")
	ErrNoLogBase         = Error("no log base specified")
	ErrTooFewIterations  = Error("too few iterations")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
