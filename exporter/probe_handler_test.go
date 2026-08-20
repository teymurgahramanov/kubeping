package main

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func doProbe(t *testing.T, handler http.HandlerFunc, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/probe", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler(rec, req)
	return rec
}

func decodeProbeResponse(t *testing.T, body io.Reader) probeResponse {
	t.Helper()
	var resp probeResponse
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return resp
}

// --- Tests ---

func TestHandleProbe_Success(t *testing.T) {
	orig := probeModules["tcp"]
	defer func() { probeModules["tcp"] = orig }()
	probeModules["tcp"] = func(address string, timeout int, insecureTLS bool) (bool, error) {
		if address != "127.0.0.1:80" {
			t.Errorf("unexpected address: %s", address)
		}
		return true, nil
	}

	rec := doProbe(t, handleProbe(newTestLogger(), 5), `{"module":"tcp","address":"127.0.0.1:80"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	resp := decodeProbeResponse(t, rec.Body)
	if !resp.Result {
		t.Error("expected result true")
	}
	if resp.Error != "" {
		t.Errorf("expected empty error, got %q", resp.Error)
	}
}

func TestHandleProbe_ProbeFailure(t *testing.T) {
	orig := probeModules["http"]
	defer func() { probeModules["http"] = orig }()
	probeModules["http"] = func(address string, timeout int, insecureTLS bool) (bool, error) {
		return false, errors.New("connection refused")
	}

	rec := doProbe(t, handleProbe(newTestLogger(), 5), `{"module":"http","address":"http://example.com"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	resp := decodeProbeResponse(t, rec.Body)
	if resp.Result {
		t.Error("expected result false")
	}
	if resp.Error != "connection refused" {
		t.Errorf("expected error %q, got %q", "connection refused", resp.Error)
	}
}

func TestHandleProbe_ProbeFailureNoError(t *testing.T) {
	orig := probeModules["icmp"]
	defer func() { probeModules["icmp"] = orig }()
	probeModules["icmp"] = func(address string, timeout int, insecureTLS bool) (bool, error) {
		return false, nil
	}

	rec := doProbe(t, handleProbe(newTestLogger(), 5), `{"module":"icmp","address":"127.0.0.1"}`)

	resp := decodeProbeResponse(t, rec.Body)
	if resp.Result {
		t.Error("expected result false")
	}
	if resp.Error != "probe failed" {
		t.Errorf("expected error %q, got %q", "probe failed", resp.Error)
	}
}

func TestHandleProbe_UnknownModule(t *testing.T) {
	rec := doProbe(t, handleProbe(newTestLogger(), 5), `{"module":"foo","address":"127.0.0.1"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleProbe_InvalidJSON(t *testing.T) {
	rec := doProbe(t, handleProbe(newTestLogger(), 5), `{not-json}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleProbe_WrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/probe", nil)
	rec := httptest.NewRecorder()
	handleProbe(newTestLogger(), 5).ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHandleProbe_DefaultTimeoutApplied(t *testing.T) {
	orig := probeModules["tcp"]
	defer func() { probeModules["tcp"] = orig }()
	var capturedTimeout int
	probeModules["tcp"] = func(address string, timeout int, insecureTLS bool) (bool, error) {
		capturedTimeout = timeout
		return true, nil
	}

	// Request omits timeout; should fall back to default (5).
	doProbe(t, handleProbe(newTestLogger(), 5), `{"module":"tcp","address":"127.0.0.1:80"}`)

	if capturedTimeout != 5 {
		t.Errorf("expected default timeout 5, got %d", capturedTimeout)
	}
}

func TestHandleProbe_ExplicitTimeoutUsed(t *testing.T) {
	orig := probeModules["tcp"]
	defer func() { probeModules["tcp"] = orig }()
	var capturedTimeout int
	probeModules["tcp"] = func(address string, timeout int, insecureTLS bool) (bool, error) {
		capturedTimeout = timeout
		return true, nil
	}

	doProbe(t, handleProbe(newTestLogger(), 5), `{"module":"tcp","address":"127.0.0.1:80","timeout":10}`)

	if capturedTimeout != 10 {
		t.Errorf("expected explicit timeout 10, got %d", capturedTimeout)
	}
}
