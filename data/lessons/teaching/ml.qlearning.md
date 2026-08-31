# Q-Learning

**Q-learning:** Off-policy model-free, learns $Q(s,a)$ via $Q(s,a)←Q(s,a)+α[r+γ\max_{a'}Q(s',a')-Q(s,a)]$, converges with infinite exploration, model-free, overestimation bias via max.

## Optimality

### Bellman
$Q^*$ satisfies Bellman optimality; $Q$-learning approximates.

### Bias
Max causes overestimation; Double Q-learning mitigates.

## Example

Gridworld: $Q(s,a)$ estimates discounted return; policy $\pi(s)=\arg\max_a Q(s,a)$.
