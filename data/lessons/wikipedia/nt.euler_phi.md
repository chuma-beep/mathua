> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Euler%27s_totient_function) — CC BY-SA 4.0

# Euler's totient function

In number theory, **Euler's totient function** counts the positive integers up to a given integer \(n\) that are relatively prime to \(n\). It is written using the Greek letter phi as \(\varphi(n)\) or \(\phi(n)\), and may also be called **Euler's phi function**. In other words, it is the number of integers \(k\) in the range \(1\leq k\leq n\) for which the greatest common divisor \(\gcd(n,k)\) is equal to 1. The integers \(k\) of this form are sometimes referred to as totatives of \(n\).

For example, the totatives of \(n=9\) are the six numbers 1, 2, 4, 5, 7 and 8. They are all relatively prime to 9, but the other three numbers in this range, 3, 6, and 9 are not; since \(\gcd(9,3)=\gcd(9,6)=3\) and \(\gcd(9,9)=9\). Therefore, \(\varphi(9)=6\). As another example, \(\varphi(1)=1\) since for \(n=1\) the only integer in the range from 1 to \(n\) is 1 itself, and \(\gcd(1,1)=1\).

Euler's totient function is a multiplicative function, meaning that if two numbers \(m\) and \(n\) are relatively prime, then \(\varphi(mn)=\varphi(m)\varphi(n)\).
This function gives the order of the multiplicative group of integers modulo (the group of units of the ring \(\Z/n\Z\)). It is also used for defining the RSA encryption system.

## History, terminology, and notation
Leonhard Euler introduced the function in 1763. However, he did not at that time choose any specific symbol to denote it. In a 1784 publication, Euler studied the function further, choosing the Greek letter \(\pi\) to denote it: he wrote \(\pi D\) for "the multitude of numbers less than \(D\), and which have no common divisor with it". This definition varies from the current definition for the totient function at \(D=1\) but is otherwise the same. The now-standard notation \(\varphi(A)\) comes from Gauss's 1801 treatise *Disquisitiones Arithmeticae*, although Gauss did not use parentheses around the argument and wrote \(\varphi A\). Thus, it is often called **Euler's phi function** or simply the **phi function**.

In 1879, J. J. Sylvester coined the term **totient** for this function, so it is also referred to as **Euler's totient function**, the **Euler totient**, or **Euler's totient**. Jordan's totient is a generalization of Euler's.

The **cototient** of \(n\) is defined as \(n-\varphi(n)\). It counts the number of positive integers less than or equal to \(n\) that have at least one prime factor in common with \(n\).

## Computing Euler's totient function
There are several formulae for computing \(\varphi(n)\).

### Euler's product formula
It states
\(\varphi(n) =n \prod_{p\mid n} \left(1-\frac{1}{p}\right),\)
where the product is over the distinct prime numbers dividing. An equivalent formulation is
\(\varphi(n) = p_1^{k_1-1}(p_1{-}1)\,p_2^{k_2-1}(p_2{-}1)\cdots p_r^{k_r-1}(p_r{-}1),\)
where \(n = p_1^{k_1} p_2^{k_2} \cdots p_r^{k_r}\) is the prime factorization of \(n\) (that is, \(p_1, p_2,\ldots, p_r\) are distinct prime numbers).

The proof of these formulae depends on two important facts.

#### Phi is a multiplicative function
This means that if \(\gcd(m,n) = 1\), then \(\varphi(m) \varphi(n) = \varphi(mn)\). *Proof outline:* Let \(A,B,C\) be the sets of positive integers which are coprime to and less than , respectively, so that \(|A| = \varphi(m)\), etc. Then there is a bijection between \(A\times B\) and by the Chinese remainder theorem.

#### Value of phi for a prime power argument
If is prime and \(k\geq1\), then

\(\varphi \left(p^k\right) = p^k-p^{k-1} = p^{k-1}(p-1) = p^k \left( 1 - \tfrac{1}{p} \right).\)

*Proof*: Since is a prime number, the only possible values of \(\gcd(p^{k},m)\) are \(1,p,p^{2},\dots,p^{k}\), and the only way to have \(\gcd(p^{k},m)>1\) is if is a multiple of , that is, \(m \in \{ p, 2p, 3p, \ldots, p^{k-1} p= p^{k}\}\), and there are \(p^{k-1}\) such multiples not greater than \(p^{k}\). Therefore, the other \(p^{k}-p^{k-1}\) numbers are all relatively prime to \(p^{k}\).

#### Proof of Euler's product formula
The fundamental theorem of arithmetic states that if *n* > 1 there is a unique expression \(n = p_1^{k_1} p_2^{k_2} \cdots p_r^{k_r},\) where *p*\(_{1}\) < *p*\(_{2}\) <.. . < *p*\(_{*r*}\) are prime numbers and each *k*\(_{*i*}\) ≥ 1. (The case 1=*n* = 1 corresponds to the empty product.) Repeatedly using the multiplicative property of and the formula for *φ*(*p*\(^{*k*}\)) gives

\(\begin{array} {rcl}
\varphi(n)&=& \varphi(p_1^{k_1})\, \varphi(p_2^{k_2})
\cdots\varphi(p_r^{k_r})\\[.1em]
&=& p_1^{k_1} \left(1- \frac{1}{p_1} \right) p_2^{k_2} \left(1- \frac{1}{p_2} \right) \cdots p_r^{k_r}\left(1- \frac{1}{p_r} \right)\\[.1em]
&=& p_1^{k_1} p_2^{k_2} \cdots p_r^{k_r} \left(1- \frac{1}{p_1} \right) \left(1- \frac{1}{p_2} \right) \cdots \left(1- \frac{1}{p_r} \right)\\[.1em]
&=&n \left(1- \frac{1}{p_1} \right)\left(1- \frac{1}{p_2} \right) \cdots\left(1- \frac{1}{p_r} \right).
\end{array}\)

This gives both versions of Euler's product formula.

An alternative proof that does not require the multiplicative property instead uses the inclusion-exclusion principle applied to the set \(\{1,2,\ldots,n\}\), excluding the sets of integers divisible by the prime divisors.

#### Example
\(\varphi(20)=\varphi(2^2 5)=20\,(1-\tfrac12)\,(1-\tfrac15)
=20\cdot\tfrac12\cdot\tfrac45=8.\)

In words: the distinct prime factors of 20 are 2 and 5; half of the twenty integers from 1 to 20 are divisible by 2, leaving ten; a fifth of those are divisible by 5, leaving eight numbers coprime to 20; these are: 1, 3, 7, 9, 11, 13, 17, 19.

The alternative formula uses only integers:
\[
\varphi(20) = \varphi(2^2 5^1)= 2^{2-1}(2{-}1)\,5^{1-1}(5{-}1) = 2\cdot 1\cdot 1\cdot 4 = 8.
\]

### Fourier transform
The totient is the discrete Fourier transform of the gcd, evaluated at 1. Let

\(\mathcal{F} \{ \mathbf{x} \}[m] = \sum\limits_{k=1}^n x_k \cdot e^*q*
