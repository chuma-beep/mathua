# Ordering a mixed list of fractions and decimals

A list that mixes fractions and decimals is ordered by converting every entry to one form. Decimals suit a quick sort, while exact fractions settle values that look equal.

## The rule

Convert each fraction to a decimal by dividing or convert each decimal to a fraction by placing it over a power of ten. Then compare like with like. Keep one more decimal place than the closest pair needs, because two entries can differ by less than a hundredth. The fraction \(\frac{1}{4}\) is 0.25, so it sorts before 0.3.

## A worked example

Order \(\frac{3}{5}\), 0.58, \(\frac{7}{12}\) and 0.62. Dividing gives \(\frac{3}{5} = 0.6\) and \(\frac{7}{12} = 0.5833\) to four places. The list becomes 0.6, 0.58, 0.5833 and 0.62. In increasing order they are 0.58, then \(\frac{7}{12}\), then \(\frac{3}{5}\), then 0.62. The two middle values are close, so the extra place is what separates them.

## Checking the result

Check the close pair exactly. The difference is \(\frac{7}{12} - \frac{58}{100} = \frac{700}{1200} - \frac{696}{1200} = \frac{4}{1200}\), which is positive, so \(\frac{7}{12}\) is indeed above 0.58. The value 0.62 exceeds 0.6 by 0.02, so its place at the top is secure. A repeating decimal such as 0.5833 is finite in print but not in value, so rounding must not reverse a close comparison.

One form for every entry is the whole method and an exact check protects the close pairs.
