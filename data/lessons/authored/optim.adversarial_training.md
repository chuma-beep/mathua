# Adversarial training

Adversarial training minimises the loss against the worst input perturbation inside a small set. The inner problem is an attack and the outer problem is training, so the method is a minimax problem in the input space.

## The inner attack

A few projected gradient steps find a perturbation that maximises the loss while staying inside the allowed ball. The number of steps trades attack strength against training time.

## The outer objective

The training gradient is taken at the perturbed input rather than at the clean input. The model therefore learns a decision boundary that has margin in the perturbed directions.

## The trade

Robust accuracy on perturbed inputs usually costs clean accuracy. The size of the perturbation set determines both the guarantee and the price, and neither is free.

To check your understanding, allow a small ball, train at the perturbed input and measure clean accuracy.

Adversarial training is a minimax problem whose inner step is a local search, and the perturbation budget is the design variable that sets the trade.
