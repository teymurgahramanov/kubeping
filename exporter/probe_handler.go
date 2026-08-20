package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// probeFunc is the signature shared by modules.ProbeTCP/ProbeHTTP/ProbeICMP.
type probeFunc func(address string, timeout int, insecureTLS bool) (bool, error)

// probeModules maps a module name to its probe implementation.
var probeModules = map[string]probeFunc{
	"tcp":  probeTCP,
	"http": probeHTTP,
	"icmp": probeICMP,
}

// Thin wrappers so tests can override probeModules without importing the
// modules package directly here.
func probeTCP(address string, timeout int, insecureTLS bool) (bool, error) {
	return modulesProbeTCP(address, timeout)
}
func probeHTTP(address string, timeout int, insecureTLS bool) (bool, error) {
	return modulesProbeHTTP(address, timeout, insecureTLS)
}
func probeICMP(address string, timeout int, insecureTLS bool) (bool, error) {
	return modulesProbeICMP(address, timeout)
}

// handleProbe returns an http.HandlerFunc that executes the probe selected by
// the request's "module" field and writes a JSON probeResponse.
func handleProbe(logger *slog.Logger, defaultTimeout int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var request probeRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		timeout := request.Timeout
		if timeout == 0 {
			timeout = defaultTimeout
		}

		probe, ok := probeModules[request.Module]
		if !ok {
			logger.Error("Unknown module")
			http.Error(w, "Unknown module", http.StatusBadRequest)
			return
		}

		result, probeErr := probe(request.Address, timeout, request.InsecureTLS)

		var response probeResponse
		if result {
			logger.Info("Probe successful")
			response.Result = true
		} else {
			if probeErr != nil {
				logger.Error("probe failed", slog.String("err", probeErr.Error()))
				response.Error = probeErr.Error()
			} else {
				response.Error = "probe failed"
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
