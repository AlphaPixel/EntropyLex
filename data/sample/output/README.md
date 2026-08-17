# Sample outputs

Encoded forms of the payloads in `../input/`, for eyeballing and for documentation examples.

These are illustrative, not authoritative. The authoritative round-trip data is the test vector set in `tests/vectors/` — that is what implementations must agree on. Sample outputs exist so a human can look at a real phrase and see what the encoding actually feels like to read aloud.

Nothing here yet: sample outputs cannot be generated until a dictionary exists.

## Naming

```
<input-basename>.<profile>.<dictionary-id>.txt
```

For example `cat_8x8_4-60bytes.el14.en-14-v1.txt`. The dictionary identity must be part of the filename, because the same payload under the same profile produces a completely different phrase under a different dictionary.

## Expected token counts for the current sample input

`cat_8x8_4-60bytes.gif` is 60 bytes = 480 bits:

| Profile | Tokens | Final remainder `r` |
|---|---|---|
| EL-8  | 60 | 0 |
| EL-12 | 40 | 0 |
| EL-14 | 35 | 4 |
| EL-16 | 30 | 0 |

Note that this payload lands on `r = 0` for every profile except EL-14, so it exercises trim-token logic only under EL-14. It is not a substitute for the test vectors, which cover every remainder class.
