# Huffman coding

Huffman coding builds an optimal prefix code by repeatedly merging the two least frequent symbols. The merged node's frequency is their sum, and the final tree is read as binary code words.

## The greedy merge

Start with one leaf per symbol weighted by frequency. Merge the two smallest weights, insert the parent, and repeat until one root remains.

## Prefix codes

Leaves carry symbols, so no code word is a prefix of another. A stream can therefore be decoded without separators.

## Optimality

An exchange argument shows that the two least frequent symbols may be placed at maximum depth, which justifies each greedy merge.

To check your understanding, take two smallest weights, place symbols at leaves and exchange the deepest pair.

Huffman's tree is optimal among prefix codes for the given frequencies. The algorithm is a clean example of a greedy choice proved by exchange.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
