"""Generic EntropyLex stream wrappers."""

from typing import BinaryIO, Protocol


class Encoding(Protocol):
    """Operations required by an EntropyLex stream wrapper."""

    def decode(self, destination: bytearray, source: bytes) -> int:
        """Decode source into destination and return the number of bytes written."""
        ...

    def encode(self, destination: bytearray, source: bytes) -> None:
        """Encode source into destination."""
        ...


def new_encoder(encoding: Encoding, writer: BinaryIO) -> BinaryIO:
    """Return a stream that encodes data with encoding before writing to writer."""
    raise NotImplementedError


def new_decoder(encoding: Encoding, reader: BinaryIO) -> BinaryIO:
    """Return a stream that decodes data read from reader with encoding."""
    raise NotImplementedError
