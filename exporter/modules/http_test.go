package modules

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProbeHTTP_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ok, err := ProbeHTTP(srv.URL, 2, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Fatal("expected true, got false")
	}
}

func TestProbeHTTP_NonOKStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	ok, err := ProbeHTTP(srv.URL, 2, false)
	if ok {
		t.Fatal("expected false, got true")
	}
	if err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}
}

func TestProbeHTTP_Unreachable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	addr := srv.URL
	srv.Close()

	ok, err := ProbeHTTP(addr, 1, false)
	if ok {
		t.Fatal("expected false, got true")
	}
	if err == nil {
		t.Fatal("expected error for unreachable server, got nil")
	}
}
