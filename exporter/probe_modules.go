package main

import (
	"github.com/teymurgahramanov/KubePing/exporter/modules"
)

// modulesProbeTCP/HTTP/ICMP delegate to the real implementations in the
// modules package. They are package-level vars so tests can swap them out.
var (
	modulesProbeTCP  = modules.ProbeTCP
	modulesProbeHTTP = modules.ProbeHTTP
	modulesProbeICMP = modules.ProbeICMP
)
