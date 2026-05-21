> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Ordinal_number) — CC BY-SA 4.0

# Ordinal numbers: first, second, third

In set theory, an **ordinal number**, or **ordinal**, is a generalization of ordinal numerals (first, second, th, etc.) aimed to extend enumeration to infinite sets. Usually Greek letters are used for ordinal number variables to help distinguish them from natural number variables.

A finite set can be enumerated by successively labeling each element with the least natural number that has not been previously used. To extend this process to various infinite sets, ordinal numbers are defined more generally as a linearly ordered class of numbers that include the natural numbers and have the property that every non-empty collection (set or proper class) of ordinals has a least or "smallest" element (this is needed for giving a meaning to "the least unused element"). This more general definition allows us to define an ordinal number \(\omega\) (omega) to be the least element that is greater than every natural number, along with ordinal numbers , etc., which are even greater than. The Zermelo–Fraenkel set theory asserts that, for *any* set of ordinals, there exists another ordinal greater than all of them. The answer to the question "What if that set is the set of all ordinals?" (the Burali-Forti paradox) is that the collection of all ordinals is not a set, but a proper class.

A linear order such that every non-empty subset has a least element is called a well-order. The axiom of choice implies that every set can be well-ordered. Given two well-ordered sets, one is isomorphic to an initial segment of the other, and the isomorphism is unique. This allows a unique ordinal to be associated with each well-ordered set, known as its order type.

Ordinal numbers are distinct from cardinal numbers, which measure the size of sets. Although the distinction between ordinals and cardinals is not this apparent on finite sets (one can go from one to the other just by counting labels), they are very different in the infinite case, where different infinite ordinals can correspond to sets having the same cardinal. Like other kinds of numbers, ordinals can be added, multiplied, and exponentiated, although none of these operations are commutative.

Ordinals were introduced by Georg Cantor in 1883 to accommodate infinite sequences and classify derived sets, which he had previously introduced in 1872 while studying the uniqueness of trigonometric series.

## Motivation

A natural number (which, in this context, includes the number 0) can be used for two purposes: to describe the *size* of a set, or to describe the *position* of an element in a sequence. When generalized to infinite sets, the notion of size leads to cardinal numbers, and the notion of position leads to the ordinal numbers described here.

In a broader mathematical sense, counting can be viewed as the instantiation of mathematical induction. To enumerate a well-ordered set is effectively to verify a property for its elements sequentially. For the natural numbers, this is standard induction: if a property holds for 0, and its truth for implies its truth for , then it holds for all natural numbers. This process corresponds to the first infinite ordinal,. Mathematical contexts often require iterating beyond a single infinite limit. The ordinal (represented in the figure) exemplifies the concept of **nested induction**. It consists of a sequence of distinct copies of the natural numbers ordered one after another. To verify a property for all ordinals less than , one performs an "inner" induction (counting through ), establishes the limit at , and then proceeds to the next sequence (). This structure parallels a nested loop in computer programming (e.g., iterating through pairs of natural numbers ordered lexicographically). Ordinals allow the definition of processes of arbitrary complexity, such as (triple nesting) or (induction over the depth of nested induction).

The validity of inductive counting rests on the property of well-foundedness, specifically the requirement that every process can be traced back to a "foundational" element. A linear order that exhibits this well-foundedness is termed a well-order. The existence of a "least" or minimal element in every non-empty subset of a well-ordered set grounds the principle of transfinite induction, generalizing standard induction by ensuring that if a property fails to hold, there exists a specific least counterexample.

Ordinals serve as the canonical abstractions of these well-ordered structures. A fundamental theorem in set theory establishes that any two well-ordered sets are comparable: given two well-orders, either they are isomorphic, or one is isomorphic to a proper initial segment of the other. This uniqueness implies that well-orders can be classified by their structure alone, independent of specific representations. Consequently, **ordinal numbers** are defined as the representative forms of these isomorphism classes.

## Definitions
### Well-ordering
The construction procedure of transfinite sequences implies that the ordinals used to label their elements are well-ordered. This means that:
Every non-empty collection of ordinals has a unique least element.
Intuitively, the least ordinal in is the "next" ordinal introduced after all ordinals strictly less than every ordinal in have been used. In contrast, need not have a greatest element, for example when is the set of all natural numbers.

Being well-ordered is stronger than merely being linearly ordered. For example, the real numbers are naturally linearly ordered, but any open interval has no least element.

The importance of well-ordering is that it enables transfinite induction. If a statement about an ordinal is not universally true for all , i.e., if it has any counterexamples, then it must have a least counterexample. Conversely, if it can be proven that is true whenever is true for all , then it cannot have any least counterexample, and thus it cannot have any counterexample at all, i.e., it is indeed universally true.

