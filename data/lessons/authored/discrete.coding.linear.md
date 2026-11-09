# Linear codes

A linear code is a vector subspace of a finite vector space. The Hamming code encodes four message bits into seven transmitted bits and corrects one error.

## The Hamming code

The Hamming [7,4] code has four message positions and three parity positions. Each parity bit checks a fixed subset of the message bits. The resulting codewords are the binary vectors whose parity checks vanish.

## Distance and error correction

The Hamming distance between two words is the number of positions in which they differ. A code with minimum distance three corrects one error. The nearest codeword to a received word is the unique codeword at distance at most one.

## Encoding messages

Write the four message bits in the first four positions. Compute each parity bit as the sum modulo two of its assigned message positions. The seven bits together form the codeword.

The three sections form one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section so the three moves stay connected.
