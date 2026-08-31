# Transformers and Attention

**Attention:** For queries $Q$, keys $K$, values $V$,

$$\text{Attention}(Q,K,V)=\text{softmax}(QK^T/\sqrt{d_k})V$$

Multi-head runs $h$ attentions in parallel, concatenated.

## Architecture

### Self-Attention
Without positional encoding, self-attention is permutation equivariant; sinusoidal encodings $PE_{(pos,2i)}=\sin(pos/10000^{2i/d})$, $PE_{(pos,2i+1)}=\cos(...)$ add order.

### Complexity
Scaled dot-product is $O(n^2)$ in sequence length $n$ due to $n×n$ attention matrix; decoder uses causal mask to be autoregressive.

## Example

Sequence length $4$, $d_k=2$: $QK^T$ is $4×4$, softmax rows give attention weights summing to $1$ per query.
