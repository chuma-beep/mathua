> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Cardinality) — CC BY-SA 4.0

# Last number equals total count

In mathematics, **cardinality** is an inherent property of sets, roughly meaning the number of individual objects they contain, which may be infinite. The concept is understood through one-to-one correspondences between sets. That is, if their objects can be paired such that each object has a pair, and no object is paired more than once.

Two sets are said to be ***equinumerous*** or *have the same cardinality* if there exists a one-to-one correspondence between them. Otherwise, one is said to be *strictly larger* or *strictly smaller* than the other. Using this concept, it is possible to show there are different sizes of infinity.

A set is countably infinite if it can be placed in one-to-one correspondence with the set of natural numbers \(\{1,2,3,4,\cdots\}\). For example, the set of even numbers \(\{2,4,6,\cdots\}\) and the set of rational numbers are countable. Uncountable sets are those strictly larger than the set of natural numbers. The set of all real numbers and the powerset of the set of natural numbers are proven to be uncountable by so-called diagonal arguments. Cantor's theorem generalizes these arguments to show there is an infinite hierarchy of infinities.

For finite sets, cardinality recovers the usual concept of size as "number of elements." However, it is more often difficult to ascribe "sizes" to infinite sets. A system of cardinal numbers can be developed to extend the role of natural numbers in answering "how many". Most commonly, the Aleph numbers \(\aleph_0, \aleph_1, \aleph_2,... \aleph_\omega, \aleph_{\omega+1}...\) are used; since their definition naturally extends the process of counting, and it can be shown that every infinite set has cardinality equivalent to some Aleph.

The set of natural numbers has cardinality. The question of whether the real numbers have cardinality is known as the continuum hypothesis, which has been shown to be both unprovable and undisprovable in standard set theories such as Zermelo–Fraenkel set theory. Alternative set theories and additional axioms give rise to different properties and have often strange or unintuitive consequences. However, every theory of cardinality using standard logical foundations of mathematics admits Skolem's paradox.

The basic concepts of cardinality go back as early as the 6th century BCE, and there are several close encounters with it throughout history, however, the results were generally dismissed as paradoxical. It is considered to have been first introduced formally to mathematics by Georg Cantor at the turn of the 20th century. Cantor's theory of cardinality was then formalized, popularized, and explored by many influential mathematicians of the time, and has since become a fundamental concept of mathematics.

## Basics
### Definition
Cardinality is an inherent property of sets which defines their size, roughly meaning the number of individual objects they contain. Fundamentally however, it is different from the concepts of number or counting as the cardinalities of two sets can be compared without referring to their number of elements, or defining number at all. For example, in the image above, a set of apples is compared to a set of oranges such that every fruit is used exactly once which shows these two sets have the same cardinality, even if one doesn't know how many of each there are. Thus, cardinality is measured by putting sets in one-to-one correspondence. If it is possible, the sets are said to have the *same cardinality*, and if not, one set is said to be *strictly larger* or *strictly smaller* than the other.

### Sets and functions
| not surjective | surjective |
| --- | --- |
| not injective **general** **function** **surjective** **only** |  |
| injective **injective** **only** **bijective** |  |

The basic concepts of cardinality are developed in terms of sets and functions, which are somewhat more abstract than their counterparts outside of mathematics. Informally, a set can be understood as any collection of objects, usually represented with curly braces. For example, \(S = \{1,2,3\}\) specifies a set, called , which contains the numbers 1, 2, and 3. The symbol represents set membership, for example says "1 is a member of the set " which is true by the definition of above. Here is finite, but that is not a requirement in general. The only requirement for a set is that it is well-defined. That is, for any object, one can determine whether belongs to that set , or does not belong to that set. One example of an infinite set is the set of all natural numbers \(\{ 1,2,3,\cdots \}\).

A function, or correspondence, maps members of one set to the members of another, often represented with an arrow diagram. For example, the adjacent table depicts several functions which map sets of natural numbers to sets of letters. If a function does not map two members to the same place, it is called injective. If a function covers every member in the output set, it is called surjective. If a function is both injective and surjective, it is called bijective or a one-to-one correspondence. Functions are not limited to those one can draw an arrow diagram for, so long as the function is well-defined. That is, for each possible input, one can determine the output. For example, one may define a function on the natural numbers by multiplying by two:

