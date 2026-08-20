#!/usr/bin/env python3
"""Derive the 7x7 ASCII forms from the same module logic as the SVG."""

from pathlib import Path

N = 7

def grid():
    g = [[0]*N for _ in range(N)]
    for r in range(5): g[r][0] = 1                 # E stem, col 0, rows 0-4
    for r in (0, 2, 4):
        for c in range(5): g[r][c] = 1             # E arms, cols 0-4
    for r in range(1, 7): g[r][6] = 1              # L stem, col 6, rows 1-6
    for c in range(1, 7): g[6][c] = 1              # L foot, row 6, cols 1-6
    return g

def invert(g): return [[1-v for v in row] for row in g]

def render_solid(g, ch="#"):
    return [ "".join(ch if v else " " for v in row) for row in g ]

def render_line(g):
    out = []
    for r in range(N):
        row = ""
        for c in range(N):
            if not g[r][c]:
                row += " "; continue
            up    = r > 0     and g[r-1][c]
            down  = r < N-1   and g[r+1][c]
            left  = c > 0     and g[r][c-1]
            right = c < N-1   and g[r][c+1]
            h, v = left or right, up or down
            row += "+" if (h and v) else ("-" if h else ("|" if v else "+"))
        out.append(row)
    return out

def frame(lines):
    w = max(len(l) for l in lines)
    body = [" " + l.ljust(w) + " " for l in lines]
    edge = " " * (w + 2)
    return "\n".join([edge] + body + [edge]) + "\n"

g = grid(); ng = invert(g)
files = {
    "entropylex-7x7-hash-positive.txt": frame(render_solid(g)),
    "entropylex-7x7-hash-negative.txt": frame(render_solid(ng)),
    "entropylex-7x7-ascii-positive.txt": frame(render_line(g)),
    "entropylex-7x7-ascii-negative.txt": frame(render_line(ng)),
}
OUT = Path(__file__).parent

for name, text in files.items():
    # newline="" keeps the committed files LF on every platform; writing beside
    # the script keeps them out of whatever directory it was invoked from.
    (OUT / name).write_text(text, encoding="utf-8", newline="")
    print("=== " + name)
    print(text.replace(" ", "\u00b7") if False else text, end="")
    print("--- widths:", {len(l) for l in text.rstrip("\n").split("\n")})
