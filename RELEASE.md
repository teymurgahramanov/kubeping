# KubePing 1.2.1

KubePing 1.2.1 makes the HTTP probe sensitive to SSL/TLS certificate errors.

The HTTP probe previously skipped TLS certificate verification (`InsecureSkipVerify`), so certificate problems were never reported. Verification is now enabled, and all SSL/TLS errors — expired certificates, hostname mismatches, untrusted or unknown certificate authorities, and handshake failures — are surfaced with their original Go TLS error messages.