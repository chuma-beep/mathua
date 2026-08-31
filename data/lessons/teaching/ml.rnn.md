# Recurrent Neural Networks

**RNN:** Processes sequences by maintaining a hidden state $h_t$ updated as $h_t = \tanh(W_{hh}h_{t-1}+W_{xh}x_t+b)$, sharing $W_{hh},W_{xh}$ across time steps.

## Memory Over Time

### Shared Weights and Unrolling
Unrolling $T$ steps yields a deep feedforward network with tied weights, handling variable-length sequences.

### Vanishing Gradients
Backprop through time multiplies many Jacobians; gradients vanish, hindering long dependencies. LSTMs use input/forget/output gates to preserve gradients.

## Example

Sequence $x_1,x_2,x_3$: $h_1=f(x_1,h_0)$, $h_2=f(x_2,h_1)$, $h_3=f(x_3,h_2)$ — same $f$ reused, memory persists via $h_t$.
