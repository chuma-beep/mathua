# The Radon-Nikodym theorem

The Radon-Nikodym theorem gives a density for one measure with respect to another. It is the measure-theoretic form of a derivative.

## Absolute continuity

A measure \(\nu\) is absolutely continuous with respect to \(\mu\) when \(\mu(A) = 0\) implies \(\nu(A) = 0\), written \(\nu \ll \mu\). If \(\nu\) is given by a density \(f\), then \(\nu(A) = \int_A f \, d\mu\), and the implication is immediate.

## The density

Under absolute continuity and sigma-finiteness, a density exists. Take \(\nu([a,b]) = \int_a^b 3t^2\,dt\) on \([0,1]\). The density is \(f(t) = 3t^2\), and \(\nu([0,1]) = 1\). The density is unique up to agreement almost everywhere, so \(f\) and \(f + 1_{\{0\}}\) give the same measure.

## Conditional expectation

Conditional expectation is a Radon-Nikodym derivative. For a uniform variable on \([0,1]\) and the partition into \([0,1/2)\) and \([1/2,1]\), the conditional expectation is \(1/4\) on the first atom and \(3/4\) on the second. Each value is the average of the variable over its atom.

The theorem names the density that absolute continuity promises, and conditional expectation is its most used instance.
