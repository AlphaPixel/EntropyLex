# EntropyLex specification

This document defines EntropyLex, a reversible encoding that converts a sequence of bytes into human-language tokens and converts those tokens back into exactly the original bytes. It defines the bit-level mapping, profile sizes, dictionary structure, written phrase rules, validation requirements, dictionary construction process, and considerations for languages other than English.

EntropyLex is an encoding, not encryption. It does not make a payload secret. It also does not include a checksum to detect arbitrary changes or an error-correcting mechanism that can reconstruct damaged data. Its human-usability goals depend on dictionary quality; the exact conversion between bytes and token indices does not.

Nothing in this document is final. Items marked **(provisional)** are working proposals that may change after measurements and initial implementations. In requirement statements, **must** means required for compatibility, **should** means recommended unless a documented reason justifies another choice, and **may** means optional.

---

## 1. Purpose

EntropyLex represents binary data as tokens intended to be easier for people to read, write, speak, and remember than strings such as hexadecimal, Base32, or Base58. Example payloads include keys, identifiers, and hash or checksum values. Every payload bit must be recoverable exactly.

EntropyLex is intended to:

- Represent any sequence of bytes without ambiguity
- Represent `w` payload bits with each normal token, where `w` is fixed by the selected profile
- Allow decoding without an explicit length header
- Reconstruct a whole number of bytes
- Guarantee a deterministic round trip: the same input produces the same tokens, and decoding those tokens reproduces the input
- Use one set of rules for profiles with different token widths

---

## 2. Core Concepts

The following terms are used throughout this specification:

- **Payload:** the input bytes being represented.
- **Bit:** a binary digit, either 0 or 1.
- **Byte:** eight bits. EntropyLex accepts only whole bytes as input.
- **Token:** one item emitted by the encoder. An English token is a word; another dictionary may use a character or a fixed multi-character item.
- **Dictionary:** an ordered list that maps numerical indices to tokens, plus the rules used to recognize those tokens in writing.
- **Index:** a token's numerical position in the dictionary. Encoding maps bits to an index and then to a token; decoding maps the token back to the same index and bits.
- **Profile:** the named set of numerical parameters for an EntropyLex variant. Its main parameter is normal-token width `w`.
- **Normal token:** a token that represents exactly `w` bits and may appear anywhere in a phrase.
- **Trim token:** a token that represents fewer than `w` bits and may appear only at the end of a phrase.
- **Canonical:** the one required form used when output must be identical across implementations. A decoder may accept additional written forms where this specification explicitly permits them.
- **Deterministic:** producing the same output whenever all inputs and settings are the same.
- **UTF-8:** the standard byte encoding used here for Unicode text.
- **Hash:** a fixed-size value calculated from arbitrary bytes. EntropyLex uses SHA-256 hashes for file checksums and dictionary fingerprints; a hash is an identifier, not encryption.

### 2.1 Input data
EntropyLex accepts a sequence of bytes. Every payload length is therefore a multiple of eight bits.

Let the input size be `N` bytes. The payload then contains `8N` bits.

An empty payload (N = 0) is legal and encodes to an empty token sequence.

### 2.2 Output tokens
EntropyLex produces a sequence of dictionary tokens. There are two kinds:

1. Normal tokens
2. Trim tokens (used only for the final token)

A token is the smallest item that a dictionary maps to a value. An English dictionary uses one word per token. A character-based dictionary may use one character or a fixed multi-character item; see section 12.

### 2.3 Token width
Each normal token represents exactly `w` bits, where `w` is fixed by the chosen profile. There are `2^w` possible `w`-bit values, so the normal part of the dictionary must contain exactly `2^w` distinct tokens.

### 2.4 Byte alignment constraint
Dividing `8N` payload bits into groups of `w` may leave fewer than `w` bits at the end. Only certain remainder lengths can occur.

Let `g = gcd(8, w)`, the greatest common divisor of 8 and `w`. Let `r = 8N mod w`, where `mod` means the remainder after division. Every possible `r` is a multiple of `g`:

```
r ∈ { 0, g, 2g, ..., w - g }
```

The profile needs trim tokens only for the nonzero values in this set. Those values therefore determine the size of the trim part of the dictionary.

For the reference profile, `w = 14`, `g = gcd(8, 14) = 2`, so `r ∈ {0, 2, 4, 6, 8, 10, 12}`.

---

## 3. Profiles

EntropyLex defines several profiles that differ in normal-token width `w`. All profiles use the same bit order, trim-token method, and encoding and decoding algorithms. A profile is named `EL-<w>`; for example, EL-12 uses `w = 12`.

### 3.1 General profile arithmetic

For token width `w` and eight-bit bytes:

```
g            = gcd(8, w)
remainders R = { 0, g, 2g, ..., w - g }
normal count = 2^w
trim count   = Σ 2^r  for r ∈ R, r > 0
             = (2^w - 2^g) / (2^g - 1)
total count  = 2^w + trim count
tokens for N bytes = ceil(8N / w)
```

`Σ` means “sum the following expression over the stated values.” `ceil(x)` means round `x` up to the next integer. Trimming disappears only when `w` divides eight evenly, as EL-8 does. A width larger than eight can still leave a remainder even when it is a multiple of eight: EL-16 leaves eight bits when the payload contains an odd number of bytes.

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

The arithmetic also permits odd widths such as 11 or 13. For an odd width, `g = 1`, so every remainder from 0 through `w-1` can occur and each nonzero length needs its own trim group. No odd-width profile is currently defined because this produces the largest possible trim portion for a given `w`, without an identified benefit that offsets that cost.

### 3.3 Token count comparison

Tokens required for representative payloads:

| Payload | EL-8 | EL-12 | EL-14 | EL-16 |
|---|---|---|---|---|
| 16 bytes (128 bit key)      | 16 | 11 | 10 | 8  |
| 32 bytes (256 bit key)      | 32 | 22 | 19 | 16 |
| 60 bytes (8x8 4-color GIF)  | 60 | 40 | 35 | 30 |

### 3.4 EL-8: one byte per token

EL-8 is the simplest defined profile:

- Token width equals byte width. Each byte maps to one token and no token crosses a byte boundary.
- `r` is always 0, so the trim dictionary is empty and the trim steps in sections 4.2, 6.4, and 7.2 do nothing.
- The dictionary contains 256 entries. This permits strict requirements for differences in spelling and pronunciation, and complete human review is practical.
- Decoding always produces a whole number of bytes. The byte-alignment check in section 7.3 therefore detects no EL-8 errors; see section 13.5.

