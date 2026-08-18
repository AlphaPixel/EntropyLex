# EntropyLex dictionary files

EntropyLex can store one dictionary in two forms:

- **LXJ** (`.lxj`) is the required JSON form. LXJ version 1 is defined in this document and in [`lxj-v1.schema.json`](lxj-v1.schema.json).
- **LXB** (`.lxb`) will be an optional binary form generated from LXJ. Its byte layout is not yet defined.

An implementation can therefore begin with LXJ alone. LXB should be added only if measurements show that JSON loading is a meaningful cost.

## 1. LXJ version 1 at a glance

An LXJ file contains:

- the file-format version;
- a stable dictionary identifier and release number;
- language, writing system, and EntropyLex profile;
- exact rules for turning written input into tokens;
- one token array, in numerical index order;
- a fingerprint of everything that affects phrase interpretation; and
- records of the source material and selection settings.

The array position is the token's index. The file does not repeat an index beside each token, because two copies of the same fact could disagree.

This abbreviated example omits most tokens but shows the structure. The omission makes it explanatory JSON, not a valid LXJ file.

```json
{
  "format": {
    "name": "LXJ",
    "version": 1
  },
  "dictionary": {
    "id": "entropylex-en-8-test",
    "version": 1,
    "purpose": "test",
    "name": "EntropyLex English EL-8 test dictionary",
    "language": "en",
    "script": "Latn"
  },
  "profile": {
    "name": "EL-8",
    "normal_bits": 8,
    "remainder_bits": [],
    "normal_count": 256,
    "trim_count": 0,
    "total_count": 256
  },
  "recognition": {
    "unicode_version": "15.1.0",
    "normalization": "NFC",
    "case": "ascii-lower",
    "token_text": "unicode-scalars",
    "tokenization": {
      "kind": "separator-runs",
      "canonical_separator": "U+0020",
      "separators": ["U+0009", "U+000A", "U+000D", "U+0020", "U+002D"]
    }
  },
  "indexing": "continuous-v1",
  "tokens": ["able", "about", "...254 more tokens..."],
  "fingerprint": {
    "recipe": "LXFP-1",
    "sha256": "64 lowercase hexadecimal digits"
  },
  "provenance": {
    "sources": ["structured source records"],
    "selection": "a structured selection record"
  }
}
```

The complete EL-8 example is [`../../tests/fixtures/dict/entropylex-en-8-test-v1.lxj`](../../tests/fixtures/dict/entropylex-en-8-test-v1.lxj).

## 2. Which rules are authoritative

Three layers work together:

1. [`lxj-v1.schema.json`](lxj-v1.schema.json) checks the JSON structure, required fields, basic types, known profile values, and array size.
2. This document defines checks that JSON Schema cannot express reliably, such as Unicode normalization, duplicate object names, token sorting, and fingerprint calculation.
3. [`../../SPEC.md`](../../SPEC.md) defines what token indices mean during EntropyLex encoding and decoding.

A file is valid only if it satisfies all three. If they conflict, that is a specification defect; implementations should report it rather than guessing.

## 3. Field reference

Unknown properties are invalid in LXJ v1. This catches spelling mistakes and prevents one implementation from silently ignoring information another implementation treats as meaningful.

### 3.1 `format`

| Property | Required value | Meaning |
|---|---|---|
| `name` | `"LXJ"` | Identifies the file as an EntropyLex JSON dictionary. |
| `version` | `1` | Selects this exact JSON structure and its validation rules. |

An implementation must reject a format version it does not support.

### 3.2 `dictionary`

| Property | Rule | Meaning |
|---|---|---|
| `id` | 1–64 lowercase ASCII letters, digits, or internal hyphens | Stable machine-readable name shared by its releases; do not put the release number in this value. |
| `version` | integer from 1 through 4,294,967,295 | Release number for that identifier. |
| `purpose` | `"production"` or `"test"` | States whether the word list is intended as a reviewed release or only as test data. |
| `name` | 1–128 Unicode characters | Human-readable name. |
| `language` | lowercase language tag | Language identity, such as `"en"`. Use two or three lowercase ASCII letters, followed only when needed by hyphen-separated parts of two through eight lowercase letters or digits. |
| `script` | four-letter writing-system code | Use the ISO 15924 form: one uppercase and three lowercase ASCII letters, such as `"Latn"` for Latin script. |

