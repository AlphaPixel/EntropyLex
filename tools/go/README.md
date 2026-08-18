This folder will contain EntropyLex dictionary-selection and verification tools written in Go.

Each independently written implementation belongs in a subfolder named for its author, for example `satoshi-nakamoto`. The `src/` tree uses the same organization.

The planned sequence is `eldict-ingest`, `eldict-filter`, `eldict-score`, `eldict-select`, `eldict-emit`, `eldict-compile`, and `eldict-verify`. The main tools README explains each stage. An implementation must state which stages it provides and the exact source dataset versions used to test it.

Independently written implementations must produce the same mapping, fingerprint, and canonical LXJ v1 bytes from byte-for-byte identical source files and settings. Generated LXB files must also match byte for byte after that layout is defined. Comparing results checks both the implementations and whether the specification is precise enough.

See ../README.md for the stage descriptions and ../../SPEC.md section 11 for the requirements implementations must follow.