The tradeoff is phrase length: an EL-8 phrase for a 256-bit key contains 32 tokens, compared with 19 under EL-14.

EL-8 is the recommended first implementation target. It exercises dictionary loading, dividing written phrases into tokens, written-form normalization, and encode-then-decode tests without groups that cross byte boundaries or trim-token handling. The Pretty Good Privacy (PGP) Word List is an earlier system that also maps byte values to words and provides a useful comparison.

### 3.5 EL-12: whole- or half-byte boundaries

EL-12 is the intermediate step:

- `g = 4`, so a token boundary occurs only at a whole byte or after four bits of a byte. Four bits are sometimes called a **nibble**. Because one hexadecimal digit also represents four bits, intermediate values are straightforward to compare with hexadecimal output.
- Only two nonzero remainders exist, `r ∈ {4, 8}`, so the trim portion contains 272 entries rather than EL-14's 5,460.
- The total dictionary contains 4,368 entries. Reviewed source-list sizes suggest that familiar, short words should be available at this scale, but the selection tools described in section 11 must confirm that enough remain after all filters.
- Each EL-12 normal token represents 12 bits, or about 86 percent of EL-14's 14 bits, while its complete dictionary is one-fifth the size.

EL-12 is the recommended second implementation target. It introduces groups that cross byte boundaries, remainder calculation, trim-token lookup, and validation that the result contains a whole number of bytes. EL-14 uses the same mechanisms with a larger dictionary and more possible remainder lengths.

### 3.6 EL-14: the reference profile

EL-14 is the target that minimizes phrase length among the English profiles currently considered plausible. It is the default profile when a discussion uses “EntropyLex” without further qualification. Its suitability depends on finding enough familiar English base words after all selection filters; that has not yet been established. See section 15 for the implementation order.

EL-14 has six groups of trim tokens, one for each possible nonzero remainder length. It requires 5,460 trim tokens and 21,844 tokens in total.

### 3.7 EL-16 and wider

EL-16 is not currently practical for single-word English tokens because it needs 65,792 distinct tokens. Its bit handling is simpler than EL-12 or EL-14: normal tokens contain exactly two bytes, and it needs only one trim group for an eight-bit remainder. Languages or writing systems with larger token inventories may provide more candidates, but section 12 explains why familiarity remains a limiting factor. No practical EL-16 dictionary has been demonstrated.

### 3.8 Profile and dictionary identification

A token sequence does not describe itself. Correct decoding requires knowing both:

1. The profile (`w`), and
2. The exact dictionary used

The phrase itself does not normally contain either value. The surrounding application, file, or protocol must provide them separately, just as it must say whether a character string uses Base32 or Base58. A dictionary file therefore declares its profile, and implementations must report its dictionary fingerprint; see section 11.7. The fingerprint is a SHA-256 identifier calculated from the ordered token mapping and the rules used to recognize written tokens. Matching fingerprints mean that two dictionary files interpret phrases the same way, even if one is JSON and the other is binary.

Implementations must not attempt to auto-detect the profile from a token sequence. A sequence valid under one dictionary may be valid but different under another.

One exception is under consideration. Scheme C in section 11.8 would use completely separate token sets for different profiles. Any token could then identify its profile within that particular family of dictionaries. The empty phrase would remain ambiguous because it contains no token. This decision has not been made, so implementations cannot rely on it. Even if adopted, a token would identify only the profile, not the exact index mapping or written-token rules; checking the expected fingerprint would still be necessary.

---

## 4. Dictionary Structure

A dictionary contains a normal part followed by a trim part. “Normal dictionary” and “trim dictionary” below refer to these two parts of one ordered dictionary, not separate files.

### 4.1 Normal dictionary
- Size: 2^w tokens
- Index range: 0 to 2^w - 1
- Each normal token represents one `w`-bit value
- Used for all tokens except, possibly, the final token

### 4.2 Trim dictionary
A trim token is used only as the final token, to encode the remaining bits. For each possible nonzero remainder length `r`, the trim part must uniquely represent every possible `r`-bit value. That requires `2^r` tokens in the group for that remainder length.

The index of each trim token must identify both:

1. The remainder size `r`
2. The value of the final `r` payload bits

Trim tokens may appear only as the final token. The normal and trim parts cannot share a token, and no token may occur twice. The decoder can therefore determine a token's kind from its dictionary index without a separate marker inside the phrase.

For EL-8 the trim part is empty, so these requirements have no effect.

### 4.3 One continuous index range **(provisional)**

LXJ and LXB store one ordered token list. Normal and trim tokens therefore use one continuous range of indices:

- Indices `0 .. 2^w - 1` are the normal tokens, index equal to the encoded value.
- Trim tokens follow, grouped from shortest to longest remainder. `trim_base(r)` is the number of trim tokens in all preceding groups. For remainder length `r` and numerical value `v`:

```
index(r, v) = 2^w + trim_base(r) + v
trim_base(r) = Σ 2^j  for j ∈ R, 0 < j < r
```

For EL-14 the trim group starting positions are:

| `r` | group size | `trim_base(r)` | index range |
|---|---|---|---|
| 2  | 4    | 0    | 16384 .. 16387 |
| 4  | 16   | 4    | 16388 .. 16403 |
| 6  | 64   | 20   | 16404 .. 16467 |
| 8  | 256  | 84   | 16468 .. 16723 |
| 10 | 1024 | 340  | 16724 .. 17747 |
| 12 | 4096 | 1364 | 17748 .. 21843 |

For EL-12:

| `r` | group size | `trim_base(r)` | index range |
|---|---|---|---|
| 4 | 16  | 0  | 4096 .. 4111 |
| 8 | 256 | 16 | 4112 .. 4367 |

For EL-16, the single trim group `r = 8` occupies indices 65536 .. 65791.

A decoder therefore needs one mapping from written tokens to indices. The index range tells it whether the token is normal or trim and, for a trim token, how many bits and which value it represents.

---

## 5. Bit Order and Canonical Forms

### 5.1 Bit order
Read the payload as one bit sequence in **most-significant-bit-first** order. Bit 0 of the sequence is the highest-value bit of the first byte; bit 7 is that byte's lowest-value bit; bit 8 is the highest-value bit of the second byte, and so on.

