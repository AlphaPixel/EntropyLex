This folder will contain implementations of EntropyLex in the language denoted by the folder name (e.g. "python").

Each implementation will be in a subfolder, named for the author of the implementation. E.g. ("satoshi-nakamoto").

Each implementation must state which EntropyLex profiles it supports (EL-8, EL-12, EL-14, EL-16). The expected build order is EL-8 first (no bit splitting, no trim words), then EL-12 (nibble-boundary splits, small trim dictionary), then EL-14. Implementations supporting more than one profile should use a single code path parameterized by token width, not one implementation per profile. See ../../SPEC.md sections 3 and 15.
