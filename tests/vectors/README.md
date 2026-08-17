# Conformance vectors

Language-agnostic test data that every EntropyLex implementation must pass.

## Two classes of vector

**Dictionary-independent vectors** express the expected result as a sequence of **token indices** rather than words. Because the bitstream packing, remainder arithmetic, trim classification, and validation rules are all defined over indices (SPEC.md sections 4.3, 5.1, 6, 7), these vectors test the entire encoding mechanism without reference to any particular dictionary.

This matters for sequencing: **the full conformance suite for the encoding logic can be written now, before a single word has been chosen.** Dictionary derivation and implementation are independent tracks, and neither blocks the other.

**Dictionary-dependent vectors** express the expected result as an actual phrase, and are pinned to a specific dictionary fingerprint. They test dictionary loading, canonical form, and input normalization (SPEC.md sections 5.2, 5.3, 11.7). These cannot be authored until dictionaries exist.

## Required coverage

### Remainder classes

Every reachable remainder must be exercised in every profile. There is a convenient result here:

> **Payloads of N = 0 through 7 bytes exercise every remainder class in every defined profile.**

EL-14's remainder cycles with `N mod 7` — 0, 8, 2, 10, 4, 12, 6, back to 0 — so eight payloads of ascending length cover all seven classes. EL-12 (cycle length 3) and EL-16 (cycle length 2) are subsumed. Eight small payloads give complete trim coverage across the whole family.

Worked example, payload bytes `a5 5a 3c c3 0f f0 99` truncated to length N:

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

Note the N = 7 EL-14 row: four tokens, all normal, no trim token. That is the `r = 0` case, and it is the one where a phrase legitimately ends on a normal token.

### Negative vectors

Validation rules D1 through D5 (SPEC.md section 7.3) each need at least one sequence that must be **rejected**:

- A token index outside the dictionary
- A trim token in a non-final position (D2)
- A token sequence whose total bit length is not divisible by 8 (D3) — for EL-14, any run of normal tokens whose count is not a multiple of 4
- A trim token claiming a remainder outside the profile's remainder set (D4)

An implementation that decodes these without error is not conformant. Silent acceptance of malformed input is the failure mode this format can least afford.

### Round-trip

Encode-then-decode must reproduce the input exactly, for every vector, in every supported profile. Decode-then-encode must reproduce the canonical phrase exactly.

## Format

JSON, one file per vector set, UTF-8, LF endings. Proposed shape **(provisional — settle when the first set is authored)**:

```json
{
  "vector_set": "core-roundtrip",
  "version": 1,
  "dictionary_independent": true,
  "vectors": [
    {
      "id": "el14-n02",
      "profile": "EL-14",
      "payload_hex": "a55a",
      "indices": [10582, 16386],
      "remainder": 2,
      "note": "r=2 trim class"
    }
  ]
}
```

Dictionary-dependent sets add `"dictionary": "entropylex-en-14-v1"` and `"dictionary_sha256": "..."` at the top level, and `"phrase"` per vector.

## Provenance

The index values above were computed directly from the specification (MSB-first packing, unified index space per section 4.3) and independently verified. Any implementation disagreeing with them has found either a bug in itself or an ambiguity in the spec — both are worth reporting.
