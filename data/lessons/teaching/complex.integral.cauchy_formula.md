# The Cauchy Integral Formula

## The formula

The Cauchy integral formula recovers the value of an analytic function inside a simple closed contour from its values on the contour. If f is analytic on and inside a positively oriented circle C centred at a, then f(a) is the contour integral of f(z) divided by z minus a, scaled by one over two pi i. The denominator is the only singular factor, and the contour surrounds its pole exactly once.

The formula is local as well as global. Shrinking the circle around a keeps the value f(a) unchanged, because the integral of an analytic function over an annulus between two circles is zero. That freedom lets us choose a small circle where estimates are easy, which is how the formula proves that an analytic function has derivatives of every order.

To use the formula, identify the point a, confirm that the contour winds once around it, and verify that the remaining factor is analytic inside. Then the integral reduces to evaluating the analytic factor at a. A contour that winds twice multiplies the answer by two, and a contour that does not surround a gives zero.

## A circle integral

For a monomial z to the n, the value at a is a to the n, so the formula converts a contour integral into a power. The same computation with the denominator squared gives the derivative, and repeating it gives every higher derivative. This is the standard way to evaluate integrals of the form z to the n divided by z minus a around a circle.

Worked examples usually choose a contour large enough to contain a and small enough to avoid other singularities. If the integrand has additional factors, check whether they are analytic in the enclosed region. A factor such as z minus b contributes only when b is also inside, in which case the integrand has two poles and the single-pole formula no longer applies directly.

A common error is to apply the formula when the contour is not positively oriented. Reversing the orientation changes the sign of the integral. Another error is to forget that the formula needs the point a strictly inside the contour; a point on the contour makes the integrand singular there and the hypotheses fail.

## Derivatives from the formula

Differentiating under the integral sign gives the general derivative formula: the k-th derivative at a is a contour integral of f(z) divided by z minus a to the k plus one, with a factorial and the same scaling. The formula shows that analyticity is a much stronger condition than real differentiability, because one derivative implies infinitely many.

The derivative formula also yields Cauchy's estimates. If f is bounded by M on a circle of radius r around a, then the k-th derivative at a is bounded by k factorial times M divided by r to the k. These estimates control growth and are the input to Liouville's theorem and to results about polynomial degree.

For computation, the most useful cases are the first and second derivatives. They appear when an integrand has a double or triple pole, and they let us evaluate residues without expanding a full Laurent series. The factorial in the formula is easy to forget, so check the order of the pole against the power in the denominator.

