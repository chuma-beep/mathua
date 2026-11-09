# Geometric means between terms

A geometric mean between two terms keeps a constant ratio rather than a constant difference. Inserting means is a question about equal multiplicative steps, so the arithmetic uses roots and powers.

## One mean between two terms

A single geometric mean between positive numbers \(a\) and \(c\) is \(\sqrt{ac}\). It is the number that makes the ratio from \(a\) to the mean equal to the ratio from the mean to \(c\).

For the endpoints 2 and 18, the product is 36 and the square root is 6. The sequence 2, 6, 18 has ratio 3 at each step, so the mean satisfies the defining property.

Both endpoints must be positive for the real geometric mean to exist. A negative product has no real square root, and the geometric mean of a mixed-sign pair is not a real number.

## The common ratio

When \(k\) means are inserted between \(a\) and \(b\), there are \(k + 1\) equal ratios. The common ratio is the \((k+1)\)-th root of \(b/a\), so \(r = (b/a)^{1/(k+1)}\).

For 2 means between 3 and 81, there are 3 ratios and \(r = (81/3)^{1/3} = 27^{1/3} = 3\). The sequence is 3, 9, 27, 81.

The ratio can be below one. Between 81 and 3 the ratio is \(1/3\), and the sequence decreases. The root formula handles that case without a separate rule.

## The later means

Each mean is the first term multiplied by a power of the ratio. The \(j\)-th mean is \(a r^{j}\), where \(r\) is the common ratio.

With first term 3 and ratio 3, the second mean is \(3 \times 3^{2} = 27\). The exponent counts the number of multiplicative steps from the first term, not the number of means already written.

A check is to multiply the ratio through the list and confirm that the last step lands on the final endpoint. If the final term is not reached, the root or the exponent was taken with the wrong count.
