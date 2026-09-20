package core

import (
	"context"
	"time"
)

// EventSource produces a stream of events.
type EventSource interface {
	Events(ctx context.Context) (<-chan Event, error)
	Close() error
}

// EventStore is persists events losslessly. Implementation: the Parquet store
type EventStore interface {
	Write(ctx context.Context, batch []Event) error
	Close() error
}

// Query is a structured request over persisted events. Zero-valued fields
// are ignored (n filter). Query is the ONLY thing user input is translated into
// engines never execute raw generated SQL.
type Query struct {
	Process string
	File string
	Type *EventType
	Since time.Time
	Limit int
}

// QueryEngine answers structured queries over peristed events.
type QueryEngine interface {
	Query(ctx context.Context, q Query) (<-chan Event, error)
}
