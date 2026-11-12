# Integrating differential forms

A differential form is integrated by pulling it back to a parameter domain.

## Pulling back to a parameter domain

A parameterization maps a plane domain into a surface. The pullback substitutes it into the form.

The result is a multiple of the area element, so the integral becomes an ordinary double integral. Orientation is part of the data, so reversing the parameterization changes the sign.

For the unit disk the area form pulls back to \(r\,dr\wedge d\phi\), whose integral is \(\pi\), about 3.1416.

## The integral of a two-form

On a surface in three dimensions a two-form combines \(dy\wedge dz\), \(dz\wedge dx\) and \(dx\wedge dy\).

Integrating it gives the flux of the matching vector field through the surface. On the unit sphere the area form pulls back to \(\sin\theta\,d\theta\wedge d\phi\), whose integral is \(4\pi\).

## A worked flux

Take \(\omega=x\,dy\wedge dz+y\,dz\wedge dx+z\,dx\wedge dy\) over the unit sphere, which pulls back to \(\sin\theta\,d\theta\wedge d\phi\).

The \(\theta\) integral from 0 to \(\pi\) gives 2 and the \(\phi\) integral gives \(2\pi\), so the flux is \(4\pi\). That is three times the volume of the unit ball.

The factor 3 is the divergence of the position field, which the divergence theorem confirms.
