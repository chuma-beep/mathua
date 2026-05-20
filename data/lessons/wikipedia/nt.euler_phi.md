> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Euler%27s_totient_function) — CC BY-SA 4.0

# Euler's totient function

}}


In number theory, **Euler's totient function** counts the positive integers up to a given integer \(n\) that are relatively prime to \(n\). It is written using the Greek letter phi as \(\varphi(n)\) or \(\phi(n)\), and may also be called **Euler's phi function**. In other words, it is the number of integers \(k\) in the range \(1\leq k\leq n\) for which the greatest common divisor \(\gcd(n,k)\) is equal to 1. The integers \(k\) of this form are sometimes referred to as totatives of \(n\).

For example, the totatives of \(n=9\) are the six numbers 1, 2, 4, 5, 7 and 8. They are all relatively prime to 9, but the other three numbers in this range, 3, 6, and 9 are not, since \(\gcd(9,3)=\gcd(9,6)=3\) and \(\gcd(9,9)=9\). Therefore, \(\varphi(9)=6\). As another example, \(\varphi(1)=1\) since for \(n=1\) the only integer in the range from 1 to \(n\) is 1 itself, and \(\gcd(1,1)=1\).

Euler's totient function is a multiplicative function, meaning that if two numbers \(m\) and \(n\) are relatively prime, then \(\varphi(mn)=\varphi(m)\varphi(n)\).
This function gives the order of the multiplicative group of integers modulo  (the group of units of the ring \(\Z/n\Z\)). It is also used for defining the RSA encryption system.

## History, terminology, and notation
Leonhard Euler introduced the function in 1763. However, he did not at that time choose any specific symbol to denote it. In a 1784 publication, Euler studied the function further, choosing the Greek letter \(\pi\) to denote it: he wrote \(\pi D\) for "the multitude of numbers less than \(D\), and which have no common divisor with it". This definition varies from the current definition for the totient function at \(D=1\) but is otherwise the same. The now-standard notation \(\varphi(A)\) comes from Gauss's 1801 treatise *Disquisitiones Arithmeticae*, although Gauss did not use parentheses around the argument and wrote \(\varphi A\). Thus, it is often called **Euler's phi function** or simply the **phi function**.

In 1879, J. J. Sylvester coined the term **totient** for this function, so it is also referred to as **Euler's totient function**, the **Euler totient**, or **Euler's totient**. Jordan's totient is a generalization of Euler's.

The **cototient** of \(n\) is defined as \(n-\varphi(n)\). It counts the number of positive integers less than or equal to \(n\) that have at least one prime factor in common with \(n\).

## Computing Euler's totient function
There are several formulae for computing \(\varphi(n)\).

### Euler's product formula
It states
\(\varphi(n) =n \prod_{p\mid n} \left(1-\frac{1}{p}\right),\)
where the product is over the distinct prime numbers dividing .

An equivalent formulation is
\(\varphi(n) = p_1^{k_1-1}(p_1{-}1)\,p_2^{k_2-1}(p_2{-}1)\cdots p_r^{k_r-1}(p_r{-}1),\)
where \(n = p_1^{k_1} p_2^{k_2} \cdots p_r^{k_r}\) is the prime factorization of \(n\) (that is, \(p_1, p_2,\ldots, p_r\) are distinct prime numbers).

The proof of these formulae depends on two important facts.

#### Phi is a multiplicative function
This means that if \(\gcd(m,n) = 1\), then \(\varphi(m) \varphi(n) = \varphi(mn)\). *Proof outline:* Let \(A,B,C\) be the sets of positive integers which are coprime to and less than , , , respectively, so that \(|A| = \varphi(m)\), etc. Then there is a bijection between \(A\times B\) and  by the Chinese remainder theorem.

#### Value of phi for a prime power argument
If  is prime and \(k\geq1\), then

\(\varphi \left(p^k\right) = p^k-p^{k-1} = p^{k-1}(p-1) = p^k \left( 1 - \tfrac{1}{p} \right).\)

