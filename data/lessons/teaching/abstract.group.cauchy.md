# Cauchy's theorem for groups

Cauchy's theorem for groups is part of the abstract algebra curriculum. This lesson develops the idea in three sections, from the defining picture through the main computation to the interpretation.

## The Theorem

Cauchy's theorem states that if a finite group G has order n and a prime p divides n, then G contains an element of order p. The element need not be unique, and the subgroup it generates has order p by Lagrange's theorem. The result is a converse of the elementary fact that the order of an element divides the order of the group.

The theorem is weaker than the full Sylow theorems but much easier to state. It guarantees at least one element of each prime order, not a subgroup of the largest possible prime power. For a group of order 12, for example, it guarantees elements of orders 2 and 3, while Sylow theory describes the subgroups of orders 4 and 3.

The proof can be built from the class equation or from a group action on ordered tuples. The action approach counts tuples whose entries multiply to the identity and shows that the number of fixed tuples is divisible by p. A fixed tuple with all entries equal gives the required element.

## Prime Divisors

The number of distinct prime divisors of the group order is the number of different prime orders that Cauchy guarantees. For n = 60 = 2^2 times 3 times 5, the theorem guarantees elements of orders 2, 3 and 5. It does not guarantee elements of orders 4 or 6, because those are not prime.

The largest prime divisor gives the largest prime order that must appear. This is often useful in a quick existence question. A group of order 100 = 2^2 times 5^2 must contain an element of order 5, and the largest guaranteed prime order is 5, not 25.

Counting prime divisors is a matter of factoring n and discarding repeated factors. The exponent of a prime does not change the count, because the theorem asks for an element of order p rather than p to a power. That distinction is the most common source of errors in exercises.

## Prime Power Groups

A group whose order is p^k for a prime p is called a p-group. In such a group the only prime divisor is p, so Cauchy guarantees an element of order p. The centre of a nontrivial p-group is also nontrivial, which follows from the class equation and is a separate but related result.

Not every element of a p-group has order p. An element can have order p^2 or higher, up to p^k. The theorem only promises that at least one element has the smallest nontrivial prime order. For a cyclic group of order p^k there is exactly one subgroup of order p, and its non-identity elements all have order p.

If the order is written as p^k, the exponent k counts repeated prime factors. Finding it is the same as finding the exponent in the prime factorisation. For n = 81 = 3^4, the group is a 3-group and Cauchy guarantees an element of order 3, while the exponent of the group order is 4.