A token's numerical value is formed in stream order: the first bit read has the highest place value. The `r` bits represented by a trim token are interpreted the same way, producing a value from 0 through `2^r - 1`. “Right-aligned” means that if this value is stored in a wider integer, it occupies the lowest-value `r` bit positions and all higher positions are zero.

This ordering applies to all profiles and is not configurable.

### 5.2 Canonical phrase form
The canonical written form is the one form an encoder must produce. It places the tokens in order with one ordinary space (ASCII U+0020) between them and no space before the first or after the last. Tokens use the letter case and Unicode Normalization Form C (NFC) recorded by the dictionary. NFC gives canonically equivalent Unicode character sequences one standard representation.

### 5.3 Input normalization for decoding **(provisional)**
For human-entered phrases, a dictionary may permit more than the canonical single-space form. The dictionary must record the exact recognition steps and their order; implementations must not rely on whatever a programming language happens to classify as whitespace, a hyphen, or a letter-case equivalent.

The intended English behavior is:

1. Convert the complete input to Unicode NFC.
2. Apply the letter-case conversion named by the dictionary. The first English dictionary is expected to use ASCII lowercase, avoiding language-dependent Unicode case rules.
3. Treat any nonempty run of declared separator characters as one token boundary. The initial separator set is expected to include ASCII space, tab, carriage return, line feed, and hyphen-minus; the final set remains to be decided.
4. Ignore separators before the first token and after the last token.
5. Look up each resulting token without further alteration.

The final format must list separator code points explicitly and name an exact case-conversion algorithm. A token cannot contain a character that the same dictionary treats as a separator.

Decoders must not silently correct spelling or replace an unknown token with the closest known token. If a token is absent from the dictionary, decoding fails. A separate, clearly labeled suggestion feature is permitted, but the normal decoding operation must never apply a suggestion automatically.

For a dictionary whose tokens each contain the same fixed number of written units (see section 12), a phrase may omit separators. The dictionary must say whether one unit means one Unicode code point or one extended grapheme cluster, which is a base character together with any combining marks displayed as one user-perceived character. Its decoder divides the phrase by that declared unit and count rather than relying on spaces. The specification must settle this choice before a separator-free dictionary is published.

### 5.4 Empty payload
A zero-byte payload encodes to a zero-token phrase, whose canonical written form is the empty string. A zero-token phrase decodes to a zero-byte payload. If a decoder accepts leading and trailing separators under section 5.3, input containing only those separators also normalizes to the legal empty phrase.

---

## 6. Encoding Process

### 6.1 Read input bytes as one bit sequence
Read the input bytes consecutively as one sequence of `8N` bits, using the most-significant-bit-first order from section 5.1.

### 6.2 Emit complete `w`-bit tokens
While at least `w` bits remain:

- Read the next `w` bits
- Interpret them as a value between 0 and 2^w - 1
- Map the value to a normal token and append it to the output sequence

### 6.3 Determine remainder
After consuming all complete `w`-bit groups, calculate the number of bits left:

```
r = (8N mod w)
```

`r` is always a member of the profile's remainder set `R`.

### 6.4 Emit final trim token
If `r = 0`:

- No trim token is required
- Encoding ends on the last normal token

If `r > 0`:

- Read the remaining `r` bits
- Use the pair `(r, final_value)` to select the appropriate trim token
- Append the trim token as the final token

Encoding ends after emitting the final trim token.

In EL-8, `r` is always 0 and this step never executes. More generally, trimming is unnecessary when `w` divides eight evenly. EL-16 does not meet that condition: an odd number of input bytes leaves an eight-bit remainder.

---

## 7. Decoding Process

Given a sequence of tokens, first handle the empty case: return a zero-byte payload without attempting to find a final token. Otherwise continue with sections 7.1 through 7.4.

### 7.1 Decode all but the last token
For each token except the final one:

- Look up its index in the complete dictionary
- Reject the sequence if that index belongs to the trim range
- Otherwise append the normal token's `w` bits to the reconstructed bit sequence

This enforces D2 below: a trim token cannot appear before the final position.

### 7.2 Decode the final token
If the final token is a normal token:

- It contributes `w` bits

If the final token is a trim token:

- Identify its remainder size `r`
- Extract the `r` bits encoded by that token
- Append those `r` bits to the reconstructed bit sequence

### 7.3 Validation rules
A decoder must reject a sequence unless all of the following hold:

- **D1** Every token is present in the dictionary for the declared profile.
- **D2** No trim token appears in any position other than the final one.
- **D3** The total reconstructed bit length is divisible by 8.
- **D4** If the final token is a trim token with remainder `r`, then `r ∈ R` and `r > 0`.
- **D5** An empty sequence is valid and yields an empty payload.

D3 is the encoding's only check that can detect some changes to an otherwise valid token sequence. It is not a checksum and does not detect every error; see section 13.5.

### 7.4 Convert the reconstructed bits to bytes
Group bits into bytes, MSB-first, and output the original binary data.

---

## 8. Trimming Behavior Summary

The input contains whole eight-bit bytes and the encoder emits `w`-bit normal tokens, so:

- Remainder sizes are always multiples of `g = gcd(8, w)`
- Maximum remainder is `w - g`
- Remainders not divisible by `g` never occur
- Only the nonzero members of `R` require trim tokens

| Profile | Remainder step `g` | Max remainder | Trim groups | Trim tokens |
|---|---|---|---|---|
| EL-8  | 8 | —  | 0 | 0    |
| EL-12 | 4 | 8  | 2 | 272  |
| EL-14 | 2 | 12 | 6 | 5460 |
| EL-16 | 8 | 8  | 1 | 256  |

Trim tokens must be distinct from the normal tokens in the same dictionary.

The design considered storing the payload length in a required header. Even a one-byte header would add data to every payload. Trim tokens instead increase dictionary size and identify the final bit count only when needed. EntropyLex currently chooses that tradeoff.

---

## 9. Dictionary Size Summary

| Profile | Normal | Trim | Total | English feasibility |
|---|---|---|---|---|
| EL-8  | 256   | 0    | 256   | Far below reviewed source-list sizes; complete human review is practical |
| EL-12 | 4096  | 272  | 4368  | Below reviewed common-word source-list sizes; final survivors remain to be measured |
| EL-14 | 16384 | 5460 | 21844 | Raw open sources are large enough; survival after familiarity and distance filters remains to be measured |
| EL-16 | 65536 | 256  | 65792 | Not feasible for single-word English tokens; see section 12 |

