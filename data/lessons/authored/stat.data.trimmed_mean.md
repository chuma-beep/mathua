# The trimmed mean

A trimmed mean discards a fixed share of the largest and smallest observations before averaging what remains, so a few extreme values cannot drag the average far.

## Why trim

An ordinary mean gives every observation the same leverage, so one outlier can move it a long way. Nine measurements near 20 with a single reading of 900 have a mean of 108, while the median is 20. A trimmed mean sits between the two, keeping more information than the median because it averages the values near the centre.

## Removing the extremes

A 10 percent trimmed mean removes the same count from both ends, not 10 percent of each tail independently. With 20 ordered values, 10 percent is 2 values, so drop the smallest 2 and the largest 2 and average the remaining 16. With 25 values the calculation gives 2.5, which must be rounded; the rounding rule changes the answer.

## The trimmed average

Take the ordered values 4, 8, 9, 11, 12, 13, 15 and 60. Dropping the smallest and the largest leaves 8, 9, 11, 12, 13 and 15, whose sum is 68 and whose mean is \(68/6 = 11.33\). The untrimmed mean is \(132/8 = 16.5\), pulled upward by the 60, so reporting both states how much the extreme value mattered.

Trimming is a decision about which data to ignore, so the share trimmed belongs in the report rather than in a footnote.