`id`, `version`, `purpose`, and `name` label or describe the artifact. They do not change how a phrase decodes, so they are not part of the dictionary fingerprint. `language` and `script` are part of the fingerprint: changing either declares a different dictionary identity even if the token strings happen to match.

### 3.3 `profile`

LXJ v1 accepts the four profiles currently defined by the main specification. Every value on a row must match; a file cannot, for example, claim `EL-8` while declaring 12-bit normal tokens.

| `name` | `normal_bits` | `remainder_bits` | `normal_count` | `trim_count` | `total_count` |
|---|---:|---|---:|---:|---:|
| `EL-8` | 8 | `[]` | 256 | 0 | 256 |
| `EL-12` | 12 | `[4, 8]` | 4,096 | 272 | 4,368 |
| `EL-14` | 14 | `[2, 4, 6, 8, 10, 12]` | 16,384 | 5,460 | 21,844 |
| `EL-16` | 16 | `[8]` | 65,536 | 256 | 65,792 |

The repeated counts make the file easy to inspect. A loader must still compare them with the required values above and with the actual length of `tokens`.

### 3.4 `recognition`

These fields tell a decoder exactly how to recognize written input. LXJ v1 deliberately supports one simple form of tokenization: tokens separated by declared characters. A later format version can add separator-free character dictionaries without changing this definition.

`unicode_version` must be `"15.1.0"`. It pins the Unicode character data used by recognition so implementations do not silently rely on different library versions.

`normalization` must be `"NFC"`. Normalize the complete input using Unicode 15.1.0 Normalization Form C before any other recognition step. Unicode normalization stability means later conforming Unicode versions give the same result for characters assigned by Unicode 15.1.0, but a loader still validates the declared version rather than substituting an arbitrary one.

`case` is one of:

- `"none"`: do not change letter case.
- `"ascii-lower"`: replace each ASCII capital letter `A` through `Z` (U+0041 through U+005A) with its lowercase counterpart `a` through `z` by adding 32 to its code point. Leave every other code point unchanged. This is not locale-sensitive Unicode case folding.

`token_text` must be `"unicode-scalars"`. Each token must be a nonempty sequence of Unicode scalar values: Unicode code points other than the surrogate range U+D800 through U+DFFF. Its UTF-8 encoding must contain no more than 255 bytes. Each stored token must already be in NFC and unchanged by the declared `case` operation.

`tokenization` contains:

- `kind`, which must be `"separator-runs"`;
- `canonical_separator`, which must be `"U+0020"`, the ordinary ASCII space; and
- `separators`, a nonempty list of `U+` code-point labels in increasing numerical order.

Every separator label uses uppercase hexadecimal. Values through U+FFFF use exactly four digits. Larger values use five or six digits without a leading zero. A listed value must not exceed U+10FFFF and must be a Unicode scalar value. `U+0020` must be present. No stored token may contain a listed separator.

To divide input into tokens after normalization and case conversion:

1. Treat any nonempty run of listed separator characters as one boundary.
2. Ignore a separator run before the first token or after the last token.
3. Look up each remaining token exactly. Do not correct spelling or replace unknown words.

The canonical phrase produced by an encoder joins stored tokens with one U+0020 space and has no leading or trailing space. The empty token sequence has the empty string as its canonical phrase. Input containing only separators also becomes the empty token sequence.

The first English dictionary uses these separators:

| Label | Character |
|---|---|
| `U+0009` | horizontal tab |
| `U+000A` | line feed |
| `U+000D` | carriage return |
| `U+0020` | space |
| `U+002D` | hyphen-minus |

This list is exact. A decoder must not treat every character that its programming language calls whitespace as a separator.

### 3.5 `indexing` and `tokens`

`indexing` must be `"continuous-v1"`. The array position is the numerical token index:

- positions `0` through `normal_count - 1` are normal tokens;
- trim groups follow in the order given by `remainder_bits`; and
- a trim group for `r` remaining bits contains `2^r` tokens.

The normal part and the complete trim part are each independently sorted in increasing Unicode code-point order. Compare the first code point that differs; the smaller value sorts first. If one token is a complete prefix of another, the shorter token sorts first. The remainder-group boundaries divide the already sorted trim part into index ranges; they do not restart sorting.

The complete token array must contain no empty value and no duplicate. Each token must meet the `recognition` rules. LXJ v1 permits at most 65,792 entries and at most 255 UTF-8 bytes per token.

