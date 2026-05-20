> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Inverse_trigonometric_functions) — CC BY-SA 4.0

# Inverse trig functions

In mathematics, the **inverse trigonometric functions** (occasionally also called *antitrigonometric*, *cyclometric*, or *arcus* functions) are the inverse functions of the trigonometric functions, under suitably restricted domains. Specifically, they are the inverses of the sine, cosine, tangent, cotangent, secant, and cosecant functions, and are used to obtain an angle from any of the angle's trigonometric ratios. Inverse trigonometric functions are widely used in engineering, navigation, physics, and geometry.

## Notation


Several notations for the inverse trigonometric functions exist. The most common convention is to name inverse trigonometric functions using an arc- prefix: arcsin(*x*), arccos(*x*), arctan(*x*), etc. (This convention is used throughout this article.)  This notation arises from the following geometric relationships:
when measuring in radians, an angle of  radians will correspond to an arc whose length is , where  is the radius of the circle. Thus in the unit circle, the cosine of x function is both the arc and the angle, because the arc of a circle of radius 1 is the same as the angle. Or, "the arc whose cosine is " is the same as "the angle whose cosine is ", because the length of the arc of the circle in radii is the same as the measurement of the angle in radians. In computer programming languages, the inverse trigonometric functions are often called by the abbreviated forms , , .

The notations sin<sup>−1</sup>(*x*), cos<sup>−1</sup>(*x*), tan<sup>−1</sup>(*x*), etc., as introduced by John Herschel in 1813, are often used as well in English-language sources, much more than the also established sin<sup>[−1]</sup>(*x*), cos<sup>[−1]</sup>(*x*), tan<sup>[−1]</sup>(*x*) – conventions consistent with the notation of an inverse function, that is useful (for example) to define the multivalued version of each inverse trigonometric function: \(\tan^{-1}(x) = \{\arctan(x) + \pi k \mid k \in \mathbb Z\} ~.\) However, this might appear to conflict logically with the common semantics for expressions such as sin<sup>2</sup>(*x*) (although only sin<sup>2</sup> *x*, without parentheses, is the really common use), which refer to numeric power rather than function composition, and therefore may result in confusion between notation for the reciprocal (multiplicative inverse) and inverse function.

The confusion is somewhat mitigated by the fact that each of the reciprocal trigonometric functions has its own name — for example, (cos(*x*))<sup>−1</sup> . Nevertheless, certain authors advise against using it, since it is ambiguous. Another precarious convention used by a small number of authors is to use an uppercase first letter, along with a “−1” superscript: Sin<sup>−1</sup>(*x*), Cos<sup>−1</sup>(*x*), Tan<sup>−1</sup>(*x*), etc. Although it is intended to avoid confusion with the reciprocal, which should be represented by sin<sup>−1</sup>(*x*), cos<sup>−1</sup>(*x*), etc., or, better, by sin<sup>−1</sup> *x*, cos<sup>−1</sup> *x*, etc., it in turn creates yet another major source of ambiguity, especially since many popular high-level programming languages (e.g. Mathematica and MAGMA) use those very same capitalised representations for the standard trig functions, whereas others (e.g. Python, Matlab, MAPLE) use lower-case.

Hence, since 2009, the ISO 80000-2 standard has specified solely the "arc" prefix for the inverse functions.

## Basic concepts


### Principal values
Since none of the six trigonometric functions are one-to-one, they must be restricted in order to have inverse functions. Therefore, the result ranges of the inverse functions are proper (i.e. strict) subsets of the domains of the original functions.

For example, using  in the sense of multivalued functions, just as the square root function \(y = \sqrt{x}\) could be defined from \(y^2 = x,\) the function \(y = \arcsin(x)\) is defined so that \(\sin(y) = x.\) For a given real number \(x,\) with \(-1 \leq x \leq 1,\) there are multiple (in fact, countably infinitely many) numbers \(y\) such that \(\sin(y) = x\); for example, \(\sin(0) = 0,\) but also \(\sin(\pi) = 0,\) \(\sin(2 \pi) = 0,\) etc. When only one value is desired, the function may be restricted to its principal branch. With this restriction, for each \(x\) in the domain, the expression \(\arcsin(x)\) will evaluate only to a single value, called its principal value. These properties apply to all the inverse trigonometric functions.

The principal inverses are listed in the following table.

{| class="wikitable" style="text-align:center"
|-
! scope="col" | Name
! scope="col" | Usual notation
! scope="col" | Definition
! scope="col" | Domain of  for real result
! scope="col" | Range of usual principal value (radians)
! scope="col" | Range of usual principal value (degrees)
|-
! scope="row" | arcsine
| 1= *y* = arcsin(*x*) || 1=*x* = sin(*y*) || −1 ≤ *x* ≤ 1 || −}} || −90° ≤ *y* ≤ 90°
|-
! scope="row" | arccosine
| 1= *y* = arccos(*x*) || 1=*x* = cos(*y*) || −1 ≤ *x* ≤ 1 || 0 ≤ *y* ≤ π || 0° ≤ *y* ≤ 180°
|-
! scope="row" | arctangent
| 1= *y* = arctan(*x*) || 1=*x* = tan(*y*) || all real numbers || −}} || −90° < *y* < 90°
|-
! scope="row" | arccotangent
| 1= *y* = arccot(*x*) || 1=*x* = cot(*y*) || all real numbers || 0 < *y* < π || 0° < *y* < 180°
|-
! scope="row" | arcsecant
| 1= *y* = arcsec(*x*) || 1=*x* = sec(*y*) || |*x*| ≥ 1 || 0 ≤ *y* <  or  || 0° ≤ *y* < 90° or 90° < *y* ≤ 180°
|-
! scope="row" | arccosecant
| 1= *y* = arccsc(*x*) ||1=*x* = csc(*y*) || |*x*| ≥ 1 || − or 0 < *y* ≤  || −90° ≤ *y* < 0 or 0° < *y* ≤ 90°
|-
|}
Note: Some authors define the range of arcsecant to be {{nowrap|(\(0 \leq y < \frac{\pi}{2}\)}} or {{nowrap|\(\pi \leq y < \frac{3 \pi}{2}\) ),}} because the tangent function is nonnegative on this domain. This makes some computations more consistent. For example, using this range, \(\tan(\arcsec(x)) = \sqrt{x^2 - 1},\) whereas with the range {{nowrap|(\(0 \leq y < \frac{\pi}{2}\)}} or {{nowrap|\(\frac{\pi}{2} < y \leq \pi\)),}} we would have to write \(\tan(\arcsec(x)) = \pm \sqrt{x^2 - 1},\) since tangent is nonnegative on \(0 \leq y < \frac{\pi}{2},\) but nonpositive on \(\frac{\pi}{2} < y \leq \pi.\) For a similar reason, the same authors define the range of arccosecant to be \(( - \pi < y \leq - \frac{\pi}{2}\) or \(0 < y \leq \frac{\pi}{2} ) .\)

#### Domains
If  is allowed to be a complex number, then the range of  applies only to its real part.


### Solutions to elementary trigonometric equations
Each of the trigonometric functions is periodic in the real part of its argument, running through all its values twice in each interval of \(2 \pi:\)

* Sine and cosecant begin their period at \(2 \pi k-\frac{\pi}{2}\) (where \(k\) is an integer), finish it at \(2 \pi k+\frac{\pi}{2},\) and then reverse themselves over \(2 \pi k+\frac{\pi}{2}\) to \(2 \pi k+\frac{3\pi}{2}.\)
* Cosine and secant begin their period at \(2 \pi k,\) finish it at \(2 \pi k+\pi.\) and then reverse themselves over \(2 \pi k+\pi\) to \(2 \pi k+2 \pi.\)
* Tangent begins its period at \(2 \pi k-\frac{\pi}{2},\)  finishes it at \(2 \pi k+\frac{\pi}{2},\) and then repeats it (forward) over \(2 \pi k+\frac{\pi}{2}\) to \(2 \pi k+\frac{3 \pi}{2}.\)
* Cotangent begins its period at \(2 \pi k,\) finishes it at \(2 \pi k+\pi,\) and then repeats it (forward) over \(2 \pi k+\pi\) to \(2 \pi k+2 \pi.\)

