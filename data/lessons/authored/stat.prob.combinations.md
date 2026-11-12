# Counting unordered selections

A combination counts the subsets of a given size drawn from a larger set, so two selections that differ only in order are counted once.

## Order does not matter

A hand of cards is the standard case. The five cards dealt are the same hand whatever order they arrive in, so the count must not multiply by the orders. Choosing 3 of 7 books for a bag gives 35 selections rather than the 210 ordered sequences, because each selection has \(3! = 6\) orders.

## Dividing out the orderings

The formula is \(\binom{n}{r} = n!/(r!(n-r)!)\). The numerator counts ordered selections and the denominator divides by the r! orders within a selection and by the \((n-r)!\) orders of the objects left behind. For \(n=7\) and \(r=3\) the value is \(5040/(6 \times 24) = 35\). Row 7 of Pascal's triangle reads 1, 7, 21, 35, 35, 21, 7, 1.

## A worked count

A lottery draws 6 numbers from 50 without regard to order. The number of tickets is \(\binom{50}{6} = 50!/(6! \times 44!)\), which equals 15,890,700, so a single ticket has probability one in 15,890,700. An ordered draw would give \(50 \times 49 \times \cdots \times 45\), larger by a factor of 720.

The formula is small but the arithmetic is not; the safe habit is to write the falling product first and divide by the factorial afterwards.