### Well-ordered sets
In the usual Zermelo–Fraenkel (ZF) formalization of set theory, a well-ordered set is a totally ordered set such that every non-empty subset has a least element. Here “set” means a collection that is itself an object of ZF, as opposed to a proper class.

A well-ordered set can be written as a transfinite sequence by assigning ordinal labels to its elements in a way that respects ordering: if and only if. (Such a one-to-one correspondence between two ordered sets is called an order isomorphism.) The labels used in such an enumeration form an initial segment of the ordinals, in the sense that if a label is used then every smaller ordinal label is also used. By well-orderedness, there is a unique least ordinal that is not used as a label; this ordinal determines the "length" of the sequence, and is called the order type of the well-ordering. Thanks to 0-based indexing, the order type coincides with the cardinality for finite sets, including the empty set. However, for infinite sets, different well-order relations can have different order types.

### Definition of an ordinal as an equivalence class
Order types can be defined without a preexisting notion of ordinals, by working directly with order isomorphisms between general well-ordered sets, just as cardinality can be defined through bijections between general (unordered) sets. As with bijections, being order-isomorphic is an equivalence relation on well-ordered sets; its equivalence classes correspond to order types (ordinals).

In the *Principia Mathematica* approach, the order type of a well-ordered set is *identified with* its isomorphism class, i.e., the set of all well-ordered sets order-isomorphic to. Since elements of are allowed to be *anything*, this definition has the “set of all sets” flavor, and in ZF such collections are generally too large to be sets. This definition can still be used in type theory and in Quine's axiomatic set theory New Foundations and related systems.

In ZF and related systems of axiomatic set theory, these equivalence classes are generally too large to form sets. Consequently, it is necessary to select a unique, canonical representative from each class—a single set that embodies the structure of the well-ordering. Since the fundamental relation in set theory is set membership (), the ideal representation is one where the abstract order relation is translated directly into the membership relation. ## Von Neumann definition of ordinals
The **von Neumann representation** provides this canonical form. It relies on the observation that any well-founded relation satisfying certain properties can be mapped to a specific set where the relation becomes set membership. This mapping is known as the Mostowski collapse lemma.

When applied to a well-ordering, the Mostowski collapse yields a specific set where the order relation is true if and only if. The resulting sets have the property that they are transitive: every element of is also a subset of (i.e., the union of the set is contained within the set). In this representation, each ordinal is identified with the set of all preceding ordinals.

Thus the finite **von Neumann ordinals** are defined recursively as , \(1 = \{0\}\), \(2 = \{0,1\}\), etc. The first infinite ordinal is represented by the set of all finite ordinals, i.e., the set of von Neumann natural numbers \(\mathbb{N} = \{0, 1, 2, \ldots\}\). Then \(\omega + 1 = \{0, 1, 2, \ldots, \omega\} = \mathbb{N} \cup \{\omega\}\), and so on.

Informally, one may define an ordinal recursively as a downward closed set of ordinals. Such a recursive definition is usually justified with the transitive closure. However, the von Neumann ordinals are already transitive sets, allowing them to be formally defined by a concise statement:
A set \(S\) is an ordinal if and only if \(S\) is transitive (every element of \(S\) is a subset of \(S\)) and strictly well-ordered.by set membership (\(\in\)).
The set \(\omega \equiv \mathbb{N}\) is usually defined as the smallest inductive set (containing and closed under successor). The "smallest" constraint guarantees that each element of is either zero or the successor of another element of , which allows induction to demonstrate that indeed satisfies the formal definition of ordinals stated above.

### Basic properties
Defining the strict order relation as the membership relation restricted to the class of all ordinals, the recursive characterization \(\gamma = \{\alpha \mid \alpha < \gamma\}\) is reduced to the following statement:

* *If and is an ordinal, then is an ordinal*: transitivity of as a set is found by finding that all the relevant ordinals are elements of and then using the transitivity of the relation within ; well-orderedness of follows straightforwardly from well-orderedness of. The non-strict order relation has an alternative characterization: if and only if for ordinals. The "if" direction follows from transitivity of , and the "only if" direction follows from:
* *If are both ordinals and , then *: let \(\gamma = \min\{\beta \setminus \alpha\}\). is transitive so \(\alpha = \{\xi \in \beta \mid \xi < \gamma\} = \gamma \in \beta\).
This implies that is a partial order. It is in fact a total order, and a well-order:
* *If and are both ordinals, then or *: is an ordinal, so or , or else , contradicting irreflexivity.
* If is a nonempty set of ordinals, then must be in by similar logic.

So all ordinal numbers form a well-ordered class \(\mathrm{ON}\), and thus any nonempty set of ordinals equipped with is a well-ordered set.

