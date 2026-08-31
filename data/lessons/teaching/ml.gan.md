# Generative Adversarial Networks

**GAN:** Generator $G(z)$ maps noise $z$ to data space; discriminator $D(x)$ predicts real vs fake. Training is minimax:

$$\min_G \max_D \mathbb{E}_{x\sim p_{data}}[\log D(x)] + \mathbb{E}_{z}[\log(1-D(G(z)))]$$

## Adversarial Training

### Equilibrium
Nash equilibrium where $G$ matches $p_{data}$ and $D=0.5$ everywhere; $D$ provides gradient for $G$ via backprop.

### Failure Modes
Mode collapse: $G$ outputs few modes; Wasserstein GAN replaces JS with Earth mover distance for stability.

## Example

$G$ generates $28×28$ digits, $D$ classifies real MNIST vs fake; after training $G$'s samples look like digits.
