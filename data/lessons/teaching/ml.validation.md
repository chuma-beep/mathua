# Train/Validation Split and Cross-Validation

**Model selection needs three splits:** training (fit), validation (tune hyperparameters, pick model), and test (final unbiased estimate). Reusing test for tuning leaks information.

## Reliable Estimates

### Hold-Out Validation
Split data into train/val/test (e.g. 70/15/15). Fit on train, tune on val, report on test exactly once.

### k-Fold Cross-Validation
Split data into $k$ folds (e.g. $k=5$). Rotate: train on $k-1$ folds, validate on the remaining, average $k$ scores. More reliable than a single split, especially on small data.

## Example

Dataset of 1000 points, $k=5$: each fold has 200 points. Fold 1 validates on points 0–199, trains on 200–999, and so on. Average validation error across 5 folds estimates test error with lower variance than one 80/20 split.
