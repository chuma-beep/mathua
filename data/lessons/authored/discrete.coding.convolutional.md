# Convolutional codes

A convolutional code emits output bits from a sliding window of recent input bits. The encoder is a shift register, and the decoder searches a trellis for the most likely path.

## The shift register

Each input bit shifts into a register. Output bits are linear combinations of the register contents, so the code is linear but not block-aligned.

## The trellis

States are register contents. Each input bit labels a transition, and a received sequence is a path through the trellis.

## Viterbi decoding

Dynamic programming keeps the best path into each state at every time step, so the search cost grows linearly with the sequence length.

To check your understanding, shift in one bit, label states by register and keep the best path per state.

Convolutional codes spread redundancy across time rather than within a block. Viterbi decoding turns maximum-likelihood decoding into a shortest-path problem on a trellis.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
