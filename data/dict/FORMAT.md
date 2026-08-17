# EntropyLex dictionary files

EntropyLex uses two representations of the same dictionary:

- **LXJ** (`.lxj`) is the official, human-readable JSON representation.
- **LXB** (`.lxb`) is an optional, faster-loading binary representation generated from LXJ.

LXJ is the source of truth. An LXB file must be reproducible from its LXJ file and must not contain a different mapping or different decoding rules. Implementations are required to support LXJ. LXB support is optional until its byte layout is finalized.

## One dictionary, two files

For example:

```text
entropylex-en-14-v1.lxj
entropylex-en-14-v1.lxb
```

The two files have different bytes and therefore different ordinary file checksums. They carry the same **dictionary fingerprint**, because they describe the same token mapping and decoding behavior.

## LXJ contents

The final JSON schema is not yet settled. At minimum, an LXJ file must record:

- LXJ format version and dictionary release version
- Dictionary name, profile, and token width
- Language and script labels
- Rules that affect recognition of a written phrase: Unicode normalization, case handling, and token segmentation
- Index-assignment scheme and reachable trim sizes
- Normal, trim, and total token counts
- One flat token array, in index order
- Dictionary fingerprint and the version of the fingerprint recipe
- Structured source records and the selection-configuration checksum

Array position is the token index. Entries must not repeat their own indices: duplicated index fields would add size and create another value that could disagree with array order.

Language and script labels are useful identification and provenance. They affect the fingerprint only when they select actual parsing or decoding behavior. The behavior itself is what must be identified.

Source records must be machine-readable objects rather than a single prose string. Each selected source should record its role, immutable version or commit, download location, license, downloaded-file SHA-256 checksum, and the relevant record count.

## Dictionary fingerprint

A dictionary fingerprint is a SHA-256 identifier that answers:

> Will these two files interpret every valid phrase in exactly the same way?

It covers the facts that can change that answer:

- The fingerprint-recipe version
- Profile and index-assignment rules
- Written-token recognition rules
- Every token, in index order

It does not cover comments, JSON indentation, property order, timestamps, source URLs, or quality reports. Those facts do not change how a phrase decodes.

The exact bytes passed to SHA-256 are called the **fingerprint input**. The final specification must define those bytes independently of JSON and LXB serialization so all languages calculate the same fingerprint. Each variable-length string will be preceded by its byte length. That keeps the lists `["ab", "c"]` and `["a", "bc"]` unambiguously different.

The fingerprint is not proof that a file is official or trustworthy. It detects a mismatch only when an application already knows the expected fingerprint from a trusted release or configuration. An ordinary per-file checksum may additionally be published to detect damage to one particular LXJ or LXB file.

## LXB goals

LXB exists to reduce startup work where measurements show that LXJ loading matters. Version 1 should remain deliberately small and predictable:

- Fixed magic bytes and a format version
- Fixed byte order
- Counts and the decoding fields needed without parsing JSON
- The shared dictionary fingerprint
- Canonical UTF-8 token bytes
- Bounds information sufficient to reject truncated or malformed files
- Optional provenance metadata that a runtime may skip

Two token layouts remain candidates: a sequence of length-prefixed strings, and an offset table followed by one string block. The offset table makes index-to-token access immediate but costs four bytes per token. The length-prefixed form is smaller but must be scanned once while loading. This choice will be made by a benchmark rather than assumed in advance.

LXB version 1 will not contain a portable hash table, trie, compressed string block, or minimal-perfect-hash structure. Those add implementation and compatibility costs to a dictionary that is only tens of thousands of short tokens. A loader may build its own language-native lookup table.

## What can be verified

The dictionary file alone contains enough information for **structural verification**:

- Supported format and fingerprint recipe
- Required counts for the declared profile
- Valid UTF-8 and required Unicode normalization
- No empty or duplicate tokens
- Token-character and separator rules
- Correct normal and trim ordering
- Matching calculated and stored fingerprints
- For LXB, valid lengths or offsets and no truncated sections
- Exact LXJ/LXB agreement when both are supplied

Selection quality requires more than the final dictionary. Checking familiarity, pronunciation distance, semantic distance, sensitive-term exclusions, and reproducibility also requires the pinned source datasets and selection configuration. A quality report records those checks; it is evidence about the derivation, not something the dictionary can prove by itself.

## Open work

Before either representation becomes final:

1. Settle the allowed token characters and separators.
2. Settle the JSON schema and publish a small LXJ example.
3. Define the exact fingerprint input.
4. Benchmark LXJ, length-prefixed LXB, and offset-table LXB loaders.
5. Define the LXB byte layout from the winning simple design.
6. Publish malformed-file and LXJ/LXB equivalence vectors.
