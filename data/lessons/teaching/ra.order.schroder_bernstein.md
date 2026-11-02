# The Schroder-Bernstein Theorem

## Injections and surjections

Two sets have the same cardinality when a bijection exists between them. Schroder-Bernstein reduces that requirement to two one-way maps: an injection from A into B and an injection from B into A are enough to produce a bijection. The theorem is useful because one-way embeddings are often easy to write down.

For finite sets the statement is immediate from counting. The content lies in the infinite case, where an injection can miss elements and a surjection can repeat them. The proof must manufacture a bijection without assuming that either injection is surjective.

## The Cantor-Bernstein argument

Compose the two injections to obtain a map from A into itself. The image misses some part of A, and iterating the missed part produces a countable chain of subsets. Splitting A into the part that survives the iteration and the part that does not lets the proof use one injection on the first piece and the other on the second.

The construction is sometimes described with a back-and-forth bookkeeping rule. Each element is assigned a direction according to how many inverse steps it can take. Elements with an even history go one way and elements with an odd history go the other, and the resulting map is both injective and surjective.

## Consequences

The theorem shows that the closed interval and the whole real line have the same cardinality, since each embeds into the other. It also gives a convenient way to prove equipotence between sets of very different geometric shape without writing an explicit formula for a bijection.

Schroder-Bernstein is a statement about cardinality, not about order or measure. The interval and the line have equal cardinality even though the interval is bounded. Keeping those notions separate prevents the common mistake of treating a bijection as if it had to preserve distance.
