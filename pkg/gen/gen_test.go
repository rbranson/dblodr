package gen

import (
	"testing"
	"time"
)

func TestJitterDuration(t *testing.T) {
	tests := []struct {
		name   string
		d      time.Duration
		jitter float64
	}{
		{"NoJitter", time.Second, 0},
		{"SmallJitter", time.Second, 0.1},
		{"MediumJitter", time.Second, 0.5},
		{"LargeJitter", time.Second, 1},
	}

	const iterations = 1000

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < iterations; i++ {
				got := JitterDuration(tt.d, tt.jitter)
				if tt.jitter == 0 && got != tt.d {
					t.Errorf("Expected %v, got %v", tt.d, got)
				}
				if tt.jitter != 0 {
					min := tt.d - time.Duration(float64(tt.d)*tt.jitter/2)
					max := tt.d + time.Duration(float64(tt.d)*tt.jitter/2)
					if got < min || got > max {
						t.Errorf("Expected duration between %v and %v, got %v", min, max, got)
					}
				}
			}
		})
	}
}
