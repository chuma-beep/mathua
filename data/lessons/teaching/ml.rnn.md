# Recurrent Neural Networks

**RNN:** For a sequence, keep a hidden state updated as $h_t = \tanh(W_{hh}h_{t-1} + W_{xh}x_t + b)$: the same weights $W_{hh}, W_{xh}$ reused every step. Memory persists because each state feeds the next.

## Worked: unrolling three steps

Take inputs $x_1, x_2, x_3$ with initial state $h_0 = 0$, $W_{hh} = 0.5$, $W_{xh} = 1$, $b = 0$:
1. $h_1 = \tanh(0.5\cdot 0 + 1\cdot x_1)$: with $x_1 = 1$, $h_1 = \tanh(1) \approx 0.76$.
2. $h_2 = \tanh(0.5\cdot 0.76 + x_2)$: with $x_2 = 0$, $h_2 = \tanh(0.38) \approx 0.36$ — the past lingers, halved.
3. $h_3 = \tanh(0.5\cdot 0.36 + x_3)$: with $x_3 = -1$, $h_3 = \tanh(-0.82) \approx -0.68$ — new input plus fading memory.

So unrolling turns recurrence into a deep feedforward net with tied weights: same function $f$ applied at every step, memory carried by $h_t$.

## Worked: gradients vanish

Take backpropagation through 10 steps where each Jacobian scales gradients by 0.5:
1. One step back multiplies by 0.5; two steps by $0.5^2 = 0.25$.
2. Ten steps give $0.5^{10} \approx 0.001$ — the first input's influence nearly gone.
3. Long dependencies therefore train glacially: the signal decays exponentially with distance.

So vanilla RNNs remember recent steps and forget distant ones — the vanishing gradient, quantified as repeated multiplication below 1.

## When gates are needed

LSTMs fix the decay with input, forget, and output gates that let gradients flow unscaled down a cell state. Use gated cells (LSTM, GRU) for long-range dependencies; plain recurrence suffices only for short histories.
