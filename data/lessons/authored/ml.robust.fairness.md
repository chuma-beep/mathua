# Fairness metrics

Fairness metrics compare model behaviour across groups. Different definitions cannot all hold at once when base rates differ.

## Demographic parity

The positive rate must be equal across groups. The difference is the largest gap between those rates.

## Equalised odds

True positive and false positive rates must match across groups. This is stricter than parity because it conditions on the true label.

## The trade-off

An imperfect predictor cannot satisfy every fairness definition simultaneously. Choosing one is a policy decision, not a mathematical one.

To check your understanding, compare positive rates, match true positive rates and name the conflict.

Fairness is a family of constraints, and the choice among them is explicit.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
