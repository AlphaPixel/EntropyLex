# hanson-prototype

EL-8 encoder/decoder. Shows the EL-8 path works end to end. Not meant to be maintained, future implementation will succeed this one. No dependencies, Python 3.8+.

Does: **EL-8 only**, LXJ v1 dictionaries.
Doesn't: LXFP-1 fingerprint check, schema validation, LXB, anything wider than 8 bits, or anything else.

```
python el8.py encode -d ../../../tests/fixtures/dict/entropylex-en-8-test-v1.lxj -i key.bin -o phrase.txt
python el8.py decode -d ../../../tests/fixtures/dict/entropylex-en-8-test-v1.lxj -i phrase.txt -o key.bin
```

Options follow the common CLI in the top-level README (`-d`, `-i`, `-o`). Encode/decode is picked with the first argument instead of being two separate programs to keep it compact and simple.
