# EntropyLex Tools

This folder holds the preprocessing and support tooling for EntropyLex — principally the **dictionary derivation pipeline** described in SPEC.md section 11. None of it exists yet; the subfolders are placeholders that establish the layout.

Dictionaries are not written by hand. They are derived from a master word list by a deterministic, re-runnable pipeline, so that the selection criteria are auditable and the artifacts are reproducible. Twenty-one thousand words chosen by human judgment would be neither.

## Layout

Tools follow the same convention as `src/`:

```
tools/<language>/<implementor-name>/
```

- The language folder denotes the implementation language (e.g. `python`).
- Each implementation lives in a subfolder named for its author (e.g. `satoshi-nakamoto`), exactly as in `src/`.

Multiple independent implementations of the same tool are expected and welcome. They are a cross-check: two derivation pipelines run against the same pinned inputs and configuration must produce the same ordered mapping and dictionary fingerprint. The canonical LXJ formatter and LXB compiler must also produce byte-identical files after their formats are finalized. A disagreement means that either the specification is incomplete or one implementation is wrong, and both outcomes are worth finding.

## The pipeline

| Stage | Tool | Responsibility |
|---|---|---|
| 1 | `eldict-ingest`  | Load master word list plus frequency, pronunciation, and part-of-speech data into a normalized candidate table |
| 2 | `eldict-filter`  | Apply hard disqualifiers; emit surviving candidates with rejection reasons logged |
| 3 | `eldict-score`   | Compute orthographic, phonetic, and semantic feature vectors for each candidate |
| 4 | `eldict-select`  | Build the confusability graph and select the final token set to a target count |
| 5 | `eldict-emit`    | Assign canonical indices, partition normal versus trim, and write the canonical LXJ file |
| 6 | `eldict-compile` | Compile LXJ into optional LXB without changing the mapping or fingerprint |
| 7 | `eldict-verify`  | Verify structure from dictionary files and quality/reproducibility from pinned inputs; emit a JSON report |

Every stage is deterministic given its inputs and a recorded configuration file. Source datasets must be pinned to immutable versions, checksummed, and recorded as structured source objects in LXJ. The reviewed source candidates and their license conditions are listed in `../data/dict/SOURCES.md`.

Stage 3 is deliberately metric-pluggable: a Chinese dictionary swaps the distance families for visual and romanization metrics and leaves stages 1, 2, 4, 5, 6, and 7 untouched. See SPEC.md section 12.4.

## Outputs

Derived LXJ, LXB, and verification reports belong in `data/dict/`, not here. Tools are inputs to the repository; dictionaries are outputs. LXJ is authoritative and LXB is reproducibly generated from it.

## Open design question

The composition of the EL-8, EL-12, and EL-14 dictionaries relative to one another is **not yet decided** — nested subset, independent optimization, or disjoint partition. That choice changes what stages 4 and 5 must do, so it should be settled before `eldict-select` is written. See SPEC.md section 11.8.

## Reference

- SPEC.md section 10 — word selection criteria
- SPEC.md section 11 — dictionary derivation, stage by stage
- SPEC.md section 11.7 — LXJ/LXB formats, provenance, and verification boundaries
- SPEC.md section 11.8 — the undecided cross-profile composition question
- SPEC.md section 12 — non-English and ideographic dictionaries
- `../data/dict/FORMAT.md` — working file-format design
- `../data/dict/SOURCES.md` — candidate datasets, licenses, and counts
