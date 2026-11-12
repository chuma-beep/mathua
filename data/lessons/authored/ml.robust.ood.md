# Out-of-distribution detection

Out-of-distribution detection flags inputs that do not belong to the training distribution. A score summarises unfamiliarity, and a threshold turns the score into a decision.

## The score

A confidence score, an energy or a distance in feature space can serve. The choice fixes which unfamiliar inputs are detectable.

## The threshold

A threshold trades false alarms against missed detections. Moving it changes both rates in opposite directions.

## The evaluation

AUROC summarises the ranking across all thresholds. It is invariant to the base rate of unfamiliar inputs.

To check your understanding, pick a score, move the threshold and compute AUROC.

OOD detection is a ranking problem with a threshold on top.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