Specific ordinals can be constructed explicitly with the following principles:
* is an ordinal.
* For any ordinal , \(\mathrm{succ}\;\alpha = \alpha \cup \{\alpha\}\) is an ordinal and \(\mathrm{succ}\;\alpha = \min\{\beta \mid \beta > \alpha\}\).
* If A is a set of ordinals, then \(\sup A = \bigcup A\) is an ordinal.

The explicit forms of \(\mathrm{succ}\) and imply that for any set of ordinals, there exists another ordinal greater than all of them. In other words,
* (Burali-Forti paradox) The class of all ordinals \(\mathrm{ON}\) is not a set; otherwise \(\mathrm{succ} \sup \mathrm{ON}\) would be an ordinal not in \(\mathrm{ON}\).

### Order types
Every well-ordered set \(S\) is order-isomorphic to exactly one ordinal, known as its *order type*. Uniqueness is guaranteed because a well-ordered set cannot be isomorphic to a proper initial segment of itself, preventing isomorphism to two distinct ordinals. The existence of this ordinal is proven by defining a map pairing each \(x \in S\) with the ordinal representing the order type of the initial segment \(S_x = \{y \in S \mid y < x\}\). By the Axiom schema of replacement, the range of this map is a set of ordinals. Because this range is downward closed (the order type of a segment of a segment is a smaller ordinal), the range is itself an ordinal \(\gamma\). The domain of the isomorphism must be all of \(S\); otherwise, the least element \(z \in S\) outside the domain would imply \(S_z \cong \gamma\), which would effectively include \(z\) in the domain, a contradiction. Thus, \(S \cong \gamma\).

## Successor and limit ordinals
Every ordinal number is one of three types: the ordinal zero, a successor ordinal, or a limit ordinal.
* *Zero*: The ordinal \(0 = \emptyset\) is the least ordinal.
* *Successor ordinals*: An ordinal \(\alpha\) is a successor if \(\alpha = S(\beta) = \beta \cup \{\beta\}\) for some ordinal \(\beta\). In this case, \(\beta\) is the maximum element of \(\alpha\).
* *Limit ordinals*: An ordinal \(\lambda\) is a limit ordinal if \(\lambda \neq 0\) and \(\lambda\) is not a successor ordinal.

There is variation in the definition of limit ordinals regarding the inclusion of zero. Some texts, such as *Introduction to Cardinal Arithmetic* by Holz et al., define a limit ordinal as a non-zero ordinal that is not a successor. In contrast, other standard set theory texts, including Jech's *Set Theory* and Just and Weese's *Discovering Modern Set Theory*, define a limit ordinal simply as any ordinal that is not a successor, which implies that 0 is a limit ordinal. When the topological definition is used (based on the order topology), 0 is not a limit ordinal because it is not a limit point of the set of smaller ordinals (which is empty); Rosenstein's *Linear Orderings* uses this definition. When 0 is included as a limit, ordinals that are strictly greater than 0 and not successors are usually referred to as "nonzero limit ordinals".

The following properties characterize nonzero limit ordinals:

\(\lambda\) is a nonzero limit ordinal if and only if \(\lambda \neq 0\) and for every ordinal \(\alpha < \lambda\), the successor \(S(\alpha)\) is also less than \(\lambda\).

This implies that a nonzero limit ordinal is equal to the supremum of all ordinals strictly less than it:

\(\lambda\) is a nonzero limit ordinal if and only if \(\lambda = \sup \{ \alpha \mid \alpha < \lambda \} = \bigcup \lambda\) and \(\lambda \neq 0\).

For example, \(\omega\) is a limit ordinal because any natural number is less than , and the successor of any natural number is also a natural number (hence less than ). It is the least limit ordinal because each natural number is either zero or a successor.

## Termination of decreasing sequences
Any strictly decreasing sequence of ordinals must be finite. This follows directly from the natural ordering of ordinals being a well-order: if there existed such an infinite decreasing sequence, then the set \(\{\alpha_i \mid i \in \mathbb{N}\}\) would be a set of ordinals without a least element. By the same argument, a well-ordered set has no infinite strictly descending chains; in fact, assuming the axiom of dependent choice, any total order that satisfies this condition is a well-order, giving an alternative characterization of well-ordered sets. The fact this is true for natural numbers is the basis of Fermat's method of proof by infinite descent, which can be generalized to ordinals and other well-ordered classes too, as a special case of transfinite induction where the proof of any only requires for at most one specific. This property may be surprising when the initial value is an infinite ordinal. Indeed, for sequences starting from a natural number , the longest sequence is always one that decreases by 1 every step, leading to a sequence with steps ( elements). This strategy of descending to the immediate predecessor remains valid for successor ordinals. However, a limit ordinal has no immediate predecessor to descend to, so any next term must jump to some , skipping infinitely many ordinals strictly between and. For example, when descending from , one must choose a finite natural number, and thus "commit" to the number of maximum remaining steps. Descending from allows one to make such a commitment times, and descending from allows one to commit to a finite value of. Larger ordinals may allow more complicated decision structures, but the number of descending steps remains unbounded but finite.

