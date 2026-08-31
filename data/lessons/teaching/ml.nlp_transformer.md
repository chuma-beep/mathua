# Transformers for NLP

**Transformer:** Self-attention $ \text{softmax}(QK^T/\sqrt d)V$ with multi-head, positional encoding, $O(n^2)$ complexity, causal mask for decoder autoregressive.

## NLP Use

### Encoding
Without positional encoding, self-attention permutation equivariant; sinusoidal encodings add order.

### Heads
Multi-head captures different relations.

## Example

BERT base: $12$ layers, $768$ hidden, $12$ heads; attention matrix $n×n$.
