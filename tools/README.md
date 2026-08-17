# EntropyLex Tools

This folder is reserved for programs that select, generate, and verify EntropyLex dictionaries. SPEC.md section 11 defines the planned sequence of tools. No tools have been implemented yet; the current subfolders only define where implementations will live.

Except for possible manual review of EL-8, dictionaries will not be selected entirely by hand. Given byte-for-byte identical source files and settings, the tools must produce the same ordered mapping. They must also record every exclusion and every change to an allowed similarity limit so reviewers can inspect how the result was produced. Purely manual selection cannot provide those properties consistently for more than 21,000 words.

## Layout

Tools follow the same convention as `src/`:

```
tools/<language>/<implementor-name>/
```

- The language folder names the programming language used to write the tool, such as `python`.
- Each implementation lives in a subfolder named for its author (e.g. `satoshi-nakamoto`), exactly as in `src/`.

More than one author may implement the same tool independently. Comparing their output checks both the code and the specification. When given identical source files and settings, implementations must produce the same ordered mapping and dictionary fingerprint. Once the LXJ and LXB layouts are final, the generated files must also match byte for byte. A disagreement indicates an implementation error or a rule that the specification has not defined precisely enough.

## The tool sequence

The tools form a pipeline: each numbered stage consumes recorded inputs and produces recorded output for the next stage.

| Stage | Tool | Responsibility |
|---|---|---|
| 1 | `eldict-ingest`  | Load word, frequency, pronunciation, and grammatical data from different sources into one consistently formatted candidate table |
| 2 | `eldict-filter`  | Apply rules that always exclude a candidate and record the reason for every removal |
| 3 | `eldict-score`   | Calculate numerical spelling, pronunciation, and meaning comparisons between candidates |
| 4 | `eldict-select`  | Mark pairs that are too similar and select the required number of nonconflicting tokens |
| 5 | `eldict-emit`    | Assign normal and trim roles and numerical indices, then write LXJ |
| 6 | `eldict-compile` | Generate optional LXB from LXJ without changing the mapping, rules, or fingerprint |
| 7 | `eldict-verify`  | Check file structure directly; when exact source files and settings are available, also repeat quality checks and write a JSON report |

Every stage must be deterministic: identical inputs and settings produce identical output. Each source dataset must be fixed to an unchangeable release or exact repository commit, identified by a SHA-256 value calculated from every downloaded file byte, and recorded in separate machine-readable LXJ fields. Reviewed candidates and their license conditions are listed in `../data/dict/SOURCES.md`.

Stage 3 allows each writing system to define appropriate similarity measures. A Chinese dictionary, for example, can compare visual character shape and pronunciation written in the Latin alphabet instead of using the English spelling measures. Stages 1, 2, 4, 5, 6, and 7 retain the same roles. See SPEC.md section 12.4.

## Outputs

Generated LXJ and LXB files and their verification reports belong in `data/dict/`, not here. LXJ is authoritative; LXB is generated from LXJ and must contain the same mapping and rules.

## Open design question

Whether EL-8, EL-12, and EL-14 share any words is **not yet decided**. The alternatives are nested sets, separate selection with uncontrolled overlap, and completely separate token sets. This choice changes stages 4 and 5 and must be settled before `eldict-select` is implemented. See SPEC.md section 11.8.

## Reference

- SPEC.md section 10 — word selection criteria
- SPEC.md section 11 — dictionary derivation, stage by stage
- SPEC.md section 11.7 — LXJ/LXB formats, source history, and what each verification type can establish
- SPEC.md section 11.8 — the undecided question of whether profiles share tokens
- SPEC.md section 12 — dictionaries for other languages and writing systems
- `../data/dict/FORMAT.md` — working file-format design
- `../data/dict/SOURCES.md` — candidate datasets, licenses, and counts
