# Transformers for NLP

**Transformer for text:** Self-attention over token embeddings plus positional encodings, stacked in layers with multiple heads. Order comes only from the encodings — the attention itself is permutation equivariant.

## Worked: why order needs encoding

Take the sentences "dog bites man" and "man bites dog" with identical token embeddings:
1. Without positional encoding, self-attention sees the same multiset of tokens in both — outputs reorder with inputs, so the two sentences get the same representation up to order.
2. Add sinusoidal encodings: position 1, 2, 3 markers differ, so "dog" as subject versus object enters different queries.
3. The model now separates agent from patient — word order becomes visible to attention.

So encodings carry all sequential information: remove them and the transformer is a bag-of-words model.

## Worked: BERT base dimensions

Take BERT base: 12 layers, hidden width 768, 12 heads:
1. Per-head dimension: $768/12 = 64$ — each head attends in a 64-dimensional subspace.
2. One layer's attention matrix on 512 tokens is $512\times 512$: about 262000 weights per head per layer.
3. Twelve heads times 12 layers re-derive relations (syntax, coreference, agreement) in parallel subspaces.

So width splits across heads, depth stacks relations: dimensions budget expressiveness, heads divide the labor.

## When the decoder masks

Generation must not see the future: the decoder applies a causal mask zeroing attention to later positions, making prediction autoregressive. Encoders (BERT-style) skip the mask and read both directions — bidirectional understanding, but no generation.
