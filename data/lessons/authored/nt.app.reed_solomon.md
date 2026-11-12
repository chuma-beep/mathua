# Reed-Solomon error correction

Reed-Solomon codes evaluate polynomials at many points. Redundancy lets the receiver correct errors rather than merely detect them.

## The code

A message of k symbols becomes a polynomial of degree less than k. The codeword evaluates it at n distinct field elements, so n must not exceed the field size.

## The redundancy

The difference between n and k is the number of parity symbols. More redundancy means more errors can be corrected or more erasures recovered.

## Error correction

A code with n minus k parity symbols corrects half that many errors. The bound comes from needing two errors to move a word to another codeword.

To check your understanding, choose the field size, subtract dimension and halve the redundancy.

Reed-Solomon coding is interpolation with slack, and the slack is exactly the distance between codewords.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
