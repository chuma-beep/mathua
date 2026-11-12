# The Pareto distribution

A Pareto distribution describes a positive quantity whose tail decays as a power of the value, so a small share of cases can hold a large share of the total.

## The tail index

The survival function is \(S(x) = (x_m/x)^\alpha\) for x at least the minimum \(x_m\), with tail index \(\alpha\). Small values of \(\alpha\) mean heavy tails. With \(x_m\) of 1000 and \(\alpha\) of 2, survival beyond 2000 is \((1/2)^2 = 0.25\), so a quarter of the cases exceed twice the minimum.

## The mean

The mean is \(\alpha x_m/(\alpha-1)\); it exists only when \(\alpha\) exceeds 1. With \(\alpha\) of 2 and \(x_m\) of 1000 the mean is 2000, twice the minimum. As \(\alpha\) approaches 1 the mean grows without bound; at or below 1 there is no finite mean at all. Near that boundary the sample mean wanders without settling.

## A worked example

Suppose incomes above 20,000 follow a Pareto tail with \(\alpha=1.5\). The mean of the tail is \(1.5 \times 20000/0.5 = 60000\). The share held by those above 80,000 is \(4^{1-1.5} = 4^{-0.5} = 0.5\), so half of the tail's income sits with the top quarter.

A power-law tail is a statement about relative proportions; it stays true when the units change.
