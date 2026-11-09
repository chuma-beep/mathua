# Multiple comparisons and false discovery rate

Multiple testing raises the chance of a false positive when many hypotheses are tested at once. Bonferroni controls the family-wise error rate, and Benjamini-Hochberg controls the false discovery rate.

## Bonferroni correction

Bonferroni divides the target error rate by the number of tests. A p-value is rejected when it is at most alpha divided by m. The rule is simple and controls the family-wise error rate without assumptions about dependence.

## Family-wise error

The family-wise error rate is the probability of at least one false rejection. The Bonferroni bound is alpha times m before correction. It is conservative when many tests are correlated.

## Benjamini-Hochberg

Order the p-values from smallest to largest. Compare the i-th p-value with i times alpha divided by m and take the largest i that passes. This step-up rule controls the false discovery rate under independence or positive dependence.

The three sections form one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section so the three moves stay connected.
