# EntropyLex trademarks and the EntropyLex mark

[`LICENSE`](LICENSE) places the EntropyLex **software and documentation** under
the MIT License. That grant is deliberately broad: copy it, change it, ship it
in a commercial product, no permission needed.

It does **not** grant rights in the EntropyLex name or logo. Trademarks answer a
different question from copyright — not *may I use this work?* but *may I put
this label on my thing?* An open licence that also handed over the name would
leave nobody able to tell an official release from anyone else's fork, which
helps no one. This document says what you may do with the marks without asking.

## 1. What this covers

The **Marks** are:

- the word mark **EntropyLex**, and confusingly similar variants;
- the **E/L logo** — the mirrored-L mark described in
  [`assets/graphics/logos/DESIGN.md`](assets/graphics/logos/DESIGN.md), in every
  variant and every format, raster or vector; and
- the ASCII transcriptions of that logo in `assets/graphics/logos/`.

`entropylex_mark.py` and `ascii_mark.py` are source code and are MIT-licensed
like everything else. The MIT grant covers the *scripts*; it does not convert
their output into an unrestricted work. Running a generator produces the Mark,
and the Mark is governed by this document.

## 2. What you may do without asking

- **Say what your software does.** "Implements EntropyLex", "EntropyLex-compatible",
  "an EntropyLex encoder for Rust", "supports EntropyLex EL-8 and EL-12". Accurate
  statements of fact need no permission, and a conformance claim backed by the
  published test vectors is exactly such a statement.
- **Name the project in prose** — articles, talks, papers, documentation,
  courseware, comparisons, and critical or unfavourable commentary alike.
- **Use the unmodified logo to refer to this project**, at a size and clearance
  consistent with `DESIGN.md`, in a context where it plainly identifies
  EntropyLex rather than you.
- **Redistribute the repository** with the Marks intact, including forks. If you
  publish a modified fork, see section 3.
- **Link to the project** using the name or the logo.

## 3. What requires written permission

- Using a Mark, or anything confusingly similar, as **the name or logo of your
  own** product, service, company, organisation, app, or package — including as
  a prefix or suffix (`EntropyLex Pro`, `EntropyLexHub`, `entropylex-cloud`).
- Any use implying **endorsement, affiliation, sponsorship, certification, or
  official status** that does not exist.
- **Modifying the logo**: recolouring outside the declared palette, skewing,
  rotating, redrawing, animating, adding effects, changing its proportions, or
  combining it with other elements into a new mark.
- Applying a Mark to **merchandise** offered for sale.
- Registering a Mark, a confusingly similar mark, or a corresponding domain name
  or social-media handle.

Requests are welcome and are usually granted for reasonable uses. Contact
AlphaPixel LLC. https://alphapixeldev.com

## 4. Modified versions

If you distribute a modified version, you may state factually what it is
derived from — "a fork of EntropyLex", "based on EntropyLex" — but you may not
name it in a way that suggests it is the official project, and you should not
ship the logo as the identity of your fork. Choose your own name and mark, and
credit EntropyLex in your description.

This matters most for the wire format. A dictionary or profile that diverges
from [`SPEC.md`](SPEC.md) but still calls itself EntropyLex makes the name
useless as a compatibility signal, which is the one job it does.

## 5. Using the logo correctly

`DESIGN.md` section 9 states the construction rules. Where you have permission
to use the logo, follow them: never recolour the two letters differently, never
place anything inside the clear-space ring, never use a framed variant below
48 px, and never skew or add stripes to the E.

Ship a self-contained variant — the flood chips or the badge — anywhere the
background is not known. The transparent and `currentColor` variants disappear
against the wrong surface, and a broken mark serves nobody.

## 6. Reservation of rights

AlphaPixel LLC reserves all rights in the Marks not expressly granted here. This
document may be updated; the version in the repository at the time you obtain it
governs your use. Nothing here limits fair use, nominative use, or equivalent
rights under applicable law, and nothing here restricts your rights under the
MIT License to the software itself.
