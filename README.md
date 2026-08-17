# EntropyLex: A Human-Readable Encoding for Binary Data

EntropyLex is a reversible encoding: it converts a sequence of bytes into dictionary tokens and converts those tokens back into exactly the original bytes. The first planned dictionaries use common English words. Possible payloads include cryptographic keys, device identifiers, authentication tokens, and hash values.

In this document, **high-entropy data** means data whose bits are difficult to predict, as cryptographic keys should be. EntropyLex preserves those bits, but it does not encrypt them, hide them, or add error correction. Anyone who can read the phrase and has the correct dictionary can recover the bytes.

Hexadecimal, Base32, and Base58 represent bytes with letters and digits. They are compact and convenient for software, but long strings in those formats can be difficult to remember, distinguish visually, or communicate aloud. EntropyLex instead uses carefully selected words intended to be familiar and difficult to confuse in spelling or speech.

The encoding rules do not depend on English. English is simply the first planned language. EntropyLex defines several **profiles**, each of which specifies how many bits one ordinary token represents. Wider profiles use fewer words per payload but require larger dictionaries and more complicated handling of the final token. See [Profiles](#profiles).

## Why It Exists

EntropyLex solves three persistent problems:

1. Long strings of hexadecimal digits or apparently random characters are difficult to memorize or communicate accurately. EntropyLex does not depend on a particular estimate of human memory capacity; it depends only on words being easier to handle than dense character strings in some situations.

2. Many existing word-based systems belong to a particular application, such as a cryptocurrency wallet recovery format. Their words may include checksums or feed a later key-generation process. EntropyLex instead defines a direct, application-independent mapping between bytes and tokens.

3. The same encoding can represent any whole number of bytes. The application does not need a special EntropyLex variant for keys, identifiers, images, or other byte sequences.

Typical uses include memorized master keys, offline recovery secrets, pairing codes, identifiers, hash values, and transfers in which a person must read, write, or speak the data. EntropyLex can encode files too, although word phrases become unwieldy as the payload grows. For scale, the 60-byte sample GIF under `data/sample/input/` would require 35 EL-14 tokens.

## How It Works - Conceptual Overview

EntropyLex accepts a sequence of bytes. A byte contains eight bits, and a bit is a binary digit with value 0 or 1. The encoder reads those bits in a defined order and replaces groups of bits with dictionary tokens.

The sections below first describe EL-14, in which an ordinary token represents 14 bits. It is the target that uses the fewest English words among the profiles currently considered practical. The simpler EL-8 and EL-12 profiles are described under [Profiles](#profiles) and are the first implementation targets.

### 1. Normal tokens

The normal part of an EL-14 dictionary contains 16,384 carefully selected English words. This count is required because 14 bits have `2^14 = 16,384` possible values. Each dictionary position, called an **index**, corresponds to one of those values:

```
14-bit value → dictionary index → word
```

The encoder places the payload bytes end to end as one bit sequence and divides that sequence into 14-bit groups. It interprets each complete group as a number from 0 through 16,383 and emits the word at that dictionary index. The decoder performs the reverse lookup.

Why 14 bits? Published common-word lists and lists reduced to base word forms contain roughly 25,000 to 84,000 entries, while broader lexical databases contain more. A 15-bit normal dictionary alone would need 32,768 tokens. A 14-bit normal dictionary needs 16,384, leaving more candidates to discard because they are rare, offensive, difficult to spell, or too similar to another word. The derivation experiment still has to establish whether enough suitable words remain.

### 2. Controlled Ending on Byte Boundaries

The number of payload bits is not usually divisible by 14. EntropyLex therefore uses **trim tokens** for the 2 to 12 bits that may remain after the last complete 14-bit group. A trim token can appear only at the end of a phrase, and its dictionary position tells the decoder both how many bits it represents and what those bits are.

The encoder:

- Reads each byte from its highest-value bit to its lowest-value bit. This is commonly called most-significant-bit-first order.
- Emits 14-bit words until fewer than 14 bits remain.
- Maps the remaining 2–12 bits to the appropriate trim token. Only even remainder lengths can occur, as explained below. If no bits remain, the phrase ends with a normal word and no trim token is added.

This ensures that decoding reconstructs exactly the original byte-aligned payload without needing a separate length field.

### 3. Decoding

Decoding reverses the process:

- Convert each word back into its bit sequence.  
- For all but the last word, read 14 bits.  
- For the final word, read only the number of bits it indicates.  
- Reassemble the bits into bytes.

The transformation is reversible and deterministic: the same bytes and dictionary always produce the same phrase, and decoding that phrase reproduces the bytes exactly.

### 4. Trimming tokens

Without a trim token, the encoder would have to add meaningless padding bits to fill the final 14-bit group. The decoder would not know which final bits belonged to the payload. A separate payload-length field could solve that problem, but EntropyLex instead records the remaining bit count in the final token.

Payload lengths are multiples of 8 bits, while normal EL-14 tokens consume 14 bits. The greatest common divisor of 8 and 14 is 2: it is the largest integer that divides both numbers. Consequently, the possible nonzero remainder lengths are 2, 4, 6, 8, 10, and 12 bits. No odd remainder can occur.

A required length field would add overhead to every encoded payload. Trim tokens instead increase the dictionary size while adding a final token only when bits remain. That trade is part of the EntropyLex design.

Representing every possible value at each remainder length requires `2^12 + 2^10 + 2^8 + 2^6 + 2^4 + 2^2`, or `4,096 + 1,024 + 256 + 64 + 16 + 4 = 5,460`, trim words. Together with the 16,384 normal words, an EL-14 dictionary contains 21,844 tokens. Open sources contain more raw entries, but the derivation pipeline still has to show that 21,844 familiar and sufficiently distinct words remain after filtering. Less-frequent surviving words can be assigned to trim positions because a phrase contains at most one trim token.

### 5. Summary

EntropyLex provides an exact, word-based representation for arbitrary short byte sequences. Its dictionaries are intended to make phrases easier to remember, distinguish, and communicate than conventional base-encoded strings. Whether a particular phrase is practical to memorize or speak depends on its length, dictionary quality, and user.

## Profiles

The 14-bit design above is one member of a family. A **profile** fixes the normal-token width `w`: the number of payload bits represented by each normal token. Profiles use the name `EL-<w>`, such as EL-12. They otherwise share the same bit order, trim-token method, and encoding and decoding rules.

The table uses `g = gcd(8, w)`, the greatest common divisor of the eight-bit byte width and the token width. Every possible remainder is a multiple of `g`, so trim tokens are needed only for lengths `g, 2g, …, w-g`.

| Profile | Bits/word | `g` | Remainders | Normal words | Trim words | Total dictionary |
|---|---|---|---|---|---|---|
| EL-8  | 8  | 8 | none                    | 256    | 0    | 256    |
| EL-12 | 12 | 4 | 4, 8                    | 4,096  | 272  | 4,368  |
| EL-14 | 14 | 2 | 2, 4, 6, 8, 10, 12      | 16,384 | 5,460| 21,844 |
| EL-16 | 16 | 8 | 8                       | 65,536 | 256  | 65,792 |

Words needed for a payload:

| Payload | EL-8 | EL-12 | EL-14 | EL-16 |
|---|---|---|---|---|
| 16 bytes (128-bit key)     | 16 | 11 | 10 | 8  |
| 32 bytes (256-bit key)     | 32 | 22 | 19 | 16 |
| 60 bytes (the sample GIF)  | 60 | 40 | 35 | 30 |

### EL-8 — one byte per token

An eight-bit token represents exactly one byte. No byte crosses a token boundary, there are no remaining bits, and EL-8 needs no trim tokens. Encoding is a direct substitution: byte → word; decoding reverses it.

The dictionary contains 256 words. The small count permits strict selection for short, common words that differ clearly in spelling and pronunciation, and it makes complete human review practical. The Pretty Good Privacy (PGP) Word List is an earlier byte-to-word system and provides a useful comparison.

The tradeoff is phrase length: a 256-bit key needs 32 EL-8 words instead of 19 EL-14 words.

EL-8 is the first implementation target. It tests dictionary loading, splitting a phrase into tokens, applying written-form rules, and verifying encode-then-decode behavior without requiring groups of bits to cross byte boundaries.

### EL-12 — token boundaries occur at whole or half bytes

At 12 bits, `g = 4`, so a token boundary occurs only at a whole byte or after four bits of a byte. A four-bit half-byte is sometimes called a **nibble**. These boundaries are easy to compare with hexadecimal notation, where each digit also represents four bits.

Only two nonzero remainder lengths exist, 4 and 8 bits, so the trim portion contains 272 words instead of EL-14's 5,460. The whole EL-12 dictionary contains 4,368 words. Source lists appear large enough to support a very familiar vocabulary at this size, but the derivation pipeline must confirm the result after all filters are applied.

An EL-12 normal word represents 12 bits, about 86 percent of EL-14's 14 bits, while the complete dictionary is one-fifth the size.

EL-12 is the second implementation target. It adds every mechanism EL-14 needs: grouping bits across byte boundaries, calculating the remaining-bit count, using trim tokens, and rejecting results that do not reconstruct a whole number of bytes. Its intermediate values remain small enough to check by hand.

### EL-14 — the reference profile

EL-14 is the target that minimizes the number of English tokens among the profiles currently considered plausible. It is also the default profile when a discussion says “EntropyLex” without naming one. It needs six remainder-length groups, 5,460 trim words, and 21,844 words in total. Whether enough suitably familiar and distinct English words exist remains an experimental question.

### EL-16 and wider

EL-16 requires 65,792 tokens, more than a usable single-word English vocabulary can supply. Its bit handling is nevertheless simple: it has one eight-bit trim group, and normal tokens contain exactly two bytes. Languages or scripts with larger token inventories may offer more candidates, but familiarity and distinctness still limit what is usable.

### The decoder needs the profile and dictionary

A phrase normally does not state which profile or dictionary produced it. The surrounding application or message must provide that information separately, just as it must say whether a character string uses Base32 or Base58. Implementations must report a **dictionary fingerprint**, a SHA-256 identifier calculated from the ordered token mapping and the rules used to recognize written tokens. An application that already knows the expected fingerprint can reject the wrong dictionary instead of silently decoding the phrase to different bytes. The authoritative LXJ file and its optional LXB form share one fingerprint because they describe the same behavior.

One proposal would make the profile identifiable from the words: assign completely separate word sets to EL-8, EL-12, and EL-14, with no word shared between profiles. This remains undecided; see [the open question below](#open-question-whether-profiles-share-words).

## Other Languages and Character-Based Dictionaries

Nothing in the encoding rules is English-specific. Bit grouping, trim tokens, and index assignment remain the same; only the token inventory and the rules for recognizing written tokens change.

The largest practical `w` depends on how many suitable tokens a language or writing system can provide. A source database may contain many entries that are too rare, visually similar, difficult to pronounce, or otherwise unsuitable. EL-14 is therefore a target to test for English, not a proven capacity. A writing system with many familiar characters may support a different profile. EL-16, for example, needs a much larger dictionary than EL-14 but has simpler remainder handling because it needs only one trim group.

**Chinese with one character per token.** A Han character—the character family used in written Chinese—can serve as one token. Such a phrase is visually compact and may omit spaces. A finalized dictionary must define exactly what the decoder counts, such as Unicode code points, rather than relying on a programming language's informal idea of a character. The Table of General Standard Chinese Characters contains 8,105 characters. That is more than EL-12's required 4,368 tokens but less than EL-14's 21,844. EL-12 is therefore the largest defined profile that could use only characters from that standard, before filtering for familiarity and visual distinctness.

**Japanese has a smaller widely known kanji inventory.** *Kanji* are the Han-derived characters used in Japanese. The jōyō list—the government guide to characters for general use—contains 2,136. Even the 3,000 kanji allowed in personal names, consisting of the jōyō list plus 864 additional jinmeiyō name characters, falls short of EL-12's 4,368-token requirement. The often-cited figure of about 10,000 comes from Japanese Industrial Standard (JIS) X 0213, a computer character-encoding standard; it describes what software can represent, not what readers commonly know. A practical Japanese dictionary restricted to single kanji therefore appears limited to EL-8.

Chinese routinely writes most content words and many other elements with Han characters, while Japanese also uses hiragana and katakana—two scripts that represent sounds and are collectively called **kana**. Consequently, Japanese text can rely on a smaller set of commonly known kanji. Based on the reviewed raw counts, Chinese can be investigated at EL-12 while Japanese single-kanji tokens do not reach EL-12.

**Multi-character tokens for EL-16 and wider profiles.** Reaching EL-16's 65,792 tokens would require words or compounds containing more than one character. A standard Chinese dictionary, 现代汉语词典 (*Xiandai Hanyu Cidian*), contains about 69,000 entries. Using nearly all of it would leave almost no entries available to discard for rarity, visual similarity, or other problems. A practical EL-16 dictionary would need a much larger starting list. It would also need unambiguous token boundaries: fixed two-code-point tokens could be counted directly, while variable-length tokens would require separators. Current evidence supports investigating single-character Chinese at EL-12; it does not yet establish a practical EL-16 dictionary.

Full citations for all of these counts are in [SPEC.md section 12.6](SPEC.md).

**Selection criteria depend on the writing system.** English selection gives substantial weight to spoken confusion caused by homophones and accent differences. Mandarin and Japanese contain many characters or words with the same pronunciation, so their dictionaries should give greater weight to visual differences. Candidate measures include stroke count, shared character components such as radicals, and overall shape. A secondary measure can compare **romanizations**, which write pronunciations in the Latin alphabet: pinyin with tone for Mandarin and rōmaji for Japanese. Meaning still matters because neighboring characters can accidentally form a familiar compound that a person may remember as a unit rather than as two independent tokens.

## Where the Dictionaries Come From

Except for possible manual review of the 256-word EL-8 dictionary, the dictionaries will not be selected entirely by hand. Human selection of 21,844 words would be difficult to repeat consistently or inspect systematically. Instead, software will derive each dictionary from recorded source datasets and selection settings.

The objective is to choose the required number of familiar words while making the easiest-to-confuse pair as different as possible in spelling, pronunciation, and meaning. “Familiarity floor” below means the lowest acceptable usage-frequency score.

The intended pipeline—an ordered sequence in which each tool passes recorded output to the next—is:

1. **Load and combine sources.** Load a starting word list, usage-frequency data derived from bodies of text, a pronunciation dictionary, data that connects word forms such as “walked” to the base form “walk,” and a numerical model of similarity in meaning. This stage converts the different source formats into one candidate table. “Common” should be supported by measured usage, not personal judgment.
2. **Apply required exclusions.** Remove candidates that violate firm rules before comparing word pairs. Proposed rules exclude words outside 3–8 characters in the basic ASCII Latin character set, proper names, redundant grammatical forms, offensive or specialized terms, homophones (different words pronounced the same), words with multiple unrelated pronunciations, unpredictable spelling, and words below the familiarity floor. Record every removal and its reason so reviewers can inspect the result.
3. **Measure similarity.** Compare each surviving pair in three ways:
   - *Spelling* — count insertions, deletions, replacements, and adjacent-letter swaps; give extra weight to mistakes involving nearby keyboard keys; and detect long shared beginnings or endings.
   - *Sound* — compare sequences of speech sounds, treating sounds formed similarly in the mouth as more alike; also compare syllable count and which syllables are stressed.
   - *Meaning* — use a numerical language model to identify related words, such as “big” and “large.” Keeping such pairs apart reduces accidental substitution from memory and reduces the chance that random tokens resemble a meaningful sentence.
4. **Choose the words.** Treat two words as conflicting if any comparison says they are too similar. Consider candidates from most to least familiar and keep a word only when it does not conflict with one already chosen. If this does not produce the required count, change the similarity limits by recorded amounts and repeat. The quality report must show every such change.
5. **Assign roles and indices.** Give the most familiar selected words the normal-token positions because normal tokens appear throughout a phrase. Give the remaining words the trim-token positions, which can occur only once at the end. Sort each group by written form before assigning numerical indices so the result does not depend on source ordering.
6. **Write, generate the binary form, and verify.** Write the authoritative, human-readable JavaScript Object Notation (JSON) dictionary in LXJ (`.lxj`) format and optionally produce the faster-loading LXB (`.lxb`) binary form. Confirm that both files describe the same index-to-token mapping and recognition rules. A file can be checked by itself for structure, counts, valid tokens, and its fingerprint. Claims about familiarity, pronunciation, meaning, and repeatable derivation require the recorded source datasets and settings; those results belong in a separate quality report.

Every stage must be deterministic: the same exact input files and settings must produce the same ordered mapping and fingerprint. Once the file layouts are final, they must also produce exactly the same LXJ and LXB bytes. Step 3 can use measures appropriate to each writing system—for example, visual-shape and romanized-pronunciation comparisons for Chinese—without changing the other stages.

These tools will live under [`tools/`](tools/), organized first by programming language and then by author name, as `src/` is. Independently written pipelines provide a useful check: when given byte-for-byte identical input files and settings, they must produce the same mapping and fingerprint. Once the LXJ and LXB layouts are final, their output files must also match byte for byte. A disagreement reveals either an implementation error or a rule that the specification has not defined precisely enough.

The working file-format design is in [`data/dict/FORMAT.md`](data/dict/FORMAT.md). Reviewed source candidates, license conditions, and published or measured counts are in [`data/dict/SOURCES.md`](data/dict/SOURCES.md).

### Open question: whether profiles share words

**This is not decided yet.** When EL-8, EL-12, and EL-14 dictionaries are selected from the same starting vocabulary, how should their word sets relate? Three alternatives remain under consideration:

**A — nested sets.** Every EL-8 word also appears in EL-12, and every EL-12 word also appears in EL-14 (`EL-8 ⊂ EL-12 ⊂ EL-14`). This needs only 21,844 distinct words and is simple to document. However, an EL-8 phrase also consists entirely of valid EL-14 words, so choosing the wrong profile can silently produce different bytes. The EL-8 words would also be inherited from the larger selection rather than optimized as the clearest possible set of 256.

**B — independent selection.** Select each profile separately for its required count. A 256-word dictionary can enforce much larger differences between words than a 21,844-word dictionary can. However, words may overlap unpredictably, so some phrases would reveal their profile and others would not.

**C — separate sets.** Allow no word to appear in more than one profile. Any nonempty phrase would then identify the profile that produced it because any one of its words would do so. The empty phrase contains no word and would remain valid under every profile. A decoder using another profile would reject a nonempty phrase instead of returning incorrect bytes. This would also give EL-8 some ability to detect a profile mismatch, although it would still not detect one valid EL-8 word being replaced by another.

The cost of Scheme C is the number of distinct suitable words required. `2^14 + 2^12 + 2^8 = 20,736` counts only normal tokens. Trim tokens must also be separate from every normal and trim set, adding 5,732:

| | EL-8 | EL-12 | EL-14 | Total |
|---|---|---|---|---|
| Normal words | 256 | 4,096 | 16,384 | 20,736 |
| Trim words   | 0   | 272   | 5,460  | 5,732  |
| **Per profile** | **256** | **4,368** | **21,844** | **26,468** |

The total of 26,468 exceeds the reviewed 25,000-entry 12dicts core list, although broader reusable sources contain many more base word forms. Only the 20,736 normal words require the highest familiarity; the 5,732 trim words may use less-frequent survivors because a trim word appears at most once per phrase. These counts make Scheme C worth testing, but do not prove that enough acceptable words exist. The derivation pipeline must report how many candidates survive and whether similarity limits had to be weakened.

A partial alternative would keep normal-word sets separate but share trim words. The shared trim set would still have to avoid all three normal sets, so it would save only 272 words. It would also fail to identify the profile for a one-byte payload, which is represented by one trim word in both EL-12 and EL-14. For those reasons, this alternative is not recommended.

Giving EL-16 its own completely separate word set as well would require 92,260 English words in total and is not considered practical.

The decision must be made before `eldict-select` is implemented. It determines whether selection occurs once followed by assignment to profiles, or occurs separately for each profile while excluding words already used elsewhere.

## Implementation Order

The order below introduces complexity gradually:

1. **EL-8** — load a dictionary, separate written phrases into tokens, apply written-form rules, verify encode-then-decode behavior, and run shared examples. No token crosses a byte boundary and no trim token is needed.
2. **EL-12** — group bits across byte boundaries, calculate remaining bits, use trim tokens, and verify that decoding produces a whole number of bytes. This introduces every EL-14 mechanism with a dictionary one-fifth the size.
3. **EL-14** — the reference profile. No new mechanism, only scale.
4. **EL-16 and non-English dictionaries** — add token-recognition rules needed by other writing systems, such as splitting a separator-free phrase into individual characters.

Implementations should state which profiles they support. Code that supports multiple profiles should take `w` and the associated profile values as data rather than duplicate the encoding algorithm for each profile.

See [SPEC.md](SPEC.md) for the rules implementations must follow, including bit order, index assignment, decoder rejection conditions, and dictionary derivation requirements.

