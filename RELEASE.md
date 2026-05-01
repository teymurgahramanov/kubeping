# KubePing 1.1.0

Kubeping 1.1.0 introduces a redesigned web interface with a cleaner layout, and minor tweaks improving overall usability. Deployment has also been enhanced with more flexible Helm configuration, including configurable ports and TLS support.

❗**Breaking change**: Prometheus metric name has been updated from `probe_result` to `kubeping_probe_result`. Update your dashboards and alerts accordingly.