This periodicity is reflected in the general inverses, where \(k\) is some integer.

The following table shows how inverse trigonometric functions may be used to solve equalities involving the six standard trigonometric functions.
It is assumed that the given values \(\theta,\) \(r,\) \(s,\) \(x,\) and \(y\) all lie within appropriate ranges so that the relevant expressions below are well-defined.
Note that "for some \(k \in \Z\)" is just another way of saying "for some integer \(k.\)"

The symbol \(\,\iff\,\) is logical equality and indicates that if the left hand side is true then so is the right hand side and, conversely, if the right hand side is true then so is the left hand side (see this footnote for more details and an example illustrating this concept).

{| class="wikitable" style="border: none;"
|+
|-
! Equation !! if and only if !! colspan="7" | Solution
|-
| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\sin \theta = y\)
| style="text-align: center;" |\(\iff\)
| style='border-style: solid none solid none; text-align: left; padding-left: 2em;' |\(\theta =\,\)
| style='border-style: solid none solid none; text-align: right;' |\((-1)^k\)
| style='border-style: solid none solid none; text-align: left;' |\(\arcsin (y)\)
| style='border-style: solid none solid none;' |\(+\)
| style='border-style: solid none solid none;' |
| style='border-style: solid none solid none; padding-right: 2em;' |\(\pi k\)
| style="text-align: center; padding-left: 1em; padding-right: 1em;" |for some \(k \in \Z\)

|-

| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\csc \theta = r\)
| style="text-align: center;" |\(\iff\)
| style='border-style: solid none solid none; text-align: left; padding-left: 2em;' |\(\theta =\,\)
| style='border-style: solid none solid none; text-align: right;' |\((-1)^k\)
| style='border-style: solid none solid none; text-align: left;' |\(\arccsc (r)\)
| style='border-style: solid none solid none;' |\(+\)
| style='border-style: solid none solid none;' |
| style='border-style: solid none solid none; padding-right: 2em;' |\(\pi k\)
| style="text-align: center; padding-left: 1em; padding-right: 1em;" |for some \(k \in \Z\)

|-

| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\cos \theta = x\)
| style="text-align: center;" |\(\iff\)
| style='border-style: solid none solid none; text-align: left; padding-left: 2em;' |\(\theta =\,\)
| style='border-style: solid none solid none; text-align: right;' |\(\pm\,\)
| style='border-style: solid none solid none; text-align: left;' |\(\arccos(x)\)
| style='border-style: solid none solid none;' |\(+\)
| style='border-style: solid none solid none;' |\(2\)
| style='border-style: solid none solid none; padding-right: 2em;' |\(\pi k\)
| style="text-align: center; padding-left: 1em; padding-right: 1em;" |for some \(k \in \Z\)

|-

| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\sec \theta = r\)
| style="text-align: center;" |\(\iff\)
| style='border-style: solid none solid none; text-align: left; padding-left: 2em;' |\(\theta =\,\)
| style='border-style: solid none solid none; text-align: right;' |\(\pm\,\)
| style='border-style: solid none solid none; text-align: left;' |\(\arcsec (r)\)
| style='border-style: solid none solid none;' |\(+\)
| style='border-style: solid none solid none;' |\(2\)
| style='border-style: solid none solid none; padding-right: 2em;' |\(\pi k\)
| style="text-align: center; padding-left: 1em; padding-right: 1em;" |for some \(k \in \Z\)
|-

|-

| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\tan \theta = s\)
| style="text-align: center;" |\(\iff\)
| style='border-style: solid none solid none; text-align: left; padding-left: 2em;' |\(\theta =\,\)
| style='border-style: solid none solid none; text-align: right;' |
| style='border-style: solid none solid none; text-align: left;' |\(\arctan (s)\)
| style='border-style: solid none solid none;' |\(+\)
| style='border-style: solid none solid none;' |
| style='border-style: solid none solid none; padding-right: 2em;' |\(\pi k\)
| style="text-align: center; padding-left: 1em; padding-right: 1em;" |for some \(k \in \Z\)

|-

| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\cot \theta = r\)
| style="text-align: center;" |\(\iff\)
| style='border-style: solid none solid none; text-align: left; padding-left: 2em;' |\(\theta =\,\)
| style='border-style: solid none solid none; text-align: right;' |
| style='border-style: solid none solid none; text-align: left;' |\(\arccot (r)\)
| style='border-style: solid none solid none;' |\(+\)
| style='border-style: solid none solid none;' |
| style='border-style: solid none solid none; padding-right: 2em;' |\(\pi k\)
| style="text-align: center; padding-left: 1em; padding-right: 1em;" |for some \(k \in \Z\)

|}

where the first four solutions can be written in expanded form as:

{| class="wikitable" style="border: none;"
|+
|-
! Equation !! if and only if !! colspan="7" | Solution
|-
| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\sin \theta = y\)
| style="text-align: center;" |\(\iff\)
| style='border-style: solid none solid none; text-align: left; padding-left: 2em;' |\(\theta = \;\;\;\,\arcsin(y)+2 \pi k\) or \(\theta =-\arcsin(y)+2 \pi k+\pi\)
| style="text-align: center; padding-left: 1em; padding-right: 1em;" |for some \(k \in \Z\)

|-

| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\csc \theta = r\)
| style="text-align: center;" |\(\iff\)
| style='border-style: solid none solid none; text-align: left; padding-left: 2em;' |\(\theta = \;\;\;\,\arccsc(r)+2 \pi k\) or \(\theta =-\arccsc(r)+2 \pi k+\pi\)
| style="text-align: center; padding-left: 1em; padding-right: 1em;" |for some \(k \in \Z\)

|-

| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\cos \theta = x\)
| style="text-align: center;" |\(\iff\)
| style='border-style: solid none solid none; text-align: left; padding-left: 2em;' |\(\theta = \;\;\;\,\arccos(x)+2 \pi k\) or \(\theta =-\arccos(x)+2 \pi k\)
| style="text-align: center; padding-left: 1em; padding-right: 1em;" |for some \(k \in \Z\)

|-

| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\sec \theta = r\)
| style="text-align: center;" |\(\iff\)
| style='border-style: solid none solid none; text-align: left; padding-left: 2em;' |\(\theta = \;\;\;\,\arcsec(r)+2 \pi k\) or \(\theta =-\arcsec(r)+2 \pi k\)
| style="text-align: center; padding-left: 1em; padding-right: 1em;" |for some \(k \in \Z\)

|}

For example, if \(\cos \theta = -1\) then \(\theta = \pi+2 \pi k = -\pi+2 \pi (1+k)\) for some \(k \in \Z.\) While if \(\sin \theta = \pm 1\) then \(\theta = \frac{\pi}{2}+\pi k =-\frac{\pi}{2}+\pi (k+1)\) for some \(k \in \Z,\) where \(k\) will be even if \(\sin \theta = 1\) and it will be odd if \(\sin \theta = -1.\) The equations \(\sec \theta = -1\) and \(\csc \theta = \pm 1\) have the same solutions as \(\cos \theta = -1\) and \(\sin \theta = \pm 1,\) respectively. In all equations above  for those just solved (i.e. except for \(\sin\)/\(\csc \theta = \pm 1\) and \(\cos\)/\(\sec \theta =-1\)), the integer \(k\) in the solution's formula is uniquely determined by \(\theta\) (for fixed \(r, s, x,\) and \(y\)).

