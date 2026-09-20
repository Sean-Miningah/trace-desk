# TraceDesk -- build/test/bench/lint/gen
BIN := bin
GOFLAGS ?=  -trimpath
LDFLAGS ?= -s -w

.PHONY: build test bench lint gen run clean

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

gen: ## regenerate eBF Go bindings
	go generate ./...

run: build
	./$(BIN)/tracedesk run --source=fake --duration=5s

clean:
	rm -rf $(BIN)
