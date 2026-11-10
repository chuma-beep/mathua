# Voting classifiers

## Votes

A voting classifier combines the labels predicted by several models. In hard voting each model contributes one vote. The label with the most votes wins.

The vote count is not the same as a probability. A model that predicts a label contributes one vote regardless of its confidence.

## Majority

A majority requires more than half of the votes. With five voters, three votes are a majority. With an even number of voters, a tie is possible.

The winning margin is the difference between the winning and losing vote counts. A larger margin means the ensemble is less divided.

## Ties

A tie has no strict majority. The tie rule decides the output, for example by choosing the smallest label or by using a fallback model. The rule must be fixed before evaluation.

Report the tie rule with the result. A tie broken arbitrarily makes the ensemble difficult to reproduce.

A final check is to restate the answer in the units and context of the original question. If the result is plausible and the units match, the calculation is ready.
