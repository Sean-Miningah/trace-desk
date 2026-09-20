// Package provides a fake source for development and testing purposes.
package fake

import (
	"context"
	"math/rand"
	"net/netip"
	"time"

	"github.com/sean-miningah/trace-desk/internal/core"
)

type Source struct {
	Rate time.Duration // delay between events; 0 means as fast as possible
	rng *rand.Rand
}

func New(rate time.Duration) *Source {
	return &Source{
		Rate: rate,
		rng:  rand.New(rand.NewSource(1)),
	}
}

func (s *Source) Events(ctx context.Context) (<-chan core.Event, error) {
	out := make(chan core.Event)
	go func() {
		defer close(out)
		var id uint64
		for {
			if s.Rate > 0 {
				select {
					case <- ctx.Done():
						return
					case <- time.After(s.Rate):
				}
			} else if ctx.Err() != nil {
				return
			}
			id++
			ev := s.sample(id)
			select {
				case <- ctx.Done():
					return
				case out <- ev:
			}
		}
	}()
	return out, nil
}

// Close impliments core.EventSource
func (s *Source) Close() error { return nil }

var procs = []string{"firefox", "curl", "bash", "sshd", "cron"}

func (s *Source) sample(id uint64) core.Event {
	name := procs[s.rng.Intn(len(procs))]
	ev := core.Event{
		ID: id,
		Timestamp: time.Now(),
		PID: 1000 + uint32(s.rng.Intn(500)),
		PPID:  1,
		ProcessName: name,
	}
	switch s.rng.Intn(4) {
		case 0:
			ev.Type = core.ProcessExec
			ev.Data = core.ProcessData{
				Command: name,
			}
		case 1:
			ev.Type = core.FileOpen
			ev.Data = core.FileData{Path: "etc/passwd"}
		case 2:
			ev.Data = core.NetworkData{
				Proto: "tcp",
				RemoteAddr: netip.MustParseAddrPort("93.184.216.34:443"),
			}
		default:
			ev.Type = core.ProcessExit
			ev.Data = core.ProcessData{Command: name, ExitCode: 0}
	}
	return ev
}
