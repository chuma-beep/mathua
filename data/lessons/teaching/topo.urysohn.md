# Urysohn's Lemma

**Urysohn's lemma:** In a normal space $X$, for any disjoint closed sets $A,B$ there exists a continuous $f:X\to[0,1]$ with $f|_A=0$ and $f|_B=1$.

## Separating Closed Sets

### Normality
Normal means disjoint closed sets have disjoint open neighborhoods. Metric spaces are normal, so the lemma applies broadly.

### The Separating Function
$f$ is $0$ on $A$, $1$ on $B$, and transitions continuously between. It proves normal spaces have rich continuous functions.

## Example

In $\mathbb R$ with $A=\{0\}$ and $B=\{1\}$, $f(x)=x$ on $[0,1]$ (extended constantly outside) separates them. In any metric space, $f(x)=d(x,A)/(d(x,A)+d(x,B))$ does the job.
