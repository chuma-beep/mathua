# Cross-Entropy Loss

**Cross-entropy loss:** For true label $y$ in $\{0, 1\}$ and predicted probability $\hat{y}$, $\text{CE} = -y\log\hat{y} - (1-y)\log(1-\hat{y})$: the negative log-probability assigned to the truth. For one-hot $y$ over classes, $\text{CE} = -\sum_c y_c \log \hat{y}_c$.

## Worked: one confident prediction

Take true label 1, prediction 0.8:
1. Only the $y=1$ term survives: $\text{CE} = -\log 0.8$.
2. Evaluate: $-\log 0.8 \approx 0.22$.
3. Sharpen to prediction 0.95 and the cost falls to about 0.05 — reward for calibrated confidence.

So cross-entropy reads off surprise: likely truth costs little, unlikely truth costs much.

## Worked: confident and wrong

Take true label 1, prediction 0.1:
1. Score it: $\text{CE} = -\log 0.1 \approx 2.30$ — ten times the cost of predicting 0.8.
2. Push to 0.01 and the cost doubles again to about 4.61; as $\hat{y} \to 0$ the loss grows without bound.
3. Compare MSE on the same pair: $(0.1-1)^2 = 0.81$, bounded no matter how confident the error.

So cross-entropy punishes confident mistakes far harder than squared error, which is why classification trains on it: probabilities, not distances.

## When cross-entropy misleads

On noisy or mislabeled data the unbounded penalty chases bad labels: one flipped label with prediction 0.99 pays about 4.61 and drags the boundary. Label smoothing or noise-robust losses cap that pull — use cross-entropy where labels are trusted.
