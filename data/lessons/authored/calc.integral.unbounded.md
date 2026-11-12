# Improper integrals with unbounded integrands

## An unbounded integrand

An integrand can be unbounded at an endpoint and still have finite area. For \(\int_0^1 x^{-1/2} dx\) the antiderivative is \(2\sqrt{x}\), so the integral from \(a\) to 1 is \(2 - 2\sqrt{a}\), which tends to 2 as \(a\) falls to zero. The area under the spike is finite.

## The comparison test

The comparison test settles most cases. If \(f\) lies between 0 and \(g\) near the bad point and the integral of \(g\) converges, so does the integral of \(f\). On the interval from 0 to 1 the bound \(x^{-1/2}/(1+x) \le x^{-1/2}\) holds. The larger integral is 2, so the smaller one converges too.

## When the integral converges

For \(x^{-p}\) at zero the test is \(p < 1\). The cases differ sharply: \(\int_0^1 x^{-1/2} dx = 2\) converges, while \(\int_0^1 x^{-1} dx\) diverges, since the value from \(a\) to 1 is \(-\ln a\), which is 6.9078 at \(a = 0.001\). The boundary \(p = 1\) is the divergent one.

A complete answer names the relation, substitutes the data and states the result in the language of the problem. When a second method is available, compare the two results before moving on.
