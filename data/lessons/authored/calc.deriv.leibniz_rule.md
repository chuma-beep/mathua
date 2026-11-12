# The Leibniz rule for higher derivatives

The Leibniz rule expands the \(n\)-th derivative of a product into a binomial sum. Each term pairs a derivative of one factor with a complementary derivative of the other.

## The term count

The \(n\)-th derivative of a product has \(n+1\) terms. The indices run from zero to \(n\), one for each way to split the order between the factors.

## The coefficient

The coefficient of the term with \(k\) derivatives on the second factor is \(\binom{n}{k}\). The binomial coefficients are symmetric, so the two ends of the sum match.

## The second derivative

For \(n=2\) the rule gives three terms: \(f''g+2f'g'+fg''\). The middle coefficient is two, which is the only case most courses write out in full.

To check your understanding, count the splits, choose the order split and expand three terms.

The product rule generalises to a binomial sum over derivative orders.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
