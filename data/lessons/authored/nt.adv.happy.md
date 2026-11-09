# Happy numbers

A happy number reaches 1 under the map that sums the squares of the digits. Unhappy numbers eventually enter a cycle.

## The digit-square map

Replace a number by the sum of the squares of its decimal digits. Repeat the operation. The map is deterministic, so the sequence is fixed once the starting number is chosen.

## The happiness test

A number is happy when its digit-square sequence reaches 1. It is unhappy when the sequence enters the cycle that contains 4. Every starting value does one or the other.

## Searching forward

To find the next happy number after a bound, test successive integers with the digit-square map. The first one that reaches 1 is the answer. The search is finite and usually short.

The three sections are one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section, so the three moves stay connected.

A complete answer names the relation, substitutes the given values and states the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question.
