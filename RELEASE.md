# KubePing 1.2.0

- 📡 **HTTP and ICMP probes:** Added HTTP and ICMP modules alongside TCP, selectable from a new module dropdown in the web UI.
- 🎯 **Target visibility:** The requested target is displayed as the results table heading.
- 🔐 **TLS validation:** HTTP probes now verify certificates and surface original SSL/TLS errors, including expired certificates, hostname mismatches, untrusted CAs, and handshake failures.
- 📜 **Trusted CA roots:** CA certificates are bundled into the exporter image so HTTPS probes trust roots such as Let's Encrypt / ISRG Root X1.
- ✨ **UI cleanup:** The app version moved from the footer to the header next to the KubePing logo. Header and footer dividers and bold target-heading styling were removed.
- 🐳 **Docker Hub namespace:** Web and exporter images are now published as `kubeping/kubeping-web` and `kubeping/kubeping-exporter` instead of under the per-user `github.repository` namespace.
- 🛠️ **Release workflow:** Fixed `APP_VERSION` updates in `config.py` by matching its actual file format.
- ✅ **Test coverage:** Added unit tests for the TCP, HTTP, and ICMP probe modules and the `/probe` HTTP handler.
