# Argument of a Complex Number

**Argument:** For $z \ne 0$, $\arg(z)$ is the counter-clockwise angle from the positive real axis to $z$. Axis points sit at 0, 90, 180, and 270 degrees.

## Worked: argument of 1+i

Find $\arg(1+i)$:
1. Locate the quadrant: $a=1 > 0$ and $b=1 > 0$, so $z$ is in the first quadrant.
2. Take the ratio: $b/a=1$, and $\arctan(1)=45$ degrees.
3. Confirm on the circle: $|z|=\sqrt{2}$ and $z=\sqrt{2}(\cos 45+i\sin 45)$.

So in the first quadrant the argument is just $\arctan(b/a)$; the other quadrants shift this by 180 or 360 degrees.

## Worked: arguments on the axes

Read the four axis arguments directly off the plane:
1. $z=1$ points along the positive real axis, so $\arg(1)=0$ degrees.
2. $z=i$ points straight up, so $\arg(i)=90$ degrees.
3. $z=-1$ points left, so $\arg(-1)=180$ degrees.

So the axes give the reference angles; every other argument is measured against them counter-clockwise.

## Worked: arguments add under multiplication

Multiply $z_1=1+i$ by $z_2=1+i$ and track the angle:
1. Each factor has argument 45 degrees, summing to 90 degrees.
2. Multiply out: $(1+i)(1+i)=1+2i+i^2=2i$.
3. Check: $\arg(2i)=90$ degrees, matching the sum.

So $\arg(z_1 z_2)=\arg(z_1)+\arg(z_2)$ modulo 360 degrees: moduli multiply while angles add.
