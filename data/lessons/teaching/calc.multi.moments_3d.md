# Moments and Centers of Mass in Three Dimensions

## Density and Mass

A solid with density \(\rho(x,y,z)\) has mass equal to the triple integral of density over its volume. When the density is constant, the mass is density times volume. A box with side lengths \(a\), \(b\) and \(c\) therefore has mass \(\rho abc\), which is the simplest case and a useful check on a more complicated integral.

Density can vary from point to point, so the mass is not always the product of a single density and a volume. The integral adds the contributions of small pieces, each with mass approximately density times volume. This is the same limiting process that defines the triple integral itself.

## First Moments

The first moments measure where the mass is located. The moment with respect to the \(xy\)-plane is the integral of \(z\rho\), and the centre of mass has \(z\)-coordinate equal to that moment divided by the total mass. The other two coordinates are defined analogously with \(x\) and \(y\).

For a box with constant density, the centre of mass is the geometric centre. The \(z\)-coordinate is half the height, so the first moment is the mass times half the height. This agrees with the integral and gives a quick way to check a computation before moving to a more general solid.

## Moments of Inertia

The moment of inertia about an axis measures resistance to rotation. About the \(z\)-axis it is the integral of \((x^2+y^2)\rho\), because each piece contributes in proportion to its squared distance from the axis. For a constant-density box, the result is \(\rho abc(a^2+b^2)/12\).

The factor twelve comes from averaging the squared distance over the box. The moment of inertia grows with the square of the side lengths, so doubling every dimension multiplies the moment by four when the density is fixed. These scaling checks catch many algebra errors in the final answer.
