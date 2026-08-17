# LXJ dictionary fixtures

`entropylex-en-8-test-v1.lxj` is the complete, valid LXJ v1 example and the bootstrap dictionary for the first EL-8 implementation.

It contains 256 lowercase English words in numerical index order. Encoding and decoding examples belong to the shared test cases, not to the dictionary artifact description; see [`../../vectors/el8-test-v1.json`](../../vectors/el8-test-v1.json) and [`../../vectors/README.md`](../../vectors/README.md).

The companion `.tokens.txt` and `.settings.json` files are the exact inputs named in the LXJ provenance fields. Their line endings are fixed to LF by the repository's `.gitattributes`, so their recorded SHA-256 values remain stable across operating systems.

This is a test fixture, not the first production English dictionary. Its words were chosen manually to exercise loading and conversion. They have not passed the familiarity, pronunciation, spelling-similarity, meaning-similarity, or sensitive-language review required for a production dictionary. The LXJ file therefore declares `"purpose": "test"`.

Expected identifiers:

| Value | SHA-256 |
|---|---|
| Dictionary fingerprint (`LXFP-1`) | `191b8d6b489c54f99dfdab3d09de03b3f6f64d9859fc668fba1da48f6d52fd7f` |
| Token source file | `60090a007f06bd35143197fd8e53929bfca34b07ab111ad248110f75037b4302` |
| Selection settings file | `eade2b844783214e9d0a1028492f1305e686e64341b57e0a0e272507ce0eddee` |
