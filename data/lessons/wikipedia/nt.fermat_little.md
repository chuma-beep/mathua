> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Fermat%27s_little_theorem) — CC BY-SA 4.0

# Fermat's little theorem

In number theory, **Fermat's little theorem** states that if is a prime number, then for any integer , the number *a*\(^{*p*}\) − *a* is an integer multiple of. In the notation of modular arithmetic, this is expressed as

\[
a^p \equiv a \pmod p.
\]

For example, if *a* and *p* , then 2\(^{7}\) , and 128 − 2 7 × 18is an integer multiple of 7.

If is not divisible by ; that is, if is coprime to , then Fermat's little theorem is equivalent to the statement that *a*\(^{*p* − 1}\) − 1 is an integer multiple of , or in symbols:

\[
a^{p-1} \equiv 1 \pmod p.
\]

For example, if *a* and *p* , then 2\(^{6}\) , and 64 − 1 7 × 9is a multiple of 7.

Fermat's little theorem is the basis for the Fermat primality test and is one of the fundamental results of elementary number theory. The theorem is named after Pierre de Fermat, who stated it in 1640. It is called the "little theorem" to distinguish it from Fermat's Last Theorem.

## History

Pierre de Fermat first stated the theorem in a letter dated October 18, 1640, to his friend and confidant Frénicle de Bessy. His formulation is equivalent to the following:

If is a prime and is any integer not divisible by , then *a*\(^{ *p* − 1}\) − 1 is divisible by. Fermat's original statement was

This may be translated, with explanations and formulas added in brackets for easier understanding, as:

Every prime number [] divides necessarily one of the powers minus one of any [geometric] progression [*a*, *a*\(^{2}\), *a*\(^{3}\), …] [that is, there exists such that divides *a\(^{t}\)* − 1], and the exponent of this power [] divides the given prime minus one [divides *p* − 1]. After one has found the first power [] that satisfies the question, all those whose exponents are multiples of the exponent of the first one satisfy similarly the question [that is, all multiples of the first have the same property].

Fermat did not consider the case where is a multiple of nor prove his assertion, only stating:

(And this proposition is generally true for all series [*sic*] and for all prime numbers; I would send you a demonstration of it, if I did not fear going on for too long.)

Euler provided the first published proof in 1736, in a paper titled "Theorematum Quorundam ad Numeros Primos Spectantium Demonstratio" (in English: "Demonstration of Certain Theorems Concerning Prime Numbers") in the *Proceedings* of the St. Petersburg Academy, but Leibniz had given virtually the same proof in an unpublished manuscript from sometime before 1683.

The term "Fermat's little theorem" was probably first used in print in 1913 in *Zahlentheorie* by Kurt Hensel:

