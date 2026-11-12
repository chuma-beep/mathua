# Sample complexity of learning

Sample complexity is the number of examples needed to meet an error and confidence target. It converts a learning guarantee into a data requirement.

## The growth

The bound grows with the capacity and with \(1/\epsilon\). The confidence enters through a logarithm of \(1/\delta\).

## Finite classes

For a class of size \(N\) the bound is proportional to \(\log N\). Enumerating the class is not required, only counting it.

## Infinite classes

Infinite classes use VC dimension or Rademacher complexity in place of the logarithm. The shape of the bound stays the same.

To check your understanding, invert epsilon, count the class and replace the count.

Sample complexity is capacity divided by accuracy, with confidence on top.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
