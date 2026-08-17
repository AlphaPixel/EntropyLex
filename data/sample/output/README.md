# Sample outputs

This folder will contain EntropyLex phrases generated from the sample input files in `../input/`.

These files are illustrations, not compatibility tests. The shared test cases in `tests/vectors/` provide exact inputs and expected outputs that implementations must match. Sample phrases instead let readers see the phrase length and judge how an actual output looks and sounds.

No phrase files exist yet because no completed dictionary exists.

## Naming

```
<input-basename>.<profile>.<dictionary-id>.txt
```

For example: `cat_8x8_4-60bytes.el14.en-14-v1.txt`. The filename includes the dictionary identifier because changing the index-to-word mapping changes the phrase even when the payload and profile remain the same. The exact dictionary fingerprint must be recorded in an accompanying machine-readable record or document; the short filename identifier is not a substitute for it.

## Expected token counts for the current sample input

`cat_8x8_4-60bytes.gif` is 60 bytes = 480 bits:

| Profile | Tokens | Final remainder `r` |
|---|---|---|
| EL-8  | 60 | 0 |
| EL-12 | 40 | 0 |
| EL-14 | 35 | 4 |
| EL-16 | 30 | 0 |

Here `r` is the number of payload bits left after all complete normal tokens. For this sample, `r = 0` in every profile except EL-14, so only EL-14 ends with a trim token. The shared tests use several payload lengths to cover every possible remainder group.
