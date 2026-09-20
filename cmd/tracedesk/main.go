package main

import (
	"log/slog"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		slog.Info("tracedesk", "version", version)
		return
	}
	slog.Info("tracedesk", "version", version)
}