This property is useful for proving termination for any procedure. If the states of a computation (computer program or game) can be well-ordered—in such a way that each step is followed by a "lower" step—then the computation will terminate.

## Transfinite sequence
If \(\alpha\) is any ordinal and \(X\) is a set, an \(\alpha\)-indexed sequence of elements of \(X\) is a function from \(\alpha\) to \(X\). This concept, a **transfinite sequence** (if \(\alpha\) is infinite) or *ordinal-indexed sequence*, is a generalization of the concept of a sequence. An ordinary sequence corresponds to the case \(\alpha = \omega\), while a finite \(\alpha\) corresponds to a tuple, a.k.a. string.

While a sequence indexed by a specific ordinal \(\alpha\) is a set, a sequence indexed by the class of all ordinals is a proper class. The Axiom schema of replacement guarantees that any initial segment of such a class-sequence (the restriction of the function to some specific ordinal \(\delta\)) is a set.

When \(\langle x_{\iota} \mid \iota < \lambda \rangle\) is a transfinite sequence of ordinals indexed by a limit ordinal \(\lambda\) and the sequence is *increasing* (i.e. \(\iota < \rho \implies x_{\iota} < x_{\rho}\)), its *limit* is defined as the least upper bound of the set \(\{ x_{\iota} \mid \iota < \lambda \}\).

A transfinite sequence \(f\) mapping ordinals to ordinals is said to be *continuous* (in the order topology) if for every limit ordinal \(\lambda\) in its domain,
* if *f* (λ) is a limit ordinal and for every ε < *f* (λ) there exists a δ < λ such that for every γ, if δ < γ < λ, then ε < *f* (γ) ≤ *f* (λ), and
*if *f* (λ) is *not* a limit ordinal, there exists a δ < λ such that for every γ, if δ < γ < λ, then *f* (γ) = *f* (λ).
A sequence is called **normal** if it is both strictly increasing and continuous. If a sequence *f* is increasing (not necessarily strictly) and continuous and λ is a limit ordinal, then \(f(\lambda) = \bigcup_{\beta < \lambda} f(\beta)\).

## Transfinite induction

Transfinite induction holds in any well-ordered set, but it is so important in relation to ordinals that it is worth restating here.

Any property that passes from the set of ordinals smaller than a given ordinal α to α itself, is true of all ordinals.

That is, if *P*(α) is true whenever *P*(β) is true for all β < α, then *P*(α) is true for *all* α. Or, more practically: in order to prove a property *P* for all ordinals α, one can assume that it is already known for all smaller β < α.

### Transfinite recursion
Transfinite induction can be used not only to prove theorems but also to define functions on ordinals. This is known as *transfinite recursion*.

Formally, a function is defined by transfinite recursion on the ordinals if, for every ordinal , the value is specified using the set of values \(\{ F(\beta) \mid \beta < \alpha \}\).

Very often, when defining a function *F* by transfinite recursion on all ordinals, the definition is separated into cases based on the type of the ordinal:
# *Base case*: Define \(F(0)\).
# *Successor step*: Define \(F(\alpha+1)\) assuming \(F(\alpha)\) is defined.
# *Limit step*: For a limit ordinal \(\lambda\), define \(F(\lambda)\) as the limit of \(F(\beta)\) for all \(\beta < \lambda\) (either in the sense of ordinal limits or some other notion of limit if the codomain allows it).

The interesting step in the definition is usually the successor step. If \(F(\alpha)\) for limit ordinals \(\alpha\) is defined as the limsup of \(F(\beta)\) for \(\beta < \alpha\) and \(F\) takes ordinal values and is non-decreasing, the function \(F\) will be continuous as defined above. Ordinal addition, multiplication and exponentiation are continuous as functions of their second argument.

The existence and uniqueness of such a function are proven by constructing it as the union of partial approximations. The proof proceeds in three steps:
# *Local Existence:* For any specific ordinal , one proves the existence of a unique "recursion segment"—a function defined on that satisfies the recursive rule for all. # *Uniqueness and Compatibility:* One proves that any two recursion segments agree on their common domain. If is a segment on and is a segment on with , then restricted to is identical to. # *Global Definition:* The global class function is defined as the union of all such unique recursion segments. For any ordinal , the value is the value assigned to by any recursion segment defined on a domain larger than. The rigorous justification for local existence relies on the axiom schema of replacement for the step of limit ordinals in order to collect the recursion segments into a set.

