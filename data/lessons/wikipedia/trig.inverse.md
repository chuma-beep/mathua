> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Inverse_trigonometric_functions) — CC BY-SA 4.0

# Inverse trig functions

In mathematics, the **inverse trigonometric functions** (occasionally also called *antitrigonometric*, *cyclometric*, or *arcus* functions) are the inverse functions of the trigonometric functions, under suitably restricted domains. Specifically, they are the inverses of the sine, cosine, tangent, cotangent, secant, and cosecant functions, and are used to obtain an angle from any of the angle's trigonometric ratios. Inverse trigonometric functions are widely used in engineering, navigation, physics, and geometry.

## Notation

Several notations for the inverse trigonometric functions exist. The most common convention is to name inverse trigonometric functions using an arc- prefix: arcsin(*x*), arccos(*x*), arctan(*x*), etc. (This convention is used throughout this article.) This notation arises from the following geometric relationships:
when measuring in radians, an angle of radians will correspond to an arc whose length where is the radius of the circle. Thus in the unit circle, the cosine of x function is both the arc and the angle, because the arc of a circle of radius 1 is the same as the angle. Or, "the arc whose cosine is " is the same as "the angle whose cosine is ", because the length of the arc of the circle in radii is the same as the measurement of the angle in radians. In computer programming languages, the inverse trigonometric functions are often called by the abbreviated forms ,. The notations sin\(^{−1}\)(*x*), cos\(^{−1}\)(*x*), tan\(^{−1}\)(*x*), etc., as introduced by John Herschel in 1813, are often used as well in English-language sources, much more than the also established sin\(^{[−1]}\)(*x*), cos\(^{[−1]}\)(*x*), tan\(^{[−1]}\)(*x*) – conventions consistent with the notation of an inverse function, that is useful (for example) to define the multivalued version of each inverse trigonometric function: \(\tan^{-1}(x) = \{\arctan(x) + \pi k \mid k \in \mathbb Z\} ~.\) However, this might appear to conflict logically with the common semantics for expressions such as sin\(^{2}\)(*x*) (although only sin\(^{2}\) *x*, without parentheses, is the really common use), which refer to numeric power rather than function composition, and therefore may result in confusion between notation for the reciprocal (multiplicative inverse) and inverse function.

The confusion is somewhat mitigated by the fact that each of the reciprocal trigonometric functions has its own name — for example, (cos(*x*))\(^{−1}\). Nevertheless, certain authors advise against using it; since it is ambiguous. Another precarious convention used by a small number of authors is to use an uppercase first letter, along with a “−1” superscript: Sin\(^{−1}\)(*x*), Cos\(^{−1}\)(*x*), Tan\(^{−1}\)(*x*), etc. Although it is intended to avoid confusion with the reciprocal, which should be represented by sin\(^{−1}\)(*x*), cos\(^{−1}\)(*x*), etc., or, better, by sin\(^{−1}\) *x*, cos\(^{−1}\) *x*, etc., it in turn creates yet another major source of ambiguity, especially since many popular high-level programming languages (e.g. Mathematica and MAGMA) use those very same capitalised representations for the standard trig functions, whereas others (e.g. Python, Matlab, MAPLE) use lower-case.

Hence; since 2009, the ISO 80000-2 standard has specified solely the "arc" prefix for the inverse functions.

## Basic concepts

### Principal values
Since none of the six trigonometric functions are one-to-one, they must be restricted in order to have inverse functions. Therefore, the result ranges of the inverse functions are proper (i.e. strict) subsets of the domains of the original functions.

For example, using in the sense of multivalued functions, just as the square root function \(y = \sqrt{x}\) could be defined from \(y^2 = x,\) the function \(y = \arcsin(x)\) is defined so that \(\sin(y) = x.\) For a given real number \(x,\) with \(-1 \leq x \leq 1,\) there are multiple (in fact, countably infinitely many) numbers \(y\) such that \(\sin(y) = x\); for example, \(\sin(0) = 0,\) but also \(\sin(\pi) = 0,\) \(\sin(2 \pi) = 0,\) etc. When only one value is desired, the function may be restricted to its principal branch. With this restriction, for each \(x\) in the domain, the expression \(\arcsin(x)\) will evaluate only to a single value, called its principal value. These properties apply to all the inverse trigonometric functions.

The principal inverses are listed in the following table.

{| class="wikitable" style="text-align:center"
|-
! scope="col" | Name
! scope="col" | Usual notation
! scope="col" | Definition
! scope="col" | Domain of for real result
! scope="col" | Range of usual principal value (radians)
! scope="col" | Range of usual principal value (degrees)
|-
! scope="row" | arcsine
| 1= *y* = arcsin(*x*) || 1=*x* = sin(*y*) || −1 ≤ *x* ≤ 1 || −|| −90° ≤ *y* ≤ 90°
|-
! scope="row" | arccosine
| 1= *y* = arccos(*x*) || 1=*x* = cos(*y*) || −1 ≤ *x* ≤ 1 || 0 ≤ *y* ≤ π || 0° ≤ *y* ≤ 180°
|-
! scope="row" | arctangent
| 1= *y* = arctan(*x*) || 1=*x* = tan(*y*) || all real numbers || −|| −90° < *y* < 90°
|-
! scope="row" | arccotangent
| 1= *y* = arccot(*x*) || 1=*x* = cot(*y*) || all real numbers || 0 < *y* < π || 0° < *y* < 180°
|-
! scope="row" | arcsecant
| 1= *y* = arcsec(*x*) || 1=*x* = sec(*y*) || |*x*| ≥ 1 || 0 ≤ *y* < or || 0° ≤ *y* < 90° or 90° < *y* ≤ 180°
|-
! scope="row" | arccosecant
| 1= *y* = arccsc(*x*) ||1=*x* = csc(*y*) || |*x*| ≥ 1 || − or 0 < *y* ≤ || −90° ≤ *y* < 0 or 0° < *y* ≤ 90°
|-
|}
Note: Some authors define the range of arcsecant to be \(\arccot(z)\)
