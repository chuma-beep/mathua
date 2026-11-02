# Multiclass classification

## Reducing to binary problems

One-versus-rest trains one binary classifier per class, each separating that class from all the others. A problem with k classes therefore needs k classifiers.

One-versus-one trains a classifier for every pair of classes. The number of pairs is k times k minus one divided by two, and a vote combines their predictions.

## Direct multiclass models

A softmax model produces one score per class and converts the scores into probabilities that sum to one. Training minimises the cross-entropy between the predicted distribution and the true class.

For n examples and k classes the score matrix has n times k entries, and the predicted label is the column with the largest score in each row.

## Evaluation

Accuracy alone can hide poor performance on rare classes. A confusion matrix records counts for every true and predicted pair, and per-class precision and recall are read from its rows and columns.

Macro averaging treats every class equally, while micro averaging weights each class by its number of examples. The two agree only when the classes are balanced.
