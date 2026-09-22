# trace-desk

## Development setup

The eBPF programs use CO-RE and need a `vmlinux.h` header generated from the
running kernel's BTF. It's host-specific and git-ignored, so each developer
generates it locally.

Prerequisites (Linux):

- `clang` (compiles the BPF C sources)
- `bpftool` — e.g. `apt install linux-tools-common linux-tools-$(uname -r)`
- a kernel built with `CONFIG_DEBUG_INFO_BTF` (exposes `/sys/kernel/btf/vmlinux`)

Generate the header and Go bindings:

```sh
make headers   # writes internal/source/ebpf/vmlinux.h from kernel BTF
make gen       # regenerates eBPF Go bindings (runs `make headers` if missing)
```

`make gen` auto-runs `make headers` when the header is absent, so a fresh
checkout only needs `make gen`. The committed `*_bpfel.go`/`*_bpfeb.go`
bindings embed their `.o` objects, so a plain `make build` works without
regenerating.
