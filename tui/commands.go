package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/sean-miningah/trace-desk/internal/core"
)

// waitForEvent blocks on the bus subscription and turns the next event into a Msg.
func waitForEvent(sub <-chan core.Event) tea.Cmd {
	return func() tea.Msg {
		ev, ok := <-sub
		if !ok {
			return tea.QuitMsg{}
		}
		return eventMsg(ev)
	}
}

// pollDrops refreshes the drop gauge a few times a second
func pollDrops(read func() uint64) tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(time.Time) tea.Msg {
		return dropsMsg(read())
	})
}
