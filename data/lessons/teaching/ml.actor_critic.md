# Actor-Critic Methods

**Actor-critic:** Actor $\pi_\theta(a|s)$ policy, critic $V_w(s)$ value baseline; advantage $A=Q-V$, GAE $\hat{A}$, A2C synchronous, DDPG for continuous actions.

## Advantage

### GAE
Generalized Advantage Estimation balances bias-variance.

### Continuous
DDPG actor-critic for continuous control.

## Example

A2C: parallel actors collect trajectories, update both actor and critic synchronously.
