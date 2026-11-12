# Occupancy problems

Occupancy problems place balls into boxes and ask about empty boxes, collisions or the largest load.

## Empty boxes

With \(m\) balls and \(n\) boxes, the expected number of empty boxes is \(n(1-1/n)^m\).

## Poisson limit

When both counts grow together, the occupancy counts become independent Poisson variables.

## Maximum load

With as many balls as boxes, the largest load grows like \(\log n/\log\log n\).

To check your understanding, fix a box, let both counts grow and compare n balls with n boxes.

Occupancy counts balance emptiness against collisions.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