*Proof*: Since  is a prime number, the only possible values of \(\gcd(p^{k},m)\) are \(1,p,p^{2},\dots,p^{k}\), and the only way to have \(\gcd(p^{k},m)>1\) is if  is a multiple of , that is, \(m \in \{ p, 2p, 3p, \ldots, p^{k-1} p= p^{k}\}\), and there are \(p^{k-1}\) such multiples not greater than \(p^{k}\). Therefore, the other \(p^{k}-p^{k-1}\) numbers are all relatively prime to \(p^{k}\).

#### Proof of Euler's product formula
The fundamental theorem of arithmetic states that if *n* > 1 there is a unique expression \(n = p_1^{k_1} p_2^{k_2} \cdots p_r^{k_r},\) where *p*<sub>1</sub> < *p*<sub>2</sub> < ... < *p*<sub>*r*</sub> are prime numbers and each *k*<sub>*i*</sub> ≥ 1. (The case 1=*n* = 1 corresponds to the empty product.) Repeatedly using the multiplicative property of  and the formula for *φ*(*p*<sup>*k*</sup>) gives

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

\(\mathcal{F} \{ \mathbf{x} \}[m] = \sum\limits_{k=1}^n x_k \cdot e^{{-2\pi i}\frac{mk}{n}}\)

where *x<sub>k</sub>*  for {{math|*k* ∈ {1, ..., *n*}}}. Then

\(\varphi (n) = \mathcal{F} \{ \mathbf{x} \}[1] = \sum\limits_{k=1}^n \gcd(k,n) e^{-2\pi i\frac{k}{n}}.\)

The real part of this formula is

\(\varphi (n)=\sum\limits_{k=1}^n \gcd(k,n) \cos {\tfrac{2\pi k}{n}}
.\)

For example, using \(\cos\tfrac{\pi}5 = \tfrac{\sqrt 5+1}4\) and \(\cos\tfrac{2\pi}5 = \tfrac{\sqrt 5-1}4\):
\[
\begin{array}{rcl}
\varphi(10) &=& \gcd(1,10)\cos\tfrac{2\pi}{10} + \gcd(2,10)\cos\tfrac{4\pi}{10}
+ \gcd(3,10)\cos\tfrac{6\pi}{10}+\cdots+\gcd(10,10)\cos\tfrac{20\pi}{10}\\
&=& 1\cdot(\tfrac{\sqrt5+1}4)  + 2\cdot(\tfrac{\sqrt5-1}4) +
1\cdot(-\tfrac{\sqrt5-1}4) + 2\cdot(-\tfrac{\sqrt5+1}4) + 5\cdot (-1) \\
&& +\  2\cdot(-\tfrac{\sqrt5+1}4) + 1\cdot(-\tfrac{\sqrt5-1}4) + 2\cdot(\tfrac{\sqrt5-1}4) +
1\cdot(\tfrac{\sqrt5+1}4) + 10 \cdot (1) \\
&=& 4 .
\end{array}
\]
Unlike the Euler product and the divisor sum formula, this one does not require knowing the factors of . However, it does involve the calculation of the greatest common divisor of  and every positive integer less than , which suffices to provide the factorization anyway.

### Divisor sum
The property established by Gauss, that

\(\sum_{d\mid n}\varphi(d)=n,\)

where the sum is over all positive divisors  of , can be proven in several ways. (See Arithmetical function for notational conventions.)

One proof is to note that *φ*(*d*) is also equal to the number of possible generators of the cyclic group *C*<sub>*d*</sub> ; specifically, if *C*<sub>*d*</sub>  with 1=*g*<sup>*d*</sup> = 1, then *g*<sup>*k*</sup> is a generator for every  coprime to . Since every element of *C*<sub>*n*</sub> generates a cyclic subgroup, and each subgroup *C*<sub>*d*</sub> ⊆ *C*<sub>*n*</sub> is generated by precisely *φ*(*d*) elements of *C*<sub>*n*</sub>, the formula follows. Equivalently, the formula can be derived by the same argument applied to the multiplicative group of the th roots of unity and the primitive th roots of unity.

