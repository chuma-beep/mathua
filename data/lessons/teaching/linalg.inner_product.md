# Inner Products and Orthogonality

**Inner product:** $\langle u,v\rangle$ satisfies $\langle u,u\rangle\ge0$ with equality iff $u=0$, symmetry, linearity. Standard dot product $\langle u,v\rangle=u\cdot v$; weighted $\langle u,v\rangle=2u_1v_1+u_2v_2$ also works.

## Orthogonality

### Cauchy-Schwarz
$|\langle u,v\rangle|\le\|u\|\|v\|$ with equality iff $u,v$ linearly dependent.

### Gram-Schmidt
Orthogonalize via $v_2' = v_2 - \text{proj}_{v_1}v_2$ using $\langle v_2,v_1\rangle/\langle v_1,v_1\rangle$.

## Example

$u=(1,2)$, $v=(2,1)$: $\langle u,v\rangle=1·2+2·1=4$, not orthogonal.
