//go:build ignore
#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

struct event { __u32 pid; __u32 ppid; __u8 comm[16]; };

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1 << 24);
} events SEC(".maps");

SEC("tp/sched/sched_process_exec")
int handle_exec(void *ctx) {
	struct event *e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
	if (!e) return 0;
	e->pid = bpf_get_current_pid_tgid() >> 32;
	bpf_get_current_comm(&e->comm, sizeof(e->comm));
	bpf_ringbuf_submit(e, 0);
	return 0;
}
char _license[] SEC("license") = "Dual MIT/GPL";
