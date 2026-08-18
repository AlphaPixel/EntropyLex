This folder will contain EntropyLex encoders and decoders written in Go.

Each independently written implementation belongs in a subfolder named for its author, for example `satoshi-nakamoto`.

Each implementation must state which profiles it supports: EL-8, EL-12, EL-14, or EL-16. The intended order is EL-8 first, where one token represents one byte and no trim token is needed; then EL-12, where token boundaries may occur after four bits of a byte and a small trim portion is required; then EL-14. Code supporting more than one profile should take token width and related counts as parameters rather than duplicate the algorithm. See ../../SPEC.md sections 3 and 15.

Dictionary loaders must support LXJ (`.lxj`), the authoritative JSON format. LXJ v1 is defined by `../../data/dict/FORMAT.md` and `../../data/dict/lxj-v1.schema.json`; a complete EL-8 test file is at `../../tests/fixtures/dict/entropylex-en-8-test-v1.lxj`. LXB (`.lxb`) remains an optional binary form generated from LXJ. A future loader supporting LXB must verify that its fingerprint, index-to-token mapping, and written-token rules match the corresponding LXJ dictionary.
