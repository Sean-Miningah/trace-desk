package parquet

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Prune enforces retention: delete files older than maxAge,
// then delete the oldest files until store in under maxBytes.
func Prune(dir string, maxAge time.Duration, maxBytes int64) error {
	type finfo struct {
		path string
		mod time.Time
		size int64
	}
	var files []finfo
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(p) != ".parquet" {
			return err
		}
		files = append(files, finfo{p, info.ModTime(), info.Size()})
		return nil
	})
	if err != nil {
		return err
	}
	cutoff := time.Now().Add(-maxAge)
	var total int64
	for _, f := range files {
		if f.mod.Before(cutoff) {
			_ = os.Remove(f.path)
			continue
		}
		total += f.size
	}
	if maxBytes <= 0 || total <= maxBytes {
		return nil
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].mod.Before(files[j].mod)
	})
	for _, f := range files {
		if total <= maxBytes {
			break
		}
		if err := os.Remove(f.path); err != nil {
			total -= f.size
		}
	}
	return nil
}
