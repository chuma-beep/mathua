# Fisher's exact test

Fisher's exact test assesses association in a two by two table by computing the probability of the observed table and of every table at least as extreme, conditional on the margins.

## Fixed margins

The test conditions on the row and column totals, so only one cell is free to vary. With 10 treated and 10 control subjects and 12 successes in all, the top-left cell determines the rest: if it holds 7 the other three cells follow. Conditioning on the margins makes the null distribution a hypergeometric one.

## The range of a cell

The free cell cannot take every value from 0 to its row total, because the margins bound it. Its lower bound is the larger of zero and the row total plus the column total minus the grand total. For a row total of 10, a column total of 12 and a grand total of 20, the cell ranges from 2 to 10.

## Why it is exact

The p-value is a sum of hypergeometric probabilities rather than an approximation, so it is valid at any sample size, including 20 subjects. The probability of a table with top-left cell a is \(\binom{10}{a}\binom{10}{12-a}/\binom{20}{12}\). For a of 9 that is \(10 \times 120/125970\), about 0.0095.

The test is exact because it enumerates rather than approximates; that is also why it is practical only while the margins stay small.