### Etymology and related terms
The term *cardinality* originates from the post-classical Latin *cardo* ("to hinge"), which referred to something central or pivotal, both literally and metaphorically. This passed into medieval Latin and then into English, where *cardinal* came to describe things considered to be, in some sense, fundamental, such as, *cardinal sins*, *cardinal directions*, and (in linguistics) *cardinal numbers*. The last of which referred to numbers used for counting (e.g., *one*, *two*, *three*), as opposed to *ordinal numbers*, which express order (e.g., *first, second, third*), and *nominal numbers* used for labeling without meaning (e.g., jersey numbers and serial numbers).

In mathematics, the notion of cardinality was first introduced by German mathematician Georg Cantor in the late 19th century, who used the term *Mächtigkeit*, which may be translated as "magnitude" or "power", though Cantor credited the term to a work by Jakob Steiner on projective geometry. Around 1930, the terms *cardinality* and *cardinal number* were adopted from the grammatical sense, and later translations would use these terms.

## Comparing sets
### Equinumerosity

An intuitive relationship between two sets having the "same size" is that their objects can be paired one-to-one. That is, if each object in one set can be assigned a dedicated "pair" in the other, and no object from either set left unpaired, it seems reasonable to say there are the same number of objects in each set. A one-to-one pairing between two sets defines a bijective function between them by mapping each object to its pair. Similarly, a bijection between two sets defines a pairing of their elements by pairing each object with the one it maps to. Therefore, these notions of "pairing" and "bijection" are, at least intuitively, equivalent. Thus, the following definition is given:

Two sets are said to have the *same cardinality* or be *equinumerous* if their members can be paired one-to-one. That is, if there exists a function between them which is bijective. This is often written as , or. |

Alternatively, these sets may be said to be *equivalent*, *similar*, *equipotent*, or *equipollent*.
|
|

|Similar
|Equinumerous
|Equipotent
|Equipollent

For example, the set \(E = \{0, 2, 4, 6, \text{...}\}\) of even numbers has the same cardinality as the set \(\N = \{0, 1, 2, 3, \text{...}\}\) of natural numbers; since the function is a bijection from to. The property for finite sets that "the whole is greater than the part" is no longer true for infinite sets, and the existence of surjections or injections that don't work does not prove that there is no bijection. For example, the function from to , defined by multiplying by 4, is injective, but not surjective (since 2, for instance, is not mapped to). Further, the function from to , defined by rounding down to the nearest even number, \(h(n) = 2 \operatorname{floor}(n/2)\) (cf. *floor function*), is surjective, but not injective, (since 0 and 1 for instance both map to 0). Neither nor can challenge , which was established by the existence of. #### Equivalence

A fundamental result in developing a theory of cardinality is that equinumerosity forms an equivalence relationthat is, a relation satisfying the same three basic properties as equality: reflexivity, symmetry, and transitivity. 

Reflexivity, the property that every set has the same cardinality as itself , follows from the identity function: for any set , the function that maps each element to itself is a bijection from to ****. Symmetry, the property that, if has the same cardinality as , then has the same cardinality as , holds because any bijection has an inverse function \(f^{-1}:B \to A\), which is also bijective. Transitivity, the property that if and have the same cardinality , and and have the same cardinality then so do and , follows from the composition of functions: given bijections and , their composition is a bijection from to (see image).

Since equinumerosity satisfies all three properties, it is an equivalence relation. This means that it groups sets into equivalence classesgroups of sets that are all equinumerous with one anotherwhere each group defines a possible size of a set. This motivates the notion of cardinal numbers, which are representatives that stand for the "size" of each group, developed in the section below.

### Inequality

A set is not larger than a set if it can be mapped into without overlap. That is, the cardinality of is less than or equal to the cardinality of if there is an injective function from to ****. This is written , and eventually ,
|
|

and read as " is *not greater than* ," or " is *dominated by*. " If , but there is no injection from to , then is said to be *strictly* smaller than , written without the underline as or. For example, if has four elements and has five, then the following are true , and. The basic properties of an inequality are reflexivity (for any , ), transitivity (if and , then ) and antisymmetry (if and , then ). Cardinal inequality as defined above is reflexive since the identity function is injective, and is transitive by function composition. Antisymmetrythe fact that, given injections in both directions, then one can find a bijectionis non-trivial, and is the content of the so-called Schröder–Bernstein theorem. One such proof can be summarized as follows, with the help of the adjacent illustration:

