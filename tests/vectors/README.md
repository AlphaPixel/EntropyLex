# Shared specification test cases

These JSON test cases specify exact inputs and expected outputs. Every implementation must pass the cases for the profiles and file formats it claims to support.

## Two kinds of test case

**Dictionary-independent cases** express the expected result as numerical token indices rather than words. The specification defines bit grouping, the number of remaining bits, normal-versus-trim index ranges, and validation in terms of these indices (SPEC.md sections 4.3, 5.1, 6, and 7). These cases can therefore test the encoding algorithm before words have been selected.

The encoding test cases can be written before a dictionary exists. Dictionary selection and implementation can therefore proceed independently.

**Dictionary-dependent cases** include actual phrases and name the exact dictionary fingerprint they require. They test file loading, the one written form an encoder must produce, and the permitted variations a decoder accepts (SPEC.md sections 5.2, 5.3, and 11.7). Each case must pass when the dictionary is loaded from LXJ or from its matching LXB. These cases cannot be completed until dictionaries exist.

## Required coverage

### Remainder groups

Every possible remaining-bit length must appear in the tests for each profile:

> **Payloads of N = 0 through 7 bytes exercise every remainder group in every defined profile.**

For EL-14, the remainder repeats according to `N mod 7`, where `mod` means the remainder after division. For payload lengths 0 through 7, the sequence is 0, 8, 2, 10, 4, 12, 6, then 0 again. EL-12 repeats every three byte lengths and EL-16 every two, so their complete cycles also occur within 0 through 7. These eight lengths therefore cover every trim group in every defined profile.

The following example uses the byte sequence `a5 5a 3c c3 0f f0 99`. Row `N` uses only its first `N` bytes:

| N | Payload | `r` | Expected EL-14 indices |
|---|---|---|---|
| 0 | *(empty)* | 0  | `[]` |
| 1 | `a5` | 8  | `[16633]` |
| 2 | `a55a` | 2  | `[10582, 16386]` |
| 3 | `a55a3c` | 10 | `[10582, 17296]` |
| 4 | `a55a3cc3` | 4  | `[10582, 9164, 16391]` |
| 5 | `a55a3cc30f` | 12 | `[10582, 9164, 18531]` |
| 6 | `a55a3cc30ff0` | 6  | `[10582, 9164, 3135, 16452]` |
| 7 | `a55a3cc30ff099` | 0  | `[10582, 9164, 3135, 12441]` |

The same payloads under EL-12 and EL-8:

| N | EL-12 `r` | EL-12 indices | EL-8 indices |
|---|---|---|---|
| 0 | 0 | `[]` | `[]` |
| 1 | 8 | `[4277]` | `[165]` |
| 2 | 4 | `[2645, 4106]` | `[165, 90]` |
| 3 | 0 | `[2645, 2620]` | `[165, 90, 60]` |

For `N = 7` under EL-14, all four tokens are normal because no bits remain after the final complete 14-bit group. This is a required successful case; a phrase does not always end with a trim token.

### Required rejection and acceptance cases

Rules D1 through D4 in SPEC.md section 7.3 describe invalid conditions. Each needs at least one sequence that the decoder must reject:

- A token index outside the dictionary (D1)
- A trim token in a non-final position (D2)
- A token sequence whose total bit length is not divisible by 8 (D3) — for EL-14, any run of normal tokens whose count is not a multiple of 4
- A trim token claiming a remainder outside the profile's remainder set (D4)

Rule D5 is different: it requires an empty token sequence to be accepted and decoded as an empty payload. It therefore needs a successful test rather than a rejection test.

An implementation that accepts a D1–D4 case or rejects the D5 empty case does not conform to the specification.

### Exact conversion in both directions

For every successful test, encoding followed by decoding must reproduce the original bytes exactly. Decoding followed by encoding must produce the canonical phrase exactly, even when the input phrase used another permitted spacing or letter-case form.

## Format

Use JSON with one file per test set, UTF-8 text encoding, and line-feed (`LF`) characters to end lines. The proposed structure is **provisional** and must be settled when the first set is written:

```json
{
  "test_set": "core-encode-decode",
  "version": 1,
  "dictionary_independent": true,
  "cases": [
    {
      "id": "el14-n02",
      "profile": "EL-14",
      "payload_hex": "a55a",
      "indices": [10582, 16386],
      "remainder": 2,
      "note": "r=2 trim group"
    }
  ]
}
```

`payload_hex` writes each payload byte as two hexadecimal digits, so `a55a` represents the two bytes `a5 5a`. `indices` is the expected ordered list of numerical token indices.

Dictionary-dependent sets add `"dictionary": "entropylex-en-14-v1"`, `"dictionary_fingerprint": "sha256:..."`, and `"fingerprint_recipe": "..."` once for the entire set, plus `"phrase"` in each case. The fingerprint recipe identifies the version of the calculation rules. `dictionary_fingerprint` identifies the index mapping and written-token rules; it is not a checksum over all bytes of either the LXJ or LXB file.

Separate deliberately invalid files are required for both formats. LXJ cases cover missing or wrongly typed fields, incorrect counts, invalid UTF-8 or Unicode normalization, duplicate tokens, incorrect order, and fingerprint mismatch. LXB cases also cover files that end early and section lengths or starting positions that point outside the file. For every valid LXJ/LXB pair, tests must compare the token and behavior at every index.

## How the example indices were obtained

The index values above were calculated from the specification's highest-value-bit-first order and the continuous index ranges in section 4.3, then checked with a separate calculation. A disagreement may reveal an implementation error or an ambiguous rule in the specification; report enough detail to determine which.
