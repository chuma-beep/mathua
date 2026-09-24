# Policy Gradients

**Policy gradient:** For a parameterized policy $\pi_\theta$, ascend $\nabla_\theta J = \mathbb{E}[G_t \nabla\log\pi_\theta(a|s)]$ (REINFORCE): weight each action's score gradient by its Monte Carlo return $G_t$. On-policy, unbiased, and high-variance without a baseline.

## Worked: one REINFORCE step

Take learning rate $\alpha = 0.1$, return $G_t = 4$, score gradient $\nabla\log\pi = 0.5$:
1. Scale return by gradient: $4\cdot 0.5 = 2$ — good outcomes push the action up.
2. Scale by the rate: $0.1\cdot 2 = 0.2$.
3. Update: $\theta \leftarrow \theta + 0.2$, raising the probability of an action that returned 4.

So REINFORCE is reward-weighted imitation of the agent's own behavior: repeat what worked, in proportion to how well it worked.

## Worked: subtracting a baseline

Take state value $V = 5$ with action values $Q = 7$ and $Q = 3$:
1. Advantages: $7 - 5 = 2$ for the good action, $3 - 5 = -2$ for the bad one.
2. The update now pushes actions relative to average — above-average up, below-average down — instead of pushing everything up by its raw return.
3. The baseline adds no bias (it does not depend on the action) but cuts variance: learning no longer swings with the state's intrinsic goodness.

So advantage $A = Q - V$ centers the signal: the baseline absorbs luck of the state, the residual judges the action.

## When gradients need structure

Raw REINFORCE on long episodes drowns in variance; actor-critic methods learn the baseline (critic) alongside the policy (actor), and PPO clips the surrogate objective to stay near the current policy. Add entropy bonuses to keep exploration alive — pure exploitation collapses early.
