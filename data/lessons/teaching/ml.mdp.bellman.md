# The Bellman equation

## Value of a state

The value of a state under a policy is the expected discounted return starting from that state. The Bellman equation writes this value recursively as the immediate reward plus the discounted value of the next state.

For a deterministic transition with discount factor gamma, the relation is V(s) equals r plus gamma times V(s prime). Repeated substitution expands the value into the discounted sum of future rewards.

## Optimality

The Bellman optimality equation replaces the action taken by the policy with the maximum over actions. The optimal value of a state is the largest immediate reward plus discounted next value available from that state.

This equation is a fixed point condition. Value iteration applies the right-hand side repeatedly and converges to the optimal value when the discount factor is below one.

## State-action values

The action-value form records the value of taking a particular action in a state and then following the optimal policy. It is the quantity estimated by tabular Q-learning.

A tabular representation stores one number for every state-action pair, so a problem with S states and A actions needs S times A entries.
