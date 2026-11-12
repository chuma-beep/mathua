# Addition formulas for hyperbolic functions

Hyperbolic functions obey addition formulas like the circular ones, except that the sign of a product term is reversed.

## The double-angle formulas

Setting the two arguments equal collapses each addition formula. Thus
\(\sinh 2u = 2\sinh u\cosh u\) and \(\cosh 2u = \cosh^{2}u + \sinh^{2}u\). At \(u = 1\), with
\(\sinh 1 \approx 1.1752\) and \(\cosh 1 \approx 1.5431\), the first gives
\(2(1.1752)(1.5431) \approx 3.6269\), which is \(\sinh 2\). The second gives
\(2(1.5431)^{2} - 1 \approx 3.7622\), which is \(\cosh 2\). Both come out positive, unlike the
circular double angle, which can change sign.

## Adding two hyperbolic angles

The general formulas are \(\sinh(u+v) = \sinh u\cosh v + \cosh u\sinh v\) and
\(\cosh(u+v) = \cosh u\cosh v + \sinh u\sinh v\). With \(u = 1\) and \(v = 2\) the first reads
\(1.1752(3.7622) + 1.5431(3.6269) \approx 10.018\), which is \(\sinh 3\). The plus sign inside the
cosine formula is the difference from the circular case, where a minus appears. That sign is the
whole difference between the two families of formulas.

## Products of sinh and cosh

Read backwards, the addition formulas turn products into sums. The identity
\(2\sinh u\sinh v = \cosh(u+v) - \cosh(u-v)\) gives, at \(u = 2\) and \(v = 1\),
\(10.0677 - 1.5431 \approx 8.5246\), matching \(2(3.6269)(1.1752)\). The companion identity
\(2\cosh u\cosh v = \cosh(u+v) + \cosh(u-v)\) gives \(10.0677 + 1.5431 \approx 11.611\).

The formulas are worth memorising in pairs, since each has a companion. Check any hyperbolic
identity against the circular one and watch the signs.
