# KubePing 1.2.0

KubePing 1.2.0 adds HTTP and ICMP probe modules to the web UI, alongside the existing TCP probe.

A module selector in the probe form lets you choose between TCP, HTTP, and ICMP. The address input hint updates dynamically based on the selected module.

The requested target is now displayed as the heading of the results table.

The exporter already supported all three modules; the web UI now exposes them instead of TCP only.
