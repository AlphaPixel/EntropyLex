# Dictionaries

This folder will hold EntropyLex dictionary files and their verification reports. It currently contains no completed dictionaries.

Dictionaries are outputs of the selection tools under `tools/`. Except for possible manual review of EL-8, they are not selected or edited by hand. Running the recorded tool version with the exact recorded source files and settings must reproduce the same ordered mapping. A mismatch means that the file, implementation, settings, or specification needs investigation.

## Naming

```text
entropylex-<lang>-<w>-v<version>.lxj
entropylex-<lang>-<w>-v<version>.lxb
```

`<lang>` is the dictionary's language code, `<w>` is its normal-token width in bits, and `<version>` is the dictionary release number.

For example, `entropylex-en-14-v1.lxj` is the required JSON form of version 1 of the English EL-14 dictionary. `entropylex-en-14-v1.lxb` is the optional binary form generated from that LXJ file.

## Format

**LXJ** (`.lxj`) is the authoritative, human-readable JSON file. **LXB** (`.lxb`) is generated from LXJ and is intended for implementations whose measurements show that JSON loading is too slow. Implementations must support LXJ. LXB remains optional while its version 1 byte layout is being measured and finalized.

Both forms carry the same dictionary fingerprint because they assign the same token to every index and use the same written-token rules. A checksum over every file byte differs because the JSON and binary bytes are different.

See [FORMAT.md](FORMAT.md) for the working design, the plain-language fingerprint definition, what a file can verify by itself, and the remaining format decisions.

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
- `SOURCES.md` — candidate input datasets, licenses, and counts
