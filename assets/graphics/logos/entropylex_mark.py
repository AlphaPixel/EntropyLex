#!/usr/bin/env python3
"""
EntropyLex — E/L modular mark.
==============================

CONSTRUCTION
    Everything sits on a 7 x 7 grid of one module, a. The stroke weight is
    exactly one module, and so is every piece of whitespace inside the mark.

        E   upright, occupies modules 0-4 in both axes.
            Stem = column 0. Arms = rows 0, 2, 4, each running to column 4.
            The two counters are rows 1 and 3: one module tall.

        L   mirrored, hugging the right and bottom edges.
            Stem = column 6, full height. Foot = row 6.

        gap E's arms end at column 4; L's stem begins at column 6.
            The clear column 5 between them is one module wide.
            The clear row 5 between E's bottom arm and L's foot is likewise
            one module tall.

WHY IT WORKS
    Because the counters, the vertical clearance and the horizontal clearance
    are all exactly one module, the entire negative space inside the mark is a
    single ribbon of constant width. Traced from the top edge it runs down the
    channel, pulses twice to the left through the E's counters, turns along the
    bottom and exits at the lower left: a square-wave pulse train with a 1:1
    mark-to-space ratio, which is what a bytestream looks like on a scope.

    That constant width is the whole craft claim. If any clearance drifts, the
    wave stops reading and the mark just looks like two letters near each other.
    Do not adjust one gap without adjusting all of them.

NOT NEGOTIABLE
    - Never skew, italicise, or add stripes to the E. That silhouette belongs
      to a bankrupt energy company and the resemblance is instant.
    - Never colour the two letters differently. One mark, one colour.

Output: explicit closed paths. No masks, no strokes, no live text.
Usage:  python3 entropylex_mark.py
"""

import math
from pathlib import Path

INK = "#16181C"     # graphite
PAPER = "#F2F1ED"   # neutral off-white
SEAL = "#B7382C"    # accent, diagrams and the fingerprint seal only
WHITE = "#FFFFFF"   # ground for the white-flood chip ONLY, never the mark

# The black-flood chip grounds in INK, not in a separate near-black. An earlier
# hand-edited export used #231F20, Illustrator's CMYK-black conversion, which
# put an undeclared fourth colour into the identity. Grounding in INK makes the
# dark chip the exact inverse of the badge.

A = 150             # one module
R = 16              # corner fillet, ~11% of the module
OUT = Path(__file__).parent


def n(v):
    return f"{v:.2f}".rstrip("0").rstrip(".")


def fillet(pts, r):
    """Round every corner of a closed rectilinear polygon."""
    m, d = len(pts), []
    for i in range(m):
        p0, p1, p2 = pts[(i - 1) % m], pts[i], pts[(i + 1) % m]
        v1 = (p1[0] - p0[0], p1[1] - p0[1])
        v2 = (p2[0] - p1[0], p2[1] - p1[1])
        l1, l2 = math.hypot(*v1), math.hypot(*v2)
        u1, u2 = (v1[0] / l1, v1[1] / l1), (v2[0] / l2, v2[1] / l2)
        rr = min(r, l1 / 2, l2 / 2)
        a = (p1[0] - u1[0] * rr, p1[1] - u1[1] * rr)
        b = (p1[0] + u2[0] * rr, p1[1] + u2[1] * rr)
        sweep = 1 if u1[0] * u2[1] - u1[1] * u2[0] > 0 else 0
        d.append(("M " if i == 0 else "L ") + f"{n(a[0])} {n(a[1])}")
        d.append(f"A {n(rr)} {n(rr)} 0 0 {sweep} {n(b[0])} {n(b[1])}")
    return " ".join(d) + " Z"


def letters(a=A, r=R, foot_start=1, stem_start=1):
    """foot_start: column where L's foot begins. stem_start: row where L's stem
    begins. Both default to 1, which shortens the L by one module at each open
    end. That opens the negative ribbon at BOTH ends: it enters at the right
    edge on row 0, turns down the channel, pulses twice through the E's
    counters, turns along row 5 and exits at the bottom edge on column 0. Entry
    and exit are 180 degrees apart, which is the reversibility of the encoding
    stated as a shape. Set either to 0 to close that end instead."""
    S = 7 * a
    E = [(0, 0), (5 * a, 0), (5 * a, a), (a, a), (a, 2 * a), (5 * a, 2 * a),
         (5 * a, 3 * a), (a, 3 * a), (a, 4 * a), (5 * a, 4 * a), (5 * a, 5 * a), (0, 5 * a)]
    L = [(6 * a, stem_start * a), (7 * a, stem_start * a), (7 * a, 7 * a),
         (foot_start * a, 7 * a), (foot_start * a, 6 * a), (6 * a, 6 * a)]
    return S, fillet(E, r), fillet(L, r)



