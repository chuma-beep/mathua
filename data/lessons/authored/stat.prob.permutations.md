# Counting ordered arrangements

A permutation counts arrangements of a chosen number of objects in which the order of those objects is part of the arrangement.

## Order matters

Two selections are different permutations when they contain the same objects in a different order. Choosing a president and a vice president from 6 candidates gives 30 ordered pairs, while choosing two delegates gives 15 unordered pairs. The factor of 2 between them is the two orders of one pair.

## The falling product

The number of ways to fill r ordered positions from n objects is \(n(n-1)\cdots(n-r+1)\), a product of r descending factors. Each factor counts the choices left once the earlier positions are filled. For 8 runners and 3 medals the product is \(8 \times 7 \times 6 = 336\).

## A worked count

A club of 12 members elects a chair, a secretary and a treasurer. The chair has 12 choices, the secretary 11 and the treasurer 10, so there are \(12 \times 11 \times 10 = 1320\) outcomes. A committee of three with no titles would give \(1320/6 = 220\), because each committee has \(3! = 6\) assignments to the roles.

Counting ordered arrangements is mostly a matter of deciding, before multiplying, whether two arrangements that swap objects count once or twice. That decision fixes whether the count is a permutation or a combination.