The formula can also be derived from elementary arithmetic. For example, let *n*  and consider the positive fractions up to 1 with denominator 20:
\(\tfrac{ 1}{20},\,\tfrac{ 2}{20},\,\tfrac{ 3}{20},\,\tfrac{ 4}{20},\,
 \tfrac{ 5}{20},\,\tfrac{ 6}{20},\,\tfrac{ 7}{20},\,\tfrac{ 8}{20},\,
 \tfrac{ 9}{20},\,\tfrac{10}{20},\,\tfrac{11}{20},\,\tfrac{12}{20},\,
 \tfrac{13}{20},\,\tfrac{14}{20},\,\tfrac{15}{20},\,\tfrac{16}{20},\,
 \tfrac{17}{20},\,\tfrac{18}{20},\,\tfrac{19}{20},\,\tfrac{20}{20}.\)

Put them into lowest terms:
\(\tfrac{ 1}{20},\,\tfrac{ 1}{10},\,\tfrac{ 3}{20},\,\tfrac{ 1}{ 5},\,
 \tfrac{ 1}{ 4},\,\tfrac{ 3}{10},\,\tfrac{ 7}{20},\,\tfrac{ 2}{ 5},\,
 \tfrac{ 9}{20},\,\tfrac{ 1}{ 2},\,\tfrac{11}{20},\,\tfrac{ 3}{ 5},\,
 \tfrac{13}{20},\,\tfrac{ 7}{10},\,\tfrac{ 3}{ 4},\,\tfrac{ 4}{ 5},\,
 \tfrac{17}{20},\,\tfrac{ 9}{10},\,\tfrac{19}{20},\,\tfrac{1}{1}\)

These twenty fractions are all the positive  ≤ 1 whose denominators are the divisors *d*  . The fractions with 20 as denominator are those with numerators relatively prime to 20, namely , , , , , , , ; by definition this is *φ*(20) fractions. Similarly, there are *φ*(10) fractions with denominator 10, and *φ*(5) fractions with denominator 5, etc. Thus the set of twenty fractions is split into subsets of size *φ*(*d*) for each *d* dividing 20. A similar argument applies for any *n.*

Möbius inversion applied to the divisor sum formula gives
\(\varphi(n) = \sum_{d\mid n} \mu\left( d \right) \cdot \frac{n}{d}  = n\sum_{d\mid n} \frac{\mu (d)}{d},\)

where  is the Möbius function, the multiplicative function defined by \(\mu(p) = -1\) and \(\mu(p^k) = 0\) for each prime 1=*p* and 1=*k* ≥ 2. This formula may also be derived from the product formula by multiplying out \(\prod_{p\mid n} (1 - \frac{1}{p})\) to get \(\sum_{d \mid n} \frac{\mu (d)}{d}.\)

An example:
\[
\begin{align}
\varphi(20) &= \mu(1)\cdot 20 + \mu(2)\cdot 10 +\mu(4)\cdot 5 +\mu(5)\cdot 4 + \mu(10)\cdot 2+\mu(20)\cdot 1\\[.5em]
&= 1\cdot 20 - 1\cdot 10 + 0\cdot 5 - 1\cdot 4 + 1\cdot 2 + 0\cdot 1 = 8.
\end{align}
\]


## Some values
The first 100 values  are shown in the table and graph below:


