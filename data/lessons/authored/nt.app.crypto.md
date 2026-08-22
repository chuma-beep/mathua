# Public-key cryptography and RSA

How can two people who have never met exchange secrets in the open? The answer reshaped commerce: make the *encryption* method public and keep only the *undoing* of it private.

RSA is the classic scheme. Alice picks two large primes \\(p\\) and \\(q\\), publishes their product \\(N = pq\\) together with an exponent \\(e\\) chosen coprime to \\((p-1)(q-1)\\), and computes her private key \\(d\\) with \\(ed \equiv 1 \pmod{(p-1)(q-1)}\\). To send a message encoded as a number \\(m < N\\), Bob transmits

\[
c \equiv m^{e} \pmod{N},
\]

which anyone can do — \\(N\\) and \\(e\\) are public. Alice recovers \\(m\\) by computing \\(c^{d} \bmod N\\).

Why does decryption work? Euler's theorem (the generalization of Fermat's little theorem to composite moduli) guarantees \\(m^{ed} \equiv m \pmod N\\) whenever \\(\gcd(m,N)=1\\), because \\(ed = 1 + k(p-1)(q-1) = 1 + k\varphi(N)\\).

The security rests on an asymmetry of effort: multiplying \\(p\\cdot q\\) is instant, but factoring \\(N\\) back into its primes is believed intractable for thousand-bit \\(N\\). An eavesdropper knows everything except the factorization, and without it, computing \\(d\\) from \\((N,e)\\) is out of reach with known methods. Every HTTPS handshake you use borrows this idea.
