# KubePing 1.2.0

- 🌐 **HTTP and ICMP probes:** Added HTTP and ICMP probe modules alongside TCP. Select the desired module from the new dropdown in the web UI.
- ✨ **Refined web UI:** Improved the interface for a cleaner and more polished experience.
- 🔐 **Optional TLS validation:** HTTP probes can verify TLS certificates and report the original error, including expired certificates, hostname mismatches, untrusted certificate authorities, and handshake failures. Validation can also be disabled when required.
- 🐳 **Unified Docker Hub repository:** Web and exporter images are now published to `teymurgahramanov/kubeping`, using the `web-<version>` and `exporter-<version>` tag formats.
- 🚀 **Go upgrade:** Exporter upgraded to Go 1.27.0 (from 1.21.6).
- ✅ **Test coverage:** Added unit tests for the TCP, HTTP, and ICMP probe modules and the `/probe` HTTP handler.

For the complete list of changes, see the [changelog](./CHANGELOG.md).
