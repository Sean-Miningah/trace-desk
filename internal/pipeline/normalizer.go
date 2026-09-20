package pipeline

import (
	"sync/atomic"
	"time"

	"github.com/sean-miningah/trace-desk/internal/core"
)

// Ras iw source record before normalization. Real source is (eBPF) will have their own
// raw structs. Normalizer maps them onto core.Event
type Raw struct {
	Type core.EventType
	PID uint32
	PPID uint32
	UID uint32
	GID uint32
	ProcessName string
	Data core.EventData
	When time.Time
}

type Normalizer struct {
	seq atomic.Uint64
}

// NewNormalizer returns a a ready Normalizer
func NewNormalizer() *Normalizer {
	return &Normalizer{}
}

// Normalize converts a raw record intoa a core.Event
func (n *Normalizer) Normalize(r Raw) core.Event {
	ts := r.When
	if ts.IsZero() {
		ts = time.Now()
	}
	return core.Event{
		ID: n.seq.Add(1),
		Timestamp: ts,
		Type: r.Type,
		PID: r.PID,
		PPID: r.PPID,
		UID: r.UID,
		GID: r.GID,
		ProcessName: r.ProcessName,
		Data: r.Data,
	}
}