This construction allows definitions such as ordinal addition, multiplication, and exponentiation to be rigorous. For example, exponentiation is defined recursively on :
* \(\alpha^0 = 1\)
* \(\alpha^{\beta+1} = \alpha^\beta \cdot \alpha\) (for successor ordinals)
* \(\alpha^\lambda = \bigcup_{0 < \beta < \lambda} \alpha^\beta\) (for limit ordinals )

The principle of transfinite recursion also implies that recursion can be performed up to a specific ordinal \(\delta\) (defining a set rather than a proper class). This is often used to define sequences of length \(\omega\). For example, to prove that a normal function \(f\) has arbitrarily large fixed points, one constructs a sequence starting with any ordinal \(\gamma_0\) and defines \(\gamma_{n+1} = f(\gamma_n)\) by recursion on \(n < \omega\). Then the limit \(\delta = \sup_{n<\omega} \gamma_n\) is a fixed point of \(f\) because continuity ensures \(f(\delta) = f(\sup \gamma_n) = \sup f(\gamma_n) = \sup \gamma_{n+1} = \delta\).

### Indexing classes of ordinals
Any well-ordered set is similar (order-isomorphic) to a unique ordinal number \(\alpha\); in other words, its elements can be indexed in increasing fashion by the ordinals less than. This applies, in particular, to any set of ordinals: any set of ordinals is naturally indexed by the ordinals less than some. The same holds, with a slight modification, for *classes* of ordinals (a collection of ordinals, possibly too large to form a set, defined by some property): any class of ordinals can be indexed by ordinals (and, when the class is unbounded in the class of all ordinals, this puts it in class-bijection with the class of all ordinals). So the \(\gamma\)-th element in the class (with the convention that the "0-th" is the smallest, the "1-st" is the next smallest, and so on) can be freely spoken of. Formally, the definition is by transfinite induction: the \(\gamma\)-th element of the class is defined (provided it has already been defined for all \(\beta<\gamma\)), as the smallest element greater than the \(\beta\)-th element for all. This could be applied, for example, to the class of limit ordinals: the \(\gamma\)-th ordinal, which is either a limit or zero is \(\omega\cdot\gamma\) (see ordinal arithmetic for the definition of multiplication of ordinals). Similarly, one can consider *additively indecomposable ordinals* (meaning a nonzero ordinal that is not the sum of two strictly smaller ordinals): the \(\gamma\)-th additively indecomposable ordinal is indexed as. The technique of indexing classes of ordinals is often useful in the context of fixed points: for example, the \(\gamma\)-th ordinal \(\alpha\) such that \(\omega^\alpha = \alpha\) is written. These are called the "epsilon numbers".

## Arithmetic of ordinals

There are three usual operations on ordinals: addition, multiplication, and exponentiation. Each can be defined in essentially two different ways: either by constructing an explicit well-ordered set that represents the operation or by using transfinite recursion. The Cantor normal form provides a standardized way of writing ordinals. It uniquely represents each ordinal as a finite sum of ordinal powers of ω. However, this cannot form the basis of a universal ordinal notation due to such self-referential representations as ε\(_{0}\) = ω\(^{ε\(_{0}\)}\).

Ordinals are a subclass of the class of surreal numbers, and the so-called "natural" arithmetical operations for surreal numbers are an alternative way to combine ordinals arithmetically. They retain commutativity at the expense of continuity.

Interpreted as *nimbers*, a game-theoretic variant of numbers, ordinals can also be combined via nimber arithmetic operations. These operations are commutative but the restriction to natural numbers is generally not the same as ordinary addition of natural numbers.

## Ordinals and cardinals
### Initial ordinal of a cardinal
Each ordinal associates with one cardinal, its cardinality. If there is a bijection between two ordinals (e.g. ω and ω + 1 > ω), then they associate with the same cardinal. Any well-ordered set having an ordinal as its order-type has the same cardinality as that ordinal. The least ordinal associated with a given cardinal is called the *initial ordinal* of that cardinal. Every finite ordinal (natural number) is initial, and no other ordinal associates with its cardinal. But most infinite ordinals are not initial, as many infinite ordinals associate with the same cardinal. The axiom of choice is equivalent to the statement that every set can be well-ordered, i.e. that every cardinal has an initial ordinal. In theories with the axiom of choice, the cardinal number of any set has an initial ordinal, and one may employ the Von Neumann cardinal assignment as the cardinal's representation. (However, we must then be careful to distinguish between cardinal arithmetic and ordinal arithmetic.) In set theories without the axiom of choice, a cardinal may be represented by the set of sets with that cardinality having minimal rank (see Scott's trick).

