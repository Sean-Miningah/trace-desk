package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/sean-miningah/trace-desk/internal/core"
	"github.com/sean-miningah/trace-desk/internal/pipeline"
	"github.com/sean-miningah/trace-desk/internal/source/fake"
	"github.com/sean-miningah/trace-desk/tui"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	switch os.Args[1] {
	case "run":
		if err := runCmd(os.Args[2:]); err != nil {
			slog.Error("run", "error", err)
			os.Exit(1)
		}
	case "version":
		slog.Info("tracedesk", "version", version)
		return
	case "tui":
		if err := tuiCmd(os.Args[2:]); err != nil {
			slog.Error("tui", "error", err)
			os.Exit(1)
		}
	case "":
	default:
		usage()
	}
}

func usage() {
	fmt.Printf("trace-desk %s\nusage:tracedesk <run|version> [flags]\n", version)
}

func runCmd(args []string) error {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	source := fs.String("source", "fake", "event source: fake|ebpf")
	dur := fs.Duration("duration", 5*time.Second, "how long to run")
	_ = fs.Parse(args)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	ctx, cancelTimeout := context.WithTimeout(ctx, *dur)
	defer cancelTimeout()

	var src core.EventSource
	switch *source {
	case "fake":
		src = fake.New(0)
	default:
		return fmt.Errorf("unknown source: %s", *source)
	}
	defer src.Close()

	raw, err := src.Events(ctx)
	if err != nil {
		return err
	}

	norm := pipeline.NewNormalizer()
	bus := pipeline.NewBus()
	store := bus.Subscribe("store", 1024, pipeline.Block)
	ui := bus.Subscribe("ui", 64, pipeline.DropNewest)

	//Normalize raw -> canonical events, feed the bus.
	normalized := make(chan core.Event, 1024)
	go func() {
		defer close(normalized)
		for r := range raw {
			normalized <- norm.Normalize(pipeline.Raw{
				Type:        r.Type,
				PID:         r.PID,
				PPID:        r.PPID,
				UID:         r.UID,
				GID:         r.GID,
				ProcessName: r.ProcessName,
				Data:        r.Data,
				When:        r.Timestamp,
			})
		}
	}()
	go bus.Run(ctx, normalized)

	// Drain subscribers: count what each saw
	var stored, shown uint64
	doneStore := drain(store, &stored)
	doneUI := drain(ui, &shown)
	<-doneStore
	<-doneUI

	secs := dur.Seconds()
	fmt.Printf("stored=%d (%.0f/s, drops=%d) ui=%d (drops=%d)\n", stored, float64(stored)/secs, bus.Drops("store"), shown, bus.Drops("ui"))
	return nil
}

func drain(ch <-chan core.Event, counter *uint64) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for range ch {
			*counter++
		}
	}()
	return done
}

func tuiCmd(args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	src := fake.New(50 * time.Millisecond)
	raw, err := src.Events(ctx)
	if err != nil {
		return err
	}
	norm := pipeline.NewNormalizer()
	bus := pipeline.NewBus()
	sub := bus.Subscribe("tui", 128, pipeline.DropNewest)

	normalized := make(chan core.Event, 256)
	go func() {
		defer close(normalized)
		for r := range raw {
			normalized <- norm.Normalize(pipeline.Raw{
				Type:        r.Type,
				PID:         r.PID,
				PPID:        r.PPID,
				UID:         r.UID,
				GID:         r.GID,
				ProcessName: r.ProcessName,
				Data:        r.Data,
				When:        r.Timestamp,
			})
		}
	}()
	go bus.Run(ctx, normalized)
	return tui.Run(sub, func() uint64 { return bus.Drops("tui") })
}
