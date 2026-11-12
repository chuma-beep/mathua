# Shannon entropy

Shannon entropy measures the average surprise of a source. It is the expected negative logarithm of the symbol probability and bounds the best possible prefix-code length.

## The definition

For probabilities \(p_i\), entropy is \(H=-\sum p_i\log_2 p_i\) bits. A fair coin has one bit and a certain symbol has none.

## Units

Base-two logarithms measure bits, natural logarithms measure nats, and the two differ by the constant \(\ln 2\).

## The coding bound

Any prefix code has average length at least \(H\), and Huffman coding comes within one bit of it.

To check your understanding, weight by probability, choose the logarithm base and compare average length.

Entropy is a property of the source, not of a particular code. It sets the floor that every lossless code must respect.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
