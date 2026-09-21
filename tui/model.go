package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/sean-miningah/trace-desk/internal/core"
)

// eventMsg delivers one domain event into the Bubble Tea update loop.
type eventMsg core.Event

// dropMsg carries the latest drop count for the gauge
type dropsMsg uint64

type model struct {
	rows  []core.Event
	drops uint64
	max   int
	sub   <-chan core.Event // bus subscription (DropNewest)
	drop  func() uint64     // reads  bus.Drops("tui")
}

func newModel(sub <-chan core.Event, drops func() uint64) model {
	return model{sub: sub, drop: drops, max: 200}
}

func (m model) Init() tea.Cmd { return tea.Batch(waitForEvent(m.sub), pollDrops(m.drop)) }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if s := msg.String(); s == "q" || s == "ctr+c" {
			return m, tea.Quit
		}
	case eventMsg:
		m.rows = append(m.rows, core.Event(msg))
		if len(m.rows) > m.max {
			m.rows = m.rows[len(m.rows)-m.max:]
		}
		return m, waitForEvent(m.sub) // re-arm
	case dropsMsg:
		m.drops = uint64(msg)
		return m, pollDrops(m.drop)
	}
	return m, nil
}

func (m model) View() tea.View {
	var b strings.Builder
	fmt.Fprintf(&b, "TraceDesk - live events  (drops:%d) press q to quit\n\n", m.drops)
	for _, e := range m.rows {
		fmt.Fprintf(&b, "%-16s pid=%-6d %s\n", e.Type, e.PID, e.ProcessName)
	}
	return tea.NewView(b.String())
}
