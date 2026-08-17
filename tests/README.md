# Tests

Shared test data that every EntropyLex implementation can use, regardless of programming language.

The repository is intended to support independently written implementations in several programming languages. They need common examples with exact expected results so agreement can be measured rather than assumed. Those examples live here.

An implementation keeps tests of its own internal functions next to its code under `src/<language>/<author>/`. This folder instead contains the shared **conformance suite**: tests of externally visible behavior required by the specification. Every implementation must pass the applicable shared tests before it is described as conformant.

## Layout

```
tests/fixtures/dict/  Complete dictionary files and their exact construction inputs
tests/vectors/        JSON test cases for successful conversion and required rejection behavior
```

The folder retains the common name **test vector**, meaning a fixed input with its expected result. Its README generally uses the plainer term **test case**.

Some test cases use actual words from a specific dictionary. The first valid fixture is `fixtures/dict/entropylex-en-8-test-v1.lxj`, a 256-entry LXJ v1 dictionary for EL-8 loader and byte-conversion work. Each dictionary-dependent case must pass when its dictionary is loaded from authoritative LXJ or, after LXB is defined, from its matching binary form. Deliberately invalid LXJ files belong under `tests/fixtures/dict/invalid/`; corresponding LXB cases will follow after that layout is defined.

## What implementations are expected to do

Load each test case, perform the stated operation, compare the actual result with the expected result, and report pass or fail for that case. An implementation that has not run the applicable shared tests must not be described as conformant.

Test sets have version numbers. New cases may be added without changing existing cases. If an existing implementation fails a new case, determine whether the implementation is wrong or the specification was ambiguous and correct the appropriate one.

See `vectors/README.md` for the format and coverage requirements.
