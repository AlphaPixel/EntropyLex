This folder contains the Hanson-Claude Python implementation of EntropyLex. It is an independent implementation that follows the structure of the Go implementation under `go/arcadium`; it is not an Arcadium project.

The implementation is currently only a skeleton. Its common interface, EL-8 dictionary holder, tests, build commands, and continuous-integration configuration intentionally stop at the same point as the Go skeleton. Encoding, decoding, and LXJ loading are not implemented yet.

Planned profile order is EL-8, then EL-12, then EL-14. All profiles will use the same encoding and decoding functions, with token width and related counts supplied as profile parameters.

The dictionary loader will first support LXJ version 1 as defined in `../../../data/dict/FORMAT.md` and `../../../data/dict/lxj-v1.schema.json`. Its initial target is the complete test dictionary at `../../../tests/fixtures/dict/entropylex-en-8-test-v1.lxj`. Optional LXB (`.lxb`) support will follow only after its byte layout is final. Loading either form must produce the same dictionary fingerprint, index-to-token mapping, and written-token behavior.
