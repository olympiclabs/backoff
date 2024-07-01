// Copyright © 2024 Timothy E. Peoples

package rerun

import (
	"math"
	"time"
)

// LogarithmicDelay implements the Algorithm interface to generate a
// logarithmic progression of wait times defined by the function:
//
//	     A·ln(C·X + M) + V
//	W = -------------------
//	            D
//
// ...where:
//
//   - X: Iteration Number (given as the uint argument to Wait)
//   - W: Generated Wait Time (as returned by Wait)
//   - A: Amplifier Field
//   - C: Coefficient Field
//   - M: Modifier Field
//   - V: VerticalOffset Field
//   - D: Denominator Field
//
// Since none of the above fields are of type time.Duration, the Units field
// should be used to ensure the return value from Wait is interpreted at the
// intended resolution.
type LogarithmicDelay struct {
	Start time.Duration
	Units delayUnits

	Amplifier      float64
	Coefficient    float64
	Modifier       float64
	VerticalOffset float64
	Denominator    float64
}

// OK checks the validity of the receiver's fields then calculates a wait time
// for each iteration value from 1 to n inclusive. If any field has an invalid
// value or if a wait time is calculated to be less than zero, an error is
// returned.
//
// The field rules for this type are:
//
//   - The Start field cannot be less than zero.
//
//   - The Amplifier, Coefficient, and Denominator fields cannot be zero but may
//     be positive or negative. Albeit, negative values are likely to generate
//     negative wait times, which will also cause an error.
//
//   - The Modifier and VerticalOffset fields may contain any value that do not
//     result in a negative calculated wait time.
//
// OK contributes to implementing the Algorithm interface.
func (ld LogarithmicDelay) OK(n uint) error {
	if ld.Start < 0 {
		return ErrNegativeDuration
	}

	for i := uint(1); i < n; i++ {
		if ld.Wait(i) < 0 {
			return ErrNegativeDuration
		}
	}

	return nil
}

func (ld LogarithmicDelay) Warmup() time.Duration {
	return ld.Start
}

func (ld LogarithmicDelay) Wait(n uint) time.Duration {
	if n == 0 {
		return 0
	}

	d := time.Duration(ld.Amplifier*math.Log(ld.Coefficient*float64(n)+ld.Modifier) + ld.VerticalOffset)
	return time.Duration(ld.Units) * d
}
