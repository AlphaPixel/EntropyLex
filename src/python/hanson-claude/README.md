This folder is reserved for a Python reference implementation of EntropyLex written by Claude Code. Its purpose is to test the specification and provide known behavior that independently written implementations can compare against. No code has been implemented yet.

Planned profile order is EL-8, then EL-12, then EL-14. All profiles will use the same encoding and decoding functions, with token width and related counts supplied as profile parameters.

The dictionary loader will first support the authoritative JSON format, LXJ (`.lxj`). Optional LXB (`.lxb`) support will follow only after its byte layout is final. Loading either form must produce the same dictionary fingerprint, index-to-token mapping, and written-token behavior.
