# CHANGELOG

## [1.2.0]

### Added
- HTTP and ICMP probe modules selectable in the web UI alongside TCP
- CA certificates bundled into the exporter image so HTTPS probes validate against trusted roots (e.g. Let's Encrypt / ISRG Root X1)
- Module selector dropdown in the probe form
- Requested target displayed as the heading of the results table

### Changed
- Web and exporter container images are now published to the `kubeping` Docker Hub namespace (`kubeping/kubeping-web`, `kubeping/kubeping-exporter`) instead of the per-user `github.repository` namespace
- Address input hint now updates dynamically based on the selected module
- Web `/ping` endpoint reads the module from the form instead of hardcoding TCP
- HTTP probe now verifies TLS certificates instead of skipping validation, surfacing all SSL/TLS errors (expired, hostname mismatch, untrusted CA, handshake failures) with their original error messages
- App version moved from the footer to the header, displayed next to the KubePing logo
- Removed line dividers from the header and footer
- Removed bold styling from the target heading above the results table

### Fixed
- Release workflow now correctly bumps `APP_VERSION` in `config.py` (sed pattern updated to match the actual file format)

## [1.1.1] - 2026-06-27

### Changed
- Helm image tag defaults are now empty in `values.yaml`, so web and exporter images use the chart `appVersion` unless explicitly overridden
- Helm ingress TLS configuration now uses `web.ingress.host` as the single source for both ingress rules and TLS hosts
- Helm ingress TLS values now use `web.ingress.tls.enabled` and `web.ingress.tls.secretName`

## [1.1.0]

### Added
- Dark and light theme support with a toggle in the header
- Outbound public IP display in the header, fetched from a configurable source (defaults to ifconfig.me)
- `PUBLIC_IP_URL` setting in `config.py`, configurable via environment variable
- ICMP probe module using pro-bing
- Default liveness and readiness probes for both web and exporter in Helm chart
- `web.replicaCount` Helm value (templated in deployment)
- ConfigMap checksum annotation to DaemonSet for automatic rollout on config changes
- TLS support to Ingress
- Configurable service port via `service.port` for both components
- Standard labels to Role and RoleBinding
- `NOTES.txt` with post-install instructions

### Changed
- Exporter metric name changed from `probe_result` to `kubeping_probe_result`
- Redesigned web UI with a cleaner layout, dark and light theme support
- Centered form on the page with larger inputs
- Renamed `/submit` endpoint to `/ping`
- Exporter port is now read dynamically from pod container spec instead of being hardcoded
- Renamed Helm container names to `kubeping-web` and `kubeping-exporter` to avoid collision

### Fixed
- URL staying on `/ping` after form submission using Post/Redirect/Get pattern
- Helm `volumes`/`volumeMounts` defaults changed from `{}` to `[]`
- Helm `nodePort` default changed from empty string to `null`

## [1.0.0]