Normal tokens may appear many times in a phrase, while a trim token appears at most once and only at the end. The most familiar and easiest-to-distinguish selected words should therefore receive normal indices. Less-frequent but still acceptable words may receive trim indices.

---

## 10. Word Selection Criteria

Word choice is not fully specified. The governing principles are:

- Distinct pronunciation to reduce confusion in speech
- Exclude or minimize homophones (different words pronounced the same) and near-homophones
- Easy, predictable spelling
- Broad familiarity across regional and social varieties of the language
- Avoidance of culturally sensitive, offensive, or specialized terms
- Prefer short to medium length words
- Prefer pronunciation patterns that differ clearly from every other selected token
- Prefer words with unrelated meanings, so neighboring tokens are less likely to resemble a sentence and near-synonyms are less likely to be substituted from memory

The final dictionary for a given profile must contain exactly the counts in section 9 for that profile — no more and no fewer.

Section 11 describes the intended repeatable software process for applying these principles.

---

## 11. Dictionary Derivation

Dictionaries are not to be selected entirely by hand, except possibly EL-8, whose 256 words are small enough for complete manual review. Software will **derive** them from recorded source datasets. The tools do not exist yet; this section describes the required process so another implementation can repeat it, reviewers can inspect every removal and decision, and a changed rule can be applied consistently.

The objective is to select exactly the required number of words while making the easiest-to-confuse selected pair as different as possible in spelling, pronunciation, and meaning. Every selected word must also meet a minimum usage-frequency requirement, called the familiarity floor.

### 11.1 Planned tool sequence **(provisional)**

The tools form a pipeline: each numbered stage consumes recorded inputs and produces recorded output for the next stage. They will live under a `tools/` directory:

| Stage | Tool | Responsibility |
|---|---|---|
| 1 | `eldict-ingest`  | Load word, frequency, pronunciation, and grammatical data from different sources into one consistently formatted candidate table |
| 2 | `eldict-filter`  | Apply rules that always exclude a candidate and record the reason for every removal |
| 3 | `eldict-score`   | Calculate numerical spelling, pronunciation, and meaning comparisons between candidates |
| 4 | `eldict-select`  | Mark pairs that are too similar and choose the required number of nonconflicting tokens |
| 5 | `eldict-emit`    | Assign normal and trim roles and indices, then write the canonical LXJ file |
| 6 | `eldict-compile` | Convert LXJ into the optional LXB binary form without changing its mapping, rules, or fingerprint |
| 7 | `eldict-verify`  | Check file structure and, when source inputs are supplied, repeat quality checks and write a report |

Every stage must be deterministic. Byte-for-byte identical source files and identical settings must produce the same ordered mapping and dictionary fingerprint. The canonical LXJ writer and eventual LXB compiler must also produce byte-for-byte identical files.

### 11.2 Stage 1 — ingest

Each source must be fixed to an immutable release or commit so later runs use exactly the same data. LXJ records this source history, also called **provenance**, in structured fields:

- A starting word list with usage frequency measured from one or more large text collections, called corpora
- A pronunciation dictionary that represents each word as a sequence of speech sounds, called phonemes
- Data connecting grammatical variants, called inflections, to a base form, called a lemma—for example, “walks” and “walked” to “walk”
- A model that represents word meanings numerically so related words can be detected

Each candidate record contains at least the written word, its base form, usage-frequency rank, pronunciation, syllable count, and grammatical category such as noun or verb.

The final loading specification must define how records from different sources are joined. It must cover spelling variants, conflicting grammatical labels, missing pronunciations or frequency values, several pronunciations for one spelling, dialect labels, and ties. Source order or a library's default map iteration order must never decide the result implicitly.

Candidate sources, license terms, published sizes, and possible uses are tracked in [`data/dict/SOURCES.md`](data/dict/SOURCES.md). Download availability alone does not grant permission to redistribute the data or a derived dictionary. Before use, record the exact release or commit, the SHA-256 checksum of the downloaded file, all license obligations, and the number of records actually loaded.

### 11.3 Stage 2 — required exclusion rules

Apply these exclusion rules before comparing candidate pairs:

- Length outside the accepted band **(provisional: 3 to 8 characters)**
- Characters outside ASCII, the basic Latin character set, for English dictionaries
- Proper nouns, brand names, and place names
- Grammatical variants whose base form is already a candidate, such as plurals, conjugated verbs, and comparatives
- Words on an offensive, slur, or sensitive-topic blocklist
- Words on a blocklist of specialized vocabulary
- Sets of homophones—different spellings with the same pronunciation; retain at most one member or remove the entire set
- Homographs—identically spelled words—with substantially different pronunciations, such as “read,” “lead,” and “wound”
- Words whose spelling is not predictable from pronunciation, or vice versa, beyond a configured threshold
- Frequency below a familiarity floor

Record every rejected candidate and the rule that rejected it. This log is a required pipeline output so reviewers can inspect the decision, not temporary diagnostic output.

Each exclusion must become an exact rule before the first dictionary is built. For example, the recorded settings must decide whether a same-sounding group keeps its most familiar member or removes the whole group; “unpredictable spelling” must be replaced by a named measure and limit; and every blocklist must have a recorded version and matching method. The alternatives above describe design choices, not permission for implementations to choose differently without recording the choice.

### 11.4 Stage 3 — distance metrics **(provisional)**

Compare every surviving candidate with every other candidate using three groups of measures. A larger distance means the pair is less likely to be confused under that measure.

**Spelling distance**
- Damerau-Levenshtein distance: the number of single-character insertions, deletions, replacements, or adjacent-character swaps needed to change one word into the other
- A variant that treats replacements involving nearby keyboard keys as easier mistakes and therefore as a smaller distance
- Length of the shared beginning or ending, to avoid selecting many words that differ only near one end

**Pronunciation distance**
- Edit distance over speech-sound sequences, with sounds made using similar mouth and vocal-cord positions treated as closer; for example, /b/ and /p/ are closer than /b/ and /s/
- Immediate rejection when a broad pronunciation grouping such as Soundex or Metaphone places two words in the same group
- Syllable count and stressed-syllable pattern, so selected words do not all have the same rhythm

**Meaning distance**
- Cosine distance between numerical word representations. Cosine distance compares their directions rather than their absolute sizes; related words generally have smaller distance.
- Purpose: avoid near-synonyms such as “big” and “large,” which a person may substitute from memory, and reduce the chance that random neighboring tokens resemble a meaningful phrase that is remembered incorrectly.

