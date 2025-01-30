package gen

import (
	"context"
	"math/rand"
	"time"
)

// RandAlphanum generates a random alphanumeric string of length n
func RandAlphanum(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// InterruptibleSleep sleeps for d, but can be interrupted by ctx
func InterruptibleSleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// JitterDuration returns a duration d with a random jitter applied,
// controlled by the jitter parameter, which acts as a coefficient. For
// example, a jitter of 0.1 means the jitter will be 10% of the duration.
// Passing jitter as 0 means no jitter is applied.
func JitterDuration(d time.Duration, jitter float64) time.Duration {
	if jitter == 0 {
		return d
	}
	jitterVal := float64(d) * jitter
	jitterDur := time.Duration(rand.Float64()*jitterVal) - time.Duration(jitterVal/2)
	return d + jitterDur
}