One issue with Scott's trick is that it identifies the cardinal number \(0\) with \(\{\emptyset\}\), which in some formulations is the ordinal number. It may be clearer to apply Von Neumann cardinal assignment to finite cases and to use Scott's trick for sets which are infinite or do not admit well orderings. Note that cardinal and ordinal arithmetic agree for finite numbers.

The α-th infinite initial ordinal is written , it is always a limit ordinal. Its cardinality is written. For example, the cardinality of ω\(_{0}\) = ω which is also the cardinality of ω\(^{2}\) or ε\(_{0}\) (all are countable ordinals). So ω can be identified with , except that the notation \(\aleph_0\) is used when writing cardinals, and ω when writing ordinals (this is important since, for example, \(\aleph_0^2\) = \(\aleph_0\) whereas \(\omega^2 > \omega\)). Also, \(\omega_1\) is the smallest uncountable ordinal (to see that it exists, consider the set of equivalence classes of well-orderings of the natural numbers: each such well-ordering defines a countable ordinal, and \(\omega_1\) is the order type of that set), \(\omega_2\) is the smallest ordinal whose cardinality is greater than , and so on, and \(\omega_\omega\) is the limit of the \(\omega_n\) for natural numbers *n* (any limit of cardinals is a cardinal, so this limit is indeed the first cardinal after all the \(\omega_n\)).

### Cofinality
The cofinality of an ordinal \(\alpha\) is the smallest ordinal \(\delta\) that is the order type of a cofinal subset of. Notice that a number of authors define cofinality or use it only for limit ordinals. The cofinality of a set of ordinals or any other well-ordered set is the cofinality of the order type of that set.

Thus for a limit ordinal, there exists a \(\delta\)-indexed strictly increasing sequence with limit. For example, the cofinality of ω\(^{2}\) is ω, because the sequence ω·*m* (where *m* ranges over the natural numbers) tends to ω\(^{2}\); but, more generally, any countable limit ordinal has cofinality ω. An uncountable limit ordinal may have either cofinality ω as does \(\omega_\omega\) or an uncountable cofinality.

The cofinality of 0 is 0. And the cofinality of any successor ordinal is 1. The cofinality of any limit ordinal is at least. An ordinal that is equal to its cofinality is called regular and it is always an initial ordinal. Any limit of regular ordinals is a limit of initial ordinals and thus is also initial even if it is not regular, which it usually is not. If the axiom of choice holds, then \(\omega_{\alpha+1}\) is regular for each *α*. In this case, the ordinals 0, 1, and \(\omega_2\) are regular, whereas 2, 3, and ω\(_{ω·2}\) are initial ordinals that are not regular.

The cofinality of any ordinal *α* is a regular ordinal, i.e. the cofinality of the cofinality of *α* is the same as the cofinality of *α*. So the cofinality operation is idempotent.

### Closed unbounded sets and classes
The concepts of closed and unbounded sets are typically formulated for subsets of a regular cardinal \(\kappa\) that is uncountable. A subset \(C \subseteq \kappa\) is said to be *unbounded* (or cofinal) in \(\kappa\) if for every ordinal \(\alpha < \kappa\), there exists some \(\beta \in C\) such that \(\alpha < \beta\). To define the property of being closed, one first defines a limit point: a non-zero ordinal \(\delta < \kappa\) is a limit point of \(C\) if \(\sup(C \cap \delta) = \delta\). The set \(C\) is *closed* in \(\kappa\) if it contains all of its limit points below \(\kappa\). A set that is both closed and unbounded is commonly referred to as a *club set*.

Examples of club sets are fundamental to set theory. The set of all limit ordinals less than \(\kappa\) is a club set, as there is always a limit ordinal greater than any given ordinal below \(\kappa\), and a limit of limit ordinals is itself a limit ordinal. If \(\kappa\) is a limit cardinal, the set of all cardinals below \(\kappa\) is unbounded, and its set of limit points—the limit cardinals—forms a closed unbounded set. Furthermore, if \(\kappa\) is a strong limit cardinal (such as an inaccessible cardinal), the set of strong limit cardinals below \(\kappa\) is also a club set. Another significant example arises from normal functions (functions \(f: \kappa \to \kappa\) that are strictly increasing and continuous); the range of any normal function is a closed unbounded subset of \(\kappa\).

Club sets possess structural properties that allow them to generate a filter. Because \(\kappa\) is regular and uncountable, the intersection of any two club sets is also a club set. More generally, the intersection of fewer than \(\kappa\) club sets is a club set. Consequently, the collection of all subsets of \(\kappa\) that contain a club set forms a \(\kappa\)-complete non-principal filter, known as the *closed unbounded filter* (or club filter).

