# Domain of a Composite Function

## Composition Order

The composite function f(g(x)) applies g first and then f. Its domain is not simply the domain of g. A value x must be allowed by g, and the output g(x) must be allowed by f. Both conditions have to hold, so the domain of the composite is the set of x in the domain of g for which g(x) lies in the domain of f. Missing the second condition is the usual mistake.

## Inner Domain

A systematic method is to find the inner domain first, then solve the inequality or equation that keeps the inner output inside the outer domain. If g(x) = x - 4 and f(u) = sqrt(u), then g is defined everywhere but the composite requires x - 4 >= 0, so the domain is x >= 4. If the inner function already has a restriction, such as a denominator, intersect that restriction with the new one rather than replacing it.

## Outer Restrictions

Composition is not commutative in general, and the two orders can have different domains. The notation (f o g)(x) means f(g(x)), not g(f(x)). When simplifying a composite expression, keep the domain restriction even if the simplified formula appears to be defined more widely, because the simplification can cancel a restriction that was present in the original composition.
