package pipeline

import (
	"testing"

	"github.com/sean-miningah/trace-desk/internal/core"
)

// DropNewest must never block the producer, even with a full buffer
func TestDropNewestNeverBlock(t *testing.T) {
	b := NewBus()
	_ = b.Subscribe("ui", 1, DropNewest)
	for range 1000 {
		b.Publish(core.Event{Type: core.ProcessExec})
	}
	if got := b.Drops("ui"); got == 0 {
		t.Fatalf("expected drops > 0 on undrained DropNewest subscriber, got 0")
	}
}

// Block must lose nothing  when the consumer keeps up.
func TestBlockLosesNothing(t *testing.T) {
	b := NewBus()
	ch := b.Subscribe("store", 8, Block)
	const n = 100
	done := make(chan int)
	go func() {
		count := 0
		for range ch {
			count++
			if count == n {
				break
			}
		}
		done <- count
	}()
	for i := range n {
		b.Publish(core.Event{ID: uint64(i + 1), Type: core.FileOpen})
	}
	if got := <-done; got != n {
		t.Fatalf("Block subscriber received %d events, expected %d", got, n)
	}
	if d := b.Drops("store"); d != 0 {
		t.Fatalf("expected drops to be 0 after Block subscriber keeps up, got %d", d)
	}
}
