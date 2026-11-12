# ElGamal encryption

ElGamal encryption randomises a Diffie-Hellman exchange. The sender creates a fresh ephemeral exponent, and the receiver's private key removes it.

## The public key

The receiver publishes a prime, a generator and the generator raised to the private exponent. Recovering the exponent is the discrete logarithm problem.

## The ciphertext

The sender picks an ephemeral exponent k and sends the generator raised to k. The second component multiplies the message by the shared secret.

## Decryption

The receiver raises the first component to the private exponent. Dividing it out of the second component leaves the message.

To check your understanding, choose a generator, draw an ephemeral exponent and raise the first part.

ElGamal is Diffie-Hellman used as a lock: the ephemeral exponent makes each ciphertext different even for the same message.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
