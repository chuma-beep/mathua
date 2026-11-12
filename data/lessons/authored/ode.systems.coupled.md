# Coupled oscillators

Coupled oscillators turn into a linear system whose normal modes decouple the motion. Each mode oscillates at its own frequency, and the general solution is a sum of modes.

## The state

Each oscillator contributes a position and a velocity, so \(n\) masses give a first-order system of size \(2n\). The coupling appears as off-diagonal blocks in the system matrix.

## The stiffness matrix

The restoring forces form a symmetric stiffness matrix with \(n^{2}\) entries. Its eigenvectors are the mode shapes and its eigenvalues give the squared frequencies.

## Two equal masses

Two equal masses joined by a spring of stiffness \(k\) have squared normal frequencies 0 and \(2k\). The zero mode is the common translation and the other mode is the relative motion.

To check your understanding, count position and velocity, count the entries and find the zero mode.

Coupled motion is a linear system whose modes oscillate independently.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
