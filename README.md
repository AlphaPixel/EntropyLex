# EntropyLex: A Human-Readable Encoding for High-Entropy Binary Data

EntropyLex is a bidirectional encoding system that converts raw binary data into short sequences of common English words. It exists to make high-entropy information (cryptographic keys, device identifiers, authentication tokens, reference hashes, and other compact binary payloads) memorable, speakable, and recoverable by humans, without sacrificing security.

Traditional base encodings (Base32, Base58, hex) optimize for machines, specific transport encodings (ASCII) and purposes (ASCII armoring while maximizing bandwidth by limiting encoding-expansion). They are hard to memorize, error-prone when spoken aloud, and visually ambiguous. 

EntropyLex reverses the priority: it maps binary data into carefully selected English words chosen for clarity, phonetic distance, and recognizability across accents. This makes the resulting phrases suitable for memorization, transcription, voice communication, and low-tech backup scenarios where humans, not machines, are the bottleneck.

## Why It Exists

EntropyLex solves three persistent problems:

1. Humans cannot reliably handle raw entropy. Hex strings, random characters, and dense base encodings are difficult to memorize or communicate accurately. Research has shown that a 10-digit string of numbers like a NANP phone number (555-867-5309) is about all the entropy a human can be trusted to memorize.

2. Mnemonic systems lack explicit, exact binary mapping. Existing “seed phrase” approaches are designed around specific wallet protocols and often rely on indirect hashing pipelines rather than strict reversible encodings.

3. There is no general mechanism for turning any short binary payload into a compact, human-manageable phrase. EntropyLex works for *any* 8-bit–symbol binary payload, independent of application or cryptographic scheme.

Typical use cases include: memorized master keys, offline recovery secrets, secure pairing codes, fixed-entropy identifiers, reference hashes, and low-bandwidth cross-channel transfers where a human must read, write, or speak the data. You could even try to turn your favorite cat GIF into a speakable ballad. The smallest valid 1 pixel GIF is 34-bytes, which would be around 20 words, an example 8x8 4-color cat GIF is 60 bytes, or about 35 words.

## How It Works - Conceptual Overview

EntropyLex treats an input as a stream of 8-bit symbols, meaning the payload is always an integer number of bytes. It then transforms that stream into a sequence of word-tokens, each representing a fixed amount of entropy.

### 1. Fixed-Entropy Words

The core dictionary contains 16,384 carefully selected English words.  
Each word corresponds to a unique 14-bit value:

```
index(word) → 14-bit chunk
```

Encoding proceeds by concatenating the payload bits and slicing them into 14-bit groups. Each group is mapped directly to a word, producing most of the phrase.

Why 14-bits? There are about 20,000-25,000 root words (the base word without pluralization or conjugation) in common English dictionaries. 15 bits would need about 32,000 (2^15^) dictionary tokens. 14 bits (2^14^=16384) fits well within a 20k token dictionary, with some room to spare.

### 2. Controlled Ending on Byte Boundaries

Because arbitrary binary streams rarely end on a 14-bit boundary, EntropyLex uses special trimming tokens. These words represent *shorter* bit sequences (from 0 to 12 bits) and appear only as the final token.

The encoder:

- Packs the payload into a bitstream.
- Emits 14-bit words until fewer than 14 bits remain.
- Maps the remaining 0–13 bits into an appropriate special token that represents a shorter bit sequence. Special short-tokens exist for all "remnant" lengths (whatever length we might need to encode that is shorter than 14 bits). 

This ensures that decoding reconstructs exactly the original byte-aligned payload without needing a separate length field.

### 3. Decoding

Decoding reverses the process:

- Convert each word back into its bit sequence.  
- For all but the last word, read 14 bits.  
- For the final word, read only the number of bits it indicates.  
- Reassemble the bits into bytes.

The transformation is fully reversible, deterministic, and carries no lossy semantics or application-specific assumptions.

### 4. Trimming tokens

Because we use 14-bit indices to pack as much entropy er word as possible, and the input is in 8-bit symbols, it's possible that the input payload could end in a way that less than 14 bits still need to be encoded. It's not possible to just emit 14 bits and hope the decoder can either guess where the end of payload should be, or that the payload consumer can tolerate garbage padding at the end of the payload. Due to the payload symbols size (8-bit) and the message token size (14-bit) both being even, only even-bit-length mismatches are possible, nonetheless, it is still required that EntropyLex be able to safely transmit a 2, 4, 6, 8, 10, or 12 bit long trimmed token that the decoder will understand how to parse and reconstruct.

An alternative strategy would be to mandate a length header (even just 8 bits) on every message, but since this unilaterally increases every message by 8 bits, this option was discarded. The trimming token method emits no more words than would be necessary for encoding an untrimmed payload with trailing garbage, however it requires the utilization of additional dictionary words to express the shorter trimmed tokens.

The trimming token dictionary needs a set of tokens to represent all sequences of 12, 10, 8, 6, 4, or 2 bits. Therefore the additonal dictionary size is 2^12^+2^10^+2^8^+2^6^+2^4^+2^2^ or 4096+1024+256+64+16+4 = 5,460 extra dictionary words. On top of the 2^16^ base 14-bit tokens (16,384), this makes for a total dictionary size of 21,844. This fits well within dictionaries of common English words that are up to 25000 words. Because they are less frequently used, we will select these words from the least-used words within the dictionary we choose.

### 5. Summary

EntropyLex provides a general, high-density, human-friendly representation for arbitrary short binary payloads, small enough to memorize, robust enough to speak aloud, and precise enough to encode cryptographic secrets or other critical binary data without structural ambiguity.

