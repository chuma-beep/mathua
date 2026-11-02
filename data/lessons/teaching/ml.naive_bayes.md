# Naive Bayes classifiers

## Conditional independence

A naive Bayes classifier applies Bayes' rule to the class label and assumes that the features are conditionally independent given the class. The joint likelihood of a feature vector factors into one term per feature.

The assumption is rarely exactly true, but the classifier often works well because the decision only needs the largest posterior rather than a calibrated probability.

## Estimation and smoothing

Class priors come from the observed class frequencies. Feature probabilities come from counts within each class, and a small pseudocount is added so that an unseen value does not force the whole posterior to zero.

With Laplace smoothing, a count of c successes out of n trials over k possible outcomes becomes c plus one divided by n plus k. The smoothed value stays strictly between zero and one.

## Parameter counts

A multinomial model with k classes, d features and v possible values per feature stores k times d times v conditional probabilities, plus the class priors. A Bernoulli model stores two probabilities per feature per class, giving k times the quantity 2d plus one.

Because the parameter count grows linearly in the number of features, naive Bayes remains cheap to train even when the feature dimension is large.
