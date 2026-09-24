# Actor-Critic Methods

**Actor-critic:** An actor holds the policy $\pi_\theta(a|s)$; a critic holds a value baseline $V_w(s)$. Updates follow the advantage $A = Q - V$: the actor climbs advantage, the critic fits values. A2C runs actors synchronously in parallel; DDPG extends the split to continuous actions.

## Worked: one TD error

Take reward $r = 1$, discount $\gamma = 0.9$, next value $V(s') = 10$, current value $V(s) = 8$:
1. Bootstrap the target: $1 + 0.9\cdot 10 = 10$.
2. Subtract the current estimate: $10 - 8 = 2$ — the TD error, a one-sample advantage.
3. The critic steps toward 10; the actor raises the taken action's probability in proportion to 2.

So the TD error does double duty: critic loss and actor weight — one number, two updates.

## Worked: advantage over raw return

Take $Q = 7$ for the taken action with $V = 5$ for the state:
1. Advantage: $7 - 5 = 2$ — the action beat its state's average by 2.
2. Generalized Advantage Estimation blends multi-step returns to trade bias against variance around this signal.
3. A2C averages such updates over parallel actors synchronously, smoothing the gradient without stale (off-policy) data.

So the critic centers learning: raw returns carry state luck, advantages carry action quality.

## When actions go continuous

Discrete policies output probabilities per action, but continuous control needs densities or deterministic maps: DDPG pairs a deterministic actor with a Q-critic and explores through noise. Use vanilla actor-critic for discrete actions, deterministic variants where actions are continuous.
