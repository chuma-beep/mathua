# The Gudermannian angle

The Gudermannian function converts a hyperbolic angle into an ordinary angle by \(gd(u) = 2\arctan(e^{u}) - \pi/2\).

## Defining the angle

The definition is equivalent to two others that are easier to use: \(\sin(gd\,u) = \tanh u\) and
\(\tan(gd\,u) = \sinh u\). At \(u = 0\) both give 0, since the hyperbolic tangent and sine are 0
there, so the Gudermannian fixes the origin. The function is odd and increases from \(-\pi/2\) to
\(\pi/2\) as \(u\) runs over the whole line. Hyperbolic distance and circular angle are the same
quantity seen through two different measures.

## Degrees from a hyperbolic input

Every hyperbolic value produces one angle between -90 and 90 degrees. For \(u = 1\),
\(\tanh 1 \approx 0.7616\), so the angle with that sine is about 49.60 degrees. For \(u = 2\),
\(\tanh 2 \approx 0.9640\), so the angle is about 74.59 degrees. Small inputs behave almost
linearly: \(u = 0.1\) gives about 5.72 degrees, close to the 5.73 degrees of 0.1 radian, because
small hyperbolic and circular angles agree closely.

## The tangent relation

The tangent form is the quickest check. At \(u = 1\), \(\sinh 1 \approx 1.1752\), so
\(\arctan 1.1752 \approx 49.60\) degrees, matching the sine form. At \(u = 2\),
\(\sinh 2 \approx 3.6269\), so \(\arctan 3.6269 \approx 74.59\) degrees. The two routes agree
because \(\tan(gd\,u) = \sinh u\) is an identity, not an approximation.

Use the tangent form when a calculator has no inverse hyperbolic key.
