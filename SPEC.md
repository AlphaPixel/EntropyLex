# EntropyLex overview and specifications document

This document defines the current design of EntropyLex, a reversible encoding that maps byte-aligned binary payloads into sequences of English words and back. The system is designed for human usability, high entropy density, deterministic reversibility, and compatibility with arbitrary 8 bit input data.

This specification covers all known properties, constraints, terminology, encoding rules, decoding rules, token classes, dictionary sizing, and trimming behavior as currently established.

---

## 1. Purpose

EntropyLex converts binary data into word sequences that are easy for humans to read, write, speak, and memorize. It addresses use cases where high entropy data such as keys, seeds, identifiers, or checksums must be represented in a human friendly format without losing exact bit precision.

EntropyLex is intended to:

- Represent arbitrary 8 bit data without ambiguity
- Provide dense entropy per word (approximately 14 bits per token)
- Allow decoding without an explicit length header
- Guarantee byte aligned reconstruction
- Support deterministic round-trip mapping

---

## 2. Core Concepts

### 2.1 Input data
EntropyLex accepts input as a sequence of 8 bit bytes. All payloads are therefore multiples of 8 bits.

Let the input size be N bytes. Total payload bits = 8N.

### 2.2 Output tokens
EntropyLex produces a sequence of English word tokens. There are two classes:

1. Normal tokens
2. Trim tokens (used only for the final token)

### 2.3 Bit width of normal tokens
Each normal token encodes exactly 14 bits. This requires 2^14 = 16384 distinct normal dictionary entries.

### 2.4 Byte alignment constraint
Because the input is 8 bit aligned and the output tokens are 14 bit aligned, dividing an 8N bitstream into 14 bit chunks yields only certain possible remainders.

Let r = (8N mod 14). Since gcd(8, 14) = 2, only even remainders occur. Valid remainders are:

r in {0, 2, 4, 6, 8, 10, 12}

This determines the required trimming behavior for the final token.

---

## 3. Dictionary Structure

Total dictionary = normal dictionary + trim dictionary.

### 3.1 Normal dictionary
- Size: 16384 tokens
- Index range: 0 to 16383 (14 bit values)
- Each normal token represents a 14 bit binary value
- Only used for all tokens except the final token

### 3.2 Trim dictionary
A trim token is used only as the final token to encode the final remainder bits. Since r can only be {2, 4, 6, 8, 10, 12}, the trim dictionary must support these tail lengths.

For each remainder r, the final token must uniquely encode every possible r bit value. Therefore:

- Number of trim tokens needed for remainder r is 2^r

Total trim token count:

2^2 + 2^4 + 2^6 + 2^8 + 2^10 + 2^12
= 4 + 16 + 64 + 256 + 1024 + 4096
= 5460 tokens

Thus the complete dictionary requires:

16384 normal tokens  
+ 5460 trim tokens  
= 21844 total tokens

Trim tokens must be mapped so that each one unambiguously identifies:

1. The remainder size r  
2. The r bit tail payload  

Trim tokens never appear except as the final token.

---

## 4. Encoding Process

### 4.1 Convert input bytes to bitstream
Concatenate all input bytes into a linear bit sequence of length 8N bits.

### 4.2 Emit full 14 bit tokens
While at least 14 bits remain:

- Extract the next 14 bits
- Interpret them as a value between 0 and 16383
- Map the value to a normal token and append it to the output sequence

### 4.3 Determine remainder
After consuming all full 14 bit groups, the remainder bit length r is:

r = (8N mod 14)

Possible r values are 0, 2, 4, 6, 8, 10, 12.

### 4.4 Emit final trim token
If r = 0:

- No trim token is required
- Encoding ends on the last normal token

If r > 0:

- Read the final r bits from the remaining bitstream
- Use the pair (r, final_value) to select the appropriate trim token
- Append the trim token as the final word

Encoding ends after emitting the final trim token.

---

## 5. Decoding Process

Given a sequence of tokens:

### 5.1 Convert all but the last token to bits
For each token except the final one:

- Look up its 14 bit value in the normal dictionary
- Append the 14 bits to the reconstructed bitstream

### 5.2 Decode the final token
If the final token is a normal token:

- It contributes 14 bits
- Total bit length must be divisible by 8 for valid decoding

If the final token is a trim token:

- Identify its remainder size r
- Extract the r bits encoded by that token
- Append these r bits to the reconstructed bitstream

### 5.3 Verify byte alignment
EntropyLex guarantees that the reconstructed bit length must be divisible by 8. If not, the sequence is invalid or corrupted.

### 5.4 Convert bitstream back to bytes
Group bits into bytes and output the original binary data.

---

## 6. Trimming Behavior Summary

The input is 8 bit aligned. The encoder emits 14 bit tokens. Because of modular arithmetic properties:

- Remainder sizes are always even
- Maximum remainder is 12 bits
- 13 bit partial tokens never occur
- Only the six remainders 2, 4, 6, 8, 10, 12 require trim tokens

Total trim tokens required: 5460.

These tokens must be distinct from the 16384 normal tokens.

---

## 7. Dictionary Size Summary

- Normal tokens: 16384  
- Trim tokens: 5460  
- Total: 21844 tokens

If desired, trim tokens may be placed in separate ranges for clarity, but the decoder must uniquely distinguish normal versus trim tokens.

---

## 8. Word Selection Criteria

Word choice is not fully specified yet, but principles include:

- Distinct pronunciation to reduce confusion in speech
- Minimal homophones or near-homophones
- Easy spelling
- Broad familiarity across English dialects
- Avoidance of culturally sensitive, offensive, or domain specific terms
- Prefer short to medium length words
- Prefer unique phonetic structure and low pairwise similarity

No formal selection algorithm is defined here, but the dictionary must contain exactly 16384 normal words and 5460 trim words.

---

## 9. Security Considerations

### 9.1 Entropy preservation
EntropyLex preserves the entropy of the input bitstream. All mappings between binary values and words are deterministic and reversible.

### 9.2 No semantic leakage
Words are treated purely as codepoints and have no semantic relationship to the data they encode.

### 9.3 Human error
Mispronunciation, homophones, and transcription mistakes are major risk factors. Dictionary design must minimize confusable pairs.

### 9.4 Length leakage
Token count leaks the approximate payload size. This may or may not be acceptable depending on application.

---

## 10. Design Rationale Summary

- 14 bits per token provides high entropy density while still allowing selection of common English words
- Byte aligned input restricts trimming cases to even bit lengths only, reducing trim token requirements
- The absence of a length header simplifies decoding by embedding the required remainder information in the final trim token
- The dictionary size of 21844 tokens is large but still feasible for curated English word lists

---

## 11. Current Status

This draft describes the complete functional behavior of the encoding and decoding processes and the structural requirements for dictionary construction. Further work includes:

- Finalized dictionary selection
- Optional error detection schemes
- Optional interface bindings for various programming languages
- Optional compression or pre-processing guidelines