A subset \(S \subseteq \kappa\) is termed *stationary* if it has a non-empty intersection with every closed unbounded set in \(\kappa\). Intuitively, stationary sets are "large" enough that they cannot be avoided by any club set. Using the notation of filters, a set is stationary if and only if it does not belong to the dual ideal of the club filter (the ideal of non-stationary sets). While every club set is stationary, not every stationary set is a club; for instance, a stationary set \(S\) may fail to be closed. Furthermore, while the intersection of a stationary set and a club set is stationary, the intersection of two stationary sets may be empty.

The distinction between club sets and stationary sets is central to the definitions of certain large cardinals. If \(\kappa\) is the smallest inaccessible cardinal, the set of singular strong limit cardinals below \(\kappa\) forms a closed unbounded set. Because this club set contains no regular cardinals, the set of regular cardinals below the first inaccessible is not stationary. This remains true if \(\kappa\) is the \(n\)-th inaccessible cardinal for some \(n < \kappa\); the regular cardinals below it will not form a stationary set. A cardinal \(\kappa\) is defined as a Mahlo cardinal precisely when the set of regular cardinals below it is stationary. By relaxing the condition on the limit cardinals, one defines a cardinal as *weakly Mahlo* if it is weakly inaccessible and the set of regular cardinals below it is stationary.

