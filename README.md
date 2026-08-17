# EntropyLex: A Human-Readable Encoding for High-Entropy Binary Data

EntropyLex is a bidirectional encoding system that converts raw binary data into short sequences of common English words. It exists to make high-entropy information (cryptographic keys, device identifiers, authentication tokens, reference hashes, and other compact binary payloads) memorable, speakable, and recoverable by humans, without sacrificing security.

Traditional base encodings (Base32, Base58, hex) optimize for machines, specific transport encodings (ASCII) and purposes (ASCII armoring while maximizing bandwidth by limiting encoding-expansion). They are hard to memorize, error-prone when spoken aloud, and visually ambiguous. 

EntropyLex reverses the priority: it maps binary data into carefully selected English words chosen for clarity, phonetic distance, and recognizability across accents. This makes the resulting phrases suitable for memorization, transcription, voice communication, and low-tech backup scenarios where humans, not machines, are the bottleneck.

The encoding itself is language-neutral — English is the first dictionary, not a requirement of the format — and it is defined as a family of profiles that trade density against complexity. See [Profiles](#profiles).

## Why It Exists

EntropyLex solves three persistent problems:

1. Humans cannot reliably handle raw entropy. Hex strings, random characters, and dense base encodings are difficult to memorize or communicate accurately. As a rough working figure, a 10-digit string of numbers like a NANP phone number (555-867-5309) is about all the entropy a human can be trusted to memorize. *(Uncited. This is a rule of thumb, not a measured result, and the familiar 7±2 finding concerns working memory rather than long-term retention. The design does not depend on the exact number — only on the uncontroversial premise that the usable figure is small.)*

2. Mnemonic systems lack explicit, exact binary mapping. Existing “seed phrase” approaches are designed around specific wallet protocols and often rely on indirect hashing pipelines rather than strict reversible encodings.

3. There is no general mechanism for turning any short binary payload into a compact, human-manageable phrase. EntropyLex works for *any* 8-bit–symbol binary payload, independent of application or cryptographic scheme.

Typical use cases include: memorized master keys, offline recovery secrets, secure pairing codes, fixed-entropy identifiers, reference hashes, and low-bandwidth cross-channel transfers where a human must read, write, or speak the data. You could even try to turn your favorite cat GIF into a speakable ballad. The smallest valid 1 pixel GIF is around 34 bytes, which would be about 20 words; the example 8x8 4-color cat GIF in `data/sample/input/` is 60 bytes, or 35 words. *(The 34-byte figure is uncited — commonly quoted minimums vary between roughly 26 and 35 bytes depending on how strictly "valid" is read. It is an illustration, not a specification input. The 60-byte figure is measured from the file itself.)*

## How It Works - Conceptual Overview

EntropyLex treats an input as a stream of 8-bit symbols, meaning the payload is always an integer number of bytes. It then transforms that stream into a sequence of word-tokens, each representing a fixed amount of entropy.

The sections below describe the 14-bit design, which is the density-optimized target. EntropyLex is actually a *family* of encodings that differ only in how many bits each word carries — see [Profiles](#profiles) below for the 8-bit and 12-bit variants, which are simpler and are the first implementation targets.

### 1. Fixed-Entropy Words

The core dictionary contains 16,384 carefully selected English words.  
Each word corresponds to a unique 14-bit value:

```
index(word) → 14-bit chunk
```

Encoding proceeds by concatenating the payload bits and slicing them into 14-bit groups. Each group is mapped directly to a word, producing most of the phrase.

Why 14-bits? There are about 20,000-25,000 root words (the base word without pluralization or conjugation) in common English dictionaries. 15 bits would need about 32,000 (2^15^) dictionary tokens. 14 bits (2^14^=16384) fits well within a 20k token dictionary, with some room to spare.

### 2. Controlled Ending on Byte Boundaries

Because arbitrary binary streams rarely end on a 14-bit boundary, EntropyLex uses special trimming tokens. These words represent *shorter* bit sequences (2 to 12 bits) and appear only as the final token.

The encoder:

- Packs the payload into a bitstream, most significant bit first.
- Emits 14-bit words until fewer than 14 bits remain.
- Maps the remaining 2–12 bits into an appropriate special token that represents a shorter bit sequence. Special short-tokens exist for all "remnant" lengths that can actually occur (see below — only even lengths are reachable). If nothing remains, no trim token is emitted and the phrase simply ends on a normal word.

This ensures that decoding reconstructs exactly the original byte-aligned payload without needing a separate length field.

### 3. Decoding

Decoding reverses the process:

- Convert each word back into its bit sequence.  
- For all but the last word, read 14 bits.  
- For the final word, read only the number of bits it indicates.  
- Reassemble the bits into bytes.

The transformation is fully reversible, deterministic, and carries no lossy semantics or application-specific assumptions.

### 4. Trimming tokens

Because we use 14-bit indices to pack as much entropy per word as possible, and the input is in 8-bit symbols, it's possible that the input payload could end in a way that less than 14 bits still need to be encoded. It's not possible to just emit 14 bits and hope the decoder can either guess where the end of payload should be, or that the payload consumer can tolerate garbage padding at the end of the payload. Because the greatest common divisor of the payload symbol size (8 bits) and the token size (14 bits) is 2, every possible leftover is a multiple of 2 bits, so only even-bit-length mismatches are possible. Nonetheless, it is still required that EntropyLex be able to safely transmit a 2, 4, 6, 8, 10, or 12 bit long trimmed token that the decoder will understand how to parse and reconstruct.

An alternative strategy would be to mandate a length header (even just 8 bits) on every message, but since this unilaterally increases every message by 8 bits, this option was discarded. The trimming token method emits no more words than would be necessary for encoding an untrimmed payload with trailing garbage, however it requires the utilization of additional dictionary words to express the shorter trimmed tokens.

The trimming token dictionary needs a set of tokens to represent all sequences of 12, 10, 8, 6, 4, or 2 bits. Therefore the additional dictionary size is 2^12^+2^10^+2^8^+2^6^+2^4^+2^2^ or 4096+1024+256+64+16+4 = 5,460 extra dictionary words. On top of the 2^14^ base 14-bit tokens (16,384), this makes for a total dictionary size of 21,844. This fits well within dictionaries of common English words that are up to 25000 words. Because they are less frequently used, we will select these words from the least-used words within the dictionary we choose.

### 5. Summary

EntropyLex provides a general, high-density, human-friendly representation for arbitrary short binary payloads, small enough to memorize, robust enough to speak aloud, and precise enough to encode cryptographic secrets or other critical binary data without structural ambiguity.

## Profiles

The 14-bit design above is one member of a family. Everything except the token width `w` is identical across the family: the same bitstream model, the same trimming concept, the same encode and decode algorithms. A profile is named `EL-<w>`.

The width choice drives everything else. With 8-bit input symbols, let `g = gcd(8, w)`. Every possible leftover at the end of the payload is a multiple of `g`, so the trim dictionary only needs to cover remainders `g, 2g, … w-g`:

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
| 34 bytes (minimal GIF)     | 34 | 23 | 20 | 17 |
| 60 bytes (the cat GIF)     | 60 | 40 | 35 | 30 |

### EL-8 — no byte splitting at all

When the word is 8 bits wide, it is exactly one byte. Nothing is ever split, the leftover is always zero, and **there is no trim dictionary and no trim logic whatsoever**. Encoding is a straight substitution: byte → word, word → byte.

The dictionary is 256 words. At that size, word selection can be ruthless — every word can be short, common, unambiguous, and phonetically far from every other word in the set, and the whole dictionary can be reviewed by a human in an afternoon. Nothing else in the family can say that. (The PGP Word List is prior art for roughly this shape, and is a useful sanity reference.)

The price is density: 8 bits per word means a 256-bit key is 32 words instead of EL-14's 19.

EL-8 is the first implementation target. It exercises dictionary loading, tokenizing, normalization, and the round-trip test harness with none of the bit-packing complexity.

### EL-12 — splits bytes only in half

At 12 bits, `g = 4`, so the bitstream only ever breaks on **nibble boundaries**. Every word boundary lands on a whole byte or a half byte, which is easy to reason about, easy to debug, and easy to eyeball against a hex dump.

Only two leftover sizes exist, 4 and 8 bits, so the trim dictionary is 272 words instead of 5,460 — small enough to print. The whole dictionary is 4,368 words, which can be drawn entirely from high-frequency core vocabulary, so every word is one a person actually knows.

The density cost is modest: 12 bits per word is 86% of EL-14, for a fifth of the dictionary.

EL-12 is the second implementation target. It exercises the complete mechanism — bitstream packer, remainder arithmetic, trim dictionary, alignment check — that EL-14 needs, while keeping every intermediate value checkable by hand.

### EL-14 — the reference profile

The density-optimized target described in the sections above, and the default meaning of "EntropyLex" when no profile is stated. It needs the full apparatus: six trim classes, 5,460 trim words, and a 21,844-word dictionary that pushes to the edge of usable common English.

### EL-16 and wider

65,792 words is well past what single English words can supply, so EL-16 is not an English profile. Structurally, though, it is nearly the simplest one in the family: one trim class, and every boundary on a whole or half byte. Scripts with far larger token inventories bring it closer to reach — though as the section below works through, "closer" is not the same as viable.

### A phrase is not self-describing — unless we make it so

A phrase does not carry its own profile or dictionary identity. Decoding correctly requires knowing both, carried out of band, exactly the way "this is Base58, not Base32" is carried out of band. Implementations should be able to report a dictionary fingerprint so a mismatch fails loudly instead of quietly producing different, well-formed bytes.

There is, however, a way to make the profile self-evident from the words alone — by carving the three dictionaries out of the master corpus as **mutually disjoint** sets. See [the open question below](#open-question-how-the-three-dictionaries-relate).

## Other Languages, and Ideographic Dictionaries

Nothing in the encoding is English-specific. The bitstream model, trimming, and index space are language-neutral; only the dictionary changes.

What *is* English-specific is the ceiling on `w`. English gives us roughly 20,000–25,000 usable common root words, which is what caps the reference profile at 14 bits. Scripts with large character inventories lift that cap — and because EL-16 has only one trim class, going wider actually makes the *encoder* simpler, not harder.

**Chinese, one character per word.** A single Han character can be a token: one glyph carries the whole chunk, phrases are visually compact, and no separators are needed because segmentation is by glyph. The limit is the character inventory — the Table of General Standard Chinese Characters defines 8,105 characters. That is comfortable for EL-12's 4,368 tokens and impossible for EL-14's 21,844. **EL-12 is therefore the natural profile for a single-character Chinese dictionary**, the same way EL-14 is natural for single-word English — a pleasant convergence, since EL-12 is also the simplest profile that has trim logic at all.

**Japanese reaches one profile less far.** The jōyō kanji list is 2,136 characters, which supports EL-8 with room to spare and falls well short of EL-12's 4,368 — as does the full 3,000-character set of kanji legally permitted in Japanese personal names (jōyō plus 864 jinmeiyō). Beware the figure "about 10,000 kanji": that is the **JIS X 0213 encoding repertoire**, what a Japanese computer can display, not what a reader knows. Sourcing a dictionary from it would be like sourcing an English one from the full OED headword list. Measured against real literacy, Japanese single-glyph dictionaries top out at EL-8.

The gap versus Chinese is structural: Chinese writes nearly everything in hanzi, so the everyday inventory runs to ~3,500 characters, while Japanese offloads grammar and much vocabulary onto kana, leaving kanji past jōyō genuinely rare. That ~3x difference is exactly one profile step.

**Compound tokens for 16 bits and beyond.** Reaching 65,792 tokens means multi-character words rather than single glyphs, and the margin is thinner than it looks. A standard Chinese desk dictionary — 现代汉语词典, 7th ed. — holds about 69,000 entries, so EL-16 would consume 95% of it with essentially nothing discarded for familiarity or distinctness. Usable EL-16 needs a comprehensive lexicon several times that size, paying the same tail-familiarity cost English pays at EL-14. Segmentation also stops being free; fixed-width two-character tokens are the likely construction, keeping segmentation a matter of counting and phrases visually regular. Realistically, **EL-12 single-glyph Chinese is the non-English profile the evidence supports**; EL-16 is aspirational everywhere.

Full citations for all of these counts are in [SPEC.md section 12.6](SPEC.md).

**The selection criteria change with the script.** The English criteria are weighted toward speech, because English fails at homophones and accents. Mandarin and Japanese have much higher homophone density, so a CJK phrase is markedly less reliable spoken than written. Those dictionaries should optimize for **visual** distinctness first — stroke count, shared radicals, overall glyph shape are the ideographic equivalent of a near-homophone — with romanization distance (pinyin plus tone, romaji) as a secondary metric, since it governs how the phrase gets typed. Semantic distance matters more, not less: character sequences readily read as meaningful compounds, so adjacent-token pairs that form a common word should be screened out.

## Where the Dictionaries Come From

The dictionaries are not going to be written by hand — 21,844 words chosen by human judgment would be neither auditable nor reproducible. They will be **derived** from a master word list by a set of preprocessing tools that do not exist yet.

The objective, stated plainly: **pick the subset of a master vocabulary that maximizes the minimum pairwise distance in spelling, sound, and meaning, subject to a hard word count and a familiarity floor.**

The intended pipeline:

1. **Ingest.** Load a master word list together with corpus frequency data, a pronunciation lexicon, lemma/inflection data, and a semantic embedding model. "Common" should be measured, not asserted.
2. **Filter.** Apply hard disqualifiers before any distance math: length outside 3–8 characters, non-ASCII, proper nouns, inflected forms whose root is already a candidate, offensive or sensitive terms, domain jargon, homophone sets, homographs with divergent pronunciations ("read", "lead", "wound"), unpredictable spellings, and anything below the familiarity floor. Every rejection gets logged with its reason, as a reviewable artifact.
3. **Score.** Compute three independent distance families:
   - *Spelling* — Damerau-Levenshtein edit distance, a keyboard-adjacency-weighted variant, and shared prefix/suffix length so the set isn't full of words differing only in an ending.
   - *Sound* — phoneme-string edit distance weighted by articulatory features (so /b/–/p/ counts as nearer than /b/–/s/), with coarse phonetic-key collisions rejected outright, plus syllable count and stress pattern to spread the selection across prosodic shapes.
   - *Meaning* — embedding cosine distance, to keep near-synonyms like "big"/"large" from both being selected (a human will substitute one for the other from memory and never notice) and to reduce the odds that a random phrase reads as a sentence.
4. **Select.** Build a confusability graph — an edge between any two words too close in *any* of the three families — then seed with the highest-frequency candidates and greedily add words that conflict with nothing already chosen. Familiarity is the one property that cannot be recovered later, so it drives the ordering. If the target count isn't reached, thresholds relax by a recorded step and the run repeats; the relaxation trace ships with the dictionary. A dictionary that only hit its count by weakening the phonetic threshold twice should say so on its face.
5. **Partition and index.** The most frequent words become the normal dictionary; the remainder become trim words, since trim words are used far less often. Each group is then sorted lexicographically to assign indices, so the artifact stays stable and diffable.
6. **Emit and verify.** Write a UTF-8 artifact with a metadata header (profile, language, counts, pinned sources, config hash, SHA-256 fingerprint) and one word per line. A verifier then re-derives everything from the artifact alone and asserts the counts, the disjointness of normal and trim sets, the absence of duplicates, the canonical ordering, the surviving filters, the minimum distance in each family with worst offending pairs named, and a round-trip of the published test vectors.

Every stage is deterministic given its inputs and a recorded config, so re-running the pipeline produces a byte-identical dictionary. The stages are metric-pluggable by design: a Chinese dictionary swaps out step 3's distance families for visual and romanization metrics and leaves the rest of the pipeline alone.

These tools live under [`tools/`](tools/), organized by language and then by implementor name, the same convention `src/` uses. Independent implementations are welcome and useful: two pipelines written from the spec, run on the same pinned inputs, must produce byte-identical artifacts. If they don't, either the spec is underspecified or one of them is wrong.

### Open question: how the three dictionaries relate

**This is not decided yet.** When EL-8, EL-12, and EL-14 dictionaries all come out of the same master corpus, how should their word sets relate? Three schemes are live, and any of them could win:

**A — nested subset.** EL-8 ⊂ EL-12 ⊂ EL-14. Cheapest at 21,844 words total, simplest to document and verify. But every EL-8 phrase is also a well-formed run of EL-14 words, so decoding under the wrong profile silently yields different, plausible bytes. And the EL-8 set is merely the top 256 of a 21,844-word selection, not the 256 most mutually distinct words the corpus could offer.

**B — independent optimization.** Each profile selected separately against its own count. Best per-profile quality by a wide margin — at n=256 you can demand minimum distances that are simply unreachable at n=21,844. But the overlap between sets is uncontrolled, so a phrase is *sometimes* profile-identifying with no way to know when.

**C — disjoint partition.** Carve the three dictionaries as mutually exclusive sets, no word in more than one. This buys a real property: **any single word in a phrase tells you which profile produced it.** Profile becomes self-identifying, wrong-profile decoding becomes a hard failure instead of silent garbage, and EL-8 — which otherwise has *no* structural error detection whatsoever — gains at least that much.

The cost of C is corpus size, and it's worth being precise about it. The natural figure to quote is 2^14 + 2^12 + 2^8 = **20,736**, but that counts only the normal dictionaries. The trim words have to be disjoint too — from each other and from every normal set — which adds 5,732:

| | EL-8 | EL-12 | EL-14 | Total |
|---|---|---|---|---|
| Normal words | 256 | 4,096 | 16,384 | 20,736 |
| Trim words   | 0   | 272   | 5,460  | 5,732  |
| **Per profile** | **256** | **4,368** | **21,844** | **26,468** |

26,468 is past the 20,000–25,000 common-root-word estimate that drove the choice of 14 bits in the first place. The saving grace: only the 20,736 normal words need to be genuinely familiar. The 5,732 trim words come from the low-frequency tail by design, since trim words appear at most once per phrase. So the demand on *good* vocabulary does fit — but C consumes essentially all the headroom.

(A tempting middle road — disjoint normal dictionaries with one shared trim dictionary — turns out not to be worth it. The shared trim set still has to avoid all three normal sets, so it saves just 272 words, and it reintroduces ambiguity exactly where it's hardest to spot: a one-byte payload is a single trim word under both EL-12 and EL-14, carrying no identifying information at all. Full disjointness is the better buy.)

Extending disjointness to EL-16 is out of reach for English at 92,260 words.

The decision needs to be made before `eldict-select` gets written, since it determines whether the selector runs once with a partitioning step or three times with exclusion sets.

## Implementation Order

Each step adds exactly one mechanism:

1. **EL-8** — dictionary loading, tokenization, normalization, round-trip harness, test vectors. No bitstream, no trimming.
2. **EL-12** — bitstream packer, remainder arithmetic, trim dictionary, alignment validation. Every EL-14 mechanism at a fifth the scale, with hand-checkable intermediates.
3. **EL-14** — the reference profile. No new mechanism, only scale.
4. **EL-16 and non-English dictionaries** — no new mechanism beyond glyph-based segmentation.

Implementations should state which profiles they support, and should share one `w`-parameterized code path rather than writing each profile separately.

See [SPEC.md](SPEC.md) for the normative definition, including bit ordering, the unified index space, decoder validation rules, and the full dictionary derivation requirements.

