# Calibration and Reliability Diagrams

**Calibrated:** Predicted confidence matches empirical accuracy: among predictions made with confidence $p$, about fraction $p$ are correct. A reliability diagram bins predictions by confidence and plots accuracy against it — perfect calibration is the diagonal.

## Worked: one ECE computation

Take 200 predictions split into two bins of 100, with gaps $|acc - conf|$ of 0.2 and 0.0:
1. Weight each bin by its share: $100/200 = 0.5$ each.
2. Bin contributions: $0.5\cdot 0.2 = 0.10$ and $0.5\cdot 0.0 = 0$.
3. Expected Calibration Error $\text{ECE} = 0.10 + 0 = 0.10$.

So ECE is a weighted average of diagonal gaps: big confident-wrong bins dominate it.

## Worked: overconfident by twenty points

Take 100 predictions at 80 percent confidence with only 60 correct:
1. Accuracy 0.6 against confidence 0.8: the gap is 0.2 — overconfident by twenty points.
2. If this bin holds all the data, $\text{ECE} = 0.2$.
3. Temperature scaling divides logits $z/T$ and tunes $T$ on validation to close the gap — sharpening ($T < 1$) or softening ($T > 1$) without changing the predicted classes.

So miscalibration is fixable post-hoc: rescale confidence, keep decisions, remeasure ECE.

## When calibration misleads

A model predicting the base rate everywhere is perfectly calibrated and perfectly useless — calibration ignores sharpness (resolution between classes). Report ECE beside accuracy and AUC: calibration says whether to trust the numbers, not whether the model discriminates.
