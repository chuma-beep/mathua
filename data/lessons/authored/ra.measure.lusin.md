# Lusin's theorem

Lusin's theorem says a measurable function is continuous once a small set is removed. Measurability is almost continuity.

## Approximation by continuous functions

On a finite measure space, for every \(\varepsilon > 0\) there is a closed set \(F\) with \(\mu(E \setminus F) < \varepsilon\) such that the restriction of \(f\) to \(F\) is continuous. The function need not be continuous on the whole space.

## Closed sets

Take the indicator of the rationals on \([0,1]\). The irrationals contain a closed set \(F\) of measure 0.99, because the irrationals are a countable intersection of open sets. On \(F\) the indicator is identically 0, which is continuous. The discarded set has measure 0.01.

## Consequences

The theorem gives density of continuous functions in \(L^1\). Given an integrable \(f\) with \(|f| \le 1\) and \(\varepsilon = 0.01\), a continuous \(g\) with \(|g| \le 1\) agrees with \(f\) off a set of measure 0.01, so the integral of \(|f - g|\) is at most 0.02. Tietze's theorem then extends \(g\) from the closed set to the whole space.

Measurability is continuity up to a small set, and that reading is what lets continuous functions stand in for integrable ones.
