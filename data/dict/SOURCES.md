# Candidate dictionary sources

This is the working inventory of source material that may be used to derive EntropyLex dictionaries. Inclusion here means that the source appears reusable under stated terms; it does not mean that it has been selected or that its license obligations disappear from a derived dictionary.

Counts use the source publisher's unit. A “word,” pronunciation record, vector, match form, and spell-checker record are not interchangeable. Counts marked **measured** were calculated from the named downloadable file; counts marked **published** come from the project's own documentation. Every source actually used must later be pinned to an immutable release or commit, downloaded again, counted after ingestion, checksummed, and reviewed with its full license text.

Reviewed 2026-08-17. This is an engineering license screen, not legal advice.

The repository's MIT software license does not erase the licenses on input data. Published dictionaries will need their own attribution/notice file and, if a ShareAlike source is selected, an explicit decision about the dictionary artifact's license.

## Recommended English source stack

| Role | Source and version reviewed | Count | License | Assessment |
|---|---|---:|---|---|
| Primary lemmas, parts of speech, and semantic relations | [Open English WordNet 2025](https://en-word.net/downloads) core edition | 135,969 words; 107,519 synsets; 355,064 relations (published) | [CC BY 4.0](https://github.com/globalwordnet/english-wordnet/blob/main/LICENSE.md) | **Recommended primary candidate.** The 2025 core edition excludes the separate proper-name resource and supplies the lexical structure needed for lemma and synonym checks. Its published count precedes EntropyLex's single-token, length, and character filters. Attribution is required. |
| Spelling, variants, basic POS, inflections, and commonness bands | [English Speller Database / SCOWL dictionary release `2026.02.25`](https://github.com/en-wl/wordlist/releases/tag/rel-2026.02.25) | 47,551 records in `en_US.dic` (measured); affix rules generate additional surface forms | [MIT-like with retained notices](https://github.com/en-wl/wordlist/blob/v2/Copyright) | **Recommended cross-check and expansion source.** The current release is dictionary-only while the new database format stabilizes. Pin the database commit when the derivation pipeline begins; do not mistake Hunspell records for distinct expanded words. |
| Common and frequency-ranked lemmas | [12dicts 6.0.2](https://wordlist.aspell.net/12dicts/) | about 25,000 in `2+2+3cmn`, 34,000 in `2+2+3frq`, and 84,000 in `2+2+3lem` (published) | Most lists are public domain; the three AGID-derived `2+2+3` lists retain [AGID/source notices](https://wordlist.aspell.net/agid-readme/) | **Recommended seed and cross-check.** The 25,000-word core is unusually close to the EntropyLex scale, but it is not large enough to survive all filters or supply three disjoint profiles by itself. Preserve all applicable notices. |
| Pronunciation and homophone detection | [CMU Pronouncing Dictionary](https://github.com/cmusphinx/cmudict) commit [`0f8072f`](https://github.com/cmusphinx/cmudict/commit/0f8072f814306c5ee4fbf992ed853601b12c01f9) | 135,166 pronunciation records and 126,052 distinct headwords after removing numbered pronunciation suffixes (measured) | Use for research or commercial purposes is unrestricted; acknowledgment requested in the [project README](https://github.com/cmusphinx/cmudict) | **Recommended for US English.** Multiple records per headword are valuable for identifying divergent pronunciations. It is not a general dialect-coverage source. |
| Semantic-distance vectors | [Stanford GloVe](https://nlp.stanford.edu/projects/glove/) 2024 Wikipedia + Gigaword or Dolma models | 1.2 million vectors in each 2024 vocabulary (published); the older 6B model has 400,000 | [Open Data Commons PDDL 1.0](https://opendatacommons.org/licenses/pddl/1-0/) | **Recommended embedding candidate.** It has the lowest licensing friction among the reviewed pretrained vectors. The 100-dimensional 2024 Wikipedia + Gigaword model is likely sufficient, but that is a benchmark choice. |
| Sensitive-term exclusions | [LDNOOBW English list](https://github.com/LDNOOBW/List-of-Dirty-Naughty-Obscene-and-Otherwise-Bad-Words) commit [`4638b97`](https://github.com/LDNOOBW/List-of-Dirty-Naughty-Obscene-and-Otherwise-Bad-Words/commit/4638b970cb8d9d82789564fcba1f4a1eb508ff1a) | 403 entries (measured), including phrases and non-word forms | [CC BY 4.0](https://github.com/LDNOOBW/List-of-Dirty-Naughty-Obscene-and-Otherwise-Bad-Words/blob/master/LICENSE) | **Recommended supplemental blocklist.** It is intentionally subjective and too small to be the only sensitive-term review. Exact single-token matching should be separated from phrase and pattern handling. |
| Sensitive-term exclusions with severity and exceptions | [dsojevic/profanity-list](https://github.com/dsojevic/profanity-list) | 434 English concepts and 809 match forms (published) | [MIT](https://github.com/dsojevic/profanity-list/blob/master/LICENSE) | **Recommended second opinion.** Severity, tags, and exceptions are useful, but patterns must not be treated as literal dictionary tokens. |

The preferred initial experiment is the intersection and union analysis of Open English WordNet, ESDB/SCOWL, 12dicts, and CMUdict, with GloVe coverage and the two blocklists applied as annotations. No single source should be treated as the finished EntropyLex vocabulary.

## Conditional sources

These sources are technically useful and redistributable, but their ShareAlike or source-specific conditions may affect how generated dictionary artifacts can be licensed. They require an explicit project licensing decision before use in the selection score.

| Role | Source and version reviewed | Count | License | Assessment |
|---|---|---:|---|---|
| Multi-corpus frequency ranking | [`wordfreq` 3.1.1](https://pypi.org/project/wordfreq/3.1.1/) | English `small`: 28,917 tokens; English `large`: 321,180 tokens (measured from the 3.1.1 wheel) | Apache 2.0 code; packaged data is [CC BY-SA 4.0 with additional source acknowledgments](https://github.com/rspeer/wordfreq#license) | **Strong technically; conditional legally.** It combines seven English frequency families and is the best reviewed familiarity measure, but using its scores to select the published dictionary may bring attribution and ShareAlike obligations. |
| Multilingual pronunciations | [WikiPron](https://github.com/CUNY-CL/wikipron) | More than 5 million word-pronunciation pairs across 348 languages (published) | Apache 2.0 code; mined Wiktionary data has separate Wiktionary licensing terms | **Useful beyond English; conditional.** Pin both the WikiPron release/configuration and the underlying data license. CMUdict remains preferable for the first US-English experiment. |
| Multilingual semantic vectors | [fastText vectors for 157 languages](https://fasttext.cc/docs/en/crawl-vectors) | Up to 2 million English vectors in the published Common Crawl model; 1 million in the Wikipedia/news model | [CC BY-SA 3.0](https://fasttext.cc/docs/en/english-vectors.html) | **Useful multilingual fallback; conditional.** The files are frequency ordered and support broad language coverage, but GloVe is simpler for the first English dictionary. |

## Sources not accepted merely because they are downloadable

- A GitHub repository without an explicit data license is not a usable source.
- “Free for research” or noncommercial data is not suitable for a generally redistributable EntropyLex dictionary.
- Search-engine and commercial-corpus frequency lists require their own redistribution review; a downstream project mentioning them does not grant EntropyLex permission.
- Large game wordlists are not automatically useful. Counts dominated by archaic, technical, inflected, or invented forms work against the familiarity requirement.
- A permissively licensed extraction tool does not change the license of the data it extracts. Raw Wiktionary/Wikipedia material must follow its applicable terms; a separately published trained model such as GloVe follows the model license stated by its publisher.

## Count-reproduction notes

The measured counts above were obtained as follows:

- `wordfreq-3.1.1-py3-none-any.whl`, SHA-256 `4b1c6ecffc6198be3396d5cf871c4423ca71c907c231348d352dd54d62b97473`: sum the token buckets in `small_en.msgpack.gz` and `large_en.msgpack.gz`.
- `hunspell-en_US-2026.02.25.zip`, SHA-256 `ac8e73310e951d88c52c2cf2ba54ceaca34f8486a81630ac8a75dc5f931179f9`: compare the declared first-line count in `en_US.dic` with its remaining record count; both are 47,551.
- CMUdict commit `0f8072f`, `cmudict.dict` SHA-256 `81917843c7f44ce2b094ac63873c2c7a4cf802040792c455ba3ca406891c3d22`: count non-comment records, then remove terminal numbered pronunciation markers such as `(2)` before counting distinct headwords.
- LDNOOBW commit `4638b97`: count the 403 lines in the English `en` file. The count includes phrases and forms that EntropyLex's token filters will later reject.

These are source counts, not projected EntropyLex survivor counts. The ingest stage must publish counts after every filter so corpus sufficiency—especially the proposed 26,468-token disjoint family—can be assessed with evidence.