### 3.6 `fingerprint`

`recipe` must be `"LXFP-1"`. `sha256` contains exactly 64 lowercase hexadecimal digits, without a `sha256:` prefix. Section 6 defines its calculation.

### 3.7 `provenance`

Provenance means the recorded history of the word list. It does not affect decoding and is not part of the fingerprint.

Each object in `sources` has exactly these properties:

| Property | Meaning |
|---|---|
| `name` | Name of the source dataset or project-authored input. |
| `version` | Unchangeable release, exact repository commit, or local fixture version. |
| `role` | What the source contributed. |
| `location` | Download URL or path from the repository root where the exact input is obtained. |
| `license` | License identifier or concise license name. |
| `sha256` | SHA-256 of the exact source-file bytes, as 64 lowercase hexadecimal digits. |
| `record_count` | Number of relevant records in that source file. |

`selection` identifies how the sources became the ordered dictionary:

| Property | Meaning |
|---|---|
| `method` | Stable name of the selection procedure. |
| `settings_location` | Download URL or path from the repository root of the exact settings file. |
| `settings_sha256` | SHA-256 of that settings file's exact bytes. |

A valid file must contain at least one source record. These records make the claimed history precise; they do not prove that the claim is true. Reproducing the selection with the named inputs provides the stronger check.

## 4. Reading and validating an LXJ file

A conforming loader performs these checks in order:

1. Reject a file larger than 32 MiB (33,554,432 bytes).
2. Decode strict UTF-8. A byte-order mark is not permitted. Invalid UTF-8 is an error.
3. Parse exactly one JSON value. Reject duplicate property names, trailing non-whitespace data, comments, trailing commas, `NaN`, and infinities.
4. Validate the result against [`lxj-v1.schema.json`](lxj-v1.schema.json).
5. Apply the additional rules in this document: valid Unicode scalar strings, maximum UTF-8 token length, sorted separator values, required U+0020 separator, normalized tokens, no separator inside a token, no duplicate token, sorted index groups, and agreement between profile counts and token count.
6. Recalculate the LXFP-1 fingerprint and compare it with `fingerprint.sha256` using an exact byte comparison of the 32 decoded hash bytes.
7. Build the implementation's token-to-index lookup table. The file must not contain a precomputed lookup table.

Any failure rejects the whole dictionary. A loader must not repair it silently.

JSON Schema alone is not a complete LXJ validator. In particular, ordinary JSON Schema does not reliably reject duplicate object names or express the normalization, sorting, and fingerprint rules.

## 5. Canonical LXJ writing

JSON property order and whitespace do not affect loading or the dictionary fingerprint. Official EntropyLex writers nevertheless use one representation so separately written tools can produce byte-for-byte identical LXJ files.

The canonical representation is:

- strict UTF-8 without a byte-order mark;
- two spaces per indentation level;
- one property or array element per line;
- `: ` between a property name and its value;
- a comma after every property or array element except the last in its container;
- line-feed (U+000A) line endings and exactly one final line feed;
- integers written in decimal without a leading zero;
- object properties in the order shown in section 3 and in the full example; and
- source records in their listed input order; source order records history but does not affect decoding.

The property order, written from outer object to inner objects, is:

```text
top level:     format, dictionary, profile, recognition, indexing, tokens, fingerprint, provenance
format:        name, version
dictionary:    id, version, purpose, name, language, script
profile:       name, normal_bits, remainder_bits, normal_count, trim_count, total_count
recognition:   unicode_version, normalization, case, token_text, tokenization
tokenization:  kind, canonical_separator, separators
fingerprint:   recipe, sha256
provenance:    sources, selection
source:        name, version, role, location, license, sha256, record_count
selection:     method, settings_location, settings_sha256
```

For JSON strings, escape quotation mark as `\"` and reverse solidus as `\\`. Use the short escapes `\b`, `\t`, `\n`, `\f`, and `\r` for U+0008, U+0009, U+000A, U+000C, and U+000D. Write every other U+0000-through-U+001F control value as `\u00xx` with lowercase hexadecimal. Write `/` and every non-control Unicode scalar value directly as UTF-8 rather than using `\/` or `\u` escapes.

Empty objects and arrays are written as `{}` and `[]` on one line. A nonempty object or array starts its contents on the next line, indents each member by two additional spaces, and places its closing character on a new line at the parent's indentation.

