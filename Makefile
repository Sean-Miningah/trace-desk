# TraceDesk -- build/test/bench/lint/gen
BIN := bin
GOFLAGS ?=  -trimpath
LDFLAGS ?= -s -w

EBPF_DIR := internal/source/ebpf
VMLINUX  := $(EBPF_DIR)/vmlinux.h

.PHONY: build test bench lint gen headers run clean

build:  ## no-CGO agent + CLI
	CGO_ENABLED=0 go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BIN)/tracedesk ./cmd/tracedesk

query:  ## CGO duckdb query binary
	CGO_ENABLED=1 go build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BIN)/tracedesk-query ./cmd/tracedesk-query

test:
	go test $(GOFLAGS) ./...

bench:
	go test $(GOFLAGS) -bench . ./...

lint:
	go vet ./...

headers: $(VMLINUX)  ## regenerate CO-RE vmlinux.h from the running kernel's BTF

$(VMLINUX):
	@command -v bpftool >/dev/null 2>&1 || { echo "bpftool not found; install it (e.g. apt install linux-tools-common linux-tools-$$(uname -r), or bpftool)"; exit 1; }
	@test -r /sys/kernel/btf/vmlinux || { echo "kernel BTF missing at /sys/kernel/btf/vmlinux; requires a CONFIG_DEBUG_INFO_BTF kernel"; exit 1; }
	bpftool btf dump file /sys/kernel/btf/vmlinux format c > $@

gen: headers ## regenerate eBPF Go bindings (auto-generates vmlinux.h if missing)
	go generate ./...

run: build
	./$(BIN)/tracedesk run --source=fake --duration=5s

clean:
	rm -rf $(BIN)
