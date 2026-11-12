# Taylor series of sine and cosine

The Taylor series of sine and cosine around zero are power series whose coefficients come from the four-step derivative cycle.

## The sine series

Only odd powers survive, because every even derivative of sine vanishes at zero. The series is
\(\sin x = x - x^{3}/6 + x^{5}/120 - x^{7}/5040 + \cdots\). At \(x = 1\) the first four terms give
\(1 - 0.16667 + 0.00833 - 0.00020 = 0.84147\), which matches \(\sin 1\) to five decimal places.
Four terms already reach that accuracy, because the factorial in each denominator grows quickly.

## The cosine series

Cosine keeps the even powers and starts at 1. The series is
\(\cos x = 1 - x^{2}/2 + x^{4}/24 - x^{6}/720 + \cdots\). At \(x = 1\) the first four terms give
\(1 - 0.5 + 0.04167 - 0.00139 = 0.54028\), while \(\cos 1 = 0.54030\). The two agree to three
decimal places; the gap is the size of the next term.

## Reading a coefficient

A coefficient is read straight from the factorial pattern. In the sine series the term in
\(x^{5}\) has coefficient \(1/120\), while the term in \(x^{6}\) is absent, so its coefficient is
zero. In the cosine series the term in \(x^{6}\) has coefficient \(-1/720 \approx -0.00139\). In
both series the coefficient of \(x^{n}\) is the value of the \(n\) th derivative at zero divided
by \(n!\). Neither series needs any restriction on \(x\).

The factorial in the denominator is what makes both series converge everywhere.
