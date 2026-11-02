# Closures of Relations

## Reflexive and Symmetric Closure

A relation on a set may lack a property that we need. The reflexive closure adds exactly the missing pairs of the form (a, a), and no others.

The symmetric closure adds the reverse of every pair that is present. If a pair and its reverse are both already present, nothing is added for that pair.

Closures are the smallest relations containing the original and satisfying the required property, which is why no unnecessary pairs are included.

## Transitive Closure

The transitive closure adds a pair (a, c) whenever there is a path from a to c through the original relation. It is the reachability relation of the corresponding directed graph.

Computing it by repeated composition works for small relations. Warshall's algorithm is the standard systematic method, and it updates a matrix of reachable pairs one intermediate vertex at a time.

The transitive closure is reflexive only if every vertex reaches itself, so reflexive and transitive closures are different operations.

## Counting the Additions

Questions about closures usually ask how many pairs must be added, or whether a particular pair belongs. Counting is easiest by comparing the original set with the closure set.

For a reflexive closure on a set of n elements, the answer is n minus the number of diagonal pairs already present. For a symmetric closure, count the pairs whose reverse is missing.

Transitive closures can add many pairs at once, so list the paths rather than guessing.
