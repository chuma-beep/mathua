# Transformations of Random Vectors

This lesson develops the definition, the main calculation and the interpretation of the idea.

## The Jacobian

A transformation from one vector of variables to another changes density by the absolute value of the Jacobian determinant. If $u=T(x)$ is one-to-one and differentiable, then

$$f_U(u)=f_X(T^{-1}(u))\left|\det \frac{\partial x}{\partial u}\right|.$$

The determinant measures how much a small volume is stretched or compressed. A factor larger than one thins the density, and a factor smaller than one concentrates it.

For a linear map $u=Ax$, the Jacobian is the constant matrix $A$, so the density is scaled by $1/|\det A|$.

## Linear changes of variable

A shear has determinant one, so it preserves density. A diagonal scaling by $a$ and $b$ has determinant $ab$, so areas are multiplied by $ab$ and densities are divided by $ab$.

A rotation has determinant one and preserves area, which is why a rotated normal density keeps the same normalising constant. A reflection has determinant minus one, but the absolute value removes the sign.

The support must be transformed as well. A density formula without the new region of support is incomplete.

## Sums and differences

The map from $(X,Y)$ to $(S,D)=(X+Y,X-Y)$ has Jacobian determinant two in absolute value, so the joint density of the pair is half the original density expressed in the new variables.

For independent variables, the marginal density of the sum is the convolution. The transformation method recovers that result and extends it to pairs of functions that are not simply a sum.

When the map is not one-to-one, split the domain into pieces on which it is, transform each piece, and add the resulting densities.

