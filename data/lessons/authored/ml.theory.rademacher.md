# Rademacher complexity

Rademacher complexity measures how well a class correlates with random signs on a sample. It gives a data-dependent generalisation bound.

## The definition

Average the largest correlation between the class and random sign vectors. The expectation is over both the sample and the signs.

## The bound

Generalisation error is at most empirical error plus twice the complexity plus a confidence term. The confidence term shrinks like \(1/\sqrt{n}\).

## The sign vectors

With \(n\) points there are \(2^n\) sign vectors. The empirical estimate samples them instead of enumerating.

To check your understanding, draw the signs, add the complexity and count the vectors.

Rademacher complexity asks whether the class can fit noise.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
