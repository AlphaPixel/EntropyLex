# EntropyLex dictionary files

EntropyLex can store one dictionary in two file forms:

- **LXJ** (`.lxj`) is the required, authoritative, human-readable JavaScript Object Notation (JSON) representation.
- **LXB** (`.lxb`) is an optional binary representation generated from LXJ. Its purpose is faster loading when measurements show that JSON parsing is a meaningful cost.

LXJ is authoritative: changes are made in or generated into LXJ first, and LXB is produced from it. Given the same LXJ file and the same version of `eldict-compile`, LXB generation must produce the same bytes. An LXB file must assign the same token to every numerical index and carry the same written-token rules as its LXJ source. Implementations must support LXJ. LXB support remains optional until its exact byte layout is final.

## One dictionary, two files

For example:

```text
entropylex-en-14-v1.lxj
entropylex-en-14-v1.lxb
```

A checksum calculated over every file byte differs between the two files because JSON and binary have different bytes. They carry the same **dictionary fingerprint** because they describe the same behavior.

## LXJ contents

The exact JSON fields and their types are not yet settled. At minimum, an LXJ file must record:

- LXJ format version and dictionary release version
- Dictionary name, profile, and normal-token width in bits
- Language and writing-system labels
- Rules for recognizing a written phrase: the standard Unicode representation used for equivalent character sequences, uppercase/lowercase handling, separators, and division into tokens
- Rules for assigning numerical indices and the possible final-token bit lengths
- Normal, trim, and total token counts
- One flat token array, in index order
- Dictionary fingerprint and the version of the fingerprint recipe
- Structured records identifying the source datasets and a checksum of the word-selection settings

The position of a token in the JSON array is its numerical index. A token entry must not repeat that index in another field because the two copies could disagree.

Language and writing-system labels describe the dictionary and its history. A label affects the fingerprint only if it selects behavior, such as case matching or division into tokens. The behavior, not the descriptive name, determines how a phrase decodes.

Source history, also called **provenance**, must use separate machine-readable fields rather than one paragraph. Each source record should state what the source contributed, identify an unchangeable release or exact repository commit, give its download location and license, record the downloaded file's SHA-256 checksum, and state the relevant record count. SHA-256 is a standard calculation that turns arbitrary bytes into a 32-byte identifying value.

## Dictionary fingerprint

A dictionary fingerprint is a SHA-256 identifier that answers:

> Will these two files interpret every valid phrase in exactly the same way?

The fingerprint calculation includes every fact that can change that answer:

- The fingerprint-recipe version
- Profile and numerical-index assignment rules
- Rules for recognizing written tokens, including normalization, case, separators, and token boundaries
- Every token, in index order

It excludes comments, JSON indentation, JSON property order, timestamps, source URLs, and quality reports. Changing those facts does not change how a phrase decodes.

SHA-256 accepts bytes. The exact byte sequence supplied to it is called the **fingerprint input**. The final specification must define that sequence independently of either file layout so JSON formatting and binary storage choices cannot change dictionary identity. Every string will be encoded as UTF-8 and preceded by its byte count. The byte counts ensure, for example, that `["ab", "c"]` cannot be mistaken for `["a", "bc"]` after the strings are placed next to one another.

Recalculating a file's fingerprint proves only that its contents agree with the fingerprint stored inside it. To detect replacement, an application must compare the result with an expected fingerprint obtained separately from a trusted release record or application setting. The fingerprint does not cryptographically prove who published the file. A release may also publish a checksum over every byte of each LXJ and LXB file to detect accidental byte changes or replacement of that exact file.

## LXB goals

LXB exists to reduce startup work where measurements show that LXJ loading matters. Version 1 should remain deliberately small and predictable:

- A short fixed byte sequence, called **magic bytes**, that identifies the file as LXB, followed by a format version
- One declared byte order for multi-byte numbers, so every language reads them identically
- Token counts and every field needed for decoding, stored where a loader can read them directly
- The shared dictionary fingerprint
- Token text encoded as UTF-8 after applying the dictionary's required written form
- Section start positions and lengths sufficient to reject files whose data ends early or extends outside the file
- Optional source-history fields that a loader may skip

Two token layouts remain candidates:

- A **length-prefixed sequence** stores each token's UTF-8 byte count immediately before that token. It uses less index data, but a loader must scan the preceding tokens once before it can locate every index.
- An **offset table and string block** stores one starting byte position per token, followed by all token bytes in one block. It permits immediate index-to-token access but costs four additional bytes per token if 32-bit positions are used.

Measurements will compare loading time, memory use, and file size before one layout is selected.

LXB version 1 will not store a ready-made mapping from token text back to its numerical index, and it will not compress the token bytes. Either feature would make the shared file and every language's loader more complicated. After reading the tokens, a loader may build whatever in-memory lookup structure best suits its programming language.

## What can be verified

Without using the original word sources, a verifier can check the file's internal structure and decoding rules:

- Supported format and fingerprint recipe
- Required counts for the declared profile
- Every token is valid UTF-8 and uses the dictionary's required Unicode normalization
- No empty or duplicate tokens
- Token-character and separator rules
- Normal and trim tokens occupy the required index ranges and order
- Matching calculated and stored fingerprints
- For LXB, every recorded length and starting position stays within the file and no section ends early
- Exact LXJ/LXB agreement when both are supplied

The final file cannot establish that its words are familiar, sufficiently different in pronunciation or meaning, or correctly screened for sensitive terms. Those checks require the exact source datasets and selection settings used to create it. A separate quality report records that evidence and whether another run with the same inputs reproduces the result.

## Open work

Before either representation becomes final:

1. Settle the allowed token characters, exact separator code points, case-conversion algorithm, normalization order, and the unit used to divide separator-free text into tokens.
2. Settle the exact JSON fields and types and publish a small LXJ example.
3. Define the exact fingerprint input.
4. Measure loading time, memory use, and file size for LXJ, length-prefixed LXB, and offset-table LXB.
5. Define the exact LXB byte layout using the simplest design that meets the measured needs.
6. Publish deliberately invalid file cases and tests proving exact LXJ/LXB agreement.
