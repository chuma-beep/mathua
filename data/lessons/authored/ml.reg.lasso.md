# Lasso regression

Lasso regression adds an L1 penalty to least squares. The penalty shrinks coefficients and can set some of them exactly to zero.

## The L1 penalty

The L1 penalty is lambda times the sum of the absolute values of the weights. It grows linearly with each weight rather than quadratically. The penalty value is computed before the data term is combined with it.

## Soft thresholding

In one dimension the lasso coefficient is the soft threshold of the least-squares correlation. Values inside the band from minus lambda to lambda become zero. Values outside the band shrink toward zero by lambda.

## The objective

The full objective is one half of the squared residual plus the L1 penalty. Doubling it clears the fraction and gives the residual squared plus twice lambda times the absolute weight. This is the integer form used in the exercises.

The three sections form one chain: the first names the object, the second carries out the computation and the third reads the result in context.

A complete answer names the relation, substitutes the quantities in the question and states the result in the language of the concept. Where a second computation is available, compare the two results before moving on.