The final scoring specification must define each formula, its numerical precision, how missing values are handled, and the allowed limit for every measure. It must also define tie-breaking so two programming languages cannot make different choices from equal scores.

### 11.5 Stage 4 — selection

The selection goal is sometimes called **maximin**: from `m` candidates, choose `n` tokens so the least-different selected pair is as different as possible under the combined spelling, pronunciation, and meaning measures.

The intended approach **(provisional)**:

1. Mark a conflict between any two candidates that fall below the allowed distance in **any** measure. Mathematically, candidates are nodes in a **confusability graph** and each conflict is an edge connecting two nodes. Two connected words cannot both be selected.
2. Order candidates from most to least familiar. Break equal frequency ranks by the normalized written token, using the same Unicode order defined for stage 5.
3. Examine candidates in that order. Add a word only if it conflicts with no word already selected. Graph theory calls the result an independent set; the operational rule is “keep the familiar word, then skip later words that are too similar to it.” This is a practical method, not proof of the mathematically best possible set.
4. The recorded settings provide an ordered list of allowed similarity-limit combinations, from strictest to weakest. Run the same selection pass under each combination until one produces at least the target count `n`.
5. Use the first `n` selected words from the strictest successful run. If no configured combination produces `n`, fail rather than inventing an unrecorded weaker limit.

The output must record every attempted combination and its survivor count. A reviewer must be able to see, for example, that the required count was reached only after the minimum pronunciation distance was reduced twice. The exact schedule and whether limits change individually or together remain to be designed, but they must be input data rather than hidden control flow.

The tools do not have to prove that no better set exists. They must use a precisely specified, repeatable method and publish measurements of the resulting set.

### 11.6 Stage 5 — assign normal and trim roles and indices

1. Sort the selected words from most to least frequent.
2. Assign the most frequent 2^w tokens to the normal dictionary and the remainder to the trim dictionary, per section 9.
3. Within each group, first apply the dictionary's required Unicode normalization and letter case. Then compare token strings one code point at a time, using the number Unicode assigns to each character. At the first difference, the smaller number sorts first; if one token is a complete prefix of the other, the shorter token sorts first. This produces one required order across implementations.
4. Assign indices using the continuous ranges in section 4.3.

Sorting by written form rather than frequency makes changes easier to review with ordinary file-comparison tools. Adding or removing a token changes a nearby run of indices instead of reordering tokens throughout the file when frequency measurements change slightly.

### 11.7 Stages 6 and 7 — dictionary files, source history, and verification

One dictionary has two representations:

- **LXJ** (`.lxj`) is the required, human-readable JavaScript Object Notation (JSON) representation and the authoritative source from which other forms are produced.
- **LXB** (`.lxb`) is an optional binary representation generated from LXJ for implementations whose measurements show that JSON loading is too slow.

For example: `entropylex-en-14-v1.lxj` and `entropylex-en-14-v1.lxb`, stored under `data/dict/`.

An LXB file is generated only from LXJ. A **file checksum** is a SHA-256 value calculated from every byte in one particular file, so the JSON and binary files have different checksums. They must nevertheless have the same **dictionary fingerprint** and assign exactly the same token to every index because they describe the same dictionary behavior. Implementations must support LXJ. LXB support is optional until its version 1 byte layout is final.

LXJ contains one token array. Each token's array position is its index from section 4.3. The file also records the profile, token counts, rules for recognizing written tokens, structured source information, a checksum of the selection settings, and the dictionary fingerprint. It must not repeat the index inside every token entry because that second copy could disagree with the array position.

The **dictionary fingerprint** is a SHA-256 identifier answering: “Would these dictionary files interpret every valid phrase in exactly the same way?” It covers the profile, index assignment, written-form rules such as case and separator handling, and every token in index order. It excludes JSON indentation, property order, comments, timestamps, source URLs, and quality reports because changing those does not change phrase decoding.

SHA-256 operates on bytes. The exact byte sequence supplied to it is called the **fingerprint input**. This specification will define that sequence separately from both JSON and binary file layout so formatting differences do not change dictionary identity. Each string will be preceded by its UTF-8 byte count, which makes adjacent values unambiguous: `["ab", "c"]` cannot produce the same input as `["a", "bc"]`. The exact version 1 recipe remains to be specified.

LXJ records labels for language and writing system. A label by itself does not affect decoding and therefore does not affect the fingerprint. Any behavior associated with it—such as letter-case matching, division into written tokens, or logical order for a right-to-left script—does affect the fingerprint.

Recalculating a fingerprint proves only that the file is internally consistent with the fingerprint it contains. Detecting replacement requires the application to compare it with an expected fingerprint obtained separately from a trusted release record or application setting. A fingerprint does not cryptographically prove who published the file. Releases may also publish a checksum for each exact LXJ and LXB file to detect accidental byte changes or replacement.

LXB version 1 must identify its format and size limits, contain every field needed for decoding, carry the shared fingerprint, and store tokens as their canonical UTF-8 bytes. Two layouts remain candidates. A **length-prefixed sequence** stores each token's byte count immediately before that token. An **offset table** stores the starting position of every token, followed by one block containing all token bytes. Measurements must compare their file size and loading speed before the layout is chosen. Version 1 will not store a prebuilt token-to-index lookup structure and will not compress the token bytes; each loader may build the lookup structure appropriate to its programming language.

Without using the original source datasets, `eldict-verify` can check the following properties directly from an LXJ or LXB file:

- Supported file format and fingerprint recipe
- Exact token counts for the declared profile
- Valid UTF-8, required Unicode normalization, and token-character rules
- No empty or duplicate tokens
- Normal and trim sets disjoint and canonically ordered
- Calculated fingerprint matches the stored fingerprint
- Encoding and decoding published dictionary-dependent test cases produces the expected results
- Every LXB section starts and ends within the file, and, when LXJ is also supplied, both files describe exactly the same dictionary

The final dictionary cannot by itself establish that its words are familiar, sufficiently different in pronunciation or meaning, or produced by the stated selection process. Checking those claims requires the exact source datasets and settings used for selection. The quality report must identify those inputs by name and checksum, show how many candidates survived each required filter, report the smallest measured distances, and list the selected pairs closest to the allowed limits.

