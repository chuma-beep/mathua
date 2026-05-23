# Euclidean Algorithm for GCD

The **greatest common divisor (GCD)** of two integers is the largest integer that divides both without remainder. The **Euclidean algorithm** is an efficient method for finding the GCD.

**Euclidean Algorithm:**
Given integers $a$ and $b$ with $a > b > 0$:
1. Divide $a$ by $b$: $a = bq + r$ where $0 \le r < b$
2. If $r = 0$, then $\gcd(a, b) = b$
3. Otherwise, replace $(a, b)$ with $(b, r)$ and repeat

**Example 1:** Find $\gcd(252, 105)$.

Step 1: $252 = 105 \times 2 + 42$
Step 2: $105 = 42 \times 2 + 21$
Step 3: $42 = 21 \times 2 + 0$

The last nonzero remainder is 21, so $\gcd(252, 105) = 21$.

**Example 2:** Find $\gcd(48, 18)$.

Step 1: $48 = 18 \times 2 + 12$
Step 2: $18 = 12 \times 1 + 6$
Step 3: $12 = 6 \times 2 + 0$

$\gcd(48, 18) = 6$

**Why it works:** Any common divisor of $a$ and $b$ also divides $r = a - bq$, so $\gcd(a,b) = \gcd(b,r)$. The remainders get smaller each step, eventually reaching zero.

**Practice:**
- $\gcd(1071, 462)$
- $\gcd(123, 45)$
