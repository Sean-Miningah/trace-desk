// Package parquet persists events as columnar Parquet, partitioned by day
// and category.
package parquet

// ProcessRow is the flat row of process.* events.
type ProcessRow struct {
	ID          uint64 `parquet:"id"`
	TimestampNs int64  `parquet:"ts_ns,delta"`
	Type        string `parquet:"type,dict,zstd"`
	PID         uint32 `parquet:"pid"`
	PPID        uint32 `parquet:"ppid"`
	UID         uint32 `parquet:"uid"`
	ProcessName string `parquet:"proc,dict,zstd"`
	Command     string `parquet:"cmd,zstd"`
	ExitCode    int32  `parquet:"exit"`
}

// FileRow is the flat row for file.* events.
type FileRow struct {
	ID          uint64 `parquet:"id"`
	TimestampNs int64  `parquet:"ts_ns,delta"`
	Type        string `parquet:"type,dict,zstd"`
	PID         uint32 `parquet:"pid"`
	ProcessName string `parquet:"proc,dict,zstd"`
	Path        string `parquet:"path,dict,zstd"`
	Bytes       int64  `parquet:"bytes"`
}

// NetworkRow is the flat row for network.* events.
type NetworkRow struct {
	ID          uint64 `parquet:"id"`
	TimestampNs int64  `parquet:"ts_ns,delta"`
	Type        string `parquet:"type,dict,zstd"`
	PID         uint32 `parquet:"pid"`
	ProcessName string `parquet:"proc,dict,zstd"`
	Proto       string `parquet:"proto,dict,zstd"`
	RemoteAddr  string `parquet:"raddr,dict,zstd"`
	RemotePort  uint32 `parquet:"rport"`
}