A loader must accept any JSON representation that passes section 4; it must not require canonical whitespace or property order. Canonical writing is a reproducible-build rule, not an input restriction.

## 6. Dictionary fingerprint: LXFP-1

The dictionary fingerprint answers a narrow question:

> Do these artifacts declare the same language and writing system and interpret every valid phrase in the same way?

It ignores names, release labels, JSON formatting, source history, and selection reports. Those facts can change without changing token-to-index behavior. It includes the profile, language and script identity, written-input rules, index scheme, and every token in index order.

### 6.1 Small encoding rules

LXFP-1 first turns those values into one exact byte sequence, then calculates SHA-256 over that sequence.

- `U32(n)` is an unsigned 32-bit integer written as four bytes, highest-value byte first. Valid values are 0 through 4,294,967,295.
- `S(text)` is `U32(number of UTF-8 bytes)` followed by the strict UTF-8 bytes of `text`.
- A list starts with `U32(number of items)`, followed by each item in order.
- `CP(label)` parses a `U+` label and writes its numerical Unicode value as `U32`.
- `||` below means place the bytes on the right immediately after the bytes on the left.

Length prefixes make boundaries unambiguous. For example, the strings `"ab"`, `"c"` cannot be confused with `"a"`, `"bc"`.

### 6.2 Exact fingerprint input

Start with these 18 literal bytes, which are ASCII `EntropyLex-LXFP-1` followed by a zero byte:

```text
45 6e 74 72 6f 70 79 4c 65 78 2d 4c 58 46 50 2d 31 00
```

Append the following values in this exact order:

```text
S(profile.name)
U32(profile.normal_bits)
list of U32(profile.remainder_bits item)
U32(profile.normal_count)
U32(profile.trim_count)
U32(profile.total_count)
S(dictionary.language)
S(dictionary.script)
S(recognition.unicode_version)
S(recognition.normalization)
S(recognition.case)
S(recognition.token_text)
S(recognition.tokenization.kind)
CP(recognition.tokenization.canonical_separator)
list of CP(recognition.tokenization.separators item)
S(indexing)
list of S(tokens item)
```

The resulting bytes are the **fingerprint input**. Calculate SHA-256 over all of them. Write the 32-byte result as 64 lowercase hexadecimal digits in `fingerprint.sha256`.

The fingerprint field is not included in its own input. LXJ format details are also excluded, which allows a future LXB file carrying the same dictionary behavior to have the same fingerprint.

Recalculating the fingerprint proves that the file agrees with the fingerprint written inside it. It does not prove who published the file. An application that must detect replacement needs an expected fingerprint obtained separately from a trusted release record or application setting.

## 7. Relationship between LXJ and LXB

LXJ is authoritative: dictionary changes are made in or generated into LXJ first, and LXB is produced from it. A checksum over every byte differs between LXJ and LXB because their file bytes differ. Their dictionary fingerprints match because the behavior described in section 6 matches.

LXB remains open work. Its eventual first version must declare its byte order and size limits, contain every field covered by LXFP-1, carry the shared fingerprint, validate every section boundary, and store canonical UTF-8 token bytes. Measurements will decide whether it uses length-prefixed token strings or an offset table plus a string block.

## 8. What the artifact can prove by itself

Without retrieving original word sources, a verifier can check:

- all structure and value rules in sections 3 and 4;
- exact counts for the profile;
- valid and normalized token text;
- no empty or duplicate tokens;
- required index-group ordering;
- agreement between the stored and calculated fingerprint; and
- exact agreement in dictionary identity and decoding behavior between LXJ and a future LXB file.

The artifact alone cannot establish that its words are familiar, easy to pronounce, sufficiently different, or correctly screened for sensitive meanings. Those are claims about how the dictionary was selected. They require the named sources, settings, and a separate quality report.

## 9. Versioning and remaining work

LXJ v1 readers must reject unknown properties, unknown recognition algorithms, and unsupported versions. Adding an optional property to v1 would therefore be a breaking change and requires a new LXJ format version. A dictionary can receive a new `dictionary.version` without changing the LXJ format version.

The LXJ v1 structure, canonical writer, English recognition behavior, and LXFP-1 calculation are now defined. Remaining format work is:

1. publish valid and deliberately invalid conformance cases;
2. add any tokenization needed for separator-free writing systems in a later LXJ version; and
3. measure and define LXB only after the LXJ implementation supplies a baseline.
