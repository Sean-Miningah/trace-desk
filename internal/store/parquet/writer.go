package parquet

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/parquet-go/parquet-go"
	"github.com/sean-miningah/trace-desk/internal/core"
)

// Writer batches events and flushes per-day, per-category Parquet files.
// Each flust writes a new sequenced file.
type Writer struct {
	dir string
	seq map[string]int
}

// NewWriter roots the store at dir (e.g "data")
func NewWriter(dir string) *Writer { return &Writer{dir: dir, seq: map[string]int{}} }

// Write implements core.EventStore: it splits the batch by category and writes
// one file par category for the batch's day.
func (w *Writer) Write(_ context.Context, batch []core.Event) error {
	var procs []ProcessRow
	var filez []FileRow
	var nets []NetworkRow
	for _, e := range batch {
		switch d := e.Data.(type) {
		case core.ProcessData:
			procs = append(procs, ProcessRow{
				ID:          e.ID,
				TimestampNs: e.Timestamp.UnixNano(),
				Type:        e.Type.String(),
				PID:         e.PID,
				PPID:        e.PPID,
				UID:         e.UID,
				ProcessName: e.ProcessName,
				Command:     d.Command,
				ExitCode:    d.ExitCode,
			})
		case core.FileData:
			filez = append(filez, FileRow{
				ID:          e.ID,
				TimestampNs: e.Timestamp.UnixNano(),
				Type:        string(e.Type),
				PID:         e.PID,
				ProcessName: e.ProcessName,
				Path:        d.Path,
				Bytes:       d.Bytes,
			})
		case core.NetworkData:
			nets = append(nets,
				NetworkRow{
					ID:          e.ID,
					TimestampNs: e.Timestamp.UnixNano(),
					Type:        e.Type.String(),
					PID:         e.PID,
					ProcessName: e.ProcessName,
					Proto:       d.Proto,
					RemoteAddr:  d.RemoteAddr.Addr().String(),
					RemotePort:  uint32(d.RemoteAddr.Port()),
				})
		}
		day := time.Now().UTC().Format("2006-01-02")
		if err := flush(w, day, "process", procs); err != nil {
			return err
		}
		if err := flush(w, day, "file", filez); err != nil {
			return err
		}
		if err := flush(w, day, "network", nets); err != nil {
			return err
		}
	}
	return nil
}

func flush[T any](w *Writer, day, category string, rows []T) error {
	if len(rows) == 0 {
		return nil
	}
	dir := filepath.Join(w.dir, day, category)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	w.seq[category]++
	name := filepath.Join(dir, fmt.Sprintf("%06d.parquet", w.seq[category]))
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	pw := parquet.NewGenericWriter[T](f, parquet.Compression(&parquet.Zstd))
	if _, err := pw.Write(rows); err != nil {
		return err
	}
	return pw.Close() // flushes buffers + footer
}