(There is a fundamental theorem holding in every finite group, usually called Fermat's little theorem because Fermat was the first to have proved a very special part of it.)

An early use in English occurs in A.A. Albert's *Modern Higher Algebra* (1937), which refers to "the so-called 'little' Fermat theorem" on page 206.

### Further history

Some mathematicians independently made the related hypothesis (sometimes incorrectly called the Chinese hypothesis) that 2\(^{*p*}\) ≡ 2 (mod *p*) if and only if is prime. Indeed, the "if" part is true, and it is a special case of Fermat's little theorem. However, the "only if" part is false: For example, 2\(^{341}\) ≡ 2 (mod 341), but 341 = 11 × 31 is a pseudoprime to base 2. See below.

## Proofs

Several proofs of Fermat's little theorem are known. It is frequently proved as a corollary of Euler's theorem.

## Generalizations
Euler's theorem is a generalization of Fermat's little theorem: For any modulus and any integer coprime to , one has

\[
a^{\varphi (n)} \equiv 1 \pmod n,
\]

where *φ*(*n*) denotes Euler's totient function (which counts the integers from 1 to that are coprime to ). Fermat's little theorem is indeed a special case, because if is a prime number, then *φ*(*n*) = *n* − 1.

A corollary of Euler's theorem is: For every positive integer , if the integer is coprime with , then

\[
x \equiv y \pmod{\varphi(n)}\quad\text{implies}\quad a^x \equiv a^y \pmod n,
\]

for any integers and. This follows from Euler's theorem; since, if \(x \equiv y \pmod{\varphi(n)}\), then *x* = *y* + *kφ*(*n*) for some integer , and one has

\[
a^x = a^{y + \varphi(n)k} = a^y (a^{\varphi(n)})^k \equiv a^y 1^k \equiv a^y \pmod n.
\]

If is prime, this is also a corollary of Fermat's little theorem. This is widely used in modular arithmetic, because this allows reducing modular exponentiation with large exponents to exponents smaller than. Euler's theorem is used with not prime in public-key cryptography, specifically in the RSA cryptosystem, typically in the following way: if

\[
y=x^e\pmod n,
\]

retrieving from the values of , and is easy if one knows *φ*(*n*). In fact, the extended Euclidean algorithm allows computing the modular inverse of modulo *φ*(*n*), that is, the integer such that

\[
ef\equiv 1\pmod{\varphi(n)}.
\]

It follows that

\[
x\equiv x^{ef}\equiv (x^e)^f \equiv y^f \pmod n.
\]

On the other hand, if *n* = *pq* is the product of two distinct prime numbers, then *φ*(*n*) = (*p* − 1)(*q* − 1). In this case, finding from and is as difficult as computing *φ*(*n*) (this has not been proven, but no algorithm is known for computing without knowing *φ*(*n*)). Knowing only , the computation of *φ*(*n*) has essentially the same difficulty as the factorization of ; since *φ*(*n*) = (*p* − 1)(*q* − 1), and conversely, the factors and are the (integer) solutions of the equation *x* 0.

The basic idea of RSA cryptosystem is thus: If a message is encrypted as *y* = *x\(^{e}\)* (mod *n*), using public values of and , then, with the current knowledge, it cannot be decrypted without finding the (secret) factors and of. Fermat's little theorem is also related to the Carmichael function and Carmichael's theorem, as well as to Lagrange's theorem in group theory.

## Converse
The converse of Fermat's little theorem fails for Carmichael numbers. However, a slightly weaker variant of the converse is **Lehmer's theorem**:

If there exists an integer such that

\[
a^{p-1}\equiv 1\pmod p
\]

and for all primes dividing *p* − 1 one has

\[
a^{(p-1)/q}\not\equiv 1\pmod p,
\]

then is prime.

This theorem forms the basis for the Lucas primality test, an important primality test, and Pratt's primality certificate.
## Pseudoprimes

If and are coprime numbers such that *a*\(^{*p*−1}\) − 1 is divisible by , then need not be prime. If it is not, then is called a *(Fermat) pseudoprime* to base. The first pseudoprime to base 2 was found in 1820 by Pierre Frédéric Sarrus: 341 = 11 × 31.

A number that is a Fermat pseudoprime to base for every number coprime to is called a Carmichael number. Alternately, any number satisfying the equality

\[
\gcd\left(p, \sum_{a=1}^{p-1} a^{p-1}\right)=1
\]

is either a prime or a Carmichael number.

## Miller–Rabin primality test
The Miller–Rabin primality test uses the following extension of Fermat's little theorem:
If is an odd prime and *p* − 1 = 2\(^{*s*}\)*d* with s > 0 and odd > 0, then for every coprime to , either *a*\(^{*d*}\) ≡ 1 (mod *p*) or there exists such that 0 ≤ *r* < *s* and *a*\(^{2*r*}\)*d* ≡ −1 (mod *p*).

This result may be deduced from Fermat's little theorem by the fact that, if is an odd prime, then the integers modulo form a finite field, in which 1 modulo has exactly two square roots, 1 and −1 modulo. Note that *a*\(^{*d*}\) ≡ 1 (mod *p*) holds trivially for *a* ≡ 1 (mod *p*), because the congruence relation is compatible with exponentiation. And *a*\(^{*d*}\) = *a*\(^{20}\)*d* ≡ −1 (mod *p*) holds trivially for *a* ≡ −1 (mod *p*) since is odd, for the same reason. That is why one usually chooses a random in the interval 1 < *a* < *p* − 1.

The Miller–Rabin test uses this property in the following way: given an odd integer for which primality has to be tested, write *p* − 1 = 2\(^{*s*}\)*d* with s > 0 and odd > 0, and choose a random such that 1 < *a* < *p* − 1; then compute *b* = *a*\(^{*d*}\) mod *p*; if is not 1 nor −1, then square it repeatedly modulo until you get −1 or have squared *s* − 1 times. If *b* ≠ 1 and −1 has not been obtained by squaring, then is a *composite* and is a witness for the compositeness of. Otherwise, is a *strong probable prime to base a*; that is, it may be prime or not. If is composite, the probability that the test declares it a strong probable prime anyway is at most 1/4, in which case is a *strong pseudoprime*, and is a *strong liar*. Therefore after non-conclusive random tests, the probability that is composite is at most 4\(^{−*k*}\), and may thus be made as low as desired by increasing. In summary, the test either proves that a number is composite or asserts that it is prime with a probability of error that may be chosen as low as desired. The test is very simple to implement and computationally more efficient than all known deterministic tests. Therefore, it is generally used before starting a proof of primality.
