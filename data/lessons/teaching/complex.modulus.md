# Modulus of a Complex Number

**Modulus:** For $z=a+bi$, the modulus is $|z|=\sqrt{a^2+b^2}$, the distance from $z$ to the origin. It satisfies $|z|^2=z\bar{z}$.

## Worked: modulus of 3+4i

Find $|3+4i|$:
1. Square the parts: $a^2=9$ and $b^2=16$.
2. Add: $|z|^2=9+16=25$.
3. Take the root: $|z|=\sqrt{25}=5$.

So the modulus is ordinary Euclidean distance: $|a+bi|$ is the hypotenuse of the right triangle with legs $|a|$ and $|b|$.

## Worked: modulus squared via the conjugate

Check $|z|^2=z\bar{z}$ on $z=3+4i$:
1. Form the conjugate: $\bar{z}=3-4i$.
2. Multiply: $z\bar{z}=(3+4i)(3-4i)=9+12i-12i+16=25$.
3. Compare with $|z|^2=25$ from the previous computation — they agree.

So $|z|^2=z\bar{z}$ in general: the cross terms cancel, leaving $a^2+b^2$.

## Worked: triangle inequality on two vectors

Test $|z+w| \le |z|+|w|$ with $z=3+4i$ and $w=5+12i$:
1. Read the moduli: $|z|=5$ and $|w|=13$, so $|z|+|w|=18$.
2. Add: $z+w=8+16i$.
3. Take the modulus: $|z+w|=\sqrt{64+256}=\sqrt{320}$, and $\sqrt{320} < 18$ since $18^2=324$.

So $|z+w| \le |z|+|w|$ holds here; in general the direct path never exceeds the detour through $w$.
