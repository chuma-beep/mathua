# Word embeddings

## Distributed representations

A word embedding maps each token to a dense vector so that similar words have nearby vectors. The vectors are learned from co-occurrence statistics or as part of training a larger language model.

A vocabulary of V tokens with embedding dimension d requires a table of V times d numbers. The table is usually much smaller than a one-hot representation of the same vocabulary.

## Context windows

Skip-gram training predicts context words from a centre word. With window size w, each token is paired with up to w tokens on either side, so a sequence of n tokens yields about 2 times w times n training pairs before edge effects.

Negative sampling replaces the full softmax over the vocabulary with a small set of positive and negative pairs, which makes training practical for large vocabularies.

## Using embeddings

Cosine similarity compares the directions of two vectors and ignores their lengths. It is computed as the dot product divided by the product of the norms.

When two pretrained embedding spaces are concatenated, the resulting dimension is the sum of the two dimensions. Downstream layers must be sized for that combined width.