The closed unbounded filter is not an ultrafilter under the standard Zermelo–Fraenkel set theory with the Axiom of Choice (ZFC). This is because one can find two disjoint stationary sets, which precludes the filter from deciding membership for every subset. For any regular cardinal \(\kappa > \omega_1\), the set of ordinals with cofinality \(\omega\) and the set of ordinals with cofinality \(\omega_1\) are disjoint stationary subsets of \(\kappa\). In the specific case of \(\kappa = \omega_1\), the non-existence of an ultrafilter relies on the Axiom of Choice. Under ZFC, the set of limit ordinals in \(\omega_1\) can be partitioned into \(\omega_1\) disjoint stationary sets (a result related to Fodor's lemma). However, in models of set theory without the Axiom of Choice, such as those satisfying the Axiom of determinacy, the club filter on \(\omega_1\) can be an ultrafilter, a property connected to \(\omega_1\) being a measurable cardinal in those contexts.

These definitions generalize to proper classes of ordinals. A class \(C\) of ordinals is unbounded if it contains arbitrarily large ordinals, and closed if the limit of any sequence of ordinals in \(C\) is also in \(C\). This topological definition is equivalent to assuming the indexing class-function of \(C\) is continuous. Notable examples of closed unbounded classes include the class of all infinite cardinals, the class of limit cardinals, and the class of fixed points of the \(\aleph\)-function. In contrast, the class of regular cardinals is unbounded but not closed. A class is stationary if it intersects every closed unbounded class.

## Some "large" countable ordinals

As mentioned above (see Cantor normal form), the ordinal ε\(_{0}\) is the smallest satisfying the equation , so it is the limit of the sequence 0, 1, \(\omega^{\omega^\omega}\), etc. Many ordinals can be defined in such a manner as fixed points of certain ordinal functions (the \(\iota\)-th ordinal such that \(\omega^\alpha = \alpha\) is called , then one could go on trying to find the \(\iota\)-th ordinal such that , "and so on", but all the subtlety lies in the "and so on"). One could try to do this systematically, but no matter what system is used to define and construct ordinals, there is always an ordinal that lies just above all the ordinals constructed by the system. Perhaps the most important ordinal that limits a system of construction in this manner is the Church–Kleene ordinal, \(\omega_1^{\mathrm{CK\) (despite the \(\omega_1\) in the name, this ordinal is countable), which is the smallest ordinal that is not a computable ordinal (i.e., the order type of a well-ordering of the natural numbers such that the predicate is computable). Considerably large ordinals can be defined below \(\omega_1^{\mathrm{CK} }\), however, which measure the proof-theoretic strength of certain formal systems (for example, \(\varepsilon_0\) measures the strength of Peano arithmetic). Large countable ordinals such as countable admissible ordinals can also be defined above the Church-Kleene ordinal, which are of interest in various parts of logic.

## Topology and ordinals

Any ordinal number can be made into a topological space by endowing it with the order topology. This topology is discrete if and only if it is less than or equal to ω. In contrast, a subset of ω + 1 is open in the order topology if and only if either it is cofinite or it does not contain ω as an element.

See the Topology and ordinals section of the "Order topology" article.

## History
The transfinite ordinal numbers, which first appeared in 1883, originated in Cantor's work with derived sets. If *P* is a set of real numbers, the derived set ** is the set of limit points of *P*. In 1872, Cantor generated the sets *P*\(^{(*n*)}\) by applying the derived set operation *n* times to *P*. In 1880, he pointed out that these sets form the sequence *P' *⊇ ··· ⊇ *P*\(^{(*n*)}\) ⊇ *P*\(^{(*n* + 1)}\) ⊇ ···, and he continued the derivation process by defining *P*\(^{(∞)}\) as the intersection of these sets. Then he iterated the derived set operation and intersections to extend his sequence of sets into the infinite: *P*\(^{(∞)}\) ⊇ *P*\(^{(∞ + 1)}\) ⊇ *P*\(^{(∞ + 2)}\) ⊇ ··· ⊇ *P*\(^{(2∞)}\) ⊇ ··· ⊇ *P*\(^{(∞2}\)) ⊇ ···. The superscripts containing ∞ are just indices defined by the derivation process.

Cantor used these sets in the theorems:
'' is countable;
|Conversely, if ** is countable, then there is an index *α* such that *P*\(^{(*α*)}\) = ∅.

These theorems are proved by partitioning **' into pairwise disjoint sets: **'\ *P*\(^{(2)}\)) ∪ (*P*\(^{(2)}\) \ *P*\(^{(3)}\)) ∪ ··· ∪ (*P*\(^{(∞)}\) \ *P*\(^{(∞ + 1)}\)) ∪ ··· ∪ *P*\(^{(*α*)}\). For *β* < *α*: since *P*\(^{(*β* + 1)}\) contains the limit points of *P*\(^{(*β*)}\), the sets *P*\(^{(*β*)}\) \ *P*\(^{(*β* + 1)}\) have no limit points. Hence, they are discrete sets, so they are countable. Proof of first theorem: If *P*\(^{(*α*)}\) = ∅ for some index *α*, then **' is the countable union of countable sets. Therefore, **' is countable.

The second theorem requires proving the existence of an *α* such that *P*\(^{(*α*)}\) = ∅. To prove this, Cantor considered the set of all *α* having countably many predecessors. To define this set, he defined the transfinite ordinal numbers and transformed the infinite indices into ordinals by replacing ∞ with *ω*, the first transfinite ordinal number. Cantor called the set of finite ordinals the first **number class**. The second number class is the set of ordinals whose predecessors form a countably infinite set. The set of all *α* having countably many predecessors—that is, the set of countable ordinals—is the union of these two number classes. Cantor proved that the cardinality of the second number class is the first uncountable cardinality.

Cantor's second theorem becomes: If **' is countable, then there is a countable ordinal *α* such that *P*\(^{(*α*)}\) = ∅. Its proof uses proof by contradiction. Let **' be countable, and assume there is no such α. This assumption produces two cases.

* Case 1: *P*\(^{(*β*)}\) \ *P*\(^{(*β* + 1)}\) is non-empty for all countable *β*. Since there are uncountably many of these pairwise disjoint sets, their union is uncountable. This union is a subset of **, so *P' * is uncountable.
* Case 2: *P*\(^{(*β*)}\) \ *P*\(^{(*β* + 1)}\) is empty for some countable *β*. Since *P*\(^{(*β* + 1)}\) ⊆ *P*\(^{(*β*)}\), this implies *P*\(^{(*β* + 1)}\) = *P*\(^{(*β*)}\). Thus, *P*\(^{(*β*)}\) is a perfect set, so it is uncountable. Since *P*\(^{(*β*)}\) ⊆ *, the set *'' is uncountable.

In both cases, **' is uncountable, which contradicts **' being countable. Therefore, there is a countable ordinal *α* such that *P*\(^{(*α*)}\) = ∅. Cantor's work with derived sets and ordinal numbers led to the Cantor-Bendixson theorem.

Using successors, limits, and cardinality, Cantor generated an unbounded sequence of ordinal numbers and number classes. The (*α* + 1)-th number class is the set of ordinals whose predecessors form a set of the same cardinality as the *α*-th number class. The cardinality of the (*α* + 1)-th number class is the cardinality immediately following that of the *α*-th number class. For a limit ordinal *α*, the *α*-th number class is the union of the *β*-th number classes for *β* < *α*. Its cardinality is the limit of the cardinalities of these number classes.

If *n* is finite, the *n*-th number class has cardinality \(\aleph_{n-1}\). If *α* ≥ *ω*, the *α*-th number class has cardinality.. Mathematical induction proves that the *n*-th number class has cardinality \(\aleph_{n-1}\). Since the *ω*-th number class is the union of the *n*-th number classes, its cardinality the limit of the \(\aleph_{n-1}\). Transfinite induction proves that if *α* ≥ *ω*, the *α*-th number class has cardinality. Therefore, the cardinalities of the number classes correspond one-to-one with the aleph numbers. Also, the *α*-th number class consists of ordinals different from those in the preceding number classes if and only if *α* is a non-limit ordinal. Therefore, the non-limit number classes partition the ordinals into pairwise disjoint sets.
