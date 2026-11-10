# Synthetic minority oversampling

SMOTE creates synthetic minority examples by interpolating between a minority point and one of its nearest minority neighbours.

## The synthetic point

Choose a minority point and one of its k nearest minority neighbours. Pick a random interpolation weight between zero and one. The synthetic point is the weighted average of the two coordinates.

## Neighbour count

The parameter k limits the neighbours considered for each minority point. A larger k gives more variety but may cross a class boundary. The count is part of the algorithm's configuration.

## Class balance

Each synthetic example increases the minority count by one. The new total is the original count plus the number created. SMOTE changes the class ratio without collecting new data.

A complete answer names the relation, substitutes the given values and states the result in the language of the problem. Where a second method is available, compare the two results before moving on.
The same three moves apply to every generated instance of this concept. Keep the arithmetic exact and check the sign or range before submitting.
