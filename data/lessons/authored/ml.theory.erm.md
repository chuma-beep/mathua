# Empirical risk minimisation

Empirical risk minimisation chooses the hypothesis with the smallest average loss on the training sample. Generalisation depends on how tightly that average tracks the population risk.

## The average

The empirical risk is a sum of losses divided by the sample size. Every point carries the same weight unless the sample is reweighted.

## The 0-1 loss

Under the 0-1 loss the empirical risk is the misclassification rate. Multiplying by the sample size gives the error count.

## The gap

The gap between empirical and population risk is controlled by capacity. A richer class can drive the empirical risk to zero while widening the gap.

To check your understanding, sum the losses, count mistakes and measure the gap.

ERM is the direct approach, and capacity is what makes it generalise.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
