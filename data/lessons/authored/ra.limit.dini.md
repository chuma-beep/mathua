# Dini's theorem

Dini's theorem upgrades pointwise convergence to uniform convergence under a monotonicity hypothesis. It is the standard route from a pointwise limit to a uniform one on a compact interval.

## Monotone convergence of functions

Let \(f_n\) be continuous on a compact set and increase pointwise to a continuous \(f\). Then \(f_n \to f\) uniformly. On \([0, 1/2]\) the functions \(t^n\) increase to 0, and the supremum of \(t^n\) is \((1/2)^n\), which tends to 0.

## Compactness

The proof covers the compact set by finitely many intervals on which \(f - f_n\) is already small at one point. Compactness turns the pointwise control into a finite cover. At each stage the increase of \(f_n\) carries the smallness from that point to a whole neighbourhood.

## Why monotonicity is needed

Without monotonicity the conclusion fails. The functions \(g_n(t) = \max(0, 1 - |nt - 1|)\) are continuous. Each rises to height 1 near \(t = 1/n\) and tends to 0 at every point of \([0,1]\). The convergence is not uniform, because the supremum stays at 1. The limit here is continuous, so compactness alone does not help.

Dini's theorem trades a monotone hypothesis for a uniform conclusion, and the moving bump shows the trade is necessary.
