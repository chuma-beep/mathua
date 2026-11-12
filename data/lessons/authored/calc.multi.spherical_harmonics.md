# Spherical harmonics in three dimensions

Spherical harmonics are the eigenfunctions of the angular part of the Laplacian on the sphere.

## Separating the sphere

Writing \(u=R(r)Y(\theta,\phi)\) splits the Laplacian into a radial part and an angular part.

The angular equation carries the eigenvalue \(l(l+1)\) for a nonnegative integer \(l\). The radial factor then grows like \(r^{l}\) or decays like \(r^{-l-1}\).

For \(l=1\) the eigenvalue is 2 and the radial powers are \(r\) and \(r^{-2}\). For \(l=2\) the eigenvalue is 6.

## The spherical harmonics

For each \(l\) there are \(2l+1\) functions \(Y_{lm}\), with \(m\) from \(-l\) to \(l\).

They are orthonormal on the sphere, so the integral of \(Y_{lm}\) against the conjugate of \(Y_{l'm'}\) is 1 when the labels agree. The degree \(l=1\) therefore gives three functions, the coordinate functions themselves.

## The first few harmonics

The constant harmonic is \(1/\sqrt{4\pi}\), about 0.2821, whose square integrates to 1 over the sphere of area \(4\pi\).

The harmonics with \(l=1\) are proportional to \(x\), \(y\) and \(z\) on the sphere. The one for \(z\) is \(\sqrt{3/(4\pi)}\cos\theta\), about \(0.4886\cos\theta\).

At the north pole that is 0.4886. At the equator it is zero.

The degree sets the eigenvalue. The order separates the functions within a degree.
