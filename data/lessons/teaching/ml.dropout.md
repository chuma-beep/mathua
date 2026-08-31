# Dropout Regularization

**Dropout:** During training, each activation is zeroed with probability $p$ (e.g. $0.5$) and the survivors are scaled by $1/(1-p)$ (inverted dropout). At test time dropout is disabled.

## Why It Helps

### Reducing Co-Adaptation
Randomly removing neurons forces the network to not rely on any single feature; it approximates training an ensemble of $2^{n}$ subnetworks.

### Regularization Effect
Dropout adds noise, discouraging fitting to training noise, improving generalization. Expected active neurons out of $100$ with $p=0.5$ is $50$.

## Example

Network with $100$ hidden units, $p=0.5$: each forward pass uses ~$50$ random units; test uses all $100$ (or training scales by $2×$). This averages subnetwork predictions.
