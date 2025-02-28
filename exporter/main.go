package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/teymurgahramanov/KubePing/exporter/modules"
	"gopkg.in/yaml.v3"
)

type configuration struct {
	Targets  map[string]targetConfig `yaml:"targets"`
	Exporter exporterConfig          `yaml:"exporter"`
}

type targetConfig struct {
	Address  string `yaml:"address"`
	Module   string `yaml:"module"`
	Interval int    `yaml:"interval"`
	Timeout  int    `yaml:"timeout"`
}

type exporterConfig struct {
	ListenPort           int `yaml:"listenPort"`
	DefaultProbeInterval int `yaml:"defaultProbeInterval"`
	DefaultProbeTimeout  int `yaml:"defaultProbeTimeout"`
}

type probeRequest struct {
	Module  string `json:"module"`
	Address string `json:"address"`
	Timeout int    `json:"timeout"`
}

type probeResponse struct {
	Result bool   `json:"result"`
	Error  string `json:"error,omitempty"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	configFile := "config.yaml"
	var config configuration

	// **Ensure Targets map is always initialized**
	config.Targets = make(map[string]targetConfig)

	// Set default exporter values
	config.Exporter.ListenPort = 8000
	config.Exporter.DefaultProbeInterval = 30
	config.Exporter.DefaultProbeTimeout = 5

	// Try reading the config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Error(fmt.Sprintf("Config file %s not found, proceeding with defaults", configFile))
		} else {
			logger.Error(fmt.Sprintf("Failed to read config file %s: %v", configFile, err))
		}
	} else {
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to parse config file %s: %v", configFile, err))
		}
	}

	// Ensure default values are applied if missing in the config file
	if config.Exporter.ListenPort == 0 {
		config.Exporter.ListenPort = 8000
	}
	if config.Exporter.DefaultProbeInterval == 0 {
		config.Exporter.DefaultProbeInterval = 30
	}
	if config.Exporter.DefaultProbeTimeout == 0 {
		config.Exporter.DefaultProbeTimeout = 5
	}

	// **Log message to indicate program continues execution**
	logger.Info("Starting server with final configuration", slog.Any("config", config))

	// Prometheus metric setup
	var (
		probeResult = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "probe_result",
				Help: "Current status of the probe (1 for success, 0 for failure)",
			},
			[]string{"target", "module", "address"},
		)
	)

	promRegistry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = promRegistry
	prometheus.DefaultGatherer = promRegistry
	prometheus.MustRegister(probeResult)

	// HTTP handlers
	http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	http.HandleFunc("/probe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var request probeRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		timeout := request.Timeout
		if timeout == 0 {
			timeout = config.Exporter.DefaultProbeTimeout
		}

		resultHandler := func(result bool, err error) {
			var response probeResponse
			response.Result = false
			if result {
				logger.Info("Probe successful")
				response.Result = true
			} else {
				if err != nil {
					logger.Error(fmt.Sprint(err.Error()))
				}
				response.Error = err.Error()
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}

		switch request.Module {
		case "tcp":
			result, err := modules.ProbeTCP(request.Address, timeout)
			resultHandler(result, err)
		case "http":
			result, err := modules.ProbeHTTP(request.Address, timeout)
			resultHandler(result, err)
		case "icmp":
			result, err := modules.ProbeICMP(request.Address)
			resultHandler(result, err)
		default:
			logger.Error("Unknown module")
			http.Error(w, "Unknown module", http.StatusBadRequest)
			return
		}
	})

	// Start HTTP server
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", config.Exporter.ListenPort), nil)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to start HTTP server: %v", err))
			os.Exit(1)
		}
	}()

	var wg sync.WaitGroup

	// **If no targets exist, ensure the program does not deadlock**
	if len(config.Targets) == 0 {
		logger.Info("No targets configured, service running with only HTTP API")
		select {} // Keeps the program running if no targets exist
	}

	// Start probes if any targets are defined
	for key, value := range config.Targets {
		wg.Add(1)
		go func(target string, module string, address string, interval int, timeout int) {
			defer wg.Done()
			if interval == 0 {
				interval = config.Exporter.DefaultProbeInterval
			}
			if timeout == 0 {
				timeout = config.Exporter.DefaultProbeTimeout
			}
			targetLogger := logger.With(slog.String("target", target))
			resultHandler := func(result bool, err error, interval int) {
				if result {
					targetLogger.Info("Probe successful")
					probeResult.WithLabelValues(target, module, address).Set(1)
				} else {
					if err != nil {
						targetLogger.Error(fmt.Sprint(err.Error()))
					}
					probeResult.WithLabelValues(target, module, address).Set(0)
				}
				time.Sleep(time.Duration(interval) * time.Second)
			}
			switch module {
			case "tcp":
				for {
					result, err := modules.ProbeTCP(address, timeout)
					resultHandler(result, err, interval)
				}
			case "http":
				for {
					result, err := modules.ProbeHTTP(address, timeout)
					resultHandler(result, err, interval)
				}
			case "icmp":
				for {
					result, err := modules.ProbeICMP(address)
					resultHandler(result, err, interval)
				}
			default:
				targetLogger.Error("Unknown module")
			}
		}(key, value.Module, value.Address, value.Interval, value.Timeout)
	}

	wg.Wait()
}
