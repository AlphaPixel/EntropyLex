#!/usr/bin/env python3
# Quick and dirty EntropyLex EL-8 encoder/decoder.
# One byte = one word, so this is basically a 256 entry lookup table.
# Only EL-8.
#
#   python el8.py encode -d dict.lxj -i key.bin -o phrase.txt
#   python el8.py decode -d dict.lxj -i phrase.txt -o key.bin

import argparse
import json
import re
import string
import sys
import unicodedata

LOWER = str.maketrans(string.ascii_uppercase, string.ascii_lowercase)

# load dictionary. SLapdash, not much validation
def load_dict(path):
    with open(path, encoding="utf-8") as f:
        d = json.load(f)
    words = d["tokens"]
    rec = d["recognition"]
    # separators look like U+0020
    seps = "".join(chr(int(s[2:], 16)) for s in rec["tokenization"]["separators"])
    return words, seps, rec["case"]


#encode is super simple. Take byte, subscript, return.
def encode(data, words):
    return " ".join(words[b] for b in data)


def decode(text, words, seps, case):
    text = unicodedata.normalize("NFC", text)
    if case == "ascii-lower":
        text = text.translate(LOWER)
    tokens = re.split("[" + re.escape(seps) + "]+", text)
    # empty bits come from leading/trailing separators. index() blows up on unknown words
    return bytes(words.index(t) for t in tokens if t)


def main():
    p = argparse.ArgumentParser(description="EntropyLex EL-8 encoder/decoder (prototype)")
    p.add_argument("mode", choices=["encode", "decode"])
    p.add_argument("-d", "--dictionary", required=True, help="LXJ dictionary file")
    p.add_argument("-i", "--input", default="-", help="input file, - for stdin (default)")
    p.add_argument("-o", "--output", default="-", help="output file, - for stdout (default)")
    args = p.parse_args()

    words, seps, case = load_dict(args.dictionary)
    if args.input == "-":
        data = sys.stdin.buffer.read()
    else:
        with open(args.input, "rb") as f:
            data = f.read()

    if args.mode == "encode":
        result = encode(data, words).encode("utf-8")
    else:
        result = decode(data.decode("utf-8"), words, seps, case)

    if args.output == "-":
        sys.stdout.buffer.write(result)
    else:
        with open(args.output, "wb") as f:
            f.write(result)


if __name__ == "__main__":
    main()
