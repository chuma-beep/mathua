# Arithmetic means between terms

An arithmetic mean sits between two terms of an arithmetic sequence, keeping the same common difference. Inserting means is therefore a question about equal gaps: the distance between the endpoints is divided into equal steps.

## One mean between two terms

A single arithmetic mean between \(a\) and \(c\) splits the interval into two equal gaps. The common difference is \((c - a)/2\), and the mean is \(a\) plus that difference.

For the endpoints 3 and 11, the difference is \(11 - 3 = 8\), and half of it is 4. The mean is \(3 + 4 = 7\). The sequence 3, 7, 11 has the same step from term to term.

The mean is also the average of the endpoints. That shortcut works only when exactly one mean is inserted, because two gaps then make the average the midpoint.

## The common difference

When \(k\) means are inserted between two endpoints, the interval is divided into \(k + 1\) equal gaps. The common difference is \((b - a)/(k + 1)\), where \(a\) and \(b\) are the first and last terms.

For 5 means between 2 and 32, there are 6 gaps. The difference is \((32-2)/6 = 5\). The sequence runs 2, 7, 12, 17, 22, 27, 32.

The count of gaps is one more than the count of means. Confusing the two is the most common error, and it shifts every inserted term.

## The later means

Once the common difference is known, each mean is a step from the previous term. The \(j\)-th mean is \(a + jd\), where \(d\) is the common difference.

With first term 2 and difference 5, the third mean is \(2 + 3 \times 5 = 17\). Counting from the first mean, not from the first term, is what makes the index \(j\) correct.

A useful check is symmetry. The means are equally spaced, so the first and last means are the same distance from the midpoint. If the computed list is not symmetric about the average of the endpoints, an index or a gap was miscounted.
