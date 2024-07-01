// Copyright © 2024 Timothy E. Peoples

package rerun

import "time"

const (
	Fixed1s    = FixedDelay(time.Second)
	Fixed100ms = FixedDelay(100 * time.Millisecond)
	Fixed500ms = FixedDelay(500 * time.Millisecond)
)

type FixedDelay time.Duration

func (fd FixedDelay) OK(uint) error {
	if fd < 0 {
		return ErrNegativeDuration
	}
	return nil
}

func (FixedDelay) Warmup() time.Duration {
	return 0
}

func (fd FixedDelay) Wait(uint) time.Duration {
	return time.Duration(fd)
}
