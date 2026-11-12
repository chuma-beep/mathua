# McNemar's test

McNemar's test compares two matched binary measurements by looking only at the pairs that disagree, since the pairs that agree carry no information about a change.

## Paired binary data

Each subject is measured twice, before and after a treatment, with each measurement yes or no. The data form a two by two table of pairs: both yes, yes then no, no then yes and both no. With 80 subjects, suppose 30 are yes on both occasions, 8 are yes then no, 22 are no then yes and 20 are no on both. The 50 agreeing pairs say nothing about the direction of change.

## Discordant pairs

Only the 8 and the 22 are discordant. They are the whole evidence. Under the null hypothesis of no change each discordant pair is equally likely to fall in either cell, so the 30 pairs split like 30 fair coin tosses. The expected count in each discordant cell is 15, while an observed 8 against 22 is 7 away from it.

## The test statistic

The statistic is \(\chi^2 = (b-c)^2/(b+c)\), with b and c the discordant counts, compared against a chi-square distribution on 1 degree of freedom. Here \((8-22)^2/30 = 196/30\), about 6.53, which exceeds the 3.84 critical value at the 5 percent level. When the discordant total is below about 25, the exact binomial version is preferred, because the chi-square approximation is then poor.

The design is the point: pairing removes between-subject variation; the test asks only whether the disagreements lean one way.
