// Copyright © 2024 Timothy E. Peoples

package rerun

import (
	"testing"
)

func TestLogarithmicDelay(t *testing.T) {
	ld := LogarithmicDelay{
		Units:          Millisecond,
		Amplifier:      300,
		Coefficient:    20,
		Modifier:       -14,
		VerticalOffset: -400,
	}

	for i := uint(1); i < 10; i++ {
		t.Logf("%3d: %v", i, ld.Wait(i))
	}
}
