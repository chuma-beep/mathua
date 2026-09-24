# Q-Learning

**Q-learning:** Learn action values $Q(s, a)$ off-policy and model-free via the Bellman backup $Q(s,a) \leftarrow Q(s,a) + \alpha[r + \gamma\max_{a'}Q(s',a') - Q(s,a)]$. Act greedily in $Q$; learn from any sufficiently exploratory behavior.

## Worked: one Bellman update

Take $Q(s, a) = 0$, learning rate $\alpha = 0.5$, reward $r = 2$, discount $\gamma = 0.9$, best next value $\max Q(s', a') = 4$:
1. Form the target: $2 + 0.9\cdot 4 = 5.6$ — immediate reward plus discounted best future.
2. Compute the surprise: $5.6 - 0 = 5.6$.
3. Update halfway: $Q \leftarrow 0 + 0.5\cdot 5.6 = 2.8$.

So each update drags $Q$ toward a one-step-lookahead target: reward now plus the best the agent can do next.

## Worked: reading the greedy policy

Take a state with $Q$ values $[1, 5, 3]$ over three actions:
1. The greedy action is $\arg\max$: action 2 with value 5.
2. With $\epsilon = 0.1$ exploration, the agent exploits with probability 0.9 and tries a random action with probability 0.1.
3. As $Q$ converges to $Q^*$, the greedy policy converges to optimal — learning is off-policy (values assume greedy future) while behavior stays exploratory.

So policy and learning decouple: behave epsilon-greedy, evaluate as if fully greedy from the next step on.

## When the max bites back

The $\max$ in the target is optimistic: noisy overestimates get selected and propagate, causing overestimation bias. Double Q-learning separates selection from evaluation with two estimators — keep the single estimator only where noise is low.