{| class="wikitable" style="text-align: right"
|+*φ*(*n*) for 1 ≤ *n* ≤ 100
! +
! 1 || 2 || 3 || 4 || 5 || 6 || 7 || 8 || 9 || 10
|-
! 0
| 1 || 1 || 2 || 2 || 4 || 2 || 6 || 4 || 6 || 4
|-
! 10
| 10 || 4 || 12 || 6 || 8 || 8 || 16 || 6 || 18 || 8
|-
! 20
| 12 || 10 || 22 || 8 || 20 || 12 || 18 || 12 || 28 || 8
|-
! 30
| 30 || 16 || 20 || 16 || 24 || 12 || 36 || 18 || 24 || 16
|-
! 40
| 40 || 12 || 42 || 20 || 24 || 22 || 46 || 16 || 42 || 20
|-
! 50
| 32 || 24 || 52 || 18 || 40 || 24 || 36 || 28 || 58 || 16
|-
! 60
| 60 || 30 || 36 || 32 || 48 || 20 || 66 || 32 || 44 || 24
|-
! 70
| 70 || 24 || 72 || 36 || 40 || 36 || 60 || 24 || 78 || 32
|-
! 80
| 54 || 40 || 82 || 24 || 64 || 42 || 56 || 40 || 88 || 24
|-
! 90
| 72 || 44 || 60 || 46 || 72 || 32 || 96 || 42 || 60 || 40
|}

In the graph at right the top line *y*  is an upper bound valid for all  other than one, and attained if and only if  is a prime number. A simple lower bound is \(\varphi(n) \ge \sqrt{n/2}\), which is rather loose: in fact, the lower limit of the graph is proportional to .


## Euler's theorem


This states that if  and  are relatively prime then

\(a^{\varphi(n)} \equiv 1\mod n.\)

The special case where  is prime is known as Fermat's little theorem.

This follows from Lagrange's theorem and the fact that *φ*(*n*) is the order of the multiplicative group of integers modulo .

The RSA cryptosystem is based on this theorem: it implies that the inverse of the function *a* ↦ *a<sup>e</sup>* mod *n*, where  is the (public) encryption exponent, is the function *b* ↦ *b<sup>d</sup>* mod *n*, where , the (private) decryption exponent, is the multiplicative inverse of  modulo *φ*(*n*). The difficulty of computing *φ*(*n*) without knowing the factorization of  is thus the difficulty of computing : this is known as the RSA problem which can be solved by factoring . The owner of the private key knows the factorization, since an RSA private key is constructed by choosing  as the product of two (randomly chosen) large primes  and . Only  is publicly disclosed, and given the difficulty to factor large numbers we have the guarantee that no one else knows the factorization.

## Other formulae
*\(a\mid b \implies \varphi(a)\mid\varphi(b)\)
*\(m \mid \varphi(a^m-1)\)
*\(\varphi(mn) = \varphi(m)\varphi(n)\cdot\frac{d}{\varphi(d)} \quad\text{where }d = \operatorname{gcd}(m,n)\)
**In particular:
*\(\varphi(2m) = \begin{cases}
2\varphi(m) &\text{ if } m \text{ is even} \\
\varphi(m) &\text{ if } m \text{ is odd}
\end{cases}\)
*\(\varphi\left(n^m\right) = n^{m-1}\varphi(n)\)
*\(\varphi(\operatorname{lcm}(m,n))\cdot\varphi(\operatorname{gcd}(m,n)) = \varphi(m)\cdot\varphi(n)\)
Compare this to the formula \(\operatorname{lcm}(m,n)\cdot \operatorname{gcd}(m,n) = m \cdot n\)  (see least common multiple).
**φ*(*n*) is even for *n* ≥ 3. Moreover, if  has  distinct odd prime factors, 2<sup>*r*</sup>
*For any *a* > 1 and *n* > 6 such that 4 ∤ *n* there exists an *l* ≥ 2*n* such that *l* .
*\(\frac{\varphi(n)}{n}=\frac{\varphi(\operatorname{rad}(n))}{\operatorname{rad}(n)}\)
where rad(*n*) is the radical of  (the product of all distinct primes dividing ).
*\(\sum_{d \mid n} \frac{\mu^2(d)}{\varphi(d)} = \frac{n}{\varphi(n)}\) 
*\(\sum_{1\le k\le n-1 \atop gcd(k,n)=1}\!\!k = \tfrac12 n\varphi(n) \quad \text{for }n>1\)
*\(\sum_{k=1}^n\varphi(k) = \tfrac12 \left(1+ \sum_{k=1}^n \mu(k)\left\lfloor\frac{n}{k}\right\rfloor^2\right)
=\frac3{\pi^2}n^2+O\left(n(\log n)^\frac23(\log\log n)^\frac43\right)\) ( cited in)
*\(\sum_{k=1}^n\varphi(k) =\frac3{\pi^2}n^2+O\left(n(\log n)^\frac23(\log\log n)^\frac13\right)\) [Liu (2016)]
*\(\sum_{k=1}^n\frac{\varphi(k)}{k} = \sum_{k=1}^n\frac{\mu(k)}{k}\left\lfloor\frac{n}{k}\right\rfloor=\frac6{\pi^2}n+O\left((\log n)^\frac23(\log\log n)^\frac43\right)\) 
*\(\sum_{k=1}^n\frac{k}{\varphi(k)} = \frac{315\,\zeta(3)}{2\pi^4}n-\frac{\log n}2+O\left((\log n)^\frac23\right)\) 
*\(\sum_{k=1}^n\frac{1}{\varphi(k)} = \frac{315\,\zeta(3)}{2\pi^4}\left(\log n+\gamma-\sum_{p\text{ prime}}\frac{\log p}{p^2-p+1}\right)+O\left(\frac{(\log n)^\frac23}n\right)\) (where  is the Euler–Mascheroni constant).