The detailed working design and remaining decisions are in [`data/dict/FORMAT.md`](data/dict/FORMAT.md).

### 11.8 Whether profiles share tokens — UNDECIDED

Each profile is a separate encoding. An EL-8 phrase cannot be decoded as EL-12 merely because both use words.

What remains undecided is whether the EL-8, EL-12, and EL-14 dictionaries should share any tokens when they are selected from the same starting vocabulary. Three approaches remain possible.

Under every approach, one derivation run should preferably produce the complete family so one quality report can compare all profiles using the same source data and measures.

#### Scheme A — nested sets

Every EL-8 token also appears in EL-12, and every EL-12 token also appears in EL-14. In mathematical notation, `EL-8 ⊂ EL-12 ⊂ EL-14`. EL-8 uses the 256 highest-ranked members of EL-12, and EL-12 is similarly drawn from EL-14.

- Requires the fewest distinct words: 21,844, no more than EL-14 alone.
- Simplest to reason about, document, and verify.
- Greatest risk of profile confusion: every EL-8 phrase consists of valid EL-14 normal tokens. Decoding with the wrong profile may return different bytes without an error.
- The EL-8 set is inherited from a selection optimized for 21,844 tokens, rather than chosen as the 256 words that differ most clearly from one another.

#### Scheme B — independent selection

Select each profile separately using similarity limits appropriate to its required size.

- Best expected quality within each profile. A 256-token selection can require much larger differences between words than a 21,844-token selection can.
- The number of distinct words required is between 21,844 and 26,468, depending on how much the selections happen to overlap.
- Overlap is uncontrolled. Some phrases would happen to contain a word unique to one profile, while others would consist entirely of shared words. An application could not rely on the phrase to identify its profile.

#### Scheme C — separate token sets (profile-identifying phrases)

Allow no token to appear in more than one of the three profiles.

Every nonempty phrase would then identify its profile. Any one token determines its profile because it exists in exactly one profile dictionary. The empty phrase is the exception: it contains no token and is valid in every profile. For nonempty phrases, the separate profile information described in section 3.8 would become a confirmation rather than the only source of that information. A decoder using the wrong profile would reject the token instead of silently returning incorrect bytes.

The cost is the number of distinct suitable words required:

| Component | EL-8 | EL-12 | EL-14 | Total |
|---|---|---|---|---|
| Normal tokens | 256 | 4,096 | 16,384 | 20,736 |
| Trim tokens   | 0   | 272   | 5,460  | 5,732  |
| Per profile   | 256 | 4,368 | 21,844 | **26,468** |

Two observations on that total:

- `2^14 + 2^12 + 2^8 = 20,736` counts only normal tokens. The trim tokens must also be separate from one another and from every normal set. They add 5,732, bringing the complete requirement to **26,468 tokens**.
- This total exceeds the reviewed 25,000-entry 12dicts core list, although broader reusable sources contain many more base word forms and lexical entries. Only the 20,736 normal tokens need the highest familiarity; the 5,732 trim tokens may use less-frequent survivors because a phrase contains at most one trim token. Scheme C is therefore testable but not proven practical. The selection report must show survivor counts and every reduction in the similarity requirements.

A partial alternative would keep the normal sets separate but share trim tokens. The shared trim set would still have to avoid all three normal sets, so it saves only 272 tokens: 26,196 instead of 26,468. It would also fail to identify the profile for a one-byte payload, which becomes a single trim token in both EL-12 and EL-14. This alternative is therefore not recommended.

Giving EL-16 a completely separate English word set as well is not considered practical. The four profiles would require 92,260 distinct tokens in total.

#### Status

**No scheme has been selected.** Scheme A uses the fewest distinct words and is simplest. Scheme B is expected to produce the clearest words within each profile. Scheme C makes every nonempty phrase identify its profile but requires 26,468 acceptable distinct words and may require weaker similarity limits.

This decision must be made before `eldict-select` is implemented. It determines whether selection runs once followed by assignment to profiles, or runs separately for each profile while excluding words already used elsewhere. Sections 3.8 and 13.5 currently assume that tokens do not reliably identify a profile.

Scheme C's separate token sets would identify the *profile* only, and only within one family of dictionaries. They would not identify the exact index mapping or written-token rules. The fingerprint requirement in section 11.7 applies under every scheme.

---

## 12. Dictionaries for Other Languages and Writing Systems

Nothing in the bit-level encoding depends on English. Other dictionaries may change their tokens and written-token recognition rules while retaining the index assignment, trim behavior, and required counts. A dictionary for any defined profile must meet section 9's counts and section 4.2's rule that its normal and trim parts share no token.

### 12.1 How a writing system affects profile choice

Profile choice is limited by usable tokens, not by the raw number of entries in a database. Reviewed English sources range from about 25,000 common entries to more than 135,000 entries in broader lexical databases, but the count that survives familiarity and similarity filters is not yet known. EL-14 is therefore the largest English profile currently worth testing, not a demonstrated limit. A writing system with many familiar characters or compounds may support a different profile. EL-16 is relevant because, despite its much larger dictionary, it has one trim group rather than EL-14's six.

The important count is how many tokens intended users can reliably recognize, distinguish, and write from dictation. Computer character standards and comprehensive dictionaries include many rare items and therefore provide only an upper bound, not a usable token count.

A larger token inventory reduces the number of tokens per payload, but reaching that size usually requires rarer and less familiar tokens.

### 12.2 Dictionaries with one character per token

For Chinese, one Han character may serve as one token. “Han character” refers to the character family used in written Chinese and, in adapted form, in Japanese and other languages. One-character tokens make phrases visually compact and may allow a decoder to separate tokens without spaces. The dictionary must still declare the exact counting unit required by section 5.3; it cannot rely on an implementation's undefined idea of a character.

The constraint is the character inventory. The Table of General Standard Chinese Characters defines 8,105 characters (3,500 at level 1, 3,000 at level 2, 1,605 at level 3). Against the profile counts:

| Profile | Total tokens | What the raw character count permits |
|---|---|---|
| EL-8  | 256    | Well below the standard's count; filtered suitability still requires measurement |
| EL-12 | 4368   | Below the combined level 1 and 2 count; more raw candidates than required exist, but survivors are unmeasured |
| EL-14 | 21844  | Not possible with single standard characters; requires compounds |
| EL-16 | 65792  | Not possible with single characters; requires compounds |

