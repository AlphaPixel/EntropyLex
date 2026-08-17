# EntropyLex overview and specifications document

This document defines the current design of EntropyLex, a reversible encoding that maps byte-aligned binary payloads into sequences of human-language tokens and back. The system is designed for human usability, high entropy density, deterministic reversibility, and compatibility with arbitrary 8 bit input data.

This specification covers all known properties, constraints, terminology, encoding rules, decoding rules, token classes, dictionary sizing, trimming behavior, profile variants, dictionary derivation, and non-English dictionary considerations as currently established.

Nothing in this document is frozen. Items marked **(provisional)** are working proposals that are expected to change once tooling and the first implementations exist.

---

## 1. Purpose

EntropyLex converts binary data into token sequences that are easy for humans to read, write, speak, and memorize. It addresses use cases where high entropy data such as keys, seeds, identifiers, or checksums must be represented in a human friendly format without losing exact bit precision.

EntropyLex is intended to:

- Represent arbitrary 8 bit data without ambiguity
- Provide dense entropy per token (up to `w` bits per token, where `w` is the profile's token width)
- Allow decoding without an explicit length header
- Guarantee byte aligned reconstruction
- Support deterministic round-trip mapping
- Scale across a family of token widths, so that low-complexity profiles and high-density profiles share one structural definition

---

## 2. Core Concepts

### 2.1 Input data
EntropyLex accepts input as a sequence of 8 bit bytes. All payloads are therefore multiples of 8 bits.

Let the input size be N bytes. Total payload bits = 8N.

An empty payload (N = 0) is legal and encodes to an empty token sequence.

### 2.2 Output tokens
EntropyLex produces a sequence of dictionary tokens. There are two classes:

1. Normal tokens
2. Trim tokens (used only for the final token)

A *token* is whatever atomic unit the dictionary defines. For English dictionaries a token is a single word. For ideographic dictionaries a token may be a single character or a fixed multi-character compound (see section 12).

### 2.3 Token width
Each normal token encodes exactly `w` bits, where `w` is fixed by the chosen profile. This requires 2^w distinct normal dictionary entries.

### 2.4 Byte alignment constraint
Because the input is 8 bit aligned and the output tokens are `w` bit aligned, dividing an 8N bitstream into `w` bit chunks yields only certain possible remainders.

Let `g = gcd(8, w)` and `r = (8N mod w)`. Every reachable remainder is a multiple of `g`:

```
r ∈ { 0, g, 2g, ..., w - g }
```

This determines the required trimming behavior for the final token, and it is the single parameter that drives the size of the trim dictionary.

For the reference profile, `w = 14`, `g = gcd(8, 14) = 2`, so `r ∈ {0, 2, 4, 6, 8, 10, 12}`.

---

## 3. Profiles

EntropyLex is defined as a family of profiles that differ only in token width `w`. All profiles share the same bitstream model, the same trimming concept, and the same encoding and decoding algorithms. A profile is named `EL-<w>`.

### 3.1 General profile arithmetic

For token width `w` with 8 bit input symbols:

```
g            = gcd(8, w)
remainders R = { 0, g, 2g, ..., w - g }
normal count = 2^w
trim count   = Σ 2^r  for r ∈ R, r > 0
             = (2^w - 2^g) / (2^g - 1)
total count  = 2^w + trim count
tokens for N bytes = ceil(8N / w)
```

When `w` is a multiple of 8, `g = w`, the remainder set collapses to `{0}`, and the trim count is zero.

### 3.2 Defined profiles

| Profile | `w` | `g` | Remainders `r` | Normal tokens | Trim tokens | Total dictionary |
|---|---|---|---|---|---|---|
| EL-8  | 8  | 8 | {0}                     | 256    | 0    | 256    |
| EL-12 | 12 | 4 | {0, 4, 8}               | 4096   | 272  | 4368   |
| EL-14 | 14 | 2 | {0, 2, 4, 6, 8, 10, 12} | 16384  | 5460 | 21844  |
| EL-16 | 16 | 8 | {0, 8}                  | 65536  | 256  | 65792  |

Trim token derivations:

- EL-12: 2^4 + 2^8 = 16 + 256 = 272
- EL-14: 2^2 + 2^4 + 2^6 + 2^8 + 2^10 + 2^12 = 4 + 16 + 64 + 256 + 1024 + 4096 = 5460
- EL-16: 2^8 = 256

Odd widths (EL-11, EL-13, …) are structurally permitted by the general arithmetic — `g = 1` and every remainder from 0 to `w-1` becomes reachable — but they are not defined profiles here, because `g = 1` maximizes trim dictionary size for a given density and yields no compensating benefit.

### 3.3 Token count comparison

Tokens required for representative payloads:

| Payload | EL-8 | EL-12 | EL-14 | EL-16 |
|---|---|---|---|---|
| 16 bytes (128 bit key)      | 16 | 11 | 10 | 8  |
| 32 bytes (256 bit key)      | 32 | 22 | 19 | 16 |
| 34 bytes (minimal GIF)      | 34 | 23 | 20 | 17 |
| 60 bytes (8x8 4-color GIF)  | 60 | 40 | 35 | 30 |

### 3.4 EL-8: the trivial profile

EL-8 is the simplest possible member of the family:

- Token width equals symbol width, so **no bit splitting occurs at all**. Each input byte maps to exactly one token, and each token maps to exactly one byte.
- `r` is always 0, so **there is no trim dictionary and no trim logic**. Sections 4.2 (trim dictionary), 6.4 (emit trim token), and 7.2 (trim decode branch) are inert.
- The dictionary is only 256 entries, which means word selection can be extremely aggressive about phonetic and orthographic distance. This is the only profile where a hand-curated, fully audited dictionary is practical.
- Byte alignment is automatic, so the structural validity check in section 7.3 always passes and provides no error detection (see section 13.5).

The cost is density: 8 bits per token means a 256 bit key becomes 32 tokens, roughly a 1.7x expansion in token count versus EL-14.

EL-8 is the recommended first implementation target. It exercises the dictionary loading, tokenizing, normalization, and round-trip test harness with none of the bitstream or trimming complexity, and it produces a working encoder in a form that can be verified by inspection. Prior art exists for this shape — the PGP Word List maps bytes to words in a comparable way — which makes EL-8 a useful sanity reference.

### 3.5 EL-12: the half-byte profile

EL-12 is the intermediate step:

- `g = 4`, so the bitstream only ever splits bytes **on nibble boundaries**. Every token boundary falls on a whole byte or a half byte, which is easy to reason about, easy to debug, and easy to display in hex alongside the tokens.
- Only two nonzero remainders exist, `r ∈ {4, 8}`, so the trim dictionary is 272 entries rather than 5460, and the trim index table is small enough to print.
- Total dictionary is 4368 entries, which is small enough to be sourced from a high-frequency core vocabulary. Every token can be a common, short, easily spelled word.
- Density is 12 bits per token, which is 86% of EL-14's density for a fifth of the dictionary size.

EL-12 is the recommended second implementation target. It exercises the full bitstream packer, the remainder calculation, the trim dictionary, and the alignment validation — that is, every mechanism EL-14 needs — while keeping the dictionary derivation problem tractable and the intermediate values inspectable by hand.

### 3.6 EL-14: the reference profile

EL-14 is the density-optimized target profile and the default meaning of "EntropyLex" without further qualification. It maximizes bits per token subject to the constraint that the dictionary remain sourceable from common English root words. See section 15 for the rationale.

EL-14 requires the full mechanism: six trim classes, 5460 trim tokens, and a 21844 entry dictionary that pushes to the edge of the usable common-English vocabulary.

### 3.7 EL-16 and wider

EL-16 is impractical for English — 65,792 tokens far exceeds the common English root vocabulary — but is structurally the *simplest* profile after EL-8: `g = 8`, exactly one trim class, and every token boundary falls on a whole byte or a half-token boundary. Scripts with much larger token inventories bring it closer to reach, though section 12 concludes it remains aspirational there too. No EL-16 dictionary is currently considered practical.

### 3.8 Profile and dictionary identification

A token sequence does not describe itself. Correct decoding requires knowing both:

1. The profile (`w`), and
2. The exact dictionary used

Both must be carried out of band, in the same way the choice of Base32 versus Base58 is carried out of band. A dictionary file therefore declares its profile, and implementations must be able to report its dictionary fingerprint (see section 11.7). The fingerprint is a short SHA-256 identifier for the ordered token mapping and the rules needed to interpret it. Matching fingerprints mean that two files decode phrases the same way, even when one file is JSON and the other is binary.

Implementations must not attempt to auto-detect the profile from a token sequence. A sequence valid under one dictionary may be valid but different under another.

One conditional exception is under consideration. If the dictionaries for a family of profiles are derived as mutually disjoint token sets (Scheme C in section 11.8), then every token identifies its profile unambiguously, and auto-detection becomes sound *within that dictionary family*. That composition decision has not been made, so this exception is not yet in force. Even if adopted, it would identify the profile only — not the exact mapping or phrase-recognition rules — so the fingerprint requirement stands either way.

---

## 4. Dictionary Structure

Total dictionary = normal dictionary + trim dictionary.

### 4.1 Normal dictionary
- Size: 2^w tokens
- Index range: 0 to 2^w - 1
- Each normal token represents a `w` bit binary value
- Used for all tokens except, possibly, the final token

### 4.2 Trim dictionary
A trim token is used only as the final token, to encode the final remainder bits. For each reachable nonzero remainder `r`, the trim dictionary must uniquely encode every possible `r` bit value, which requires 2^r tokens for that remainder class.

Trim tokens must be mapped so that each one unambiguously identifies:

1. The remainder size `r`
2. The `r` bit tail payload

Trim tokens never appear except as the final token. The normal and trim dictionaries are disjoint sets: no token appears in both, and no token appears twice. This disjointness is what allows a decoder to determine the final token's class by lookup alone, with no in-band flag.

For EL-8 the trim dictionary is empty and this entire subsection is inert.

### 4.3 Unified index space **(provisional)**

To keep dictionary mappings as a single flat list in both LXJ and LXB, all tokens occupy one index space:

- Indices `0 .. 2^w - 1` are the normal tokens, index equal to the encoded value.
- Trim tokens follow, grouped by ascending remainder size. For remainder `r` with value `v`:

```
index(r, v) = 2^w + trim_base(r) + v
trim_base(r) = Σ 2^j  for j ∈ R, 0 < j < r
```

For EL-14 the trim bases are:

| `r` | class size | `trim_base(r)` | index range |
|---|---|---|---|
| 2  | 4    | 0    | 16384 .. 16387 |
| 4  | 16   | 4    | 16388 .. 16403 |
| 6  | 64   | 20   | 16404 .. 16467 |
| 8  | 256  | 84   | 16468 .. 16723 |
| 10 | 1024 | 340  | 16724 .. 17747 |
| 12 | 4096 | 1364 | 17748 .. 21843 |

For EL-12:

| `r` | class size | `trim_base(r)` | index range |
|---|---|---|---|
| 4 | 16  | 0  | 4096 .. 4111 |
| 8 | 256 | 16 | 4112 .. 4367 |

For EL-16, the single trim class `r = 8` occupies indices 65536 .. 65791.

A decoder therefore needs only one token-to-index map; the index alone determines class, remainder size, and value.

---

## 5. Bit Order and Canonical Forms

### 5.1 Bit order
The bitstream is packed **most significant bit first**. Bit 0 of the stream is the most significant bit of the first input byte; bit 7 is its least significant bit; bit 8 is the most significant bit of the second byte, and so on.

A token's value is formed by reading its bits in stream order, with the first bit read becoming the most significant bit of the value. A trim token's `r` bit value is likewise read MSB-first and is right-aligned in the range `0 .. 2^r - 1`.

This ordering applies to all profiles and is not configurable.

### 5.2 Canonical phrase form
The canonical written form of an EntropyLex phrase is the tokens in order, separated by a single ASCII space (U+0020), with no leading or trailing whitespace, in the exact case and Unicode normalization form (NFC) recorded in the dictionary file.

### 5.3 Input normalization for decoding
Decoders should accept a superset of the canonical form and normalize before lookup:

- Any run of whitespace, hyphens, or newlines is a token separator
- Case-insensitive matching for cased scripts
- Unicode NFC normalization before lookup
- Leading and trailing separators ignored

Decoders must not perform spelling correction or nearest-token matching silently. If a token is not in the dictionary, decoding fails. A separate, clearly labeled "suggest correction" facility is permitted but must never be part of the decode path.

For dictionaries whose tokens are single fixed-width glyphs (see section 12), a phrase may legally be written with no separators, and decoders for such dictionaries must segment by glyph count rather than by whitespace.

### 5.4 Empty payload
A zero-byte payload encodes to a zero-token phrase, and a zero-token phrase decodes to a zero-byte payload. This is legal in all profiles.

---

## 6. Encoding Process

### 6.1 Convert input bytes to bitstream
Concatenate all input bytes into a linear bit sequence of length 8N bits, MSB-first per section 5.1.

### 6.2 Emit full `w` bit tokens
While at least `w` bits remain:

- Extract the next `w` bits
- Interpret them as a value between 0 and 2^w - 1
- Map the value to a normal token and append it to the output sequence

### 6.3 Determine remainder
After consuming all full `w` bit groups, the remainder bit length is:

```
r = (8N mod w)
```

`r` is always a member of the profile's remainder set `R`.

### 6.4 Emit final trim token
If `r = 0`:

- No trim token is required
- Encoding ends on the last normal token

If `r > 0`:

- Read the final `r` bits from the remaining bitstream
- Use the pair `(r, final_value)` to select the appropriate trim token
- Append the trim token as the final token

Encoding ends after emitting the final trim token.

In EL-8, and in any profile where `w` is a multiple of 8, `r` is always 0 and this step never executes.

---

## 7. Decoding Process

Given a sequence of tokens:

### 7.1 Convert all but the last token to bits
For each token except the final one:

- Look up its value in the normal dictionary
- Append the `w` bits to the reconstructed bitstream

Any of these tokens that resolves to a trim token makes the sequence invalid (see D2 below).

### 7.2 Decode the final token
If the final token is a normal token:

- It contributes `w` bits

If the final token is a trim token:

- Identify its remainder size `r`
- Extract the `r` bits encoded by that token
- Append these `r` bits to the reconstructed bitstream

### 7.3 Validation rules
A decoder must reject a sequence unless all of the following hold:

- **D1** Every token is present in the dictionary for the declared profile.
- **D2** No trim token appears in any position other than the final one.
- **D3** The total reconstructed bit length is divisible by 8.
- **D4** If the final token is a trim token with remainder `r`, then `r ∈ R` and `r > 0`.
- **D5** An empty sequence is valid and yields an empty payload.

D3 is the only structural integrity check in the format. Its strength is discussed in section 13.5.

### 7.4 Convert bitstream back to bytes
Group bits into bytes, MSB-first, and output the original binary data.

---

## 8. Trimming Behavior Summary

The input is 8 bit aligned and the encoder emits `w` bit tokens, so:

- Remainder sizes are always multiples of `g = gcd(8, w)`
- Maximum remainder is `w - g`
- Remainders not divisible by `g` never occur
- Only the nonzero members of `R` require trim tokens

| Profile | Remainder step `g` | Max remainder | Trim classes | Trim tokens |
|---|---|---|---|---|
| EL-8  | 8 | —  | 0 | 0    |
| EL-12 | 4 | 8  | 2 | 272  |
| EL-14 | 2 | 12 | 6 | 5460 |
| EL-16 | 8 | 8  | 1 | 256  |

Trim tokens must be distinct from the normal tokens in the same dictionary.

The rejected alternative was a mandatory length header. Even an 8 bit header would enlarge every message unconditionally, while trimming costs only dictionary space and emits no more tokens than an unpadded encoding would.

---

## 9. Dictionary Size Summary

| Profile | Normal | Trim | Total | English feasibility |
|---|---|---|---|---|
| EL-8  | 256   | 0    | 256   | Trivially sourceable; can be hand-audited |
| EL-12 | 4096  | 272  | 4368  | Comfortably within reviewed high-frequency source-list sizes |
| EL-14 | 16384 | 5460 | 21844 | Raw open sources are large enough; survival after familiarity and distance filters remains to be measured |
| EL-16 | 65536 | 256  | 65792 | Not feasible for single-word English tokens; see section 12 |

Because trim tokens are used less often than normal tokens, they should be drawn from the less-frequent end of the selected vocabulary, leaving the most common and most robust words for normal tokens.

---

## 10. Word Selection Criteria

Word choice is not fully specified. The governing principles are:

- Distinct pronunciation to reduce confusion in speech
- Minimal homophones or near-homophones
- Easy, predictable spelling
- Broad familiarity across dialects
- Avoidance of culturally sensitive, offensive, or domain specific terms
- Prefer short to medium length words
- Prefer unique phonetic structure and low pairwise similarity
- Prefer low semantic similarity, so that adjacent tokens do not read as a phrase and near-synonyms cannot be substituted from memory

The final dictionary for a given profile must contain exactly the counts in section 9 for that profile — no more and no fewer.

Section 11 describes the intended mechanical process for satisfying these principles.

---

## 11. Dictionary Derivation

Dictionaries are not to be written by hand (except possibly EL-8, which is small enough to audit manually). They are to be **derived** from a master word list by a set of preprocessing tools. Those tools do not exist yet; this section defines what they must do, so that dictionary construction is reproducible, auditable, and re-runnable when criteria change.

The guiding objective: **select the subset of a master vocabulary that maximizes minimum pairwise distance in spelling, sound, and meaning, subject to a hard count and a familiarity floor.**

### 11.1 Pipeline overview **(provisional)**

Planned tools, to live under a `tools/` directory:

| Stage | Tool | Responsibility |
|---|---|---|
| 1 | `eldict-ingest`  | Load master word list plus frequency, pronunciation, and part-of-speech data into a normalized candidate table |
| 2 | `eldict-filter`  | Apply hard disqualifiers; emit surviving candidates with rejection reasons logged |
| 3 | `eldict-score`   | Compute orthographic, phonetic, and semantic feature vectors for each candidate |
| 4 | `eldict-select`  | Build the confusability graph and select the final token set to a target count |
| 5 | `eldict-emit`    | Assign canonical indices, partition normal versus trim, and write the canonical LXJ file |
| 6 | `eldict-compile` | Compile LXJ into the optional LXB runtime representation without changing its meaning or fingerprint |
| 7 | `eldict-verify`  | Verify file structure and, when source inputs are supplied, re-run quality checks and emit a report |

Every stage is deterministic given its inputs and a recorded configuration file. Re-running the pipeline on the same inputs must produce the same ordered mapping and dictionary fingerprint. The canonical LXJ formatter and the eventual LXB layout must additionally produce byte-identical files.

### 11.2 Stage 1 — ingest

Input sources, all pinned to a specific version and recorded as structured LXJ provenance records:

- A master word list with usage frequency (a corpus-derived frequency list, so that "common" is measured rather than asserted)
- A pronunciation lexicon giving a phoneme sequence per word, for phonetic distance
- Lemma and inflection data, so that inflected forms can be collapsed to their root
- A word embedding model or equivalent, for semantic distance

Candidates enter with: surface form, lemma, frequency rank, phoneme string, syllable count, part of speech.

Candidate sources and their licenses, published sizes, and suitability are tracked in [`data/dict/SOURCES.md`](data/dict/SOURCES.md). A source is not approved merely because it is downloadable. Before use, its exact release or commit, downloaded-file SHA-256 checksum, license obligations, and ingested record count must be recorded.

### 11.3 Stage 2 — hard filters

Disqualifiers, applied before any distance computation:

- Length outside the accepted band **(provisional: 3 to 8 characters)**
- Non-ASCII characters, for English dictionaries
- Proper nouns, brand names, and place names
- Inflected forms whose lemma is already a candidate (plurals, conjugations, comparatives)
- Words on an offensive, slur, or sensitive-topic blocklist
- Words on a domain-jargon blocklist
- Known homophone sets — retain at most one member, or drop the whole set
- Homographs with divergent pronunciations (e.g. "read", "lead", "wound")
- Words whose spelling is not predictable from pronunciation, or vice versa, beyond a configured threshold
- Frequency below a familiarity floor

Every rejection is logged with its reason. The rejection log is a reviewable artifact, not a side effect.

### 11.4 Stage 3 — distance metrics **(provisional)**

Three independent distance families, each computed pairwise over surviving candidates:

**Orthographic (spelling) distance**
- Damerau-Levenshtein edit distance over the surface form
- A keyboard-adjacency-weighted variant, so that transpositions and adjacent-key substitutions count as nearer
- Common-prefix and common-suffix length, to avoid clusters of words that differ only in an ending

**Phonetic (sound) distance**
- Edit distance over phoneme strings, with substitution costs weighted by articulatory feature distance, so that /b/–/p/ is nearer than /b/–/s/
- Coarse phonetic-key collision (Soundex/Metaphone class) as a hard reject rather than a soft cost
- Syllable count and stress pattern, used to spread the selection across prosodic shapes

**Semantic (meaning) distance**
- Cosine distance in an embedding space
- Purpose: prevent near-synonyms ("big"/"large") that a human may substitute from memory without noticing, and reduce the chance that a random token sequence reads as a meaningful phrase and is therefore misremembered as one

### 11.5 Stage 4 — selection

Selection is a maximin problem: choose `n` tokens from `m` candidates so that even the two most easily confused selected words are as different as possible under the combined spelling, sound, and meaning checks.

The intended approach **(provisional)**:

1. Build a *confusability graph*: an edge joins any two candidates whose distance falls below the threshold in **any** of the three families. Conflicting words must not both be selected.
2. Seed with the highest-frequency candidates, since familiarity is the property that cannot be recovered later.
3. Walk through candidates from most to least familiar. Add a word only when it has no conflict edge to a word already selected. In graph terminology this is an independent-set heuristic; operationally it is simply “keep the familiar word, then skip words too similar to it.”
4. If the target count `n` is not reached, relax thresholds by a recorded step and repeat, logging exactly which threshold was relaxed and by how much.
5. If the target count is exceeded, tighten thresholds and repeat.

The relaxation trace is part of the output. A dictionary that only reached its target count by weakening the phonetic threshold twice should say so on its face.

Exact optimization is not required; a reproducible heuristic with a published quality report is preferred over an unreproducible optimum.

### 11.6 Stage 5 — partition and index assignment

1. Sort the selected set by frequency, descending.
2. Assign the most frequent 2^w tokens to the normal dictionary and the remainder to the trim dictionary, per section 9.
3. Within each group, sort tokens by Unicode code point order (NFC, and for cased scripts, lowercase) to produce the canonical order.
4. Assign indices per the unified index space in section 4.3.

Sorting the final groups lexicographically rather than by frequency makes LXJ stable and diffable: a token added or removed shifts a contiguous run of indices instead of permuting the whole mapping.

### 11.7 Stages 6 and 7 — dictionary files, provenance, and verification

One dictionary has two representations:

- **LXJ** (`.lxj`) is the canonical, human-readable JSON representation and the source of truth.
- **LXB** (`.lxb`) is an optional compiled binary representation for implementations that benefit from faster loading.

For example: `entropylex-en-14-v1.lxj` and `entropylex-en-14-v1.lxb`, stored under `data/dict/`.

An LXB file is generated only from LXJ. The two files have different ordinary file checksums because their bytes differ, but they must have the same dictionary fingerprint and must assign exactly the same token to every index. Implementations are required to support LXJ. LXB support is optional until its version 1 byte layout is finalized.

LXJ contains a flat token array whose array position is the unified index from section 4.3. It also contains the profile, counts, written-token recognition rules, structured source records, selection-configuration checksum, and dictionary fingerprint. It must not store a redundant index on every token.

The **dictionary fingerprint** is a SHA-256 identifier for the ordered mapping and the rules that affect how phrases are recognized and decoded. It covers the profile and index scheme, normalization/case/segmentation behavior, and every token in index order. It does not cover JSON indentation, property order, comments, timestamps, source URLs, or quality reports, because those do not change decoding.

The exact bytes supplied to SHA-256 are the **fingerprint input**. They will be defined independently of LXJ and LXB serialization, with every variable-length string preceded by its byte length. This prevents serialization details from changing dictionary identity and lets every implementation calculate the same result. The exact version 1 fingerprint input remains to be specified.

Language and script labels are recorded in LXJ, but descriptive labels alone do not change the fingerprint. Any actual behavior they select — for example case handling, glyph segmentation, or right-to-left display rules — does.

A fingerprint detects a mismatch only when the expected fingerprint comes from a trusted release or configuration. It does not authenticate a file against deliberate replacement. Published releases may also provide a checksum for each exact LXJ and LXB file.

The LXB version 1 design must contain fixed identification and bounds information, the decoding fields needed at runtime, the shared fingerprint, and canonical UTF-8 token bytes. A length-prefixed token sequence and an offset table plus string block are both candidates. The choice must be benchmarked before the byte layout is frozen. Version 1 will not require a serialized hash table, trie, compressed token block, or minimal-perfect hash.

`eldict-verify` can assert the following from LXJ or LXB alone:

- Supported file format and fingerprint recipe
- Exact token counts for the declared profile
- Valid UTF-8, required Unicode normalization, and token-character rules
- No empty or duplicate tokens
- Normal and trim sets disjoint and canonically ordered
- Calculated fingerprint matches the stored fingerprint
- Round-trip of published dictionary-dependent test vectors
- LXB section bounds are valid and, when LXJ is supplied, both files describe exactly the same dictionary

The final dictionary file cannot by itself prove that the source words were familiar, phonetically distant, semantically distant, or produced by the claimed pipeline. Those quality and reproducibility checks additionally require the pinned source datasets and selection configuration. The verifier's quality report must name and checksum those inputs, report the surviving hard filters and minimum distances, and list the worst offending pairs.

The detailed working design and remaining decisions are in [`data/dict/FORMAT.md`](data/dict/FORMAT.md).

### 11.8 Cross-profile dictionary composition — UNDECIDED

Profiles are independent as an encoding matter: an EL-8 phrase is not an EL-12 phrase, and phrases are never portable between profiles.

What is **not yet determined** is how the token *sets* for EL-8, EL-12, and EL-14 should relate to one another when all three are derived from the same master corpus. This is an open design decision, not a settled one. Three schemes are on the table, and any of them could be chosen.

It is desirable in all three that the smaller profiles come out of the same derivation run, so that one quality report covers the family and the profiles feel like one system rather than three unrelated encodings.

#### Scheme A — nested subset

`EL-8 ⊂ EL-12 ⊂ EL-14`. The EL-8 dictionary is the 256 highest-ranked tokens of the EL-12 dictionary, which is in turn the top of the EL-14 dictionary.

- Cheapest: 21,844 tokens total, no more than EL-14 alone requires.
- Simplest to reason about, document, and verify.
- Worst cross-profile ambiguity: every EL-8 phrase is also a well-formed run of EL-14 normal tokens, and vice versa. Decoding under the wrong profile yields different, entirely plausible bytes.
- Weakest minimum distances at small `n`: the EL-8 set is whatever the top 256 of a 21,844-token selection happens to be, not the 256 most mutually distinct words the corpus can offer.

#### Scheme B — independent optimization

Each profile is selected independently against its own target count, with thresholds tuned for that count.

- Best per-profile quality. At `n = 256` the selector can demand enormous minimum distances that are simply unattainable at `n = 21,844`.
- Corpus cost is between 21,844 and 26,468 depending on how much the selections happen to overlap.
- Overlap is uncontrolled, so cross-profile ambiguity is partial and accidental — the worst of both worlds for identification purposes, since a phrase is *sometimes* profile-identifying and there is no way to know when.

#### Scheme C — disjoint partition (self-identifying profiles)

The three dictionaries are carved from the master corpus as mutually exclusive sets, with no token appearing in more than one profile.

This buys a genuinely useful property: **the profile becomes identifiable from the tokens themselves.** Any single token in a phrase determines which profile produced it, because that token exists in exactly one profile's dictionary. That converts the out-of-band profile declaration required by section 3.8 from a necessity into a cross-check, and it makes decoding under the wrong profile a hard failure rather than a silent production of wrong bytes.

The cost is corpus size:

| Component | EL-8 | EL-12 | EL-14 | Total |
|---|---|---|---|---|
| Normal tokens | 256 | 4,096 | 16,384 | 20,736 |
| Trim tokens   | 0   | 272   | 5,460  | 5,732  |
| Per profile   | 256 | 4,368 | 21,844 | **26,468** |

Two observations on that total:

- The often-quoted figure for this idea is 2^14 + 2^12 + 2^8 = **20,736**, but that counts only the normal dictionaries. Trim tokens must also be disjoint — both from each other and from every normal set — which adds 5,732 and brings the real requirement to **26,468 tokens**.
- 26,468 exceeds the reviewed 25,000-entry 12dicts core list, although broader permissively licensed lemmatized and lexical sources contain many more raw entries. Only the 20,736 normal tokens need the highest familiarity; the 5,732 trim tokens can come from the lower-frequency surviving tail because trim tokens occur at most once per phrase. This makes Scheme C an experiment worth running, not a feasibility result. The selector must publish survivor counts and threshold relaxations before this scheme can be called practical.

A partial variant — disjoint normal dictionaries with a single shared trim dictionary — was considered and is **not recommended**. Because the shared trim set must still avoid all three normal sets, it saves only 272 tokens (26,196 versus 26,468), and it reintroduces ambiguity in exactly the case where it is hardest to notice: a one-byte payload encodes to a single trim token under both EL-12 and EL-14, so such a phrase would carry no identifying information at all. Paying 272 tokens to eliminate that case is obviously correct.

Extending disjointness to EL-16 is not possible for English. Adding it would require 92,260 tokens.

#### Status

**No scheme has been selected.** The trade is: Scheme A is cheapest and simplest, Scheme B gives the best per-profile dictionaries, and Scheme C gives self-identifying phrases and hard failure on profile mismatch while demanding 26,468 distinct survivors and potentially weaker distance thresholds.

This decision must be made before `eldict-select` is written, since it determines whether the selector runs once with a partitioning step or three times with exclusion sets. It also interacts with section 3.8 (profile identification) and section 13.5 (error detection), both of which currently assume the conservative case.

Note also that disjointness under Scheme C identifies the *profile* only, and only within a single dictionary family. It does not identify the exact token mapping or phrase-recognition rules. The fingerprint requirement in section 11.7 stands regardless of which scheme is chosen.

---

## 12. Alternate Language and Ideographic Dictionaries

Nothing in the encoding is English-specific. The bitstream model, trimming, and index space are language neutral; only the dictionary changes. A dictionary in any language, for any defined profile, is a valid EntropyLex dictionary provided it meets the counts in section 9 and the disjointness requirement in section 4.2.

### 12.1 Why other scripts change the profile calculus

Profile choice for English is bounded by usable vocabulary, not raw database size. A reviewed common-word list contains about 25,000 entries, while broader lemmatized and lexical sources contain 84,000 to more than 135,000; the usable count after familiarity and distance filters is not known yet. EL-14 is the practical ceiling to test. Scripts with large character inventories or large compound-word lexicons may shift that ceiling, and it is worth asking how far — particularly since EL-16 is structurally *simpler* than EL-14, with one trim class instead of six.

The answer, worked through below, is less than one might hope: the binding constraint is not how many characters or words a script possesses, but how many a competent reader can reliably recognize, distinguish, and write from dictation. Those two numbers differ by an order of magnitude in every script examined, and only the second one is a valid dictionary source.

The trade is between token inventory and per-token complexity: a larger inventory buys density, but each token becomes rarer and therefore less familiar to the human in the loop.

### 12.2 Single-glyph token dictionaries

For Chinese, a token may be a single Han character. This is attractive: one glyph carries `w` bits, phrases are visually compact, and no separators are needed because segmentation is by glyph.

The constraint is the character inventory. The Table of General Standard Chinese Characters defines 8,105 characters (3,500 at level 1, 3,000 at level 2, 1,605 at level 3). Against the profile counts:

| Profile | Total tokens | Single Han character feasibility |
|---|---|---|
| EL-8  | 256    | Trivial — can be drawn from the most common characters alone |
| EL-12 | 4368   | Comfortable — fits within levels 1 and 2, with selection headroom |
| EL-14 | 21844  | Not possible with single standard characters; requires compounds |
| EL-16 | 65792  | Not possible with single characters; requires compounds |

**EL-12 is therefore the natural profile for a single-character Chinese dictionary**, in the same way EL-14 is the natural profile for single-word English. This is a useful convergence: EL-12 is also the profile whose trim logic is simplest to implement, so a Chinese single-glyph dictionary and an early implementation target coincide.

#### Japanese is substantially more constrained than Chinese

The relevant question for any script is not "how many characters exist" but "how many characters can the human in the loop reliably read, write, and tell apart." For Japanese those numbers diverge sharply, and it is easy to reach for the wrong one.

| Japanese character set | Count | Nature | Largest single-glyph profile |
|---|---|---|---|
| Kyōiku kanji | 1,026 | Taught in primary school; subset of jōyō | EL-8 (4.0x headroom) |
| Jōyō kanji | 2,136 | Cabinet-notified guide for general use, 2010 | EL-8 (8.3x headroom) |
| Jōyō + jinmeiyō | 3,000 | All kanji legally permitted in personal names | EL-8 only — short of EL-12's 4,368 |
| JIS X 0208 kanji | 6,355 | Encoding repertoire | EL-12 arithmetically, not in practice |
| JIS X 0213 kanji | ~9,980 | Encoding repertoire | EL-12 arithmetically, not in practice |

The figure most often cited as "about 10,000 kanji" is the **JIS X 0213 encoding repertoire**, not a list of characters anyone is expected to know. It is the set a Japanese computer can display, which includes a long tail of characters most native readers cannot read and would not attempt to write from dictation. Using it as a dictionary source would be the equivalent of sourcing an English dictionary from the OED's full headword list rather than from a common-word frequency list — arithmetically satisfying and practically useless.

Measured against the sets that reflect actual literacy, **Japanese single-glyph dictionaries top out at EL-8.** Even the entire legally-permitted personal-name inventory, 3,000 characters, falls short of EL-12's 4,368-token requirement, and that is before any selection headroom for visual distinctness. A Japanese EL-12 dictionary must therefore extend past jōyō into characters of declining familiarity, mix in kana-written vocabulary items, or abandon single glyphs for multi-character words.

#### Why Chinese reaches further

The gap is structural, not incidental. Chinese writes essentially everything in hanzi, so near-complete text coverage requires an everyday inventory on the order of 3,500 characters. Japanese offloads grammar and a large share of vocabulary onto kana, so kanji past the jōyō set are genuinely rare in running text. The result is roughly a 3x difference in usable single-glyph inventory, which is exactly one profile step: Chinese reaches EL-12, Japanese reaches EL-8.

### 12.3 Compound-token dictionaries

To reach EL-16's 65,792 tokens, a dictionary must use multi-character words rather than single glyphs. The token count is reachable in principle, but the margin is much thinner than it first appears.

A standard Chinese desk dictionary, the 7th edition of 现代汉语词典 (Xiandai Hanyu Cidian, 2016), contains approximately 69,000 entries. EL-16 would consume 95% of it — every entry, including the rare and the archaic, with essentially nothing discarded for familiarity, visual distinctness, or offensiveness. That is not a dictionary anyone can use. Reaching EL-16 with usable tokens requires a comprehensive lexicon several times that size, and pays the same familiarity cost at the tail that English pays at EL-14, only more so.

Practically, then, EL-16 is aspirational for CJK as well. EL-12 single-glyph Chinese is the profile the evidence actually supports. Additional considerations for compound tokens:

- Segmentation is no longer free. Either the tokens are fixed-width (e.g. all exactly two characters, so segmentation is again by count) or the phrase requires explicit separators.
- Fixed-width two-character tokens are the recommended construction **(provisional)**: it preserves separator-free writing, keeps a phrase visually regular, and makes the decoder's segmentation trivial.
- Familiarity drops as the lexicon is exhausted, exactly as it does for English at EL-14.

### 12.4 Selection criteria differ by script

The English criteria in section 10 are weighted toward speech, because English's failure modes are homophones and accent variation. Other scripts fail differently:

- **Mandarin and Japanese have high homophone density.** Spoken transmission of a CJK phrase is substantially less reliable than written transmission, tone and pitch-accent distinctions notwithstanding. CJK dictionaries should therefore optimize primarily for **visual** distinctness and treat phonetic distance as secondary.
- **Visual distance** replaces orthographic edit distance: stroke count, radical/component overlap, and overall glyph shape. Characters sharing a radical and stroke count are the ideographic equivalent of a near-homophone.
- **Romanization distance** remains useful as a secondary metric — pinyin plus tone for Mandarin, romaji or kana for Japanese — since it governs how the phrase is typed and how it is spoken when speech is unavoidable.
- **Semantic distance** applies unchanged, and matters more: character sequences in CJK readily read as meaningful compounds, so semantic screening should also reject adjacent-token pairs that form a common word.

The derivation pipeline in section 11 is intended to be metric-pluggable for exactly this reason: stage 3's three distance families are replaced per script, while stages 1, 2, 4, 5, 6, and 7 are unchanged.

### 12.5 Additional requirements for non-ASCII dictionaries

- Tokens are stored and compared in Unicode NFC. The normalization and matching behavior is recorded explicitly in LXJ and compiled into LXB.
- For Chinese, the dictionary must declare simplified or traditional; variant characters must not both appear, and a simplified-to-traditional mapping is not a substitute for a separate dictionary artifact with its own fingerprint.
- Case folding does not apply; decoders must not apply cased-script normalization to these dictionaries.
- Dictionaries for right-to-left scripts must specify the writing order of tokens explicitly, distinct from the logical token order used by the encoder.
- Every dictionary records its language and script. The fingerprint covers the actual behavior needed to recognize and decode its phrases, such as normalization, case handling, glyph segmentation, and logical token order. A descriptive language label alone is not treated as decoding behavior.

### 12.6 References for the figures in this section

Character inventory counts cited above, with sources. Note the distinction that governs all of them: **standards of general use** describe what people are expected to know, while **encoding repertoires** describe what a computer can represent. Only the former are valid dictionary sources.

| Figure | Value | Type | Source |
|---|---|---|---|
| Table of General Standard Chinese Characters (通用规范汉字表), 2013 | 8,105 total — Tier 1: 3,500, Tier 2: 3,000, Tier 3: 1,605 | General use | State Council of the PRC, promulgated June 2013 — https://en.wikipedia.org/wiki/List_of_Commonly_Used_Standard_Chinese_Characters ; full list at http://hanzidb.org/character-list/general-standard ; scan at https://archive.org/details/tongyongguifan_hanzibiao_2013 |
| Jōyō kanji (常用漢字) | 2,136 | General use | Japanese Cabinet Notification, 30 November 2010 (added 196, removed 5, from the 1,945 of 1981) — https://en.wikipedia.org/wiki/J%C5%8Dy%C5%8D_kanji ; https://www.sljfaq.org/afaq/jouyou-kanji.html |
| Kyōiku kanji (教育漢字) | 1,026 | General use, subset of jōyō | Taught grades 1–6; the remaining 1,110 jōyō are taught grades 7–12 — https://en.wikipedia.org/wiki/J%C5%8Dy%C5%8D_kanji |
| Jinmeiyō kanji (人名用漢字) | 864 | General use, names only | Most recent addition 26 June 2026 — https://en.wikipedia.org/wiki/Jinmeiy%C5%8D_kanji |
| Kanji legally permitted in personal names | 3,000 | General use | 2,136 jōyō + 864 jinmeiyō — https://en.wikipedia.org/wiki/Jinmeiy%C5%8D_kanji |
| JIS X 0208 kanji | 6,355 | Encoding repertoire | Levels 1 and 2 — https://en.wikipedia.org/wiki/JIS_X_0208 |
| JIS X 0213 kanji | ~9,980 | Encoding repertoire | 6,355 plus roughly 3,625 added; the standard's own accounting varies slightly (3,625 vs 3,695) depending on how JIS X 0212 overlap is counted. Total characters in the standard including kana, Latin, Greek, Cyrillic and symbols is 11,233 — https://en.wikipedia.org/wiki/JIS_X_0213 |
| 现代汉语词典 (Xiandai Hanyu Cidian), 7th ed., 2016 | ~69,000 entries | Lexicon | Characters, words, expressions and idioms — https://en.wikipedia.org/wiki/Xiandai_Hanyu_Cidian |

The commonly cited "about 10,000 kanji" is the JIS X 0213 row of this table — an encoding repertoire, roughly 4.7x the jōyō list. It is not a measure of what a Japanese reader knows and must not be used to size a dictionary.

Sources retrieved 2026-08-17. Any figure used to actually build a dictionary must be re-verified against the primary standard document and pinned in the LXJ source records, per section 11.7.

---

## 13. Security Considerations

### 13.1 Entropy preservation
EntropyLex preserves the entropy of the input bitstream exactly. All mappings between binary values and tokens are deterministic and reversible. No profile adds, removes, or reshapes entropy.

### 13.2 No semantic leakage
Tokens are treated purely as codepoints and have no semantic relationship to the data they encode.

### 13.3 Human error
Mispronunciation, homophones, and transcription mistakes are the dominant risk. Dictionary design must minimize confusable pairs; this is the purpose of section 11. Lower-`w` profiles are meaningfully safer here, because a smaller dictionary permits much larger minimum pairwise distances.

### 13.4 Length leakage
Token count leaks the payload size. In profiles with trim tokens the leak is approximate; in EL-8 and any profile where `w` is a multiple of 8, the token count reveals the payload length **exactly**. This may or may not be acceptable depending on application.

### 13.5 Error detection is weak and profile dependent
EntropyLex has no checksum. The only structural check is D3, the byte-alignment test, and its power depends on the profile:

- A substitution of one valid token for another valid token of the same class is **never** detected in any profile.
- An inserted or deleted normal token changes the bit count by `w`. This is detected only when `w mod 8 ≠ 0`. EL-12 (`w mod 8 = 4`) and EL-14 (`w mod 8 = 6`) therefore detect every single-token insertion or deletion; **EL-8 and EL-16 detect none.**
- A deleted trim token is detected unless its remainder `r` is itself a multiple of 8 — in EL-12 and EL-14, the `r = 8` case escapes detection.

Applications requiring detection must add their own checksum to the payload before encoding. Note that a checksum is payload and consumes tokens; an 8 bit checksum costs roughly one extra token in EL-8 and a fraction of one elsewhere.

If the disjoint-partition composition of section 11.8 is adopted, a token from the wrong profile's dictionary becomes an immediate hard error. That is not a checksum and does not detect substitution *within* a profile, but it is the only structural error detection EL-8 would have at all, which is a point in that scheme's favor.

### 13.6 Cross-dictionary confusion
Because a phrase does not carry its own profile or dictionary identity, decoding with the wrong dictionary can silently produce different, well-formed output. Checking the dictionary fingerprint (section 11.7) against a trusted expected value is the mitigation and should be considered a requirement rather than an option for any application where the dictionary may vary. Merely recalculating the fingerprint stored inside an untrusted replacement file detects damage but not deliberate substitution.

Profile confusion specifically — decoding an EL-8 phrase as EL-14, say — is a subset of this risk, and its severity depends on the undecided composition question in section 11.8. Under a nested-subset composition it is maximally likely, since every EL-8 phrase is a well-formed EL-14 token run. Under a disjoint partition it becomes impossible.

---

## 14. Design Rationale Summary

- A single parameterized structure covers all profiles, so the low-complexity profiles are not a separate format but the same format with `w` chosen differently
- `w = 14` provides high entropy density while still allowing selection from common English root words
- Byte aligned input restricts remainders to multiples of `gcd(8, w)`, which is what makes the trim dictionary affordable
- Choosing `w` as a multiple of 8 eliminates trimming entirely, at a density cost — this is the EL-8 and EL-16 trade
- Choosing `w = 12` eliminates all but nibble-boundary splits and shrinks the trim dictionary by a factor of 20, at a 14% density cost
- The absence of a length header simplifies decoding by embedding the required remainder information in the final trim token
- Mechanically derived dictionaries make the selection criteria auditable and the artifacts reproducible, which hand curation cannot provide at 21,844 tokens

---

## 15. Implementation Order

The intended sequence, chosen so that each step adds exactly one mechanism:

1. **EL-8** — dictionary loading, tokenization, normalization, round-trip harness, test vectors. No bitstream, no trimming.
2. **EL-12** — bitstream packer, remainder arithmetic, trim dictionary, alignment validation. Every EL-14 mechanism, at a fifth the dictionary size and with hand-checkable intermediate values.
3. **EL-14** — the reference profile. No new mechanism, only scale.
4. **EL-16 and non-English dictionaries** — no new mechanism beyond glyph-based segmentation.

Implementations should state which profiles they support. An implementation that supports EL-12 and EL-14 should share one code path parameterized by `w`, not two.

---

## 16. Current Status

This draft describes the complete functional behavior of the encoding and decoding processes, the profile family, and the structural requirements for dictionary construction. Further work includes:

- The dictionary derivation tools described in section 11, which do not yet exist
- Finalized dictionary selection for each profile
- Published test vectors per profile, including empty payload, every remainder class, and known-invalid sequences. These live in `tests/vectors/` as language-agnostic JSON, shared by all implementations.

  Two points make this cheaper than it looks. First, vectors expressed as **token index sequences** rather than phrases are dictionary-independent, since every rule in sections 4.3, 5.1, 6 and 7 is defined over indices — so the entire conformance suite for the encoding mechanism can be authored before any dictionary exists, and implementation need not wait on derivation. Second, **payloads of N = 0 through 7 bytes exercise every remainder class in every defined profile**: EL-14's remainder cycles with `N mod 7` (0, 8, 2, 10, 4, 12, 6), and the EL-12 and EL-16 cycles are subsumed. Eight small payloads give complete trim coverage across the family.

  Negative vectors are required as well — at least one rejected sequence per validation rule D1 through D5. Silent acceptance of malformed input is the failure mode this format can least afford.
- Optional error detection schemes
- Optional interface bindings for various programming languages
- Optional compression or pre-processing guidelines
- Final LXJ schema, fingerprint input, and benchmark-selected LXB byte layout; see `data/dict/FORMAT.md`
- License decision and pinned releases for the candidate derivation inputs recorded in `data/dict/SOURCES.md`
- **Resolution of the cross-profile dictionary composition question in section 11.8** — nested subset, independent optimization, or disjoint partition. This is the largest undecided item in the specification. It must be settled before `eldict-select` is written, and it changes what sections 3.8, 13.5, and 13.6 are allowed to promise.
