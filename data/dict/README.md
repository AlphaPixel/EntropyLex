# Dictionaries

This folder holds derived EntropyLex dictionary artifacts. It is currently empty of dictionaries — none have been built yet.

Dictionaries are **outputs**, produced by the derivation pipeline under `tools/`. They are not written or edited by hand. If a dictionary here does not match what the pipeline produces from its recorded sources and configuration, one of the two is wrong.

## Naming

```
entropylex-<lang>-<w>-v<version>.txt
```

For example `entropylex-en-14-v1.txt` for the English EL-14 dictionary. The language tag and the profile width are both part of the artifact's identity — see SPEC.md section 12.5 for scripts, and section 3.8 for why identity matters at decode time.

## Format

UTF-8, LF line endings, NFC normalized. A header block of `#`-prefixed lines carries the metadata; every line after it is one token, with its index implied by line order. See SPEC.md section 11.7 for the full definition.

The `sha256` field in the header is the **dictionary fingerprint**. Implementations should be able to report it, so that an encoder and decoder working from different dictionaries fail loudly rather than producing plausible garbage.

## Expected contents, once built

| File | Profile | Normal | Trim | Total |
|---|---|---|---|---|
| `entropylex-en-8-v1.txt`  | EL-8  | 256    | 0    | 256    |
| `entropylex-en-12-v1.txt` | EL-12 | 4,096  | 272  | 4,368  |
| `entropylex-en-14-v1.txt` | EL-14 | 16,384 | 5,460| 21,844 |

How these three relate to one another — nested subset, independent optimization, or disjoint partition — is **not yet decided**. See SPEC.md section 11.8.

## Verifier reports

`eldict-verify` output belongs alongside the artifact it describes, named to match (`entropylex-en-14-v1.report.txt`). A dictionary without a report should be treated as unverified.

## Reference

- SPEC.md section 9 — required token counts per profile
- SPEC.md section 10 — word selection criteria
- SPEC.md section 11 — derivation pipeline and artifact format
- `tools/README.md` — the tooling that produces these files
