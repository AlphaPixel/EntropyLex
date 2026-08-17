# Tests

Shared, language-agnostic conformance data for EntropyLex.

The premise of this repository is multiple independent implementations, in multiple languages, by multiple authors. That is only meaningful if there is something concrete for them to agree on. This folder is that something.

Implementations keep their own unit tests next to their own code, under `src/<language>/<implementor>/`. What lives here is the shared conformance suite that **every** implementation must pass, regardless of language.

## Layout

```
tests/vectors/     Round-trip and validation vectors, as JSON
```

Dictionary-dependent vectors cover both canonical LXJ loading and matching LXB loading. Malformed-file fixtures and LXJ/LXB equivalence cases belong under `tests/vectors/` once the schemas and binary layout are finalized.

## What implementations are expected to do

Load the vectors, run them, and report pass or fail per vector. An implementation that has not run the shared vectors should not be described as conformant.

Vectors are versioned and additive. Adding a vector that an existing implementation fails is a legitimate outcome — it means either the implementation or the specification needed the correction, and both are worth discovering.

See `vectors/README.md` for the format and coverage requirements.
