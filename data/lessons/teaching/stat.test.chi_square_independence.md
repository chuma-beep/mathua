# The chi-square test of independence

The chi-square test of independence asks whether two categorical variables are associated in a population.

## Two categorical variables

Arrange the counts in a contingency table. Rows are one variable and columns are the other. Each cell contains the number of cases in that row-column combination.

The hypotheses are:

- $H_0$: the two variables are independent.
- $H_a$: the two variables are associated.

The test compares the observed counts with the counts we would expect if the variables were independent.

## Expected counts and degrees of freedom

For a cell in row $i$ and column $j$,

$$E_{ij} = \frac{(\text{row total}_i)(\text{column total}_j)}{\text{grand total}}.$$

A table with $r$ rows and $c$ columns has

$$df = (r-1)(c-1)$$

degrees of freedom. A $3 \times 2$ table therefore has $(3-1)(2-1)=2$ degrees of freedom.

## Cell contributions

The test statistic is

$$\chi^2 = \sum \frac{(O-E)^2}{E}.$$

Each cell contributes its squared gap from expectation, divided by the expected count. Large contributions are evidence against independence.

The test requires expected counts that are not too small. A common rule is that every expected count should be at least 5. When that fails, combine categories or use an exact test.
