# Dictionaries

This folder will hold production EntropyLex dictionary files and their verification reports. It currently contains no completed production dictionary. A complete test-only EL-8 fixture lives under [`../../tests/fixtures/dict/`](../../tests/fixtures/dict/) so early implementations have a stable file without implying that its words passed production review.

Dictionaries are outputs of the selection tools under `tools/`. Except for possible manual review of EL-8, they are not selected or edited by hand. Running the recorded tool version with the exact recorded source files and settings must reproduce the same ordered mapping. A mismatch means that the file, implementation, settings, or specification needs investigation.

## Naming

```text
entropylex-<lang>-<w>-v<version>.lxj
entropylex-<lang>-<w>-v<version>.lxb
```

`<lang>` is the dictionary's language code, `<w>` is its normal-token width in bits, and `<version>` is the dictionary release number.

For example, `entropylex-en-14-v1.lxj` is the required JSON form of version 1 of the English EL-14 dictionary. `entropylex-en-14-v1.lxb` is the optional binary form generated from that LXJ file.

## Format

**LXJ** (`.lxj`) is the authoritative, human-readable JSON file. LXJ version 1 now has exact fields, validation rules, canonical writing rules, and fingerprint recipe LXFP-1. **LXB** (`.lxb`) will be generated from LXJ and is intended for implementations whose measurements show that JSON loading is too slow. Implementations must support LXJ. LXB remains optional while its version 1 byte layout is measured and defined.

Both forms carry the same dictionary fingerprint because they assign the same token to every index and use the same written-token rules. A checksum over every file byte differs because the JSON and binary bytes are different.

See [FORMAT.md](FORMAT.md) for the normative LXJ v1 definition, the plain-language fingerprint explanation, the exact LXFP-1 byte recipe, and the remaining LXB decisions. [`lxj-v1.schema.json`](lxj-v1.schema.json) provides the machine-readable structural rules.

## Expected contents, once built

| File | Profile | Normal | Trim | Total |
|---|---|---|---|---|
| `entropylex-en-8-v1.lxj`  | EL-8  | 256    | 0    | 256    |
| `entropylex-en-12-v1.lxj` | EL-12 | 4,096  | 272  | 4,368  |
| `entropylex-en-14-v1.lxj` | EL-14 | 16,384 | 5,460| 21,844 |

Whether these profiles share words is **not yet decided**. The alternatives are nested sets, separate selection with uncontrolled overlap, and completely separate token sets. See SPEC.md section 11.8.

Each LXJ file may have a matching `.lxb` file with the same filename before the extension and the same dictionary fingerprint.

## Verifier reports

The JSON report from `eldict-verify` belongs next to the dictionary and uses the same base name, for example `entropylex-en-14-v1.report.json`. It separates checks possible from the dictionary file alone—such as counts, token validity, and fingerprint—from checks that require the exact source datasets and selection settings. Without a report, neither the file's structure nor the evidence supporting its word quality has been independently recorded.

## Candidate sources

[SOURCES.md](SOURCES.md) records reviewed source candidates for English and other languages, their license terms, their published or measured sizes, and unresolved conditions. A completed LXJ file must identify the exact source release or repository commit and the downloaded file's checksum; readers must not infer actual inputs from the candidate list.

## Reference

- SPEC.md section 9 — required token counts per profile
- SPEC.md section 10 — word selection criteria
- SPEC.md section 11 — dictionary selection, file generation, source history, and verification
- `tools/README.md` — the tooling that produces these files
- `FORMAT.md` — LXJ, LXB, fingerprints, and what each verification type can establish
- `lxj-v1.schema.json` — machine-readable LXJ v1 structural validation
- `SOURCES.md` — candidate input datasets, licenses, and counts
