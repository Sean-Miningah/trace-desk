// Package wires sources to consumers
// Normalizes raw events and fans each event to every subscriber under
// a per-subscriber delivery policy
package pipeline

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/sean-miningah/trace-desk/internal/core"
)

// Policy controls what happens when a subscriber's buffer is full
type Policy int

const (
	Block Policy = iota
	DropNewest
)

type subscriber struct {
	name string
	ch chan core.Event
	policy Policy
	drops atomic.Uint64
}

func (s *subscriber) deliver (ev core.Event) {
	switch s.policy {
		case Block:
			s.ch <- ev
		case DropNewest:
			select {
			case s.ch <- ev:
			default:
				s.drops.Add(1)
			}
	}
}

// Bus i an in-process channel fan-out.
type Bus struct {
	mu sync.RWMutex
	subs []*subscriber
}

func NewBus() *Bus { return &Bus{} }

// Register a named subscriber with its own delivery policy, returning and recieve end.
func (b *Bus) Subscribe(name string, buffer int, policy Policy) <- chan core.Event {
	s := &subscriber{
		name: name,
		ch: make(chan core.Event, buffer),
		policy: policy,
	}
	b.mu.Lock()
	b.subs = append(b.subs, s)
	b.mu.Unlock()
	return s.ch
}

// Drops returns the number of dropped events for the given subscriber.
func (b *Bus) Drops(name string) uint64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, s := range b.subs {
		if s.name == name {
			return s.drops.Load()
		}
	}
	return 0
}

// Publish sends an event to all subscribers.
func (b *Bus) Publish(ev core.Event){
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, s := range b.subs {
		s.deliver(ev)
	}
}

func (b *Bus) Run(ctx context.Context, in <-chan core.Event) {
	defer b.closeAll()
	for {
		select {
			case <-ctx.Done():
				return
			case ev, ok := <-in:
				if !ok {
					return
				}
				b.Publish(ev)
		}
	}
}

func (b *Bus) closeAll() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, s := range b.subs {
		close(s.ch)
	}
	b.subs = nil
}
