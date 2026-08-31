# Calibration and Reliability Diagrams

**Calibrated:** $P(y=1 | \hat{p}=p)≈p$; confidence matches accuracy. **Reliability diagram:** bins predictions by confidence, plots accuracy vs confidence; diagonal is perfect.

## Measuring Calibration

### ECE
Expected Calibration Error $\text{ECE}=∑_b (|B_b|/N) |acc(B_b)-conf(B_b)|$. Overconfident models have high ECE.

### Fixing Miscalibration
Temperature scaling: softmax with temperature $T$ scales logits $z/T$, tuning $T$ on validation to minimize ECE without changing accuracy.

## Example

Model predicts $100$ samples at $80\%$ confidence: if $80$ correct, calibrated; if $60$ correct, overconfident by $20\%$ ECE contribution $0.2×(100/N)$.
