# Data poisoning

Data poisoning corrupts training data so that the learned model misbehaves. The attacker controls a fraction of the training set rather than the input at test time.

## The poisoning fraction

If an attacker controls \(k\) of \(n\) training points, the fraction is \(k/n\). Defences often bound that fraction explicitly.

## The influence

Influence functions estimate how one training point changes a prediction. A large influence marks a point worth auditing.

## The target

An attack can target availability by degrading overall accuracy. A targeted attack changes one chosen prediction instead.

To check your understanding, count the controlled points, estimate the influence and degrade accuracy.

Poisoning attacks the training set, so the model learns the wrong thing.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
