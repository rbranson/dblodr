package gen

import (
	"sync"
	"sync/atomic"
	"time"
)

// Stats is for generic thread-safe stats tracking in-memory by key
type Stats struct {
	mu       sync.RWMutex
	counters map[string]*atomic.Int64

	// Parent will be receive all the same actions as this instance
	Parent *Stats
}

// requires mu.RLock
func (s *Stats) checkAndInitializeCounter(key string) {
	var hasKey bool
	if s.counters != nil {
		_, hasKey = s.counters[key]
	}
	if s.counters == nil || !hasKey {
		s.mu.RUnlock()
		s.mu.Lock()
		if s.counters == nil {
			s.counters = make(map[string]*atomic.Int64)
		}
		s.counters[key] = new(atomic.Int64)
		s.mu.Unlock()
		s.mu.RLock()
	}
}

func (s *Stats) counterIncr(key string, delta int64) {
	s.mu.RLock()
	s.checkAndInitializeCounter(key)
	s.counters[key].Add(delta)
	s.mu.RUnlock()

	if s.Parent != nil {
		s.Parent.counterIncr(key, delta)
	}
}

// CounterIncr adds delta to the counter at key
func (s *Stats) CounterIncr(key string, delta int) {
	s.counterIncr(key, int64(delta))
}

// CounterIncrDur adds delta to the counter at key, treating the delta based on unit
func (s *Stats) CounterIncrDur(key string, delta time.Duration, unit time.Duration) {
	s.counterIncr(key, int64(delta/unit))
}

// ResetCounter sets the counter at key to 0
func (s *Stats) ResetCounter(key string) {
	s.mu.RLock()
	s.checkAndInitializeCounter(key)
	s.counters[key].Store(0)
	s.mu.RUnlock()

	if s.Parent != nil {
		s.Parent.ResetCounter(key)
	}
}

// ReadCounters returns a copy of the current counters
func (s *Stats) ReadCounters() map[string]int64 {
	cloned := make(map[string]int64, len(s.counters))
	s.mu.Lock() // use lock instead of RLock to ensure a consistent snapshot
	defer s.mu.Unlock()
	if len(s.counters) == 0 {
		return cloned
	}
	for k, v := range s.counters {
		cloned[k] = v.Load()
	}
	return cloned
}

// ReadCounter returns the current value of the counter at key, returns 0 if
// the key doesn't exist.
func (s *Stats) ReadCounter(key string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if counter, ok := s.counters[key]; ok {
		return counter.Load()
	}
	return 0
}
