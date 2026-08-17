# Dictionaries

This folder holds derived EntropyLex dictionary artifacts. It is currently empty of dictionaries — none have been built yet.

Dictionaries are **outputs**, produced by the derivation pipeline under `tools/`. They are not written or edited by hand. If a dictionary here does not match what the pipeline produces from its recorded sources and configuration, one of the two is wrong.

## Naming

```text
entropylex-<lang>-<w>-v<version>.lxj
entropylex-<lang>-<w>-v<version>.lxb
```

For example, `entropylex-en-14-v1.lxj` is the canonical JSON representation of the English EL-14 dictionary and `entropylex-en-14-v1.lxb` is its optional compiled binary companion.

## Format

**LXJ** (`.lxj`) is the official, human-readable JSON file and the source of truth. **LXB** (`.lxb`) is generated from LXJ for faster loading. Implementations must support LXJ; LXB remains optional while its version 1 layout is being benchmarked and finalized.

Both representations carry the same dictionary fingerprint because they describe the same ordered mapping and decoding behavior. Their ordinary file checksums differ because JSON and binary files have different bytes.

See [FORMAT.md](FORMAT.md) for the working design, the plain-language fingerprint definition, what a file can verify by itself, and the remaining format decisions.

## Expected contents, once built

| File | Profile | Normal | Trim | Total |
|---|---|---|---|---|
| `entropylex-en-8-v1.lxj`  | EL-8  | 256    | 0    | 256    |
| `entropylex-en-12-v1.lxj` | EL-12 | 4,096  | 272  | 4,368  |
| `entropylex-en-14-v1.lxj` | EL-14 | 16,384 | 5,460| 21,844 |

How these three relate to one another — nested subset, independent optimization, or disjoint partition — is **not yet decided**. See SPEC.md section 11.8.

Each LXJ file may have a matching `.lxb` file with the same basename and dictionary fingerprint.

## Verifier reports

`eldict-verify` output belongs alongside the dictionary it describes, named to match (`entropylex-en-14-v1.report.json`). The report distinguishes structural checks performed from LXJ/LXB alone from quality checks that used the pinned derivation sources. A dictionary without a report should be treated as structurally unverified and as having no reproducible evidence for its word-selection quality.

## Candidate sources

[SOURCES.md](SOURCES.md) records the reviewed English and multilingual source candidates, their licenses, their published or measured sizes, and the conditions that must be resolved before selection. Actual derivation inputs must be pinned and checksummed in LXJ rather than inferred from this candidate list.

## Reference

- SPEC.md section 9 — required token counts per profile
- SPEC.md section 10 — word selection criteria
- SPEC.md section 11 — derivation pipeline and dictionary formats
- `tools/README.md` — the tooling that produces these files
- `FORMAT.md` — LXJ, LXB, fingerprints, and verification boundaries
- `SOURCES.md` — candidate input datasets, licenses, and counts
