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

Multiple independent implementations of the same tool are expected and welcome. They are a cross-check: two derivation pipelines written from this spec, run against the same pinned inputs and the same configuration, must produce **byte-identical** dictionary artifacts. If they do not, either the spec is underspecified or one implementation is wrong, and both outcomes are worth finding out.

## The pipeline

| Stage | Tool | Responsibility |
|---|---|---|
| 1 | `eldict-ingest`  | Load master word list plus frequency, pronunciation, and part-of-speech data into a normalized candidate table |
| 2 | `eldict-filter`  | Apply hard disqualifiers; emit surviving candidates with rejection reasons logged |
| 3 | `eldict-score`   | Compute orthographic, phonetic, and semantic feature vectors for each candidate |
| 4 | `eldict-select`  | Build the confusability graph and select the final token set to a target count |
| 5 | `eldict-emit`    | Assign canonical indices, partition normal versus trim, write the dictionary artifact |
| 6 | `eldict-verify`  | Re-derive metrics from the artifact, assert all invariants, emit a quality report |

Every stage is deterministic given its inputs and a recorded configuration file. Source corpora must be pinned to a specific version and recorded in the artifact's provenance header.

Stage 3 is deliberately metric-pluggable: a Chinese dictionary swaps the distance families for visual and romanization metrics and leaves stages 1, 2, 4, 5, and 6 untouched. See SPEC.md section 12.4.

## Outputs

Derived dictionary artifacts belong in `data/dict/`, not here. Tools are inputs to the repository; dictionaries are outputs.

## Open design question

The composition of the EL-8, EL-12, and EL-14 dictionaries relative to one another is **not yet decided** — nested subset, independent optimization, or disjoint partition. That choice changes what stages 4 and 5 must do, so it should be settled before `eldict-select` is written. See SPEC.md section 11.8.

## Reference

- SPEC.md section 10 — word selection criteria
- SPEC.md section 11 — dictionary derivation, stage by stage
- SPEC.md section 11.7 — artifact format, provenance, and the verifier's required assertions
- SPEC.md section 11.8 — the undecided cross-profile composition question
- SPEC.md section 12 — non-English and ideographic dictionaries