EL-12 is therefore the largest defined profile whose raw count fits within the standard Chinese character inventory. Whether 4,368 sufficiently familiar and visually distinct characters survive must still be tested. EL-12 is also an early implementation target because its trim handling is simpler than EL-14's.

#### Japanese has a smaller widely known kanji inventory

*Kanji* are the Han-derived characters used in Japanese. Japanese standards distinguish characters taught or recommended for general use from the much larger set that computer encodings can represent. The general-use counts are relevant to human reliability; encoding-repertoire counts are not.

| Japanese character set | Count | Nature | Largest one-character profile |
|---|---|---|---|
| Kyōiku kanji | 1,026 | Taught in primary school; subset of jōyō | EL-8; 4.0 times its required count |
| Jōyō kanji | 2,136 | Government guide for general use, 2010 | EL-8; 8.3 times its required count |
| Jōyō + jinmeiyō | 3,000 | All kanji legally permitted in personal names | EL-8 only — short of EL-12's 4,368 |
| JIS X 0208 kanji | 6,355 | Encoding repertoire | Raw count reaches EL-12, but familiarity does not |
| JIS X 0213 kanji | ~9,980 | Encoding repertoire | Raw count reaches EL-12, but familiarity does not |

The figure “about 10,000 kanji” refers to the **Japanese Industrial Standard (JIS) X 0213 encoding repertoire**: a standardized set of characters that software can represent. It is not a list that readers are expected to know. It includes many characters that would fail EntropyLex's familiarity requirement, just as a comprehensive historical English dictionary includes many words unsuitable for a common-word encoding.

Using the sets intended to reflect ordinary literacy, a Japanese dictionary with one kanji per token appears limited to EL-8. Even the full 3,000-character personal-name inventory falls short of EL-12's 4,368-token requirement before visually similar or unfamiliar characters are removed. A Japanese EL-12 dictionary would need rarer kanji, tokens written in the phonetic hiragana or katakana scripts, or multi-character tokens.

#### Why Chinese reaches further

Chinese routinely writes most content words and many other elements with Han characters. Japanese also uses hiragana and katakana—two scripts that represent sounds and are collectively called **kana**—for grammar and much vocabulary. It can therefore function with a smaller commonly known kanji inventory. On the reviewed raw counts, Chinese can be investigated at EL-12 while Japanese single-kanji tokens do not reach EL-12.

### 12.3 Dictionaries with multi-character tokens

EL-16's 65,792-token requirement exceeds the reviewed single-character inventories, so it would require words or compounds containing multiple characters. A database may contain enough entries numerically while still leaving too few after suitability filters.

A standard Chinese dictionary, the 7th edition of 现代汉语词典 (*Xiandai Hanyu Cidian*, 2016), contains approximately 69,000 entries. EL-16 would require about 95 percent of them, leaving almost no entries available to discard for rarity, visual similarity, sensitivity, or other problems. A practical EL-16 dictionary would therefore need a much larger starting list and measured evidence that enough familiar entries survive.

No practical EL-16 dictionary has been established by the Chinese and Japanese counts reviewed here. Current counts support investigating one-character Chinese at EL-12. Multi-character tokens add these requirements:

- The decoder needs an unambiguous token boundary. Tokens may all contain the same number of declared Unicode units—for example, exactly two code points—or the written phrase must contain separators.
- Fixed two-unit tokens are the current recommendation **(provisional)** because a decoder can divide a separator-free phrase by count. The visual width may still vary by font and character, so fixed unit count does not promise equal display width.
- Larger target counts force selection further into the rare entries of the starting dictionaries, just as EL-14 may for English.

### 12.4 Selection criteria differ by writing system

The English criteria in section 10 give substantial weight to spoken confusion caused by homophones and accent variation. Other languages and writing systems have different sources of confusion:

- **Mandarin and Japanese contain many words or characters with the same pronunciation.** Their dictionaries should therefore give visual distinction more weight and spoken distinction less weight than the English dictionary. Tone in Mandarin and pitch accent in Japanese provide some distinction but do not remove the problem.
- **Visual distance** may use stroke count, shared components such as radicals, and overall character shape. A radical is a conventional component used to classify or construct a character.
- **Romanization distance** compares pronunciation written in the Latin alphabet—for example, pinyin with tone for Mandarin and rōmaji for Japanese. It remains relevant to typing and spoken communication.
- **Meaning distance** remains relevant. Neighboring characters can form an ordinary compound word, so selection should also reject token pairs that would encourage a reader to combine or replace them from memory.

For this reason, stage 3 of the derivation pipeline must allow each dictionary to define suitable comparison measures. The source loading, required filtering, selection process, index assignment, file generation, and verification stages remain the same.

### 12.5 Additional requirements for non-ASCII dictionaries

- Tokens are stored and compared in Unicode Normalization Form C (NFC), which gives canonically equivalent character sequences one standard representation. LXJ explicitly records the normalization and matching rules, and LXB carries the same behavior.
- A Chinese dictionary must declare whether it uses simplified or traditional characters. Two variant forms of the same character must not both appear. Converting one writing system to the other requires a separately selected dictionary and fingerprint because conversion can change distinctions between tokens.
- Letter-case conversion does not apply to writing systems without uppercase and lowercase. Decoders must apply only the case behavior declared by the dictionary.
- A dictionary for a right-to-left writing system must define its phrase as a logical sequence of tokens in encoding order. Display software may render that sequence from right to left, but a decoder must recover the logical order and must not simply reverse what appears on screen. Any required direction-control characters and their treatment must be specified before such a dictionary is published.
- Every dictionary records its language and writing system. The fingerprint covers the actual behavior needed to recognize and decode its phrases, such as normalization, case handling, division into character tokens, and logical token order. A descriptive language label alone is not treated as decoding behavior.

### 12.6 References for the figures in this section

The following sources support the raw counts above. A **general-use standard** describes characters people are expected to encounter or learn. An **encoding repertoire** lists characters a computer standard can represent, including rare characters. General-use standards are better starting points for EntropyLex, although their entries must still pass the dictionary filters.

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

The commonly cited “about 10,000 kanji” is the JIS X 0213 row: a computer encoding repertoire about 4.7 times the size of the jōyō general-use list. It does not measure how many characters a reader knows and cannot by itself establish a usable dictionary size.

Sources retrieved 2026-08-17. Before a figure is used to build a dictionary, it must be checked again against the primary standard document and its exact source version must be recorded in LXJ as required by section 11.7.

