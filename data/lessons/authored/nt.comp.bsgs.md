# Baby-step giant-step

Baby-step giant-step solves a discrete logarithm by trading memory for time. It meets in the middle at the square root of the group order.

## The split

Write the unknown exponent as i m plus j, where m is the ceiling of the square root of the order. The equation separates the two unknowns.

## The tables

The baby steps store small powers in a table. The giant steps multiply by the inverse of the generator raised to m and look for a match.

## The cost

Both the time and the memory are proportional to the square root of the order. That square root is the benchmark every discrete logarithm method is measured against.

To check your understanding, take the square root, build baby steps and count time.

Baby-step giant-step is a meet-in-the-middle search, and the square root in its name is the number of entries on each side.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
