"""EL-8 dictionary-backed encoding."""


class Encoding:
    """Hold the byte-indexed token dictionary used by EL-8."""

    def __init__(self, dictionary: dict[int, str]) -> None:
        self._dictionary = dictionary


def new_encoding(dictionary: dict[int, str]) -> Encoding:
    """Return an EL-8 encoding backed by dictionary."""
    return Encoding(dictionary)
