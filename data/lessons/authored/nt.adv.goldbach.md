# Goldbach partitions

A Goldbach partition writes an even number as a sum of two primes. Counting these partitions is a concrete way to explore the additive structure of primes.

## The partition definition

For an even n, a Goldbach partition is an unordered pair of primes p and q with p plus q equal to n. The order of the two primes does not create a second partition. The count is therefore the number of pairs with p at most q.

## Smallest and largest

The smallest prime in a partition is the first p that works as p increases. The largest is the last q in the ordered list. Either end can be read from the list of valid pairs.

## Counting unordered pairs

Enumerate p from 2 to n divided by 2. Keep p when both p and n minus p are prime. Each surviving p gives exactly one unordered partition.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.

A complete answer names the relation, substitutes the given values and states the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question.