### Menon's identity

In 1965 P. Kesava Menon proved
\(\sum_{\stackrel{1\le k\le n}{ \gcd(k,n)=1}} \!\!\!\! \gcd(k-1,n)=\varphi(n)d(n),\)
where \varphi(n)\) holds for almost all \(n\), meaning for all but \(o(x)\) values of \(n\le x\) as \(x\rightarrow\infty\).

This is an elementary consequence of the fact that the sum of the reciprocals of the primes congruent to 1 modulo \(q\) diverges, which itself is a corollary of the proof of Dirichlet's theorem on arithmetic progressions.

## Generating functions
The [[Dirichlet series for *φ*(*n*) may be written in terms of the Riemann zeta function as:
\(\sum_{n=1}^\infty \frac{\varphi(n)}{n^s}=\frac{\zeta(s-1)}{\zeta(s)}\)
where the left-hand side converges for \(\Re (s)>2\).

The Lambert series generating function is

\(\sum_{n=1}^{\infty} \frac{\varphi(n) q^n}{1-q^n}= \frac{q}{(1-q)^2}\)

which converges for |*q*| < 1.

Both of these are proved by elementary series manipulations and the formulae for *φ*(*n*).

## Growth rate
In the words of Hardy & Wright, the order of *φ*(*n*) is "always 'nearly '."

First

\(\lim\sup \frac{\varphi(n)}{n}= 1,\)

but as *n* goes to infinity, for all *δ* > 0

\(\frac{\varphi(n)}{n^{1-\delta}}\rightarrow\infty.\)

These two formulae can be proved by using little more than the formulae for *φ*(*n*) and the divisor sum function *σ*(*n*).

In fact, during the proof of the second formula, the inequality

\(\frac {6}{\pi^2} < \frac{\varphi(n) \sigma(n)}{n^2} < 1,\)

true for *n* > 1, is proved.

We also have

\(\lim\inf\frac{\varphi(n)}{n}\log\log n = e^{-\gamma}.\)

Here  is Euler's constant, *γ* , so *e<sup>γ</sup>*  and *e*<sup>−*γ*</sup> .

Proving this does not quite require the prime number theorem. Since log log *n* goes to infinity, this formula shows that

\(\lim\inf\frac{\varphi(n)}{n}= 0.\)

In fact, more is true.

\(\varphi(n) > \frac {n} {e^\gamma\; \log \log n + \frac {3} {\log \log n}} \quad\text{for } n>2\)

and

\(\varphi(n) < \frac {n} {e^{ \gamma}\log \log n} \quad\text{for infinitely many } n.\)

The second inequality was shown by Jean-Louis Nicolas. Ribenboim says "The method of proof is interesting, in that the inequality is shown first under the assumption that the Riemann hypothesis is true, secondly under the contrary assumption."

For the average order, we have

\(\varphi(1)+\varphi(2)+\cdots+\varphi(n) = \frac{3n^2}{\pi^2}+O\left(n(\log n)^\frac23(\log\log n)^\frac43\right) \quad\text{as }n\rightarrow\infty,\)
due to Arnold Walfisz, its proof exploiting estimates on exponential sums due to I. M. Vinogradov and N. M. Korobov.
By a combination of van der Corput's and Vinogradov's methods, H.-Q. Liu (On Euler's function.Proc. Roy. Soc. Edinburgh Sect. A 146 (2016), no. 4, 769–775)
improved the error term to
\(O\left(n(\log n)^\frac23(\log\log n)^\frac13\right)\)
(this is currently the best known estimate of this type). The "Big " stands for a quantity that is bounded by a constant times the function of  inside the parentheses (which is small compared to *n*<sup>2</sup>).

This result can be used to prove that the probability of two randomly chosen numbers being relatively prime is <sup>2</sup>}}.

## Ratio of consecutive values
In 1950 Somayajulu proved

\(\begin{align}
\lim\inf \frac{\varphi(n+1)}{\varphi(n)}&= 0 \quad\text{and} \\[5px]
\lim\sup \frac{\varphi(n+1)}{\varphi(n)}&= \infty.
\end{align}\)

In 1954 Schinzel and Sierpiński strengthened this, proving that the set

\(\left\{\frac{\varphi(n+1)}{\varphi(n)},\;\;n = 1,2,\ldots\right\}\)

is dense in the positive real numbers. They also proved that the set

\(\left\{\frac{\varphi(n)}{n},\;\;n = 1,2,\ldots\right\}\)

is dense in the interval (0,1).

## Totient number
A **totient number** is a value of Euler's totient function: that is, an  for which there is at least one  for which *φ*(*n*) . The *valency* or *multiplicity* of a totient number  is the number of solutions to this equation. A *nontotient* is a natural number which is not a totient number. Every odd integer exceeding 1 is trivially a nontotient. There are also infinitely many even nontotients, and indeed every positive integer has a multiple which is an even nontotient.

The first few totient numbers are \(1, 2, 4, 6, 8, 10, 12, 16, 18, 20\), see sequence  .

The number of totient numbers up to a given limit  is

\(\frac{x}{\log x}e^{ \big(C+o(1)\big)(\log\log\log x)^2 }\)

for a constant *C* .

If counted accordingly to multiplicity, the number of totient numbers up to a given limit  is

\(\Big\vert\{ n : \varphi(n) \le x \}\Big\vert = \frac{\zeta(2)\zeta(3)}{\zeta(6)} \cdot x + R(x)\)

where the error term  is of order at most  for any positive .

It is known that the multiplicity of  exceeds *m*<sup>*δ*</sup> infinitely often for any *δ* < 0.55655.

### Ford's theorem
 proved that for every integer *k* ≥ 2 there is a totient number  of multiplicity : that is, for which the equation *φ*(*n*)  has exactly  solutions; this result had previously been conjectured by Wacław Sierpiński, and it had been obtained as a consequence of Schinzel's hypothesis H. Indeed, each multiplicity that occurs, does so infinitely often.

However, no number  is known with multiplicity *k* . Carmichael's totient function conjecture is the statement that there is no such .

### Perfect totient numbers

A perfect totient number is an integer that is equal to the sum of its iterated totients. That is, we apply the totient function to a number *n*, apply it again to the resulting totient, and so on, until the number 1 is reached, and add together the resulting sequence of numbers; if the sum equals *n*, then *n* is a perfect totient number.

## Applications
### Cyclotomy


In the last section of the *Disquisitiones* Gauss proves that a regular -gon can be constructed with straightedge and compass if *φ*(*n*) is a power of 2. If  is a power of an odd prime number the formula for the totient says its totient can be a power of two only if  is a first power and *n* − 1 is a power of 2. The primes that are one more than a power of 2 are called Fermat primes, and only five are known: 3, 5, 17, 257, and 65537. Fermat and Gauss knew of these. Nobody has been able to prove whether there are any more.

Thus, a regular -gon has a straightedge-and-compass construction if *n* is a product of distinct Fermat primes and any power of 2. The first few such  are
2, 3, 4, 5, 6, 8, 10, 12, 15, 16, 17, 20, 24, 30, 32, 34, 40,... .

### Prime number theorem for arithmetic progressions


### The RSA cryptosystem


Setting up an RSA system involves choosing large prime numbers  and , computing *n*  and *k* , and finding two numbers  and  such that *ed* ≡ 1 (mod *k*). The numbers  and  (the "encryption key") are released to the public, and  (the "decryption key") is kept private.

A message, represented by an integer , where 0 < *m* < *n*, is encrypted by computing *S* .

It is decrypted by computing *t* . Euler's Theorem can be used to show that if 0 < *t* < *n*, then *t* .

The security of an RSA system would be compromised if the number  could be efficiently factored or if *φ*(*n*) could be efficiently computed without factoring .

## Unsolved problems
### Lehmer's conjecture


If  is prime, then *φ*(*p*) . In 1932 D. H. Lehmer asked if there are any composite numbers  such that *φ*(*n*)  divides *n* − 1. None are known.

In 1933 he proved that if any such  exists, it must be odd, square-free, and divisible by at least seven primes (i.e. *ω*(*n*) ≥ 7). In 1980 Cohen and Hagis proved that *n* > 10<sup>20</sup> and that *ω*(*n*) ≥ 14. Further, Hagis showed that if 3 divides  then *n* > 10<sup>1937042</sup> and *ω*(*n*) ≥ 298848.

### Carmichael's conjecture


This states that there is no number \(n\) with the property that for all other numbers \(m\), \(m\neq n\), \(\varphi(m) \neq \varphi(n)\). See Ford's theorem above.

If there is a single counterexample to this conjecture, there must be infinitely many counterexamples, and the smallest one has at least ten billion digits in base 10.

### Riemann hypothesis
The Riemann hypothesis is true if and only if the inequality
\(\frac{n}{\varphi (n)}<e^\gamma \log\log n+\frac{e^\gamma (4+\gamma-\log 4\pi)}{\sqrt{\log n}}\)
is true for all \(n\geq p_{120569}\#\) where \(\gamma\) is Euler's constant and \(p_{120569}\#\) is the product of the first 120569 primes.

## See also
*Carmichael function (λ)
*Dedekind psi function (𝜓)
*Divisor function (σ)
*Duffin–Schaeffer conjecture
*Generalizations of Fermat's little theorem
*Highly composite number
*Multiplicative group of integers modulo
*Ramanujan sum
*Totient summatory function (𝛷)

## Notes


## References


The *Disquisitiones Arithmeticae* has been translated from Latin into English and German. The German edition includes all of Gauss's papers on number theory: all the proofs of quadratic reciprocity, the determination of the sign of the Gauss sum, the investigations into biquadratic reciprocity, and unpublished notes.

References to the *Disquisitiones* are of the form Gauss, DA, art. *nnn*.

*. See paragraph 24.3.2.
*
* Dickson, Leonard Eugene, "History Of The Theory Of Numbers", vol 1, chapter 5 "Euler's Function, Generalizations; Farey Series", Chelsea Publishing 1952
*.
*
*
*
*
*
*.
*
*
*
*
*
*
*.


==External links==
*
*Euler's Phi Function and the Chinese Remainder Theorem — proof that *φ*(*n*) is multiplicative
*Euler's totient function calculator in JavaScript — up to 20 digits
*Dineva, Rosica, The Euler Totient, the Möbius, and the Divisor Functions
*Plytage, Loomis, Polhill Summing Up The Euler Phi Function

