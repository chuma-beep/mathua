# Train, Validation, and Test Splits

**Three splits:** Train to fit parameters, validate to tune hyperparameters and pick models, test to report one final unbiased number. Touching test during tuning leaks information and the final number lies.

## Worked: one 5-fold split

Take a dataset of 1000 points with $k = 5$:
1. Cut into 5 folds of 200 points each: fold 1 holds points 0–199, fold 2 holds 200–399, and so on.
2. Round 1 trains on folds 2–5 (800 points) and validates on fold 1 (200 points).
3. Rotate the held-out fold 5 times and average the 5 validation errors — that average estimates test error with lower variance than any single split.

So cross-validation buys reliability by averaging: every point validates exactly once, and no single lucky split decides.

## Worked: the cost of peeking

Take a 70/15/15 train/val/test split and a search over 50 hyperparameter settings:
1. Tuning honestly on validation picks the setting with validation error 0.20; test then reports about 0.21 — close, unbiased.
2. Tuning directly on test picks the setting with test error 0.17 — but that number was optimized, so fresh data sees about 0.22.
3. The 0.05 gap is pure leakage: the search overfit the test set itself.

So report on test exactly once, after all choices are frozen — reuse turns the unbiased estimate into another training metric.

## When a single split suffices

With very large data a single validation split is already stable: at a million points the split noise is negligible and 5-fold costs 5 full trainings for nothing. Cross-validate on small data; hold out once on big data.