With the help of integer parity
\(\operatorname{Parity}(h) =
\begin{cases}
0 & \text{if } h \text{ is even } \\
1 & \text{if } h \text{ is odd } \\
\end{cases}\)
it is possible to write a solution to \(\cos \theta = x\) that doesn't involve the "plus or minus" \(\,\pm\,\) symbol:
\(cos \; \theta = x \quad\) if and only if \(\quad \theta = (-1)^h \arccos(x) + \pi h + \pi \operatorname{Parity}(h) \quad\) for some \(h \in \Z.\)
And similarly for the secant function,
\(sec \; \theta = r \quad\) if and only if \(\quad \theta = (-1)^h \arcsec(r) + \pi h + \pi \operatorname{Parity}(h) \quad\) for some \(h \in \Z,\)
where \(\pi h + \pi \operatorname{Parity}(h)\) equals \(\pi h\) when the integer \(h\) is even, and equals \(\pi h + \pi\) when it's odd.

#### Detailed example and explanation of the "plus or minus" symbol ±
The solutions to \(\cos \theta = x\) and \(\sec \theta = x\) involve the "plus or minus" symbol \(\,\pm,\,\) whose meaning is now clarified. Only the solution to \(\cos \theta = x\) will be discussed since the discussion for \(\sec \theta = x\) is the same.
We are given \(x\) between \(-1 \leq x \leq 1\) and we know that there is an angle \(\theta\) in some interval that satisfies \(\cos \theta = x.\) We want to find this \(\theta.\) The table above indicates that the solution is

\[
\,\theta = \pm \arccos x+2 \pi k\, \quad \text{ for some }k \in \Z
\]

which is a shorthand way of saying that (at least) one of the following statement is true:
#\(\,\theta = \arccos x+2 \pi k\,\) for some integer \(k,\) or
#\(\,\theta =-\arccos x+2 \pi k\,\) for some integer \(k.\)
As mentioned above, if \(\,\arccos x = \pi\,\) (which by definition only happens when \(x = \cos \pi = -1\)) then both statements (1) and (2) hold, although with different values for the integer \(k\): if \(K\) is the integer from statement (1), meaning that \(\theta = \pi+2 \pi K\) holds, then the integer \(k\) for statement (2) is \(K+1\) (because \(\theta = -\pi+2 \pi (1+K)\)).
However, if \(x \neq -1\) then the integer \(k\) is unique and completely determined by \(\theta.\)
If \(\,\arccos x = 0\,\)  (which by definition only happens when \(x = \cos 0 = 1\)) then \(\,\pm\arccos x = 0\,\) (because \(\,+ \arccos x = +0 = 0\,\) and \(\,-\arccos x = -0 = 0\,\) so in both cases \(\,\pm \arccos x\,\) is equal to \(0\)) and so the statements (1) and (2) happen to be identical in this particular case (and so both hold).
Having considered the cases \(\,\arccos x = 0\,\) and \(\,\arccos x = \pi,\,\) we now focus on the case where \(\,\arccos x \neq 0\,\) and \(\,\arccos x \neq \pi,\,\) So assume this from now on. The solution to \(\cos \theta = x\) is still

\[
\,\theta = \pm \arccos x+2 \pi k\, \quad \text{ for some }k \in \Z
\]

which as before is shorthand for saying that one of statements (1) and (2) is true. However this time, because \(\,\arccos x \neq 0\,\) and \(\,0 < \arccos x < \pi,\,\) statements (1) and (2) are different and furthermore, *exactly one* of the two equalities holds (not both). Additional information about \(\theta\) is needed to determine which one holds. For example, suppose that \(x = 0\) and that  that is known about \(\theta\) is that \(\,-\pi \leq \theta \leq \pi\,\) (and nothing more is known). Then

\[
\arccos x = \arccos 0 = \frac{\pi}{2}
\]

and moreover, in this particular case \(k = 0\) (for both the \(\,+\,\) case and the \(\,-\,\) case) and so consequently,

\[
\theta ~=~ \pm \arccos x+2 \pi k ~=~ \pm \left(\frac{\pi}{2}\right)+2\pi (0) ~=~ \pm \frac{\pi}{2}.
\]

This means that \(\theta\) could be either \(\,\pi/2\,\) or \(\,-\pi/2.\) Without additional information it is not possible to determine which of these values \(\theta\) has.
An example of some additional information that could determine the value of \(\theta\) would be knowing that the angle is above the \(x\)-axis (in which case \(\theta = \pi/2\)) or alternatively, knowing that it is below the \(x\)-axis (in which case \(\theta =-\pi/2\)).

#### Equal identical trigonometric functions


;Set of all solutions to elementary trigonometric equations

Thus given a single solution \(\theta\) to an elementary trigonometric equation (\(\sin \theta = y\) is such an equation, for instance, and because \(\sin (\arcsin y) = y\) always holds, \(\theta := \arcsin y\) is always a solution), the set of all solutions to it are:

