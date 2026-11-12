# Weibull lifetime distributions

A Weibull distribution models a positive lifetime with a scale parameter and a shape parameter, giving a survival function that decays as a power of time.

## The survival function

The survival function is \(S(t) = \exp(-(t/\lambda)^k)\), the probability that a unit survives past time t, with scale \(\lambda\) and shape k. At \(t=\lambda\) survival is \(e^{-1}\), about 0.368, whatever the shape. With \(\lambda=1000\) hours and \(k=2\), survival at 500 hours is \(\exp(-(0.5)^2) = \exp(-0.25)\), about 0.779, so more than three quarters of the units are still running at half the scale.

## The shape parameter

The shape decides which phase of life the hazard describes. For \(k=1\) the hazard is constant and the distribution is exponential. For \(k<1\) the hazard falls, which describes early failures that remove weak units. For \(k>2\) the hazard rises, the wear-out regime. At \(k=2\) the hazard grows in proportion to time, so at 2000 hours with \(\lambda=1000\) it is twice its value at 1000 hours.

## Reading a ratio

Dividing survival at two times removes the scale. The ratio \(S(2t)/S(t)\) equals \(\exp(-(2^k - 1)(t/\lambda)^k)\), so with \(k=2\) and \(t=\lambda\) the ratio is \(\exp(-3)\), about 0.0498. Survival therefore collapses across a single scale interval once t passes the scale, which is why a fitted shape well above 3 usually points at a wear-out mechanism rather than at random failure.

The two parameters answer different questions; a fitted shape is usually the more informative of the two because it names the failure regime.
