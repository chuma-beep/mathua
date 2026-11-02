# Taylor Remainder and Error Bounds

A Taylor polynomial approximates a function near a point. The remainder measures what the polynomial leaves out, and an error bound says how large that remainder can be.

## The remainder term

Write \(f(x)=P_n(x)+R_n(x)\), where \(P_n\) is the Taylor polynomial of degree \(n\) at \(a\). The remainder \(R_n(x)\) is the exact difference between the function and the polynomial. It is not the same as the next term of the series, although the next term often gives a useful estimate.

For \(f(x)=e^x\) at \(a=0\), the degree-two polynomial is \(1+x+x^2/2\). At \(x=1\) it gives 2.5 while \(e\) is about 2.718, so the remainder is about 0.218.

## Bounding the remainder

The Lagrange form says that for some point \(c\) between \(a\) and \(x\), \(R_n(x)=f^{(n+1)}(c)(x-a)^{n+1}/(n+1)!\). If \(|f^{(n+1)}|\le M\) on the interval, then \(|R_n(x)|\le M|x-a|^{n+1}/(n+1)!\).

For \(e^x\) on \([0,1]\), the derivative is always \(e^x\), so \(M=e\). The degree-four bound at \(x=1\) is \(e/5!\), about 0.0227. A larger interval or a lower degree gives a weaker bound.

## Choosing enough terms

To guarantee an error below a tolerance, choose \(n\) so the bound is smaller than the tolerance. For \(e^x\) on \([0,1]\), the bound is \(e/(n+1)!\). At \(n=5\) it is about 0.0038, and at \(n=6\) it is about 0.00054.

For an alternating series that satisfies the Leibniz conditions, the remainder is bounded by the first omitted term. That test is often sharper than a general Taylor bound, but it applies only when the series alternates and its terms decrease to zero.
