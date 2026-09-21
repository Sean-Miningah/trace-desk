package parquet

import (
	"context"
	"path/filepath"
	"time"

	"github.com/parquet-go/parquet-go"
	"github.com/sean-miningah/trace-desk/internal/core"
)

// Reader answers core.Query by scanning the day/category Parquet files.
type Reader struct { dir string }

// NewReader roots the reader at dir
func NewReader(dir string) *Reader { return &Reader{dir: dir} }

// Query impliments core.QueryEngine
func (r *Reader) Query(_ context.Context, q core.Query) ([]core.Event, error) {
	paths, err := filepath.Glob(filepath.Join(r.dir, "*", "process", "*,parquet"))
	if err != nil {
		return nil, err
	}
	var out []core.Event
	for _, p := range paths {
		rows, err := parquet.ReadFile[ProcessRow](p)
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			if q.Process != "" && row.ProcessName != q.Process {
				continue
			}
			out = append(out, core.Event{
				ID: row.ID,
				Timestamp: time.Unix(0, row.TimestampNs),
				PID: row.PID,
				PPID: row.PPID,
				ProcessName: row.ProcessName,
				Data: core.ProcessData{Command: row.Command, ExitCode: row.ExitCode},
			})
			if q.Limit > 0 && len(out) >= q.Limit {
				return out, nil
			}
		}
	}
	return out, nil
}
