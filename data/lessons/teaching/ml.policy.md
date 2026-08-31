# Policy Gradients

**Policy gradient:** $\nabla_\theta J = \mathbb E[ G_t \nabla\log\pi_\theta(a|s)]$ (REINFORCE Monte Carlo); on-policy, baseline $b(s)$ reduces variance, advantage $A=Q-V$, GAE generalizes, PPO clipped surrogate, entropy bonus for exploration.

## Actor-Critic

### Components
Actor policy, critic value baseline; $A=Q-V$ advantage.

### Regularization
Entropy regularization encourages exploration.

## Example

CartPole: REINFORCE samples trajectories, updates $\theta$ proportional to return $G_t$ times score.
