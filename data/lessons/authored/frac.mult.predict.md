# Predicting whether a product is larger or smaller

A product can be predicted before it is computed by comparing the multiplier with 1. Multiplication by a fraction either shrinks the other factor, preserves it or enlarges it.

## The method

For a positive quantity \(q\) and a multiplier \(\frac{a}{b}\), the product is smaller than \(q\) when \(a < b\) and larger when \(a > b\). Equality holds when the multiplier is 1, the boundary between the two cases. The comparison needs the multiplier to be positive. The prediction concerns the size of the product relative to \(q\), not its size relative to 1.

## A worked product

Take \(q = \frac{7}{9}\). Multiplying by \(\frac{6}{5}\) gives \(\frac{42}{45} = \frac{14}{15}\), about 0.933, which is larger than \(\frac{7}{9}\), about 0.778. Multiplying by \(\frac{4}{5}\) gives \(\frac{28}{45}\), about 0.622, which is smaller. The order of the three values is \(\frac{28}{45}\), then \(\frac{7}{9}\), then \(\frac{14}{15}\).

## Checking the result

Divide each product by its multiplier. For the first, \(\frac{14}{15} \div \frac{6}{5} = \frac{14}{15} \times \frac{5}{6} = \frac{70}{90} = \frac{7}{9}\). For the second, \(\frac{28}{45} \div \frac{4}{5} = \frac{28}{45} \times \frac{5}{4} = \frac{140}{180} = \frac{7}{9}\). Both checks return the starting value, so both predictions hold.

The multiplier alone tells the direction and the division confirms it.