# ---------------------------------------------------------------- border
#
# The frame is a constant-width band, like everything else in this mark.
#
#   weight  b   = a/5. Thin enough that it never competes with the letters;
#                 heavy enough to survive a 1x device pixel at 48 px.
#   clear   gap = a exactly. Every other clearance in the mark is one module,
#                 so this one is too. It is also the mark's clear-space rule:
#                 nothing else may enter that ring.
#   fillet  r_in = a/3 on the inner edge, r_in + b on the outer, which makes
#                 the band a true parallel offset. The letters' own 16-unit
#                 fillet was tried here and reads as a sharp corner at frame
#                 scale, because roundness is judged against the size of the
#                 shape it sits on, not in absolute units.
#
# Floor: 48 px. Below that the band falls under one device pixel and the
# counters close up. Use the unbordered mark instead.

B_WEIGHT = A / 5
B_GAP = A
B_RADIUS = A / 3


def rrect(x, y, w, h, r):
    r = min(r, w / 2, h / 2)
    return (f"M {n(x+r)} {n(y)} H {n(x+w-r)} A {n(r)} {n(r)} 0 0 1 {n(x+w)} {n(y+r)} "
            f"V {n(y+h-r)} A {n(r)} {n(r)} 0 0 1 {n(x+w-r)} {n(y+h)} H {n(x+r)} "
            f"A {n(r)} {n(r)} 0 0 1 {n(x)} {n(y+h-r)} V {n(y+r)} "
            f"A {n(r)} {n(r)} 0 0 1 {n(x+r)} {n(y)} Z")


def bordered(a=A, r=R, b=B_WEIGHT, gap=B_GAP, r_in=B_RADIUS):
    S, pe, pl = letters(a=a, r=r)
    T = S + 2 * (gap + b)
    band = rrect(0, 0, T, T, r_in + b) + " " + rrect(b, b, T - 2 * b, T - 2 * b, r_in)
    inner = (f'    <path fill-rule="evenodd" d="{band}"/>\n'
             f'    <g transform="translate({n(gap+b)},{n(gap+b)})">'
             f'<path d="{pe}"/><path d="{pl}"/></g>')
    return T, band, inner


def bordered_circle(a=A, r=R, b=B_WEIGHT, gap=B_GAP):
    """Circular frame. Radius set so the mark's CORNERS clear the band by one
    module; the mark lands at ~57% of the diameter, normal for an avatar."""
    S, pe, pl = letters(a=a, r=r)
    ri = S * (2 ** 0.5) / 2 + gap
    T = 2 * (ri + b)
    inner = (f'    <path fill-rule="evenodd" d="M {n(T/2)} 0 A {n(T/2)} {n(T/2)} 0 1 0 {n(T/2)} {n(T)} '
             f'A {n(T/2)} {n(T/2)} 0 1 0 {n(T/2)} 0 Z M {n(T/2)} {n(b)} '
             f'A {n(ri)} {n(ri)} 0 1 1 {n(T/2)} {n(T-b)} A {n(ri)} {n(ri)} 0 1 1 {n(T/2)} {n(b)} Z"/>\n'
             f'    <g transform="translate({n(T/2 - S/2)},{n(T/2 - S/2)})">'
             f'<path d="{pe}"/><path d="{pl}"/></g>')
    return T, inner


def doc(vw, vh, inner, title, desc):
    return (f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {n(vw)} {n(vh)}" '
            f'width="{n(vw)}" height="{n(vh)}" role="img" aria-label="{title}">\n'
            f'  <title>{title}</title>\n  <desc>{desc}</desc>\n{inner}\n</svg>\n')