---

## 13. Security Considerations

### 13.1 Exact information preservation
The mapping is reversible: each valid phrase under a particular profile and dictionary corresponds to one byte sequence. For random data, a reversible change of representation preserves **information entropy**, meaning the amount of unpredictability measured in bits. This does not provide secrecy; anyone with the phrase, profile, and dictionary can recover the bytes.

### 13.2 Words do not reveal payload meaning
A token's ordinary-language meaning has no relationship to the payload bits it represents. The index mapping alone determines those bits. A phrase may accidentally resemble meaningful language, which dictionary selection attempts to reduce, but that resemblance does not describe the payload.

### 13.3 Human error
People may mispronounce a word, confuse it with a similar word, or copy it incorrectly. Section 11's dictionary-selection process aims to reduce those risks. A smaller profile dictionary can impose greater differences between selected words, although phrase length increases because each token represents fewer bits.

### 13.4 Payload length is visible
An observer can infer information about payload length from the number of tokens. In EL-8, token count equals byte count. In EL-12, EL-14, and EL-16, token count alone limits the payload to a small range of lengths. An observer who also has the dictionary can distinguish a final trim token from a normal token and calculate the exact bit length in every profile. Applications must decide whether revealing length is acceptable.

### 13.5 Error detection is limited and depends on the profile
EntropyLex does not include a checksum. Rule D3 can detect some inserted or removed tokens when the resulting bit count is not divisible by eight, but it cannot detect all changes:

- Replacing one valid token with another valid token of the same kind—normal or trim—is **never** detected in any profile.
- Inserting or deleting one normal token changes the bit count by `w`. The error is detected only when the new count is not divisible by eight. EL-12 and EL-14 therefore detect every single normal-token insertion or deletion; EL-8 and EL-16 detect none.
- A deleted trim token is detected unless its remainder `r` is itself a multiple of 8 — in EL-12 and EL-14, the `r = 8` case escapes detection.

An application that needs reliable error detection must append its own checksum to the payload before EntropyLex encoding and verify it after decoding. Those checksum bits become part of the encoded payload and may increase the phrase length.

If Scheme C from section 11.8 is adopted, a token from another profile becomes an immediate error. This detects a profile mismatch but is not a checksum and cannot detect replacement of one valid token with another from the same profile.

### 13.6 Cross-dictionary confusion
A phrase normally does not contain its profile or dictionary identity. Decoding it with the wrong mapping may silently produce different bytes. An application in which dictionaries can vary must compare the loaded dictionary's fingerprint with an expected value obtained from a trusted source. Recalculating the fingerprint stored inside an untrusted replacement file shows only internal consistency; an attacker could replace both the contents and stored fingerprint.

Using the wrong profile is one form of this problem. Under Scheme A, every EL-8 phrase uses words that also exist in EL-14, so a wrong-profile decoder may accept it. Under Scheme C, the separate token sets would make that mismatch fail immediately.

---

## 14. Design Rationale Summary

- One algorithm covers every profile by taking `w` and the associated counts as parameters; the simpler profiles are not separate formats.
- EL-14 represents more bits per English token than EL-8 or EL-12 while requiring fewer candidate words than EL-15 would. Its practical English vocabulary has not yet been demonstrated.
- Whole-byte input restricts remaining-bit lengths to multiples of `gcd(8, w)`, reducing the number of required trim groups.
- A `w` that divides eight evenly eliminates trim tokens because every normal-token boundary coincides with a byte subdivision. Among the defined profiles, only EL-8 has that property. EL-16 uses whole two-byte normal tokens but still needs one eight-bit trim group for odd-byte payload lengths.
- EL-12 permits boundaries only at multiples of four bits and needs 272 trim tokens, compared with EL-14's 5,460. It represents 12 rather than 14 bits per normal token.
- The final trim token records the remaining-bit length, so EntropyLex does not need a separate payload-length field.
- Software-derived dictionaries make selection decisions reviewable and repeatable at scales where purely manual selection would not.

---

## 15. Implementation Order

The intended sequence introduces complexity gradually:

1. **EL-8** — dictionary loading, division of written phrases into tokens, written-form normalization, encode-then-decode tests, and shared test cases. No token crosses a byte boundary and there are no trim tokens.
2. **EL-12** — bit groups that cross byte boundaries, remaining-bit calculation, trim-token lookup, and validation that decoding produces whole bytes. It introduces every EL-14 mechanism with a dictionary one-fifth the size.
3. **EL-14** — the reference profile. No new mechanism, only scale.
4. **EL-16 and non-English dictionaries** — token-recognition rules needed by other writing systems, such as dividing a separator-free phrase into fixed-length character tokens.

Implementations should state which profiles they support. An implementation supporting more than one profile should use the same encoding and decoding functions with profile parameters rather than duplicate those functions.

---

## 16. Current Status

This draft describes the current encoding and decoding behavior, profile family, and structural requirements for dictionaries. The following work remains:

- The dictionary derivation tools described in section 11, which do not yet exist
- Finalized dictionary selection for each profile
- Published test cases per profile, including an empty payload, every remainder group, and known-invalid sequences. These live in `tests/vectors/` as JSON shared by implementations in every programming language.

  Tests represented as **token index sequences** do not depend on selected words because sections 4.3, 5.1, 6, and 7 define the encoding in terms of indices. Encoding tests can therefore be written before a dictionary exists. Payload lengths from `N = 0` through `N = 7` bytes collectively exercise every possible remainder in every defined profile. For EL-14, the remainder sequence is 0, 8, 2, 10, 4, 12, 6, then back to 0; the shorter EL-12 and EL-16 cycles also occur within that range.

  Rejection tests are also required for invalid conditions D1 through D4. D5 requires a positive test showing that an empty sequence is accepted and produces an empty payload. A decoder must not silently accept malformed input.
- Optional error detection schemes
- Optional application interfaces for various programming languages
- Optional payload compression or preprocessing guidelines, performed before EntropyLex encoding and reversed after decoding
- Final LXJ JSON field structure, fingerprint input, and LXB byte layout chosen after measurements; see `data/dict/FORMAT.md`
- License decision and exact releases for the candidate source data recorded in `data/dict/SOURCES.md`
- **Resolution of the shared-token question in section 11.8**—nested sets, independent selection, or separate sets. It must be settled before `eldict-select` is implemented because it changes what sections 3.8, 13.5, and 13.6 can guarantee.
