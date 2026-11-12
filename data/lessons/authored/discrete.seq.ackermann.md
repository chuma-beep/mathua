# The Ackermann function

The Ackermann function is a recursive function that grows faster than any primitive recursive function. Its definition is simple, but the resulting values dwarf tower notation almost immediately.

## The definition

\(A(0,n)=n+1\), \(A(m+1,0)=A(m,1)\), and \(A(m+1,n+1)=A(m,A(m+1,n))\). The recursion nests one level deeper at every step.

## Growth

\(A(1,n)=n+2\), \(A(2,n)=2n+3\), and \(A(3,n)=2^{n+3}-3\). Even the fourth row is a tower of powers.

## Non-primitive recursion

Primitive recursion bounds the number of nested loops in advance. Ackermann's nesting depends on the input, which is exactly what primitive recursion forbids.

To check your understanding, apply the base cases, derive the first rows and bound primitive loops.

Ackermann's function is the standard witness that not every computable function is primitive recursive. Totality and fast growth coexist with a three-line definition.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
