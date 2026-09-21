package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/sean-miningah/trace-desk/internal/core"
)

// Run starts the TUI over a bus subscription. drops reports the tui subscriber's
// dropped-event count for the on-screen gauge.
func Run(sub <-chan core.Event, drops func() uint64) error {
	_, err := tea.NewProgram(newModel(sub, drops)).Run()
	return err
}
