package core

import (
	"testing"
)

func TestEvent(t *testing.T) {
	cases := map[EventType]string{
		ProcessExec: "process.exec",
		ProcessExit: "process.exit",
		FileOpen: "file.open",
		FileRead: "file.read",
		FileWrite: "file.write",
		FileRename: "file.rename",
		NetworkConnect: "network.connect",
		NetworkAccept: "network.accept",
	}
	for typ, want := range cases {
		if got := typ.String(); got != want {
			t.Errorf("EventType(%d).String() = %q, want %q", typ, got, want)
		}
	}
	if got := EventType(200).String(); got != "unknown" {
		t.Errorf("out-of-range String = %q, want %q", got, "unknown")
	}
}
