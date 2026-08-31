# Cross-Entropy Loss

**Cross-entropy loss:** For classification with true label $y\in\{0,1\}$ and predicted probability $\hat{y}$,

$$\text{CE}(y,\hat{y}) = -y\log\hat{y}-(1-y)\log(1-\hat{y})$$

For multi-class with true one-hot $y$ and predicted distribution $\hat{y}$,

$$\text{CE} = -\sum_c y_c \log \hat{y}_c$$

## When to Use Cross-Entropy

### Classification vs Regression
Cross-entropy measures distance between probability distributions, so it fits classification where outputs are probabilities. MSE measures Euclidean distance for continuous regression targets.

### Penalizing Confident Errors
If $y=1$ but $\hat{y}\to0$, $\text{CE}\to\infty$, while MSE stays bounded. Confident wrong predictions are strongly discouraged.

## Example

True label $1$, prediction $0.8$: $\text{CE}=-\log 0.8\approx0.22$. If prediction $0.1$: $\text{CE}=-\log0.1\approx2.30$. Poor confidence on the correct class costs far more.
