# Recurrence and Return Times

This lesson develops the definition, the main calculation and the interpretation of the idea.

## Returning to a state

A state is recurrent when a chain that starts there returns with probability one. It is transient when the return probability is less than one. Recurrence is a property of the whole chain, not of a single path, so it is computed from the transition structure.

For an irreducible finite chain, every state is recurrent. Infinite chains can have transient states even when every state communicates with every other, as the simple random walk on the integers shows.

The first return time is the number of steps until the chain comes back. Its expectation is finite for a positive recurrent state and infinite for a null recurrent state.

## Stationary probabilities and return times

For an irreducible positive recurrent chain with stationary distribution $\pi$, the expected return time to state $i$ is

$$m_i=\frac{1}{\pi_i}.$$

This is a mean over returns, not a guarantee about the next return. A state with stationary probability one quarter is visited, on average, once every four steps.

The identity connects the long-run fraction of time spent in a state with the waiting time between visits. It is the reason a stationary distribution can be estimated by counting visits in a long simulation.

## Recurrence in practice

To decide recurrence from a transition matrix, look at the return probabilities and their sum. A state is recurrent exactly when the expected number of visits is infinite.

For the simple symmetric random walk on the integers, the return probability is one, but the expected return time is infinite. The walk returns eventually, yet the average wait is unbounded. That distinction separates recurrence from positive recurrence and matters when a model is used to predict waiting times.