Given injections in both directions (, and ), take set and "color" each point blue. Then take the image of set on after the injectionthat is, the set of points that the injection projects on toand color those red. Then take the image of on this image of , and color those points blue again. By recursively alternating the images of these sets, one obtains a "coloring" of the original set (see illustration). For each point, it either stops at some finite stepand thus is definitively colored blue or redor continues switching indefinitely. Of the points that are definitively colored blue, map these points onto the next recursive image (i.e. by applying then ), leaving all other points in place. The resulting set is exactly the image of. Taking this transformation, followed by the inverse of , gives a bijection from to. A further property of cardinal inequality is totality, which says that any two sets are comparable. That is, for any and , either or. A full proof of which requires concepts introduced later, however the argument can be briefly summarized as follows. Every well-ordered set is isomorphic to a unique ordinal number, called the order type of the set. By the well-ordering theorem, every set can be well-ordered. Then, by comparing their order types, one can show that or. This fact is equivalent to the axiom of choice.

## Countability
### Countable sets
A set is called *countable* if it is finite or has a bijection with the set of natural numbers , in which case it is called *countably infinite*. The term *denumerable* is also sometimes used for countably infinite sets. For example, the set of all even natural numbers is countable, and therefore has the same cardinality as the whole set of natural numbers, even though it is a proper subset. Similarly, the set of square numbers is countable, which was considered paradoxical for hundreds of years before modern set theory (cf. **). However, several other examples have historically been considered surprising or initially unintuitive since the rise of set theory.

 to. On the left, a version for the positive rational numbers. On the right, a spiral for all pairs of integers for each fraction. The rational numbers are those which can be expressed as the quotient or fraction of two integers. The rational numbers can be shown to be countable by considering the set of fractions as the set of all ordered pairs of integers, which can be visualized as the set of all integer points on a grid. Then, an intuitive function can be described by drawing a line in a repeating pattern, or spiral, which eventually goes through each point in the grid. For example, going through each diagonal on the grid for positive fractions, or through a lattice spiral for all integer pairs. These technically over cover the rationals; since, for example, the rational number \(\textstyle\frac{1}{2}\) gets mapped to by all the fractions \(\textstyle\frac{2}{4},\, \frac{3}{6}, \, \frac{4}{8}, \, \dots\), as the grid method treats these all as distinct ordered pairs. So this function shows not. This can be corrected by "skipping over" these numbers in the grid, using the Schröder–Bernstein theorem, or by designing a function which does this naturally, for example using the Calkin–Wilf tree.

A number is called algebraic if it is a solution of some polynomial equation (with integer coefficients). For example, the square root of two is a solution to , and the rational number is the solution to. Conversely, a number which cannot be the root of any polynomial is called transcendental. Two examples include Euler's number (**') and pi (). In general, proving a number is transcendental is considered to be very difficult, and only a few classes of transcendental numbers are known. However, it can be shown that the set of algebraic numbers is countable by ordering the polynomials lexicographically (for example, see **'). Since the set of algebraic numbers is countable while the real numbers are uncountable (shown in the following section), the transcendental numbers must form the vast majority of real numbers, even though they are individually much harder to identify. That is to say, almost all real numbers are transcendental.

#### Hilbert's hotel

Hilbert's paradox of the Grand Hotel is a popular thought experiment devised by the German mathematician David Hilbert to illustrate a counterintuitive property of countably infinite sets, allowing them to have the same cardinality as a proper subset of themselves. The scenario begins by imagining a hotel with an infinite number of rooms, one for each natural number, all of which are occupied. But then a new guest walks in asking for a room. The hotel accommodates by moving the occupant of room 1 to room 2, the occupant of room 2 to room 3, room 3 to room 4, and in general, room n to room n+1. Then every guest still has a room, but room 1 is open for the new guest.

Then, the scenario continues by imagining an infinitely long bus of new guests seeking a room. The hotel accommodates by moving the person in room 1 to room 2, room 2 to room 4, and in general, room n to room 2n. Thus, all the even-numbered rooms are occupied, but all the odd-numbered rooms are vacant, leaving room for the infinite bus of new guests. The scenario continues further by assuming an infinite number of these infinite buses arrive at the hotel, and showing that the hotel is still able to accommodate. Finally, an infinite bus which has a seat for every real number arrives, and the hotel is no longer able to accommodate.

### Uncountable sets

A set is called *uncountable* if it is not countable; that is, it is infinite and strictly larger than the set of natural numbers. The usual first example of this is the set of real numbers , which can be understood as the set of all numbers on the number line. One method of proving that the reals are uncountable is called Cantor's diagonal argument, credited to Cantor for his 1891 proof, though his method differs from the more common presentation.
