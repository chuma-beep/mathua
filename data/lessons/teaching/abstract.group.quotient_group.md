# Quotient Groups and Isomorphism Theorems

**Quotient group:** If $N\trianglelefteq G$, the cosets $G/N=\{gN:g\in G\}$ form a group under $(aN)(bN)=abN$. The natural map $\pi:G\to G/N$ is a homomorphism with kernel $N$.

## Building Quotients

### The Construction
Normality ensures well-defined multiplication: $(aN)(bN)=abN$ independent of representatives. $|G/N|=|G|/|N|$ when $G$ is finite.

### First Isomorphism Theorem
If $\varphi:G\to H$ is a homomorphism, then $G/\ker\varphi\cong\operatorname{im}\varphi$. This classifies quotients: every quotient is an image, every image is a quotient.

## Example

$G=\mathbb Z$, $N=2\mathbb Z$: cosets are even and odd, so $\mathbb Z/2\mathbb Z\cong\mathbb Z_2$. Map $\varphi:\mathbb Z\to\mathbb Z_2$, $\varphi(n)=n\bmod2$, has kernel $2\mathbb Z$ and image $\mathbb Z_2$.
