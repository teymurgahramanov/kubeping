# KubePing 1.2.0

KubePing 1.2.0 adds HTTP and ICMP probe modules selectable in the web UI alongside TCP, with a module selector dropdown and a dynamically updating address input hint. The requested target is displayed as the heading of the results table.

The HTTP probe now verifies TLS certificates instead of skipping validation, surfacing all SSL/TLS errors (expired, hostname mismatch, untrusted CA, handshake failures) with their original error messages.

The app version is now shown in the header next to the KubePing logo instead of the footer. The header and footer line dividers have been removed, and the bold styling on the target heading above the results table has been removed.

This release also fixes the release workflow to correctly bump `APP_VERSION` in `config.py`.