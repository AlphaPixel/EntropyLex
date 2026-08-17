# Tests

Shared test data that every EntropyLex implementation can use, regardless of programming language.

The repository is intended to support independently written implementations in several programming languages. They need common examples with exact expected results so agreement can be measured rather than assumed. Those examples live here.

An implementation keeps tests of its own internal functions next to its code under `src/<language>/<author>/`. This folder instead contains the shared **conformance suite**: tests of externally visible behavior required by the specification. Every implementation must pass the applicable shared tests before it is described as conformant.

## Layout

```
tests/vectors/     JSON test cases for successful conversion and required rejection behavior
```

The folder retains the common name **test vector**, meaning a fixed input with its expected result. Its README generally uses the plainer term **test case**.

Some test cases will use actual words from a specific dictionary. Each such case must pass when that dictionary is loaded from authoritative LXJ or from its matching LXB. Deliberately invalid dictionary files and tests that compare LXJ with LXB belong under `tests/vectors/` once both file layouts are final.

## What implementations are expected to do

Load each test case, perform the stated operation, compare the actual result with the expected result, and report pass or fail for that case. An implementation that has not run the applicable shared tests must not be described as conformant.

Test sets have version numbers. New cases may be added without changing existing cases. If an existing implementation fails a new case, determine whether the implementation is wrong or the specification was ambiguous and correct the appropriate one.

See `vectors/README.md` for the format and coverage requirements.
