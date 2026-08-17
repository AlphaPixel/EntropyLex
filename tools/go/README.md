This folder will contain EntropyLex tooling implemented in the language denoted by the folder name (e.g. "python").

Each implementation will be in a subfolder, named for the author of the implementation. E.g. ("satoshi-nakamoto"). This mirrors the convention used under `src/`.

The primary tooling target is the dictionary derivation pipeline (`eldict-ingest`, `eldict-filter`, `eldict-score`, `eldict-select`, `eldict-emit`, `eldict-verify`). Each implementation should state which pipeline stages it covers, and which source corpora and versions it was validated against.

Independent implementations must produce byte-identical dictionary artifacts from the same pinned inputs and configuration. That equivalence is the point of allowing more than one.

See ../README.md for the pipeline overview and ../../SPEC.md section 11 for the normative requirements.
