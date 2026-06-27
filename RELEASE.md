# KubePing 1.1.1

Released: 2026-06-27

KubePing 1.1.1 refines the Helm chart defaults for easier version management and simpler ingress TLS configuration.

Image tags in `values.yaml` are now empty by default. When no explicit tag is provided, the Helm templates use the chart `appVersion` for both the web and exporter images.

Ingress TLS now uses `web.ingress.host` as the single host source. Enable TLS with `web.ingress.tls.enabled` and set the certificate secret with `web.ingress.tls.secretName`.

Upgrade note: if you previously configured `web.ingress.tls` as a list with its own `hosts`, move the host value to `web.ingress.host` and configure TLS with the new `enabled` and `secretName` fields.
