package modules

import (
	"net"
	"testing"
)

func TestProbeTCP_Success(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		conn.Close()
	}()

	ok, err := ProbeTCP(ln.Addr().String(), 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Fatal("expected true, got false")
	}
}

func TestProbeTCP_Failure(t *testing.T) {
	// Use a port that is very unlikely to be in use.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close()

	ok, err := ProbeTCP(addr, 1)
	if ok {
		t.Fatal("expected false, got true")
	}
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
