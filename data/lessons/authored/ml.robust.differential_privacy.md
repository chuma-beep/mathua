# Differential privacy

Differential privacy bounds how much one record can change the output distribution. The bound is a ratio of probabilities controlled by a privacy parameter.

## The definition

Neighbouring data sets must produce outputs whose probabilities differ by at most a factor \(e^\epsilon\). Smaller \(\epsilon\) means stronger privacy.

## The Laplace mechanism

Noise scaled by sensitivity over \(\epsilon\) gives \(\epsilon\)-differential privacy. The sensitivity is the largest change one record can cause.

## Composition

Running \(k\) mechanisms with parameters \(\epsilon_i\) costs at most their sum. Composition is why privacy budgets are spent carefully.

To check your understanding, compare neighbours, compute sensitivity and add the budgets.

Differential privacy is a worst-case guarantee per record.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
