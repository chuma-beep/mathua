# Hensel's Lemma

## Lifting a Root

Hensel's lemma lifts a solution of a polynomial congruence modulo a prime to a solution modulo higher powers of that prime. The derivative at the known root controls the lifting. When the derivative is nonzero modulo the prime, the lift is unique at each stage and the root survives to every power.

## The Correction Step

The correction step is a linear congruence. If r solves f of x congruent to zero modulo p, write the lifted root as r plus p t. Expanding f and dividing by p gives a congruence for t modulo p, and the coefficient is the derivative f prime of r. Solving that linear congruence produces the next digit of the lifted root.

## When Lifting Fails

Lifting can fail when the derivative vanishes modulo the prime. A repeated root may split, persist or disappear as the modulus grows, and the simple correction formula no longer determines a unique lift. Recognising this case matters in algorithms that solve polynomial congruences and in the factorisation of polynomials over the p-adic numbers.
