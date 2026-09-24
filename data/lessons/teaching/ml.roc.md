# ROC Curves and Evaluation Metrics

**ROC curve:** For a classifier, plot true positive rate $\text{TPR} = TP/(TP+FN)$ against false positive rate $\text{FPR} = FP/(FP+TN)$ while sweeping the decision threshold from strict to lax. Area under the curve (AUC) summarizes ranking quality in one number.

## Worked: one threshold point

Take a classifier scoring 10 positives and 90 negatives, threshold set so 8 positives and 10 negatives pass:
1. True positives 8, false negatives 2: $\text{TPR} = 8/(8+2) = 0.8$.
2. False positives 10, true negatives 80: $\text{FPR} = 10/(10+80) \approx 0.11$.
3. That threshold contributes the point (0.11, 0.8) to the ROC curve; sweeping all thresholds traces the full curve from (0, 0) to (1, 1).

So each threshold is one trade-off point: stricter cuts both rates, laxer raises both.

## Worked: reading AUC

Take the finished curve with AUC 0.85:
1. AUC equals $P(\text{score(pos)} > \text{score(neg)})$: a random positive outranks a random negative with probability 0.85.
2. AUC 1 means perfect separation — some threshold splits the classes cleanly.
3. AUC 0.5 is the diagonal: no better than random guessing; below 0.5 is worse than random (flip the prediction).

So AUC reads as ranking probability, not accuracy: it ignores threshold choice and rewards ordering positives above negatives.

## When ROC flatters

With 990 negatives and 10 positives, 99 false positives still give $\text{FPR} = 0.1$ — the curve looks fine while precision is 10/109, under 10 percent. On imbalanced data read precision-recall curves instead: they expose the false-positive flood that ROC dilutes.