{| class="wikitable" style="border: none;"
|+
|-
! If \(\theta\) solves !! then !! colspan="7" | Set of all solutions (in terms of \(\theta\))
|-
| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\;\sin \theta = y\)
| style="text-align: center;" |then
| style='border-style: solid none solid none; text-align: left; padding-left: 1em;' |\(\{\varphi:\sin \varphi=y\} =\,\)
| style='border-style: solid none solid none; text-align: right; padding: 0;' |\((\theta\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\,+\, 2\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\pi \Z)\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\,\cup\, (-\theta\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(-\pi\)
| style='border-style: solid solid solid none; text-align: left;  padding-left: 0; padding-right: 2em;' |\(+ 2 \pi \Z)\)

|-

| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\;\csc \theta = r\)
| style="text-align: center;" |then
| style='border-style: solid none solid none; text-align: left; padding-left: 1em;' |\(\{\varphi:\csc \varphi=r\} =\,\)
| style='border-style: solid none solid none; text-align: right; padding: 0;' |\((\theta\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\,+\, 2\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\pi \Z)\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\,\cup\, (-\theta\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(-\pi\)
| style='border-style: solid solid solid none; text-align: left;  padding-left: 0; padding-right: 2em;' |\(+ 2 \pi \Z)\)

|-

| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\;\cos \theta = x\)
| style="text-align: center;" |then
| style='border-style: solid none solid none; text-align: left; padding-left: 1em;' |\(\{\varphi:\cos \varphi=x\} =\,\)
| style='border-style: solid none solid none; text-align: right; padding: 0;' |\((\theta\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\,+\, 2\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\pi \Z)\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\,\cup\, (-\theta\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |
| style='border-style: solid solid solid none; text-align: left;  padding-left: 0; padding-right: 2em;' |\(+ 2 \pi \Z)\)

|-

| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\;\sec \theta = r\)
| style="text-align: center;" |then
| style='border-style: solid none solid none; text-align: left; padding-left: 1em;' |\(\{\varphi:\sec \varphi=r\} =\,\)
| style='border-style: solid none solid none; text-align: right; padding: 0;' |\((\theta\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\,+\, 2\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\pi \Z)\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\,\cup\, (-\theta\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |
| style='border-style: solid solid solid none; text-align: left;  padding-left: 0; padding-right: 2em;' |\(+ 2 \pi \Z)\)
|-

|-

| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\;\tan \theta = s\)
| style="text-align: center;" |then
| style='border-style: solid none solid none; text-align: left; padding-left: 1em;' |\(\{\varphi:\tan \varphi=s\} =\,\)
| style='border-style: solid none solid none; text-align: right; padding: 0;' |\(\theta\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\,+\,\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\pi \Z\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |
| style='border-style: solid solid solid none; text-align: left;  padding-left: 0; padding-right: 2em;' |

|-

| style="text-align: center; padding: 0.5% 2em 0.5% 2em;" | \(\;\cot \theta = r\)
| style="text-align: center;" |then
| style='border-style: solid none solid none; text-align: left; padding-left: 1em;' |\(\{\varphi:\cot \varphi=r\} =\,\)
| style='border-style: solid none solid none; text-align: right; padding: 0;' |\(\theta\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\,+\,\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |\(\pi \Z\)
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |
| style='border-style: solid none solid none; text-align: left;  padding: 0;' |
| style='border-style: solid solid solid none; text-align: left;  padding-left: 0; padding-right: 2em;' |

|}

### Transforming equations
The equations above can be transformed by using the reflection and shift identities:

{| class="wikitable" style="text-align: center;"
|+ Transforming equations by shifts and reflections
|-
! scope="col" | Argument: \(\underline{\;~~~~~~\;}=\)
! scope="col" |\(-\theta\)
! scope="col" |\(\frac{\pi}{2} \pm \theta\)
! scope="col" |\(\pi \pm \theta\)
! scope="col" |\(\frac{3\pi}{2} \pm \theta\)
! scope="col" |\(2 k \pi \pm \theta,\) \((k \in \Z)\)
|-
! scope="row" |\(\sin \underline{\;~~~~~~~~~~~~~~\;}=\)
| \(-\sin \theta\)
| \(\phantom{-}\cos \theta\)
| \(\mp\sin \theta\)
| \(-\cos \theta\)
| \(\pm\sin \theta\)
|-
! scope="row" |\(\csc \underline{\;~~~~~~~~~~~~~~\;}=\)
| \(-\csc \theta\)
| \(\phantom{-}\sec \theta\)
| \(\mp\csc \theta\)
| \(-\sec \theta\)
| \(\pm\csc \theta\)
|-
! scope="row" |\(\cos \underline{\;~~~~~~~~~~~~~~\;}=\)
| \(\phantom{-}\cos \theta\)
| \(\mp\sin \theta\)
| \(-\cos \theta\)
| \(\pm\sin \theta\)
| \(\phantom{-}\cos \theta\)
|-
! scope="row" |\(\sec \underline{\;~~~~~~~~~~~~~~\;}=\)
| \(\phantom{-}\sec \theta\)
| \(\mp\csc \theta\)
| \(-\sec \theta\)
| \(\pm\csc \theta\)
| \(\phantom{-}\sec \theta\)
|-
! scope="row" |\(\tan \underline{\;~~~~~~~~~~~~~~\;}=\)
| \(-\tan \theta\)
| \(\mp\cot \theta\)
| \(\pm\tan \theta\)
| \(\mp\cot \theta\)
| \(\pm\tan \theta\)
|-
! scope="row" |\(\cot \underline{\;~~~~~~~~~~~~~~\;}=\)
| \(-\cot \theta\)
| \(\mp\tan \theta\)
| \(\pm\cot \theta\)
| \(\mp\tan \theta\)
| \(\pm\cot \theta\)
|}

These formulas imply, in particular, that the following hold:

\(\begin{align}
\sin \theta
&= -\sin(-\theta)
&&= -\sin(\pi+\theta)
&&= \phantom{-}\sin(\pi-\theta) \\
&= -\cos\left(\frac{\pi}{2}+\theta\right)
&&= \phantom{-}\cos\left(\frac{\pi}{2}-\theta\right)
&&= -\cos\left(-\frac{\pi}{2}-\theta\right) \\
&= \phantom{-}\cos\left(-\frac{\pi}{2}+\theta\right)
&&= -\cos\left(\frac{3\pi}{2}-\theta\right)
&&= -\cos\left(-\frac{3\pi}{2}+\theta\right)
\\[0.3ex]
\cos \theta
&= \phantom{-}\cos(-\theta)
&&= -\cos(\pi+\theta)
&&= -\cos(\pi-\theta) \\
&= \phantom{-}\sin\left(\frac{\pi}{2}+\theta\right)
&&= \phantom{-}\sin\left(\frac{\pi}{2}-\theta\right)
&&= -\sin\left(-\frac{\pi}{2}-\theta\right) \\
&= -\sin\left(-\frac{\pi}{2}+\theta\right)
&&= -\sin\left(\frac{3\pi}{2}-\theta\right)
&&= \phantom{-}\sin\left(-\frac{3\pi}{2}+\theta\right)
\\[0.3ex]
\tan \theta
&= -\tan(-\theta)
&&= \phantom{-}\tan(\pi+\theta)
&&= -\tan(\pi-\theta) \\
&= -\cot\left(\frac{\pi}{2}+\theta\right)
&&= \phantom{-}\cot\left(\frac{\pi}{2}-\theta\right)
&&= \phantom{-}\cot\left(-\frac{\pi}{2}-\theta\right) \\
&= -\cot\left(-\frac{\pi}{2}+\theta\right)
&&= \phantom{-}\cot\left(\frac{3\pi}{2}-\theta\right)
&&= -\cot\left(-\frac{3\pi}{2}+\theta\right)
\\[0.3ex]
\end{align}\)

where swapping \(\sin \leftrightarrow \csc,\) swapping \(\cos \leftrightarrow \sec,\) and swapping \(\tan \leftrightarrow \cot\) gives the analogous equations for \(\csc, \sec, \text{ and } \cot,\) respectively.

So for example, by using the equality \(\sin \left(\frac{\pi}{2}-\theta\right) = \cos \theta,\) the equation \(\cos \theta = x\) can be transformed into \(\sin \left(\frac{\pi}{2}-\theta\right) = x,\) which allows for the solution to the equation \(\;\sin \varphi = x\;\) (where \(\varphi := \frac{\pi}{2}-\theta\)) to be used; that solution being:
\(\varphi  = (-1)^k \arcsin (x)+\pi k \; \text{ for some } k \in \Z,\)
which becomes:

\[
\frac{\pi}{2}-\theta ~=~ (-1)^k \arcsin (x)+\pi k \quad \text{ for some } k \in \Z
\]

where using the fact that \((-1)^{k} = (-1)^{-k}\) and substituting \(h :=-k\) proves that another solution to \(\;\cos \theta = x\;\) is:

\[
\theta ~=~ (-1)^{h+1} \arcsin (x)+\pi h+\frac{\pi}{2} \quad \text{ for some } h \in \Z.
\]

The substitution \(\;\arcsin x = \frac{\pi}{2}-\arccos x\;\) may be used express the right hand side of the above formula in terms of \(\;\arccos x\;\) instead of \(\;\arcsin x.\;\)

### Relationships between trigonometric functions and inverse trigonometric functions
Trigonometric functions of inverse trigonometric functions are tabulated below. They may be derived from the Pythagorean identities. Another way is by considering the geometry of a right-angled triangle, with one side of length 1 and another side of length \(x,\) then applying the Pythagorean theorem and definitions of the trigonometric ratios. It is worth noting that for arcsecant and arccosecant, the diagram assumes that \(x\) is positive, and thus the result has to be corrected through the use of absolute values and the signum (sgn) operation.

{|class="wikitable"
|-
!\(\theta\)
!\(\sin(\theta)\)
!\(\cos(\theta)\)
!\(\tan(\theta)\)
!Diagram
|-
!\(\arcsin(x)\)
|\(\sin(\arcsin(x)) = x\)
|\(\cos(\arcsin(x)) = \sqrt{1-x^2}\)
|\(\tan(\arcsin(x)) = \frac{x}{\sqrt{1-x^2}}\)
|
|-
!\(\arccos(x)\)
|\(\sin(\arccos(x)) = \sqrt{1-x^2}\)
|\(\cos(\arccos(x)) = x\)
|\(\tan(\arccos(x)) = \frac{\sqrt{1-x^2}}{x}\)
|
|-
!\(\arctan(x)\)
|\(\sin(\arctan(x)) = \frac{x}{\sqrt{1+x^2}}\)
|\(\cos(\arctan(x)) = \frac{1}{\sqrt{1+x^2}}\)
|\(\tan(\arctan(x)) = x\)
|
|-
!\(\arccot(x)\)
|\(\sin(\arccot(x)) = \frac{1}{\sqrt{1+x^2}}\)
|\(\cos(\arccot(x)) = \frac{x}{\sqrt{1+x^2}}\)
|\(\tan(\arccot(x)) = \frac{1}{x}\)
|
|-
!\(\arcsec(x)\)
|\(\sin(\arcsec(x)) = \sqrt{1-\frac{1}{x^2}}\)
|\(\cos(\arcsec(x)) = \frac{1}{x}\)
|\(\tan(\arcsec(x)) = \sgn(x)\sqrt{x^2-1}\)
|
|-
!\(\arccsc(x)\)
|\(\sin(\arccsc(x)) = \frac{1}{x}\)
|\(\cos(\arccsc(x)) = \sqrt{1-\frac{1}{x^2}}\)
|\(\tan(\arccsc(x)) = \frac{\sgn(x)}{\sqrt{x^2-1}}\)
|
|-
|}

### Relationships among the inverse trigonometric functions


Complementary angles:
\(\begin{align}
\arccos(x) &= \frac{\pi}{2} - \arcsin(x) \\[0.5em]
\arccot(x) &= \frac{\pi}{2} - \arctan(x) \\[0.5em]
\arccsc(x) &= \frac{\pi}{2} - \arcsec(x)
\end{align}\)

Negative arguments:
\(\begin{align}
\arcsin(-x) &= -\arcsin(x) \\
\arccsc(-x) &= -\arccsc(x) \\
\arccos(-x) &= \pi -\arccos(x) \\
\arcsec(-x) &= \pi -\arcsec(x) \\
\arctan(-x) &= -\arctan(x) \\
\arccot(-x) &= \pi -\arccot(x)
\end{align}\)

Reciprocal arguments:
\(\begin{align}
\arcsin\left(\frac{1}{x}\right) &= \arccsc(x) & \\[0.3em]
\arccsc\left(\frac{1}{x}\right) &= \arcsin(x) & \\[0.3em]
\arccos\left(\frac{1}{x}\right) &= \arcsec(x) & \\[0.3em]
\arcsec\left(\frac{1}{x}\right) &= \arccos(x) & \\[0.3em]
\arctan\left(\frac{1}{x}\right) &= \arccot(x)       &= \frac{\pi}{2} - \arctan(x) \, , \text{ if } x > 0 \\[0.3em]
\arctan\left(\frac{1}{x}\right) &= \arccot(x) - \pi &= -\frac{\pi}{2} - \arctan(x) \, , \text{ if } x < 0 \\[0.3em]
\arccot\left(\frac{1}{x}\right) &= \arctan(x)       &= \frac{\pi}{2} - \arccot(x) \, , \text{ if } x > 0 \\[0.3em]
\arccot\left(\frac{1}{x}\right) &= \arctan(x) + \pi &= \frac{3\pi}{2} - \arccot(x) \, , \text{ if } x < 0
\end{align}\)
The identities above can be used with (and derived from) the fact that \(\sin\) and \(\csc\) are reciprocals (i.e. \(\csc = \tfrac1{\sin}\)), as are \(\cos\) and \(\sec,\) and \(\tan\) and \(\cot.\)

Useful identities if one only has a fragment of a sine table:
\(\begin{align}
\arcsin(x) &= \frac{1}{2}\arccos\left(1-2x^2\right) \, , \text{ if } 0 \leq x \leq 1 \\
\arcsin(x) &= \arctan\left(\frac{x}{\sqrt{1 - x^2}}\right) \\
\arccos(x) &= \frac{1}{2}\arccos\left(2x^2-1\right) \, , \text{ if } 0 \leq x \leq 1 \\
\arccos(x) &= \arctan\left(\frac{\sqrt{1 - x^2}}{x}\right) \\
\arccos(x) &= \arcsin\left(\sqrt{1 - x^2}\right) \, , \text{ if } 0 \leq x \leq 1 \text{ , from which you get } \\
\arccos    &\left(\frac{1-x^2}{1 + x^2}\right) = \arcsin \left (\frac{2x}{1 + x^2}\right) \, , \text{ if } 0 \leq x \leq 1 \\
\arcsin    &\left(\sqrt{1 - x^2}\right) =\frac{\pi}{2}-\sgn(x)\arcsin(x) \\
\arctan(x) &= \arcsin\left(\frac{x}{\sqrt{1 + x^2}}\right) \\
\arccot(x) &= \arccos\left(\frac{x}{\sqrt{1 + x^2}}\right)
\end{align}\)

Whenever the square root of a complex number is used here, we choose the root with the positive real part (or positive imaginary part if the square was negative real).

A useful form that follows directly from the table above is

\(\arctan(x) = \arccos\left(\sqrt{\frac{1}{1+x^2}}\right)\, , \text{ if } x\geq 0\).

It is obtained by recognizing that \(\cos\left(\arctan\left(x\right)\right) = \sqrt{\frac{1}{1+x^2}} = \cos\left(\arccos\left(\sqrt{\frac{1}{1+x^2}}\right)\right)\).

From the half-angle formula, \(\tan\left(\tfrac{\theta}{2}\right) = \tfrac{\sin(\theta)}{1 + \cos(\theta)}\), we get:
\(\begin{align}
\arcsin(x) &= 2 \arctan\left(\frac{x}{1 + \sqrt{1 - x^2}}\right) \\[0.5em]
\arccos(x) &= 2 \arctan\left(\frac{\sqrt{1 - x^2}}{1 + x}\right) \, , \text{ if } -1 < x \leq  1 \\[0.5em]
\arctan(x) &= 2 \arctan\left(\frac{x}{1 + \sqrt{1 + x^2}}\right)
\end{align}\)

### Arctangent addition formula
\(\arctan(u) \pm \arctan(v) = \arctan\left(\frac{u \pm v}{1 \mp uv}\right) \pmod \pi \, , \quad u v \ne 1 \, .\)
This is derived from the tangent addition formula
\(\tan(\alpha \pm \beta) = \frac{\tan(\alpha) \pm \tan(\beta)}{1 \mp \tan(\alpha)\tan(\beta)} \, ,\)
by letting
\(\alpha = \arctan(u) \, , \quad \beta = \arctan(v) \, .\)

## In calculus
### Derivatives of inverse trigonometric functions


The derivatives for complex values of *z* are as follows:
\(\begin{align}
\frac{d}{dz} \arcsin(z) &{} = \frac{1}{\sqrt{1-z^2}} \; ;   &z &{}\neq -1, +1 \\
\frac{d}{dz} \arccos(z) &{} = -\frac{1}{\sqrt{1-z^2}} \; ;  &z &{}\neq -1, +1 \\
\frac{d}{dz} \arctan(z) &{} = \frac{1}{1+z^2} \; ;          &z &{}\neq -i, +i\\
\frac{d}{dz} \arccot(z) &{} = -\frac{1}{1+z^2} \; ;         &z &{}\neq -i, +i \\
\frac{d}{dz} \arcsec(z) &{} = \frac{1}{z^2 \sqrt{1 - \frac{1}{z^{2}}}} \; ;   &z &{}\neq -1, 0, +1 \\
\frac{d}{dz} \arccsc(z) &{} = -\frac{1}{z^2 \sqrt{1 - \frac{1}{z^{2}}}} \; ;  &z &{}\neq -1, 0, +1
\end{align}\)
Only for real values of *x*:
\(\begin{align}
\frac{d}{dx} \arcsec(x) &{} = \frac{1}{|x| \sqrt{x^2-1}} \; ;  & |x| > 1\\
\frac{d}{dx} \arccsc(x) &{} = -\frac{1}{|x| \sqrt{x^2-1}} \; ; & |x| > 1
\end{align}\)

These formulas can be derived in terms of the derivatives of trigonometric functions. For example, if \(x = \sin \theta\), then \(dx/d\theta = \cos \theta = \sqrt{1-x^2},\) so
\(\frac{d}{dx}\arcsin(x) = \frac{d \theta}{dx} = \frac{1}{dx/d\theta} = \frac{1}{\sqrt{1-x^2}}.\)

### Expression as definite integrals
Integrating the derivative and fixing the value at one point gives an expression for the inverse trigonometric function as a definite integral:
\(\begin{align}
\arcsin(x) &{}= \int_0^x \frac{1}{\sqrt{1 - z^2}} \, dz \; , & |x| &{} \leq 1\\
\arccos(x) &{}= \int_x^1 \frac{1}{\sqrt{1 - z^2}} \, dz \; , & |x| &{} \leq 1\\
\arctan(x) &{}= \int_0^x \frac{1}{z^2 + 1} \, dz \; ,\\
\arccot(x) &{}= \int_x^\infty \frac{1}{z^2 + 1} \, dz \; ,\\
\arcsec(x) &{}= \int_1^x \frac{1}{z \sqrt{z^2 - 1}} \, dz = \pi + \int_{-x}^{-1} \frac{1}{z \sqrt{z^2 - 1}} \, dz\; , & x &{} \geq 1\\
\arccsc(x) &{}= \int_x^\infty \frac{1}{z \sqrt{z^2 - 1}} \, dz = \int_{-\infty}^{-x} \frac{1}{z \sqrt{z^2 - 1}} \, dz \; , & x &{} \geq 1\\
\end{align}\)
When *x* equals 1, the integrals with limited domains are improper integrals, but still well-defined.

### Infinite series
Similar to the sine and cosine functions, the inverse trigonometric functions can also be calculated using power series, as follows. For arcsine, the series can be derived by expanding its derivative, \(\tfrac{1}{\sqrt{1-z^2}}\), as a binomial series, and integrating term by term (using the integral definition as above). The series for arctangent can similarly be derived by expanding its derivative \(\frac{1}{1+z^2}\) in a geometric series, and applying the integral definition above (see Leibniz series).

\(\begin{align}
\arcsin(z) & = z + \left( \frac{1}{2} \right) \frac{z^3}{3} + \left( \frac{1 \cdot 3}{2 \cdot 4} \right) \frac{z^5}{5} + \left( \frac{1 \cdot 3 \cdot 5}{2 \cdot 4 \cdot 6} \right) \frac{z^7}{7} + \cdots \\[5pt]
& = \sum_{n=0}^\infty \frac{(2n-1)!!}{(2n)!!}\frac{z^{2n+1}}{2n+1} \\[5pt]
& = \sum_{n=0}^\infty \frac{(2n)!}{(2^n n!)^2} \frac{z^{2n+1}}{2n+1} \, ; \qquad |z| \le 1
\end{align}\)
The Taylor series for arctangent is also known as the arctangent series or Gregory's series.
\(\arctan(z)
= z - \frac{z^3}{3} +\frac{z^5}{5} - \frac{z^7}{7} + \cdots
= \sum_{n=0}^\infty \frac{(-1)^n z^{2n+1}}{2n+1} \, ; \qquad |z| \le 1 \qquad z \neq i,-i\)

Series for the other inverse trigonometric functions can be given in terms of these according to the relationships given above.  For example, \(\arccos(x) = \pi/2 - \arcsin(x)\), \(\arccsc(x) = \arcsin(1/x)\), and so on.  Another series is given by:

\(2\left(\arcsin\left(\frac{x}{2}\right) \right)^2 = \sum_{n=1}^\infty \frac{x^{2n}}{n^2\binom {2n} n}.\)

Leonhard Euler found a series for the arctangent that converges more quickly than its Taylor series:

\(\arctan(z) = \frac z {1 + z^2} \sum_{n=0}^\infty \prod_{k=1}^n \frac{2k z^2}{(2k + 1)(1 + z^2)}.\)
(The term in the sum for *n* = 0 is the empty product, so is 1.)

Alternatively, this can be expressed as

\(\arctan(z) = \sum_{n=0}^\infty \frac{2^{2n} (n!)^2}{(2n + 1)!} \frac{z^{2n + 1}}{(1 + z^2)^{n + 1}}.\)

Another series for the arctangent function is given by

\(\arctan(z) = i\sum_{n=1}^\infty\frac{1}{2n - 1}\left(\frac{1}{(1 + 2i/z)^{2n-1}} - \frac{1}{(1 - 2i/z)^{2n - 1}}\right),\)

where \(i=\sqrt{-1}\) is the imaginary unit.

#### Continued fractions for arctangent
Two alternatives to the power series for arctangent are these generalized continued fractions:

\(\arctan(z) =
\frac z {1 + \cfrac{(1z)^2}{3 - 1z^2 + \cfrac{(3z)^2}{5 - 3z^2 + \cfrac{(5z)^2}{7 - 5z^2 + \cfrac{(7z)^2}{9-7z^2 + \ddots}}}}} =
\frac{z}{1 + \cfrac{(1z)^2}{3 + \cfrac{(2z)^2}{5 + \cfrac{(3z)^2}{7 + \cfrac{(4z)^2}{9 + \ddots}}}}}\)

The second of these is valid in the cut complex plane. There are two cuts, from −**i** to the point at infinity, going down the imaginary axis, and from **i** to the point at infinity, going up the same axis. It works best for real numbers running from −1 to 1. The partial denominators are the odd natural numbers, and the partial numerators (after the first) are just (*nz*)<sup>2</sup>, with each perfect square appearing once. The first was developed by Leonhard Euler, the second by Carl Friedrich Gauss utilizing the Gaussian hypergeometric series.

### Indefinite integrals of inverse trigonometric functions
For real and complex values of *z*:
\(\begin{align}
\int \arcsin(z) \, dz &{}= z \, \arcsin(z) + \sqrt{1 - z^2} + C\\
\int \arccos(z) \, dz &{}= z \, \arccos(z) - \sqrt{1 - z^2} + C\\
\int \arctan(z) \, dz &{}= z \, \arctan(z) - \frac{1}{2} \ln \left( 1 + z^2 \right) + C\\
\int \arccot(z) \, dz &{}= z \, \arccot(z) + \frac{1}{2} \ln \left( 1 + z^2 \right) + C\\
\int \arcsec(z) \, dz &{}= z \, \arcsec(z) - \ln \left[ z \left( 1 + \sqrt{ \frac{z^2-1}{z^2} } \right) \right] + C\\
\int \arccsc(z) \, dz &{}= z \, \arccsc(z) + \ln \left[ z \left( 1 + \sqrt{ \frac{z^2-1}{z^2} } \right) \right] + C
\end{align}\)

For real *x* ≥ 1:
\(\begin{align}
\int \arcsec(x) \, dx &{}= x \, \arcsec(x) - \ln \left( x + \sqrt{x^2-1} \right) + C\\
\int \arccsc(x) \, dx &{}= x \, \arccsc(x) + \ln \left( x + \sqrt{x^2-1} \right) + C
\end{align}\)

For all real *x* not between -1 and 1:
\(\begin{align}
\int \arcsec(x) \, dx &{}= x \, \arcsec(x) - \sgn(x) \ln\left| x + \sqrt{x^2-1}\right| + C\\
\int \arccsc(x) \, dx &{}= x \, \arccsc(x) + \sgn(x) \ln\left| x + \sqrt{x^2-1}\right| + C
\end{align}\)

The absolute value is necessary to compensate for both negative and positive values of the arcsecant and arccosecant functions. The signum function is also necessary due to the absolute values in the derivatives of the two functions, which create two different solutions for positive and negative values of x. These can be further simplified using the logarithmic definitions of the inverse hyperbolic functions:
\(\begin{align}
\int \arcsec(x) \, dx &{}= x \, \arcsec(x) - \operatorname{arcosh}(|x|) + C\\
\int \arccsc(x) \, dx &{}= x \, \arccsc(x) + \operatorname{arcosh}(|x|) + C\\
\end{align}\)

The absolute value in the argument of the arcosh function creates a negative half of its graph, making it identical to the signum logarithmic function shown above.

All of these antiderivatives can be derived using integration by parts and the simple derivative forms shown above.

#### Example
Using \(\int u \, dv = u v - \int v \, du\) (i.e. integration by parts), set

\(\begin{align}
u &= \arcsin(x) & dv &= dx \\
du &= \frac{dx}{\sqrt{1-x^2}} & v &= x
\end{align}\)

Then

\(\int \arcsin(x) \, dx = x \arcsin(x) - \int \frac{x}{\sqrt{1-x^2}} \, dx,\)

which by the simple substitution \(w=1-x^2,\ dw = -2x\,dx\) yields the final result:

\(\int \arcsin(x) \, dx = x \arcsin(x) + \sqrt{1-x^2} + C\)

## Extension to the complex plane


Since the inverse trigonometric functions are analytic functions, they can be extended from the real line to the complex plane. This results in functions with multiple sheets and branch points. One possible way of defining the extension is:
\(\arctan(z) = \int_0^z \frac{dx}{1 + x^2} \quad z \neq -i, +i\)
where the part of the imaginary axis which does not lie strictly between the branch points (−i and +i) is the branch cut between the principal sheet and other sheets. The path of the integral must not cross a branch cut. For *z* not on a branch cut, a straight line path from 0 to *z* is such a path. For *z* on a branch cut, the path must approach from Re[x] > 0 for the upper branch cut and from Re[x] < 0 for the lower branch cut.

The arcsine function may then be defined as:

\(\arcsin(z) = \arctan\left(\frac{z}{\sqrt{1 - z^2}}\right) \quad z \neq -1, +1\)
where (the square-root function has its cut along the negative real axis and) the part of the real axis which does not lie strictly between −1 and +1 is the branch cut between the principal sheet of arcsin and other sheets;
\(\arccos(z) = \frac{\pi}{2} - \arcsin(z) \quad z \neq -1, +1\)
which has the same cut as arcsin;
\(\arccot(z) = \frac{\pi}{2} - \arctan(z) \quad z \neq -i, i\)
which has the same cut as arctan;
\(\arcsec(z) = \arccos\left(\frac{1}{z}\right) \quad z \neq -1, 0, +1\)
where the part of the real axis between −1 and +1 inclusive is the cut between the principal sheet of arcsec and other sheets;
\(\arccsc(z) = \arcsin\left(\frac{1}{z}\right) \quad z \neq -1, 0, +1\)
which has the same cut as arcsec.

### Logarithmic forms
These functions may also be expressed using complex logarithms. This extends their domains to the complex plane in a natural fashion. The following identities for principal values of the functions hold everywhere that they are defined, even on their branch cuts.

\(\begin{align}
\arcsin(z) &{}= -i \ln \left( \sqrt{1-z^2} + iz \right) = i \ln \left( \sqrt{1-z^2} - iz \right) &{}= \arccsc\left(\frac{1}{z}\right) \\[10pt]
\arccos(z) &{}= -i \ln \left( z + i \sqrt {1-z^2} \right) = i \ln \left( z - i \sqrt {1-z^2} \right) = \frac{\pi}{2} - \arcsin(z) &{}= \arcsec\left(\frac{1}{z}\right) \\[10pt]
\arctan(z) &{}= -\frac{i}{2}\ln \left(\frac{i - z}{i + z}\right) = -\frac{i}{2}\ln \left(\frac{1 + iz}{1 - iz}\right) &{}= \arccot\left(\frac{1}{z}\right) \\[10pt]
\arccot(z) &{}= -\frac{i}{2}\ln\left( \frac{z + i}{z - i} \right) = -\frac{i}{2}\ln\left( \frac{iz - 1}{iz + 1} \right) &{}= \arctan\left(\frac{1}{z}\right) \\[10pt]
\arcsec(z) &{}= -i \ln \left( \frac{1}{z} + i \sqrt{1 - \frac{1}{z^2}} \right) = i \ln \left( \frac{1}{z} - i \sqrt{1 - \frac{1}{z^2}} \right) = \frac{\pi}{2} - \arccsc(z) &{}= \arccos\left(\frac{1}{z}\right) \\[10pt]
\arccsc(z) &{}= -i \ln \left( \sqrt{1 - \frac{1}{z^2}} + \frac{i}{z} \right) = i \ln \left( \sqrt{1 - \frac{1}{z^2}} - \frac{i}{z} \right) &{}= \arcsin\left(\frac{1}{z}\right)
\end{align}\)

#### Generalization
Because all of the inverse trigonometric functions output an angle of a right triangle, they can be generalized by using Euler's formula to form a right triangle in the complex plane. Algebraically, this gives us:

\(ce^{i\theta} = c\cos(\theta) + ic\sin(\theta)\)
or
\(ce^{i\theta} = a + ib\)

where \(a\) is the adjacent side, \(b\) is the opposite side, and \(c\) is the hypotenuse. From here, we can solve for \(\theta\).

\(\begin{align}
e^{\ln(c) + i\theta} & = a + ib \\
\ln c + i\theta & = \ln(a + ib) \\
\theta & = \operatorname{Im}\left( \ln(a + ib) \right)
\end{align}\)
or
\(\theta = -i\ln\left(\frac{a + ib}{c}\right)\)

Simply taking the imaginary part works for any real-valued \(a\) and \(b\), but if \(a\) or \(b\) is complex-valued, we have to use the final equation so that the real part of the result isn't excluded. Since the length of the hypotenuse doesn't change the angle, ignoring the real part of \(\ln(a+bi)\) also removes \(c\) from the equation. In the final equation, we see that the angle of the triangle in the complex plane can be found by inputting the lengths of each side. By setting one of the three sides equal to 1 and one of the remaining sides equal to our input \(z\), we obtain a formula for one of the inverse trig functions, for a total of six equations. Because the inverse trig functions require only one input, we must put the final side of the triangle in terms of the other two using the Pythagorean theorem relation
\(a^2 + b^2 = c^2\)

The table below shows the values of a, b, and c for each of the inverse trig functions and the equivalent expressions for \(\theta\) that result from plugging the values into the equations \(\theta = -i\ln\left(\tfrac{a + ib}{c}\right)\) above and simplifying.
\(\begin{align}
 & a & & b & & c && -i\ln\left(\frac{a+ib}{c}\right) && \theta && \theta_{a,b\in\R}\\
\arcsin(z)\ \ & \sqrt{1 - z^2} & & z & & 1 & & -i\ln\left( \frac{\sqrt{1 - z^2} + iz}{1} \right) && = -i\ln\left( \sqrt{1 - z^2} + iz \right) && \operatorname{Im}\left(\ln\left( \sqrt{1 - z^2} + iz \right)\right) \\
\arccos(z)\ \ & z & & \sqrt{1 - z^2} & & 1 & & -i\ln\left( \frac{z + i\sqrt{1 - z^2}}{1} \right) && = -i\ln\left( z + \sqrt{z^2 - 1} \right) && \operatorname{Im}\left(\ln\left( z + \sqrt{z^2 - 1} \right)\right) \\
\arctan(z)\ \ & 1 & & z & & \sqrt{1 + z^2} & & -i\ln\left( \frac{1 + iz}{\sqrt{1 + z^2}} \right) && = -\frac{i}{2}\ln\left( \frac{i-z}{i+z} \right) && \operatorname{Im}\left(\ln\left( 1 + iz \right)\right) \\
\arccot(z)\ \ & z & & 1 & & \sqrt{z^2 + 1} & & -i\ln\left( \frac{z + i}{\sqrt{z^2 + 1}} \right) && = -\frac{i}{2}\ln\left( \frac{z + i}{z-i} \right) && \operatorname{Im}\left(\ln\left( z + i \right)\right) \\
\arcsec(z)\ \ & 1 & & \sqrt{z^2 - 1} & & z & & -i\ln\left( \frac{1 + i\sqrt{z^2 - 1}}{z} \right) && = -i\ln\left( \frac{1}{z} + \sqrt{\frac{1}{z^2}-1} \right) && \operatorname{Im}\left(\ln\left( \frac{1}{z} + \sqrt{\frac{1}{z^2}-1} \right)\right) \\
\arccsc(z)\ \ & \sqrt{z^2 - 1} & & 1 & & z & & -i\ln\left( \frac{\sqrt{z^2 - 1} + i}{z} \right) && = -i\ln\left( \sqrt{1-\frac{1}{z^2}} + \frac{i}{z} \right) && \operatorname{Im}\left(\ln\left( \sqrt{1 - \frac{1}{z^2}} + \frac{i}{z} \right)\right) \\
\end{align}\)

The particular form of the simplified expression can cause the output to differ from the usual principal branch of each of the inverse trig functions. The formulations given will output the usual principal branch when using the \(\operatorname{Im}\left( \ln z \right) \in (-\pi,\pi]\) and \(\operatorname{Re}\left(\sqrt{z}\right) \ge 0\) principal branch for every function except arccotangent in the \(\theta\) column. Arccotangent in the \(\theta\) column will output on its usual principal branch by using the \(\operatorname{Im}\left( \ln z \right) \in [0,2\pi)\) and \(\operatorname{Im}\left(\sqrt{z}\right) \ge 0\) convention.

In this sense, all of the inverse trig functions can be thought of as specific cases of the complex-valued log function. Since these definition work for any complex-valued \(z\), the definitions allow for hyperbolic angles as outputs and can be used to further define the inverse hyperbolic functions. It's possible to algebraically prove these relations by starting with the exponential forms of the trigonometric functions and solving for the inverse function.

#### Example proof
\(\begin{align}
  \sin(\phi) &= z \\
        \phi &= \arcsin(z)
\end{align}\)

Using the exponential definition of sine, and letting \(\xi = e^{i \phi},\)

\(\begin{align}
z &= \frac{e^{i \phi} - e^{-i \phi}}{2i} \\[10mu]
2iz &= \xi - \frac{1}{\xi} \\[5mu]
0 &= \xi^2 - 2i z \xi - 1 \\[5mu]
\xi &= iz \pm \sqrt{1 - z^2} \\[5mu]
\phi &= -i \ln \left(iz \pm \sqrt{1 - z^2}\right)
\end{align}\)

(the positive branch is chosen)

\(\phi= \arcsin(z) = -i \ln \left(iz + \sqrt{1-z^2} \right)\)

{| style="text-align:center;"
 |+ Color wheel graphs of **inverse trigonometric functions in the complex plane**
 |
 |
 |
 |-
 | \(\arcsin(z)\)
 | \(\arccos(z)\)
 | \(\arctan(z)\)
 |}
{| style="text-align:center;"
 |+
 |
 |
 |
 |-
 | \(\arccsc(z)\)
 | \(\arcsec(z)\)
 | \(\arccot(z)\)
 |}

## Applications
### Finding the angle of a right triangle

Inverse trigonometric functions are useful when trying to determine the remaining two angles of a right triangle when the lengths of the sides of the triangle are known. Recalling the right-triangle definitions of sine and cosine, it follows that

\(\theta = \arcsin \left( \frac{\text{opposite}}{\text{hypotenuse}} \right) = \arccos \left( \frac{\text{adjacent}}{\text{hypotenuse}} \right) .\)

Often, the hypotenuse is unknown and would need to be calculated before using arcsine or arccosine using the Pythagorean theorem: \(a^2+b^2=h^2\) where \(h\) is the length of the hypotenuse. Arctangent comes in handy in this situation, as the length of the hypotenuse is not needed.

\(\theta = \arctan \left( \frac{\text{opposite}}{\text{adjacent}} \right) \, .\)

For example, suppose a roof drops 8 feet as it runs out 20 feet. The roof makes an angle *θ* with the horizontal, where *θ* may be computed as follows:

\(\theta
= \arctan \left( \frac{\text{opposite}}{\text{adjacent}} \right)
= \arctan \left( \frac{\text{rise}}{\text{run}} \right)
= \arctan \left( \frac{8}{20} \right) \approx 21.8^{\circ} \, .\)

### In computer science and engineering
#### Two-argument variant of arctangent


The two-argument atan2 function computes the arctangent of *y*/*x* given  and , but with a range of . In other words, atan2(*y*, *x*) is the angle between the positive -axis of a plane and the point (*x*, *y*) on it, with positive sign for counter-clockwise angles (upper half-plane, *y* > 0), and negative sign for clockwise angles (lower half-plane, *y* < 0). It was first introduced in many computer programming languages, but it is now also common in other fields of science and engineering.

In terms of the standard **arctan** function, that is with range of , it can be expressed as follows:


\[
\operatorname{atan2}(y, x) = \begin{cases}
  \arctan\left(\frac y x\right)       & \quad x > 0 \\
  \arctan\left(\frac y x\right) + \pi & \quad y \ge 0,\; x < 0 \\
  \arctan\left(\frac y x\right) - \pi & \quad y < 0,\; x < 0 \\
   \frac{\pi}{2}   & \quad y > 0,\; x = 0 \\
  -\frac{\pi}{2}   & \quad y < 0,\; x = 0 \\
  \text{undefined} & \quad y = 0,\; x = 0
\end{cases}
\]


It also equals the principal value of the argument of the complex number *x* + *iy*.

This limited version of the function above may also be defined using the tangent half-angle formulae as follows:

\[
\operatorname{atan2}(y, x) = 2\arctan\left(\frac{y}{\sqrt{x^2 + y^2} + x}\right)
\]

provided that either *x* > 0 or *y* ≠ 0. However this fails if given *x* ≤ 0 and 1= *y* = 0 so the expression is unsuitable for computational use.

The above argument order (, ) seems to be the most common, and in particular is used in ISO standards such as the C programming language, but a few authors may use the opposite convention (, ) so some caution is warranted. .)}}

#### Arctangent function with location parameter
In many applications the solution \(y\) of the equation \(x=\tan(y)\) is to come as close as possible to a given value \(-\infty < \eta < \infty\). The adequate solution is produced by the parameter modified arctangent function
\(y = \arctan_\eta(x) := \arctan(x) + \pi \, \operatorname{rni}\left(\frac{\eta - \arctan(x)}{\pi} \right)\, .\)

The function \(\operatorname{rni}\) rounds to the nearest integer.

#### Numerical accuracy
For angles near 0 and , arccosine is ill-conditioned, and similarly with arcsine for angles near −/2 and /2.  Computer applications thus need to consider the stability of inputs to these functions and the sensitivity of their calculations, or use alternate methods.

## See also

*Arcsine distribution
*Inverse exsecant
*Inverse versine
*Inverse hyperbolic functions
*List of integrals of inverse trigonometric functions
*List of trigonometric identities
*Trigonometric function
*Trigonometric functions of matrices


## Notes


## References
*


## External links
*