def main():
    S, pe, pl = letters()
    glyph = f'<path d="{pe}"/><path d="{pl}"/>'
    files = {}

    files["entropylex-mark.svg"] = doc(
        S, S, f'  <g fill="{INK}">{glyph}</g>', "EntropyLex mark",
        "E and L on a 7x7 module grid. All internal whitespace is one module wide, "
        "so the negative space reads as a constant-width square-wave pulse train "
        "that enters top right and exits bottom left.")

    Sc, pec, plc = letters(foot_start=0, stem_start=0)
    files["entropylex-mark-closed-L.svg"] = doc(
        Sc, Sc, f'  <g fill="{INK}"><path d="{pec}"/><path d="{plc}"/></g>',
        "EntropyLex mark, closed L",
        "Alternative: L runs full height and full width. Tighter silhouette, but "
        "the negative ribbon dead-ends instead of passing through.")

    pad = S * 0.26
    T = S + 2 * pad
    files["entropylex-avatar.svg"] = doc(
        T, T, f'  <circle cx="{n(T/2)}" cy="{n(T/2)}" r="{n(T/2)}" fill="{INK}"/>\n'
              f'  <g transform="translate({n(pad)},{n(pad)})" fill="{PAPER}">{glyph}</g>',
        "EntropyLex avatar", "Reversed in a circle. Mark at 66% of the diameter.")

    pad2 = S * 0.22
    T2 = S + 2 * pad2
    files["entropylex-appicon.svg"] = doc(
        T2, T2, f'  <rect width="{n(T2)}" height="{n(T2)}" rx="{n(T2*0.22)}" fill="{INK}"/>\n'
                f'  <g transform="translate({n(pad2)},{n(pad2)})" fill="{PAPER}">{glyph}</g>',
        "EntropyLex app icon", "Reversed in a rounded square, iOS/macOS corner radius.")

    # lockup ---------------------------------------------------------------
    mh = 240
    k = mh / S
    gap = mh * 0.42
    tx = S * k + gap
    fs = mh * 0.86
    stack = ("Söhne, 'ABC Diatype', 'GT America', 'Suisse Int\\'l', "
             "'Helvetica Neue', Arial, sans-serif")
    files["entropylex-lockup.svg"] = doc(
        tx + fs * 4.9, mh,
        f'  <g transform="scale({n(k)})" fill="{INK}">{glyph}</g>\n'
        f'  <!-- PLACEHOLDER TEXT: outline this with the licensed face before release -->\n'
        f'  <text x="{n(tx)}" y="{n(mh/2 + fs*0.72/2)}" font-family="{stack}" '
        f'font-size="{n(fs)}" letter-spacing="1" fill="{INK}">'
        f'<tspan font-weight="400">Entropy</tspan><tspan font-weight="600">Lex</tspan></text>',
        "EntropyLex horizontal lockup",
        "Mark plus wordmark. Cap height matches the E. The weight step from "
        "Entropy to Lex is the transcode.")

    # construction sheet ---------------------------------------------------
    a = A
    grid = "".join(
        f'<line x1="{i*a}" y1="0" x2="{i*a}" y2="{S}"/><line x1="0" y1="{i*a}" x2="{S}" y2="{i*a}"/>'
        for i in range(8))
    labels = "".join(
        f'<text x="{i*a + a/2}" y="-28" text-anchor="middle">{i}</text>'
        f'<text x="-30" y="{i*a + a/2 + 14}" text-anchor="end">{i}</text>' for i in range(7))
    m = 150
    panel2 = S + 260
    inner = (
        f'  <rect width="{n(S*2 + 260 + m*2)}" height="{n(S + m*2)}" fill="{PAPER}"/>\n'
        f'  <g transform="translate({m},{m})">\n'
        f'    <g stroke="{SEAL}" stroke-width="3" opacity="0.5">{grid}</g>\n'
        f'    <g font-family="ui-monospace, Menlo, monospace" font-size="40" fill="{SEAL}">{labels}</g>\n'
        f'    <g fill="{INK}">{glyph}</g>\n'
        f'  </g>\n'
        f'  <g transform="translate({m + panel2},{m})">\n'
        f'    <rect width="{S}" height="{S}" fill="{SEAL}"/>\n'
        f'    <g fill="{INK}">{glyph}</g>\n'
        f'  </g>\n'
        f'  <g font-family="ui-monospace, Menlo, monospace" font-size="38" fill="{INK}">\n'
        f'    <text x="{m}" y="{S + m + 90}">stroke = 1 module   every clearance = 1 module</text>\n'
        f'    <text x="{m + panel2}" y="{S + m + 90}">negative space: one constant-width ribbon</text>\n'
        f'  </g>')
    files["entropylex-construction.svg"] = doc(
        S * 2 + 260 + m * 2, S + m * 2 + 130, inner, "EntropyLex construction",
        "Left: the 7x7 module grid. Right: the negative space, showing the square-wave pulse train.")


    # bordered variants -----------------------------------------------------
    T, band, bin_ = bordered()
    files["entropylex-mark-bordered.svg"] = doc(
        T, T, f'  <g fill="{INK}">\n{bin_}\n  </g>', "EntropyLex mark, bordered",
        "Framed variant for light grounds. Band weight a/5, clear space a, "
        "inner fillet a/3.")

    files["entropylex-mark-bordered-inverse.svg"] = doc(
        T, T, f'  <g fill="{PAPER}">\n{bin_}\n  </g>',
        "EntropyLex mark, bordered, inverse", "Framed variant for dark grounds.")

    files["entropylex-mark-bordered-adaptive.svg"] = doc(
        T, T, f'  <g fill="currentColor">\n{bin_}\n  </g>',
        "EntropyLex mark, bordered, adaptive",
        "Inherits the surrounding text colour, so one file serves light and dark. "
        "Set colour on the parent element; do not add a fill attribute here.")

    files["entropylex-badge.svg"] = doc(
        T, T,
        f'  <path d="{rrect(0, 0, T, T, B_RADIUS + B_WEIGHT)}" fill="{PAPER}"/>\n'
        f'  <g fill="{INK}">\n{bin_}\n  </g>', "EntropyLex badge",
        "Opaque: carries its own ground, so it holds over photography or any "
        "unknown background.")

    # Flood chips. The ground is a rounded rect on the band's OUTER radius, so
    # the ink band paints over its edge and stays the outermost element. Beyond
    # that rect the file is transparent, which is what keeps the corners of the
    # rounded frame from flattening to a square when these are rasterised.
    ground = rrect(0, 0, T, T, B_RADIUS + B_WEIGHT)

    files["entropylex-mark-bordered-whiteflood.svg"] = doc(
        T, T,
        f'  <path d="{ground}" fill="{WHITE}"/>\n'
        f'  <g fill="{INK}">\n{bin_}\n  </g>',
        "EntropyLex mark, bordered, white flood",
        "Self-contained light chip: pure white ground inside the border, ink "
        "mark. Use where the surface is unknown or must read as true white.")

    files["entropylex-mark-bordered-inverse-blackflood.svg"] = doc(
        T, T,
        f'  <path d="{ground}" fill="{INK}"/>\n'
        f'  <g fill="{PAPER}">\n{bin_}\n  </g>',
        "EntropyLex mark, bordered, black flood",
        "Self-contained dark chip: ink ground inside the border, paper mark. "
        "The exact inverse of the badge. Default for dark-mode surfaces.")

    Tc, cin = bordered_circle()
    files["entropylex-avatar-bordered.svg"] = doc(
        Tc, Tc, f'  <g fill="currentColor">\n{cin}\n  </g>',
        "EntropyLex circular avatar, bordered",
        "Circular frame for round avatar crops. Mark at ~57% of the diameter so "
        "its corners keep one module of clear space.")

    # Circular flood chips. Same reasoning as the square pair: the ground is a
    # disc on the band's outer radius, so the ring paints over its edge. These
    # are also what an avatar upload actually wants — GitHub, Discord and the
    # like composite a transparent avatar against a surface you do not control.
    disc = (f'  <circle cx="{n(Tc/2)}" cy="{n(Tc/2)}" r="{n(Tc/2)}" '
            'fill="{}"/>\n')

    files["entropylex-avatar-bordered-whiteflood.svg"] = doc(
        Tc, Tc, disc.format(WHITE) + f'  <g fill="{INK}">\n{cin}\n  </g>',
        "EntropyLex circular avatar, bordered, white flood",
        "Self-contained light circular chip: pure white ground, ink mark.")

    files["entropylex-avatar-bordered-inverse-blackflood.svg"] = doc(
        Tc, Tc, disc.format(INK) + f'  <g fill="{PAPER}">\n{cin}\n  </g>',
        "EntropyLex circular avatar, bordered, black flood",
        "Self-contained dark circular chip: ink ground, paper mark.")

    for name, body in files.items():
        # newline="" stops the platform translating line feeds into CRLF.
        # Without it a Windows run rewrites every line of every asset and the
        # whole set shows up as modified with no visible change.
        (OUT / name).write_text(body, encoding="utf-8", newline="")
        print("wrote", name)


if __name__ == "__main__":
    main()
