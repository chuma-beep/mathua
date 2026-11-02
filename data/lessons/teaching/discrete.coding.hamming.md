# Hamming Codes

## Parity and Redundancy

A Hamming code protects m data bits by adding r parity bits. The parity bits are placed at positions that are powers of two and each covers a different pattern of positions.

The receiver recomputes the parity checks. The pattern of failures, called the syndrome, identifies the position of a single flipped bit, which can then be corrected.

Redundancy is the price of correction: the transmitted word is longer than the data it carries.

## Choosing the Number of Parity Bits

The r parity bits must be able to name every one of the m + r positions plus the no-error case, so the condition is 2 to the power r at least m + r + 1.

For small data sizes the smallest valid r can be found by testing r equal to one, two, three and so on. The resulting code length is m + r.

A code with minimum distance three corrects one error, because the received word is always closest to a unique codeword.

## Codewords and Correction

A systematic Hamming code keeps the data bits in their original order and inserts parity bits between them. A non-systematic code mixes them, but the number of codewords is still 2 to the power m.

The minimum distance of the code determines how many errors can be detected and corrected. Distance d corrects the floor of (d - 1)/2 errors.

These counts let you compare codes by rate and by protection without listing every codeword.
