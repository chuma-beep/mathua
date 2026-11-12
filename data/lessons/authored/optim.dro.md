# Distributionally robust optimisation

Distributionally robust optimisation minimises the expected loss against a set of distributions around an empirical one. The adversary chooses a distribution rather than a point, which captures ambiguity in the sample.

## The ambiguity set

The set is defined by a divergence, a moment condition or a Wasserstein ball around the empirical measure. Its radius controls how much the adversary can move probability mass.

## The dual form

For a divergence ball the inner maximisation has a closed dual, and the worst-case risk becomes a regularised empirical risk. The regulariser penalises large losses rather than large parameters.

## The interpretation

The method sits between empirical risk minimisation and robust optimisation: it is less pessimistic than a worst-case point model and more conservative than the sample average.

To check your understanding, choose a divergence ball, derive the closed dual and place it between ERM and robust.

A distributionally robust objective is empirical risk with an adversarial reweighting, and the ambiguity radius is the dial between average and worst case.
