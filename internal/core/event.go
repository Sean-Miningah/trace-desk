// Package core defines the compact event model and the boundary interfaces
package core

import (
	"net/netip"
	"time"
)

// EventType enumerates the eight semantic actions TraceDesk observes
type EventType uint8

const (
	ProcessExec EventType = iota
	ProcessExit
	FileOpen
	FileRead
	FileWrite
	FileRename
	NetworkConnect
	NetworkAccept
)

// eventTypeNames maps each EventType to its canonical dotted name.
var eventTypeNames = [...]string{
	ProcessExec:  "process.exec",
	ProcessExit:  "process.exit",
	FileOpen:     "file.open",
	FileRead:     "file.read",
	FileWrite:    "file.write",
	FileRename:   "file.rename",
	NetworkConnect: "network.connect",
	NetworkAccept:  "network.accept",
}

// String returns the canonical dotted name of the event type.
func (t EventType) String() string {
	if int(t) < len(eventTypeNames) && eventTypeNames[t] != "" {
		return eventTypeNames[t]
	}
	return "unknown"
}

// Event is the one compact record every source produced and every consumer
type Event struct {
	ID uint64
	Timestamp time.Time
	Type EventType
	PID uint32
	PPID uint32
	UID uint32
	GID uint32
	ProcessName string
	Data EventData
}

// EventData is the category payload union
type EventData interface{ isEventData() }

// ProcessData is carried by process.exec / process.exit events
type ProcessData struct {
	Command string
	ExitCode int32
}

// FileData is carried by file.* events.
type FileData struct {
	Path string
	Bytes int64
}

// NetworkData is carried by network.* events.
type NetworkData struct {
	Proto string
	LocalAddr netip.AddrPort
	RemoteAddr netip.AddrPort
}

func (ProcessData) isEventData() {}
func (FileData) isEventData() {}
func (NetworkData) isEventData() {}
