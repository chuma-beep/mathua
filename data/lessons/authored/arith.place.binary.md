# Binary place value

Binary writes a number with only zero and one, using powers of two as place values. It is the place-value system computers use.

## Powers of two as places

The places in a binary numeral are one, two, four, eight, and so on. Each digit says whether that power is present. Reading from the right, the place values double at each step: \(2^0 = 1\), \(2^1 = 2\), \(2^2 = 4\), \(2^3 = 8\), and \(2^4 = 16\).

**Example:** What is the largest power of two that is at most 21?
Double 1 repeatedly until the next double would exceed 21: \(1, 2, 4, 8, 16, 32\). Since \(16 \le 21\) and \(32 > 21\), the largest power of two is **16**.

## Converting a numeral

Multiply each bit by its place value and add the results. Reading from right to left, the exponents increase starting at \(2^0\).

**Example:** Convert \(1011_2\) to decimal.
Read the bits from right to left:
* \(1 \times 2^0 = 1\)
* \(1 \times 2^1 = 2\)
* \(0 \times 2^2 = 0\)
* \(1 \times 2^3 = 8\)

Add the results: \(1 + 2 + 0 + 8 = 11\). So \(1011_2 = 11\).

## Counting one bits

The number of one bits is the number of powers of two in the expansion. It is also called the popcount. Repeatedly halving the number and adding the remainders counts the same bits.

**Example:** How many 1 bits are in the binary expansion of 13?
Divide by 2 repeatedly and track the remainders:
* \(13 \div 2 = 6\) remainder 1
* \(6 \div 2 = 3\) remainder 0
* \(3 \div 2 = 1\) remainder 1
* \(1 \div 2 = 0\) remainder 1

Reading the remainders from bottom to top gives \(1101_2\). There are **3** ones in the expansion.

---

Practice should be answered without a calculator first. Write the relation, substitute the published values and state the integer result. Then check the result against the defining identity or a second representation before moving to the next question.
