package stream

import "context"

type Backpressure struct {
	sem chan struct{}
}

func NewBackpressure(maxInFlight int) *Backpressure {
	if maxInFlight <= 0 {
		maxInFlight = 64
	}
	return &Backpressure{sem: make(chan struct{}, maxInFlight)}
}

func (b *Backpressure) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case b.sem <- struct{}{}:
		return nil
	}
}

func (b *Backpressure) Release() {
	select {
	case <-b.sem:
	default:
	}
}
