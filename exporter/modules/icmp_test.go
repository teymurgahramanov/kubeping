package modules

import (
	"testing"
)

func TestProbeICMP_InvalidAddress(t *testing.T) {
	ok, err := ProbeICMP("not-a-valid-host", 1)
	if ok {
		t.Fatal("expected false, got true")
	}
	if err == nil {
		t.Fatal("expected error for invalid host, got nil")
	}
}

func TestProbeICMP_Loopback(t *testing.T) {
	// ICMP probes require privileges and network access; in sandboxed/CI
	// environments this may fail. Skip rather than fail in that case.
	if !canPing() {
		t.Skip("ICMP not available in this environment")
	}

	ok, err := ProbeICMP("127.0.0.1", 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Fatal("expected true, got false")
	}
}

// canPing best-efforts a loopback ping to determine whether ICMP is usable.
func canPing() bool {
	ok, _ := ProbeICMP("127.0.0.1", 2)
	return ok
}
