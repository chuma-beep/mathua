# Overfitting and Generalization

**Overfitting:** Training error far below test error — the model memorized noise instead of pattern. Well-fit means both low and close; underfit means both high. Only test error measures generalization.

## Worked: three polynomial fits

Take 20 noisy points from a cubic and fit polynomials of degree 2, 4, and 12:
1. Degree 2 underfits: training error about 1.8, test error about 1.9 — both high, the curve is too stiff.
2. Degree 12 overfits: training error near 0.05, test error about 3.4 — it threads every noisy point and swings wildly between them.
3. Degree 4 balances: training error about 0.9, test error about 1.0 — both low and close.

So capacity has a sweet spot: the test error falls while the model learns signal, then rises when it starts fitting noise.

## Worked: reading the capacity curve

Take test error plotted against model capacity:
1. On the left (tiny models) both errors are high — underfit, add capacity.
2. In the middle test error bottoms out while training error keeps falling — the gap opening past the minimum is overfitting.
3. The minimum marks the best capacity; everything right of it buys training fit with test accuracy.

So the U-shape of test error against capacity diagnoses the regime: falling means underfit, rising means overfit.

## When more data fixes it

Overfitting is a data-to-capacity mismatch, so more training samples shrink the gap: the degree-12 polynomial above generalizes once given 200 points instead of 20. When data cannot grow, shrink the model or regularize instead.
