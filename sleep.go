// Copyright © 2024 Timothy E. Peoples

package rerun

import (
	"context"
	"time"
)

func sleep(ctx context.Context, d time.Duration) error {
	if d == 0 {
		return nil
	}

	if d < 0 {
		return ErrNegativeDuration
	}

	t := newTimer(d)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C():
		return nil
	}
}

var newTimer = func(d time.Duration) timer {
	return realTimer{time.NewTimer(d)}
}

type timer interface {
	Stop() bool
	C() <-chan time.Time
}

type realTimer struct {
	*time.Timer
}

func (rt realTimer) C() <-chan time.Time {
	return rt.Timer.C
}
