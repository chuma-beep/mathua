# Domain adaptation

Domain adaptation transfers a model from a source domain with labels to a target domain without them. It assumes the two domains share something worth transferring.

## The discrepancy

The source and target feature distributions differ, and the discrepancy measures that difference. Minimising it aligns the representations.

## The sample sizes

With \(m\) source and \(n\) target points the adaptation sees both sets. The target labels are absent, so the objective uses only its features.

## The assumption

Adaptation needs the label rule to be shared or close. Without that assumption no method can identify the target rule.

To check your understanding, measure the difference, count the source and share the label rule.

Domain adaptation aligns representations so a source model transfers.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
