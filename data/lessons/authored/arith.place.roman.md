# Roman numerals

Roman numerals use seven symbols. Add the symbols from left to right, and subtract when a smaller value is written before a larger one.

## The idea

Each symbol has a fixed value.

| Symbol | Value |
| --- | --- |
| I | 1 |
| V | 5 |
| X | 10 |
| L | 50 |
| C | 100 |
| D | 500 |
| M | 1000 |

Add the values when symbols go from larger to smaller: XII is 10 + 1 + 1 = 12. Subtract when a smaller symbol comes immediately before a larger one: IV is 5 - 1 = 4. Only six subtractive pairs are used.

| Pair | Value |
| --- | --- |
| IV | 4 |
| IX | 9 |
| XL | 40 |
| XC | 90 |
| CD | 400 |
| CM | 900 |

Every other numeral adds its symbols, and no symbol repeats more than three times. So 4 is IV, not IIII, and 40 is XL, not XXXX.

## The computation

Split the numeral at its subtractive pairs, then add the values of the parts.

- XIV splits into X and IV: 10 + 4 = 14.
- XLII splits into XL, I and I: 40 + 1 + 1 = 42.
- MCMXCIV splits into M, CM, XC and IV: 1000 + 900 + 90 + 4 = 1994.

To write a number as a numeral, take the largest symbol value that fits, subtract it and repeat. For 1987: 1000 leaves 987, so write M; 900 leaves 87, so write CM; 50 leaves 37, so write L; 10 leaves 27, then another 10 leaves 17, and another 10 leaves 7, so write XXX; 5 leaves 2, so write V; then two I symbols. The numeral is MCMLXXXVII.

## Reading the result

Read from left to right. A smaller symbol before a larger one subtracts, and a smaller symbol after a larger one adds. Recompute the total from the table to check the answer. The numerals in this course run from 1 to 399, so a result outside that range is a sign that a pair or a repeated symbol was misread.
