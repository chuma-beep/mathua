> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Cardinality) — CC BY-SA 4.0

# Last number equals total count

In mathematics, **cardinality** is an inherent property of sets, roughly meaning the number of individual objects they contain, which may be infinite. The concept is understood through one-to-one correspondences between sets. That is, if their objects can be paired such that each object has a pair, and no object is paired more than once.

Two sets are said to be ***equinumerous*** or *have the same cardinality* if there exists a one-to-one correspondence between them. Otherwise, one is said to be *strictly larger* or *strictly smaller* than the other. Using this concept, it is possible to show there are different sizes of infinity.

A set is countably infinite if it can be placed in one-to-one correspondence with the set of natural numbers {{tmath|1=\{1,2,3,4,\cdots\} }}. For example, the set of even numbers {{tmath|1=\{2,4,6,\cdots\} }} and the set of rational numbers are countable. Uncountable sets are those strictly larger than the set of natural numbers. The set of all real numbers and the powerset of the set of natural numbers are proven to be uncountable by so-called diagonal arguments. Cantor's theorem generalizes these arguments to show there is an infinite hierarchy of infinities.

For finite sets, cardinality recovers the usual concept of size as "number of elements." However, it is more often difficult to ascribe "sizes" to infinite sets. A system of cardinal numbers can be developed to extend the role of natural numbers in answering "how many". Most commonly, the Aleph numbers {{tmath|1=\aleph_0, \aleph_1, \aleph_2, ... \aleph_\omega, \aleph_{\omega+1}... }} are used, since their definition naturally extends the process of counting, and it can be shown that every infinite set has cardinality equivalent to some Aleph.

The set of natural numbers has cardinality . The question of whether the real numbers have cardinality  is known as the continuum hypothesis, which has been shown to be both unprovable and undisprovable in standard set theories such as Zermelo–Fraenkel set theory. Alternative set theories and additional axioms give rise to different properties and have often strange or unintuitive consequences. However, every theory of cardinality using standard logical foundations of mathematics admits Skolem's paradox.

The basic concepts of cardinality go back as early as the 6th century BCE, and there are several close encounters with it throughout history, however, the results were generally dismissed as paradoxical. It is considered to have been first introduced formally to mathematics by Georg Cantor at the turn of the 20th century. Cantor's theory of cardinality was then formalized, popularized, and explored by many influential mathematicians of the time, and has since become a fundamental concept of mathematics.

## Basics
### Definition
Cardinality is an inherent property of sets which defines their size, roughly meaning the number of individual objects they contain. Fundamentally however, it is different from the concepts of number or counting as the cardinalities of two sets can be compared without referring to their number of elements, or defining number at all. For example, in the image above, a set of apples is compared to a set of oranges such that every fruit is used exactly once which shows these two sets have the same cardinality, even if one doesn't know how many of each there are. Thus, cardinality is measured by putting sets in one-to-one correspondence. If it is possible, the sets are said to have the *same cardinality*, and if not, one set is said to be *strictly larger* or *strictly smaller* than the other.

### Sets and functions
{| class="wikitable floatright" style="text-align:center"
!
!not surjective
!surjective
|-
!not
injective
|
**general** **function**
|
**surjective** **only**
|-
!injective
|
**injective** **only**
|
**bijective**
|}
The basic concepts of cardinality are developed in terms of sets and functions, which are somewhat more abstract than their counterparts outside of mathematics. Informally, a set can be understood as any collection of objects, usually represented with curly braces. For example, {{tmath|1=S = \{1,2,3\} }} specifies a set, called , which contains the numbers 1, 2, and 3. The symbol  represents set membership, for example  says "1 is a member of the set " which is true by the definition of  above. Here  is finite, but that is not a requirement in general. The only requirement for a set is that it is well-defined. That is, for any object, , one can determine whether  belongs to that set , or  does not belong to that set . One example of an infinite set is the set of all natural numbers {{tmath|1=\{ 1,2,3,\cdots \} }}.}}

A function, or correspondence, maps members of one set to the members of another, often represented with an arrow diagram. For example, the adjacent table depicts several functions which map sets of natural numbers to sets of letters. If a function does not map two members to the same place, it is called injective. If a function covers every member in the output set, it is called surjective. If a function is both injective and surjective, it is called bijective or a one-to-one correspondence. Functions are not limited to those one can draw an arrow diagram for, so long as the function is well-defined. That is, for each possible input, one can determine the output. For example, one may define a function on the natural numbers by multiplying by two:

### Etymology and related terms
The term *cardinality* originates from the post-classical Latin *cardo* ("to hinge"), which referred to something central or pivotal, both literally and metaphorically. This passed into medieval Latin and then into English, where *cardinal* came to describe things considered to be, in some sense, fundamental, such as, *cardinal sins*, *cardinal directions*, and (in linguistics) *cardinal numbers*. The last of which referred to numbers used for counting (e.g., *one*, *two*, *three*), as opposed to *ordinal numbers*, which express order (e.g., *first, second, third*), and *nominal numbers* used for labeling without meaning (e.g., jersey numbers and serial numbers).

In mathematics, the notion of cardinality was first introduced by German mathematician Georg Cantor in the late 19th century, who used the term *Mächtigkeit*, which may be translated as "magnitude" or "power", though Cantor credited the term to a work by Jakob Steiner on projective geometry. Around 1930, the terms *cardinality* and *cardinal number* were adopted from the grammatical sense, and later translations would use these terms.}}

## Comparing sets
### Equinumerosity


An intuitive relationship between two sets having the "same size" is that their objects can be paired one-to-one. That is, if each object in one set can be assigned a dedicated "pair" in the other, and no object from either set left unpaired, it seems reasonable to say there are the same number of objects in each set. A one-to-one pairing between two sets defines a bijective function between them by mapping each object to its pair. Similarly, a bijection between two sets defines a pairing of their elements by pairing each object with the one it maps to. Therefore, these notions of "pairing" and "bijection" are, at least intuitively, equivalent. Thus, the following definition is given:

Two sets are said to have the *same cardinality* or be *equinumerous* if their members can be paired one-to-one. That is, if there exists a function between them which is bijective. This is often written as , or .
|
}}
}} Alternatively, these sets may be said to be *equivalent*, *similar*, *equipotent*, or *equipollent*.
|
|
}}

|Similar
|Equinumerous
|Equipotent
|Equipollent
}}
}} For example, the set {{tmath|1=E = \{0, 2, 4, 6, \text{...}\} }} of even numbers has the same cardinality as the set {{tmath|1=\N = \{0, 1, 2, 3, \text{...}\} }} of natural numbers, since the function  is a bijection from  to .

The property for finite sets that "the whole is greater than the part" is no longer true for infinite sets, and the existence of surjections or injections that don't work does not prove that there is no bijection. For example, the function  from  to , defined by multiplying by 4, , is injective, but not surjective (since 2, for instance, is not mapped to). Further, the function  from  to , defined by rounding down to the nearest even number, {{tmath|1=h(n) = 2 \operatorname{floor}(n/2) }} (cf. *floor function*), is surjective, but not injective, (since 0 and 1 for instance both map to 0). Neither  nor  can challenge , which was established by the existence of .

#### Equivalence

A fundamental result in developing a theory of cardinality is that equinumerosity forms an equivalence relationthat is, a relation satisfying the same three basic properties as equality: reflexivity, symmetry, and transitivity.  }}

Reflexivity, the property that every set has the same cardinality as itself , follows from the identity function: for any set , the function that maps each element to itself is a bijection from  to ****. Symmetry, the property that, if  has the same cardinality as  , then  has the same cardinality as  , holds because any bijection   has an inverse function {{tmath|1=f^{-1}:B \to A }}, which is also bijective. Transitivity, the property that if  and  have the same cardinality , and  and  have the same cardinality  then so do  and  , follows from the composition of functions: given bijections  and , their composition  is a bijection from  to  (see image).

Since equinumerosity satisfies all three properties, it is an equivalence relation. This means that it groups sets into equivalence classesgroups of sets that are all equinumerous with one anotherwhere each group defines a possible size of a set. This motivates the notion of cardinal numbers, which are representatives that stand for the "size" of each group, developed in the section below.

### Inequality

A set  is not larger than a set  if it can be mapped into  without overlap. That is, the cardinality of  is less than or equal to the cardinality of  if there is an injective function from  to ****. This is written ,  and eventually ,
|
|
}}
}} and read as " is *not greater than* ," or " is *dominated by* ." If , but there is no injection from  to , then  is said to be *strictly* smaller than , written without the underline as  or . For example, if  has four elements and  has five, then the following are true , , and .

The basic properties of an inequality are reflexivity (for any , ), transitivity (if  and , then ) and antisymmetry (if  and , then ). Cardinal inequality  as defined above is reflexive since the identity function is injective, and is transitive by function composition. Antisymmetrythe fact that, given injections in both directions, then one can find a bijectionis non-trivial, and is the content of the so-called Schröder–Bernstein theorem. One such proof can be summarized as follows, with the help of the adjacent illustration:

Given injections in both directions (, and ), take set  and "color" each point blue.  Then take the image of set  on  after the injectionthat is, the set of points that the injection projects on toand color those red. Then take the image of  on this image of , and color those points blue again. By recursively alternating the images of these sets, one obtains a "coloring" of the original set (see illustration). For each point, it either stops at some finite stepand thus is definitively colored blue or redor continues switching indefinitely. Of the points that are definitively colored blue, map these points onto the next recursive image (i.e. by applying  then ), leaving all other points in place. The resulting set is exactly the image of . Taking this transformation, followed by the inverse of , gives a bijection from  to .

A further property of cardinal inequality is totality, which says that any two sets are comparable. That is, for any  and , either  or . A full proof of which requires concepts introduced later, however the argument can be briefly summarized as follows. Every well-ordered set is isomorphic to a unique ordinal number, called the order type of the set. By the well-ordering theorem, every set can be well-ordered. Then, by comparing their order types, one can show that  or . This fact is equivalent to the axiom of choice.

## Countability
### Countable sets
A set is called *countable* if it is finite or has a bijection with the set of natural numbers , in which case it is called *countably infinite*. The term *denumerable* is also sometimes used for countably infinite sets. For example, the set of all even natural numbers is countable, and therefore has the same cardinality as the whole set of natural numbers, even though it is a proper subset. Similarly, the set of square numbers is countable, which was considered paradoxical for hundreds of years before modern set theory (cf. **). However, several other examples have historically been considered surprising or initially unintuitive since the rise of set theory.

 to . On the left, a version for the positive rational numbers. On the right, a spiral for all pairs of integers  for each fraction .
}}

The rational numbers  are those which can be expressed as the quotient or fraction  of two integers. The rational numbers can be shown to be countable by considering the set of fractions as the set of all ordered pairs of integers, which can be visualized as the set of all integer points on a grid. Then, an intuitive function can be described by drawing a line in a repeating pattern, or spiral, which eventually goes through each point in the grid. For example, going through each diagonal on the grid for positive fractions, or through a lattice spiral for all integer pairs. These technically over cover the rationals, since, for example, the rational number {{tmath|1= \textstyle\frac{1}{2} }} gets mapped to by all the fractions {{tmath|1= \textstyle\frac{2}{4},\, \frac{3}{6}, \, \frac{4}{8}, \, \dots }}, as the grid method treats these all as distinct ordered pairs. So this function shows  not . This can be corrected by "skipping over" these numbers in the grid, using the Schröder–Bernstein theorem, or by designing a function which does this naturally, for example using the Calkin–Wilf tree.

A number is called algebraic if it is a solution of some polynomial equation (with integer coefficients). For example, the square root of two  is a solution to , and the rational number  is the solution to . Conversely, a number which cannot be the root of any polynomial is called transcendental. Two examples include Euler's number (**') and pi (). In general, proving a number is transcendental is considered to be very difficult, and only a few classes of transcendental numbers are known. However, it can be shown that the set of algebraic numbers is countable by ordering the polynomials lexicographically (for example, see **'). Since the set of algebraic numbers is countable while the real numbers are uncountable (shown in the following section), the transcendental numbers must form the vast majority of real numbers, even though they are individually much harder to identify. That is to say, almost all real numbers are transcendental.

#### Hilbert's hotel

Hilbert's paradox of the Grand Hotel is a popular thought experiment devised by the German mathematician David Hilbert to illustrate a counterintuitive property of countably infinite sets, allowing them to have the same cardinality as a proper subset of themselves. The scenario begins by imagining a hotel with an infinite number of rooms, one for each natural number, all of which are occupied. But then a new guest walks in asking for a room. The hotel accommodates by moving the occupant of room 1 to room 2, the occupant of room 2 to room 3, room 3 to room 4, and in general, room n to room n+1. Then every guest still has a room, but room 1 is open for the new guest.

Then, the scenario continues by imagining an infinitely long bus of new guests seeking a room. The hotel accommodates by moving the person in room 1 to room 2, room 2 to room 4, and in general, room n to room 2n. Thus, all the even-numbered rooms are occupied, but all the odd-numbered rooms are vacant, leaving room for the infinite bus of new guests. The scenario continues further by assuming an infinite number of these infinite buses arrive at the hotel, and showing that the hotel is still able to accommodate. Finally, an infinite bus which has a seat for every real number arrives, and the hotel is no longer able to accommodate.

### Uncountable sets

A set is called *uncountable* if it is not countable; that is, it is infinite and strictly larger than the set of natural numbers. The usual first example of this is the set of real numbers , which can be understood as the set of all numbers on the number line. One method of proving that the reals are uncountable is called Cantor's diagonal argument, credited to Cantor for his 1891 proof, though his method differs from the more common presentation.

{{Quote box
| quote = {{tmath|1=\begin{align}
1 \to & \; 0.{ \color{red} \textbf{6} }18033... \\
2 \to & \; 0.1{ \color{red}\textbf{2} }3456... \\
3 \to & \; 0.60{ \color{red}\textbf{7} }927... \\
4 \to & \; 0.222{ \color{red}\textbf{2} }22... \\
\vdots \\
& \; 0.{ \color{red}\textbf{2323}22}...
\end{align} }}
| border = 0px
| fontsize = 120%
| bgcolor = #00000000
}}

It begins by assuming, by contradiction, that there is some one-to-one mapping between the natural numbers and the set of real numbers between 0 and 1 (the interval ). Then, take the decimal representation of each real number, for example,  with a leading zero followed by any sequence of digits. The number 1 is included in this set since 1  Considering these real numbers in a column, it is always possible to create a new number such that the first digit of the new number is different from that of the first number in the column, the second digit is different from the second number in the column, and so on. The new number must also have a unique decimal representation, that is, it cannot end in repeating nines or repeating zeros. For example, if the digit isn't 2, make the digit of the new number 2, and if it was 2, make it 3. Then, this new number will be different from each of the numbers in the list by at least one digit, and therefore must not be in the list. This shows that the real numbers cannot be put into a one-to-one correspondence with the naturals, and thus must be strictly larger.

Another classical example of an uncountable set, established using a related reasoning, is the power set of the natural numbers, denoted {{tmath|1=\mathcal{P}(\N) }}. This is the set of all subsets of , including the empty set and  itself. The method is much closer to Cantor's original diagonal argument. Again, assume by contradiction that there exists a one-to-one correspondence  between  and {{tmath|1=\mathcal{P}(\N) }}, so that every subset of  is assigned to some natural number. These subsets are then placed in a column, in the order defined by  (see image). Now, one may define a subset  of  which is not in the list by taking the negation of the "diagonal" of this column as follows:

If , then , that is, if 1 is in the first subset of the list, then 1 is *not* in the subset . Further, if , then , that is if the number 2 is not in the second subset of the column, then 2 *is* in the subset . Then in general, for each natural number ,  if and only if , meaning  is put in the subset  only if the nth subset in the column does not contain the number . Then, for each natural number , , meaning,  is not the nth subset in the list, for any number , and so it cannot appear anywhere in the list defined by . Since  was chosen arbitrarily, this shows that every function from  to {{tmath|1=\mathcal{P}(\N) }} must be missing at least one element, therefore no such bijection can exist, and so {{tmath|1=\mathcal{P}(\N) }} must be not be countable.

These two sets,  and {{tmath|1=\mathcal{P}(\N) }} can be shown to have the same cardinality (by, for example, assigning each subset to a decimal expansion) called the cardinality of the continuum.

Cantor's theorem generalizes the second theorem above, showing that every set is strictly smaller than its powerset. In broad strokes, the proof goes as follows: Given a set , assume by contradiction that there is a bijection  from  to {{tmath|1=\mathcal{P}(A) }}. Then, the subset  given by taking the negation of the "diagonal", formally, {{tmath|1=T = \{ a \in A : a \notin f(a) \} }}, cannot be in the list. Therefore, every function is missing at least one element, and so {{tmath|1=\vert A \vert < \vert\mathcal{P}(A)\vert }}. Further, since {{tmath|1=\mathcal{P}(A) }} is itself a set, the argument can be repeated to show {{tmath|1=\vert A \vert < \vert \mathcal{P}(A)\vert < \vert\mathcal{P}(\mathcal{P}(A))\vert }}. Taking , this shows that {{tmath|1=\mathcal{P}(\mathcal{P}(\N)) }} is even larger than {{tmath|1=\mathcal{P}(\N) }}, which was already shown to be uncountable. Repeating this argument shows that there are infinitely many "sizes" of infinity.

## Cardinal numbers
In the above sections, "the cardinality of a set" was described relationally. In other words, one set could be compared to another, intuitively comparing their "size". Cardinal numbers are a means of measuring this "size" more explicitly. For finite sets, this is simply the natural number found by counting the elements. This number is called the *cardinal number* of that set, or simply *the cardinality* of that set. The cardinal number of a set  is generally denoted by , with a vertical bar on each side, though it may also be denoted by  , {{tmath|1=\operatorname{card}(A) }}, or .
| 2= {{tmath|1=\operatorname{card}(A) }}
| 3=
}}
}}

For infinite sets, "cardinal number" is somewhat more difficult to define formally. Cardinal numbers are not usually thought of in terms of their formal definition, but immaterially in terms of their arithmetic/algebraic properties. The assumption that there is *some* cardinal function  which satisfies , sometimes called the *axiom of cardinality* or *Hume's principle*, is sufficient for deriving most properties of cardinal numbers.

Commonly in mathematics, if a relation satisfies the properties of an equivalence relation, the objects used to materialize this relation are equivalence classes, which groups all the objects equivalent to one another. These called the *FregeRussell* cardinal numbers. However, this would mean that cardinal numbers are too large to form sets (apart from the cardinal number  whose only element is the empty set), since, for example, the cardinal number  would be the set of all sets with one element, then {{tmath|1=\{ \bold 1 \} \in \bold 1 }}, and would therefore contain itself, violating regularity. Thus, due to John von Neumann, it is more common to assign representatives of these classes.

### Finite cardinals


Given a basic sense of natural numbers, a set is said to have cardinality  if it can be put in one-to-one correspondence with the set {{tmath|\{1,\,2,\, \dots, \, n \} }}, analogous to counting its elements. For example, the set {{tmath|1=S = \{ A,B,C,D \} }} has a natural correspondence with the set {{tmath|\{1,2,3,4\} }}, and therefore is said to have cardinality 4. Other terminologies include "Its cardinality is 4" or "Its cardinal number is 4". In formal contexts, the natural numbers can be understood as some construction of objects satisfying the Peano axiomsa list of properties, such that any system satisfying these properties is, in a certain sense, just like the natural numbers.

Showing that such a correspondence exists is not always trivial. Combinatorics is the area of mathematics primarily concerned with counting, both as a means and as an end to obtaining results, and certain properties of finite structures. The notion cardinality of finite sets is closely tied to many basic combinatorial principles, and provides a set-theoretic foundation to recover them. It can be shown by induction on the possible sizes of sets that finite cardinality corresponds uniquely with natural numbers (cf. **). This is related to several other concepts, verifying Hume's principle and the basis of bijective proofs, and is equivalent to a certain formulation of the pigeonhole principle, that a finite set cannot be put in one to one correspondence with a proper subset of itself.

The addition principle asserts that given disjoint sets  and , , intuitively meaning that the sum of the parts is equal to the whole. The multiplication principle asserts that given two sets  and , , intuitively meaning that there are  ways to pair objects from these sets. Both of these can be proven by a bijective proof, together with induction. The more general result is the inclusion–exclusion principle, which defines how to count the number of elements in overlapping sets.

Naturally, a set is defined to be finite if it is empty or can be put in correspondence with the set {{tmath|\{1,\,2,\, \dots, \, n\} }}, for some natural number . However, there exist other definitions of "finite" which don't rely on a definition of "number." For example, a set is called Dedekind-finite if it cannot be put in one-to-one correspondence with a proper subset of itself, though this definition requires the axiom of choice.

### Aleph numbers

The aleph numbers are a sequence of cardinal numbers that represent the sizes of infinite sets, denoted with an aleph , the first letter of the Hebrew alphabet. The first aleph number is , called "aleph-nought", "aleph-zero", or "aleph-null", which represents the cardinality of the set of all natural numbers: {{tmath|1=\aleph_0 = \vert \N \vert = \vert\{0,1,2,3,\cdots\}\vert }}.  Then,  represents the next largest cardinality, then , and so on. The most common way this is formalized in set theory is through Von Neumann ordinals, known as Von Neumann cardinal assignment.

Ordinal numbers generalize the notion of *order* to infinite sets. For example, 2 comes after 1, denoted , and 3 comes after both, denoted . Then, one defines a new number, , which comes after every natural number, denoted . Further , and so on. More formally, these ordinal numbers can be defined as follows:

{{tmath|1=0 := \{\} }}, the empty set, {{tmath|1=1 := \{0\} }}, {{tmath|1=2 := \{0,1\} }}, {{tmath|1=3 := \{0,1,2\} }}, and so on. Then one can define {{tmath|1=m < n \text{, if } \, m \in n }}, for example, {{tmath|1=2 \in \{0,1,2\} = 3 }}, therefore .  Defining {{tmath|1=\omega := \{0,1,2,3,\cdots\} }} gives  the desired property of being the smallest ordinal greater than all finite ordinal numbers. Further, {{tmath|1=\omega+1 := \{0,1,\cdots,\omega\} }}, and so on.

Since  by the natural correspondence, one may define  as the set of all finite ordinals. That is, . Then,  is the set of all countable ordinals (all ordinals  with cardinality ), the first uncountable ordinal. Since a set cannot contain itself,  must have a strictly larger cardinality: .  Furthermore,  is the set of all ordinals with cardinality less than or equal to , and in general the successor cardinal is the set of all ordinals with cardinality up to . Put another way for infinite cardinals,  is the number of possible well-orderings on  up to order isomorphism. Proving that such a set always exists is known as Hartogs' theorem, wherein the smallest ordinal not less than or equal to than a set  is called the Hartogs number of . Then,  for a limit ordinal  is the union of all lesser alephs.

The importance of ordinal numbers here is to generalize the notion of counting to infinite sets. When counting, one implicitly assigns an order to their set of objects, but no matter what order one assigns the final result of the count is always the same, which demonstrates the connection between cardinal and ordinal numbers. Further, by the well-ordering theorem, there cannot exist any set with cardinality between  and , and every infinite set has some cardinality corresponding uniquely to some aleph , for some ordinal . This allows one to use a constructive definition of the cardinality function, by assigning each set to its equinumerous aleph.

### Cardinal arithmetic

Basic arithmetic can be done on cardinal numbers in a very natural way, by extending the theorems for finite combinatorial principles above. The intuitive principle that is  and  are disjoint then addition of these sets is simply taking their union, written as . Thus if  and  are infinite, cardinal addition is defined as  where  denotes disjoint union. Similarly, the multiplication of two sets is intuitively the number of ways to pair their elements (as in the multiplication principle), therefore cardinal multiplication is defined as , where  denotes the Cartesian product. These definitions can be shown to satisfy the basic properties of standard arithmetic:

* Associativity: , and
* Commutativity: , and
* Distributivity:

Although a lot of properties of finite arithmetic hold for infinite arithmetic, as above and in the table below, strict inequalities (such as Kőnig's theorem) are rare. For example, in finite arithmetic, for any nonzero number , . However, since both the set of even numbers  and set of odd numbers  have cardinality , it shows  thus . In fact, for any infinite cardinal, . In this way, infinite cardinal addition and multiplication are considered to be remarkably well-behaved (at least under the Axiom of Choice).

Cardinal exponentiation {{tmath|1=\vert A \vert^{\vert B \vert} }} is defined via set exponentiation, the set of all functions , that is, {{tmath|1=\vert A \vert^{\vert B \vert} := \vert A^B\vert }} which naturally extends the role of "repeated multiplication" to infinite sets. For finite sets this can be shown to coincide with standard natural number exponentiation, but includes as a corollary that zero to the power of zero is one  since there is exactly one function from the empty set to itself: the empty function. A combinatorial argument can be used to show {{tmath|1=2^{\vert A \vert} = \vert\mathcal{P}(A)\vert }}, by considering the indicator function of each subset. In general, cardinal exponentiation is not as well-behaved as addition and multiplication. For example, even though it can be proven that the expression {{tmath|1=2^{\aleph_0} }}does indeed correspond to some aleph, it is unprovable from standard set theories which aleph it corresponds to.
{| class="wikitable"
|+Basic results
!Exponentiation
!Inequality and Monotonicity
!Identity elements
!Absorption laws (for infinite A,B)
|-
|{{tmath|1=\big(\vert A \vert^{\vert B \vert} \big)^{\vert C \vert } = \vert A \vert^{\vert B \times C\vert} }}
| implies
|
|{{tmath|1=\vert A \vert+\vert B \vert = \operatorname{max}(\vert A \vert,\vert B \vert) }}
|-
|{{tmath|1=\vert A \vert^{\vert B \vert + \vert C \vert } = \vert A \vert^{\vert B \vert} \cdot \vert A \vert^{\vert C \vert } }}
| implies
|
|{{tmath|1=\vert A \vert\cdot\vert B \vert = \operatorname{max}(\vert A \vert,\vert B \vert) }}
|-
|{{tmath|1=\vert A \times B\vert ^{\vert C \vert } = \vert A \vert^{\vert C \vert } \cdot \vert B \vert^{\vert C \vert } }}
| implies  {{tmath|1=\vert A \vert^{\vert C \vert } \leq \vert B \vert^{\vert C \vert } }}
|{{tmath|1=\vert A \vert^{\bold 1} = \vert A \vert }}
| (annihilator)
|-
|
| and  implies {{tmath|1=\vert C \vert ^{\vert A \vert} \leq \vert C \vert ^{\vert B \vert} }}
|
|{{tmath|1=\bold 1^{\vert A \vert} = \bold 1 }}
|}

### Set of all cardinal numbers

The *set of all cardinal numbers* refers to a hypothetical set which contains all cardinal numbers. Such a set cannot exist, which has been considered paradoxical, and related to the Burali-Forti paradox. Using the definition of cardinal numbers as representatives of their cardinalities.  It begins by assuming there is some set {{tmath|1=S := \{ X \, \vert X \text{ is a cardinal number}\} }}. Then, if there is some largest cardinal , then the powerset  is strictly greater, and thus not in . Conversely, if there is no largest element, then the union  contains the elements of all elements of , and is therefore greater than or equal to each element. Since there is no largest element in , for any element , there is another element  such that  and . Thus, for any , , and so . Thus, the collection of all cardinal numbers is too large to form a set, and is a proper class. Georg Cantor denoted  the collection of all cardinal numbers **}}}}** (tav, the last letter of the Hebrew alphabet), and considered it to be an "inconsistent multiplicity".


## Cardinality of the continuum


The real numbers  formalize the intuitive notion of the *continuum*: the unbroken, gapless set of points on the number line. There are numerous ways to construct and analyze the continuum in set theory, for example, sets of Cauchy sequences of rational numbers, or Dedekind cuts. However, a somewhat informal definition as the set of infinite sequences of digits, e.g. , with the usual arithmetic and order is a popular and sufficient choice for basic arguments involving the continuum. The cardinality of this set, denoted "" (a lowercase fraktur script "c"), turns out to be remarkably stable under various transformations.

For example, all intervals on the real line e.g. , and , have the same cardinality as the entire set . First,  is a bijection from  to . The fact that two intervals of unequal length can be mapped one-to-one has been known, even if considered paradoxical, for thousands of years (cf. **). Further, the tangent function is a bijection from the interval {{tmath|1= \textstyle\left( \frac{-\pi}{2} \, ,  \frac{\pi}{2} \right) }} to the whole real line. A more surprising example is the Cantor set, which is defined as follows: take the interval  and remove the middle third {{tmath|1= \textstyle( \frac{1}{3}, \frac{2}{3} ) }}, then remove the middle third of each of the two remaining segments, and continue removing middle thirds (see image). The Cantor set is the collection of points that survive this process. The remaining points are exactly those whose decimal expansion can be written in ternary without a 1. Reinterpreting these decimal expansions as binary (e.g. by replacing the 2s with 1s) gives a bijection between the Cantor set and the interval . The Cantor set has a measure or "length" of zero, and yet is still equinumerous to the whole real line.

In the other direction, into higher dimensions, one can find a bijective map between a one-dimensional line and a two-dimensional square. Given a point in the unit square in  (two-dimensional space) , one can interleave their digits to get the unique number  in the unit interval .  Space-filling curves offer a more visual proof, giving a continuous map from the unit interval onto the unit square. Classical examples include the Peano curve and Hilbert curve. Although these maps are not bijective, they are indeed sufficient to show  with the converse being immediate. Both of these methods can be reused at each dimension to show that {{tmath|1= \vert \R^n \vert = \vert \R^{n+1}\vert }} and thus   for any dimension . Further, the infinite cartesian product  can also be shown to be equinumerous to  by cardinal arithmetic: {{tmath|1= \textstyle \vert \R^\infty\vert = \vert \R \vert^{\aleph_0} = (2^{\aleph_0} )^{\aleph_0} = 2^{(\aleph_0 \cdot \aleph_0)} = 2^{\aleph_0} = \vert \R \vert }}.  Thus, the real numbers, all finite-dimensional real spaces, and the countable cartesian product have the same cardinality.

As shown in , the set of real numbers is strictly larger than the set of natural numbers. Specifically, {{tmath|1=\vert \R \vert = \vert\mathcal{P}(\N)\vert  }}. The Continuum Hypothesis (CH) asserts that the real numbers have the next largest cardinality after the natural numbers, that is . As shown by Gödel and Cohen, the continuum hypothesis is independent of ZFC, a standard axiomatization of set theory; that is, it is impossible to prove the continuum hypothesis or its negation from ZFC—provided that ZFC is consistent, meaning it produces no contradictions. The Generalized Continuum Hypothesis (GCH) extends this to all infinite cardinals, stating that {{tmath|1=2^{\aleph_\alpha} = \aleph_{\alpha+1} }} for every ordinal . Research on CH and GCH continues independent of ZFC, especially in descriptive set theory and through the exploration of large cardinal axioms. Without GHC, the cardinality of  cannot be written in terms of specific alephs. The Beth numbers (  , the second letter of the Hebrew alphabet) provide a concise notation for powersets of the real numbers starting from , then {{tmath|1=\beth_1 = 2^{\beth_0} =\vert \R \vert }}, and {{tmath|1=\beth_2 = \vert\mathcal{P}(\R)\vert = 2^{\beth_1} }}, and in general {{tmath|1=\beth_{n+1} = 2^{\beth_n} }} and {{tmath|1=\beth_{\lambda} = \bigcup_{\alpha<\lambda} \beth_\alpha }} if  is a limit ordinal.

## Alternative and additional axioms
Around the turn of the 20th century, set theory turned to an axiomatic approach to avoid rampant foundational issues related to its naive study (cf. **). The most common axiomatic set theory used today is Zermelo–Fraenkel set theory (ZFC). In this system, the relevant axioms include: the Axiom of Infinity, stating roughly "there is an infinite set", specifically, a set with cardinality of the natural numbers ; the Axiom of power set, which says that, for any set , the powerset  also exists; and the Axiom of Choice, explained below. ZFC has been criticized both for being too strong, and too weak. Similarly, there are many "natural extensions" of ZFC studied by set theorists. Thus, there exist many alternative systems of axioms each of which has implications to the standard theory of cardinality discussed above.

### Without the axiom of choice

The Axiom of Choice (AC) is a fundamental principle in the foundations of mathematics which has seen much controversy in its history. Informally, it states that given any collection of non-empty sets, it is possible to construct a new set by choosing one element from each set, even if the collection is infinite. The controversy over AC had mostly to do with the nature of how it chose these elements. Whereas most axioms seem to describe what sets are or allow one to construct explicit sets, AC does not tell you *which* set you have constructed, just that such a choice function exists, leaving the constructed set implicit. It is hardly controversial to modern mathematicians, however, because of its unique historical controversy it is often given special treatment not given to other axioms that basic proofs which use it ought to call it out.

The Axiom of Choice is very closely related to the nature of cardinality. AC imposes a strict structure on what cardinality can be and what it means to compare to sets. Assuming AC is false, cardinality has a much more intricate structure and resists the kind of linear order AC offers. For example, cardinal inequality does not treat injection and surjection equivalently. Specifically, there exists sets, such that there is a surjection from  onto , but no injection from  into , thus injection is a strictly stronger notion. Similarly, although cardinal inequality maintains a partial order, there exist incomparable sets under cardinal inequality, that is, the law of trichotomy fails to hold, and there are sets such that none of , ,  hold. Both of these (trichotomy and that surjection implies injection) are equivalent to AC.

Because trichotomy does not hold and the aleph sequence has a natural well-order, there exist sets whose cardinality does not correspond to any aleph. As such, the cardinality function  becomes somewhat more difficult to define. In fact, a cardinality function that maps each set to a unique "representative" of that cardinalityi.e. satisfying Hume's principle, and idempotence ( or )is impossible without AC. Set theories without AC address this by adopting the so-called *FregeRussellScott* definition, introduced by Dana Scott, reminiscent of the *FregeRussell* cardinal numbers. Under this definition, one considers the "Set of all sets" equinumerous with , but applies Scott's trick to regularize these classes. Specifically, it reduces this class to only those sets of minimal rank; that is, those sets which appear earliest in the von Neumann hierarchy. Since the von Neumann hierarchy is well-ordered, there is a least rank, and since these sets all share that rank, their collection is bounded in the hierarchy and thus constitutes a set rather than a proper class.

Cardinal arithmetic is much more involved, losing many of its simplest identities. The identity {{tmath|1=\vert A\vert + \vert B\vert = \operatorname{max}(\vert A\vert, \vert B\vert)}} requires a well-defined notion of "max" which requires all sets to be comparable, and therefore does not hold. Similarly, the statement that  for every infinite cardinal , and even the statement that the arbitrary product of non-zero cardinal numbers always non-zerocalled the *Multiplicative axiom* by Bertrand Russellare both equivalent to AC. Further, although the Generalized Continuum Hypothesis (GCH, that  for every ) is independent of ZFC, it can be shown ZF+GCH is sufficient to derive AC, and therefore, if AC is false, it follows that GCH does not hold.

### Proper classes
Some set theories include classes, which are collections of sets, which allows theories to discuss any collection of sets without running into issues self-reference (e.g. the set of all sets that don't contain themselves.). A class is called a proper class when, at least intuitively, it is "too large" to form a set. For example, the Universe of all sets, the class of all cardinal numbers, and the class of all ordinal numbers are proper classes. Such set theories include Von Neumann–Bernays–Gödel set theory (NBG), and Morse–Kelley set theory (MK). Cantor originally called the size of these "too large" sets "Absolute infinite", separating it from the transfinite. The former he characterized by its "inconsistency", causing paradoxes, and associated the concept with God.

Proper classes can, in a sense, be assigned cardinalities. The first to formally distinguish between sets and classes was John von Neumann, who formalized the notion of "too large to form sets". More accurately, he defined a class to be "too big" (a proper class) if and only if it was equinumerous with the whole Universe of sets (cf. *Axiom of limitation of size*; used as an axiom in MK set theory, and derivable as a theorem in NBG). Thus, all proper classes have the same "size". The axiom has several implications, mostly relating to the limitation of size principles of early set theory. It implies the axiom of specification, the axiom of replacement, the axiom of union, and the axiom of global choice.

### Large cardinals

Large cardinal axioms assert the existence of cardinal numbers which, as the name suggests, are very largeso large that they cannot be proven to exist within ZFC. For example, an inaccessible cardinal is, roughly, a cardinal that cannot be approached from below using basic set-theoretic operations such as unions, limits, and powersets (more formally, a regular, limit cardinal greater than ). Large cardinals are understood in terms of the Von Neumann hierarchy, denoted  (for some ordinal ), which can be understood as the sets that can be obtained from the empty set, followed by recursively applying the powerset  times. Specifically, , {{tmath|1=V_{\alpha+1} = \mathcal P(V_\alpha) }} and {{tmath|1=V_\lambda = \bigcup_{\alpha < \lambda} V_\alpha }} for a limit ordinal .

There exist many known properties defining large cardinals, which come seemingly in a linear hierarchy, in terms of consistency strength. As an analogy, in ZFC without the Axiom of Infinity can only prove the existence of finite sets. Therefore {{tmath|1=V_\omega ( = \mathcal{P}(\varnothing) \cup \mathcal{P}(\mathcal{P}(\varnothing)) \cup \cdots ) }}, whose existence is provable in usual ZFC, can serve as a model of ZFCInfinity, and thus if ZFC is consistent, ZFCInfinity is consistent. Analogously, ZFC+"There exists an inaccessible cardinal" implies the consistency of ZFC, since if  is inaccessible,  can serve as a model of ZFC (cf. *Grothendieck universe*). Stronger and stronger large cardinal axioms assert the existence of larger and larger cardinals, each of which proves the consistency—or inconsistency—of the weaker systems.

Large cardinals are at the forefront of set-theoretic research for both practical and philosophical reasons. In a practical sense, it is often the case that unproven or unprovable conjectures can be settled by sufficiently strong large cardinal axioms. For example, the existence of a measurable cardinal is inconsistent with Gödel's constructibility axiom. Similarly, the existence of a Reinhardt cardinal is inconsistent with the axiom of choice, which makes it a much more controversial large cardinal axiom. In a philosophical sense, according to a platonic view among set-theorists such as W. Hugh Woodin, these axioms simply extend the system to include sets that are "supposed" to be considered. That is, there is some fundamental universe of sets, which these axioms grant further access to. For this reason, large cardinal axioms are generally given preference compared to other possible axioms of set theory. This view is controversial among a competing philosophy, sometimes called pluralism, which posits that set theory should be understood as a multiverse of set theories, but no "absolute" or "true" model.

### Determinacy

The Axiom of Determinacy (AD) asserts that certain kinds of mathematical games on the natural numbers are determined; that is, one player will always have a guaranteed winning strategy. The first serious study of the consequences of AD started during the 1960s in descriptive set theorywhich, roughly, studies the definable sets of the real numbersafter it was noticed that it leads to very nice regulatory properties of the real numbers. Specifically, it implies the perfect set property, the property of Baire, and that all sets of real numbers are Lebesgue measurable. However, it was shown that this axiom is inconsistent with AC (cf. **), and thus was never taken as a fundamental axiom of set theory. However, its relationship with large cardinals and the continuum hypothesis is still of heavy interest among set theorists such as Donald A. Martin, John R. Steel, and W. Hugh Woodin.

Some implications of replacing the Axiom of Choice with the Axiom of Determinacy include:

* The Banach–Tarski paradox is a paradox arising from AC which allows one to disassemble a ball into finitely many pieces, and re-arrange them into two balls identical to the original. This, and related paradoxes, become impossible under AD.
* It is possible to partition the set of real numbers in such a way that there are *strictly* more nonempty sets of real numbers than there are real numbers. That is, if  denotes a general partition of , there is a partition such that .

## Skolem's paradox


In model theory, a model corresponds to a specific interpretation of a formal language or theory. It consists of a domain (a set of objects) and an interpretation of the symbols in the language, such that the axioms of the theory are satisfied within this structure. In first-order logic, the Löwenheim–Skolem theorem states that if a theory has an infinite model, then it also has models of every other infinite cardinality. Applied to set theory, it asserts that Zermelo–Fraenkel set theory, which proves the existence of uncountable sets such as , nevertheless has a countable model. Thus, Skolem's paradox was posed as follows: how can it be that there exists a domain of set theory which only contains countably many objects, but is capable of satisfying the statement "there exists a set with uncountably many elements"?

Skolem's paradox does not only apply to ZFC, but *any* first-order set theory, if it is consistent, has a model which is countable. A mathematical explanation of the paradox, showing that it is not a true contradiction in mathematics, was first given in 1922 by Thoralf Skolem, who asserted it as a reason against founding mathematics on first order logic. He explained that the countability or uncountability of a set is not absolute, but relative to the model in which the cardinality is measured. This is because, for example, if the set  is countable in a model of set theory then there is a bijection . But a submodel containing  which excludes all such functions would thus contain no bijection between  and , and therefore  would be uncountable. In second-order and higher-order logics, the Löwenheim–Skolem theorem does not hold. This is due to the fact that second-order logic quantifies over all subsets of the domain. Skolem's work was harshly received by Ernst Zermelo, who argued against the limitations of first-order logic and Skolem's notion of "relativity", but the result quickly came to be accepted by the mathematical community.

## History
### Ancient history

From the 6th century BCE, the writings of Greek philosophers, such as Anaximander, discuss infinite sets or objects, however, it was generally viewed as paradoxical and imperfect (cf. *Zeno's paradoxes*). Aristotle distinguished between the notions of actual and potential infinity, arguing that Greek mathematicians understood the difference, and that they "do not need the [actual] infinite and do not use it." The Greek notion of number (*αριθμός*, *arithmos*) was used exclusively for a definite number of definite objects (i.e. finite numbers). This would be codified in Euclid's *Elements*, where the fifth common notion states "The whole is greater than the part", called the *Euclidean principle*. This principle would be the dominating philosophy in mathematics until the 19th century.

Around the 4th century BCE, Jaina mathematics would be the first to discuss different sizes of infinity. They defined three major classes of number: enumerable (finite numbers), unenumerable (*asamkhyata*, roughly, countably infinite), and infinite (*ananta*, roughly, the continuum). Then they had five classes of infinite numbers: infinite in one direction, infinite in both directions, infinite in area, infinite everywhere, and infinite perpetually.

One of the earliest explicit uses of a one-to-one correspondence is recorded in Aristotle's *Mechanics* (), known as Aristotle's wheel paradox. The paradox can be briefly described as follows: A wheel is depicted as two concentric circles. The larger, outer circle is tangent to a horizontal line (e.g. a road that it rolls on), while the smaller, inner circle is rigidly affixed to the larger.  Assuming the larger circle rolls along the line without slipping (or skidding) for one full revolution, the distances moved by both circles are the same: the circumference of the larger circle. Further, the lines traced by the bottom-most point of each is the same length. Since the smaller wheel does not skip any points, and no point on the smaller wheel is used more than once, there is a one-to-one correspondence between the two circles.

### Pre-Cantorian set theory
 (left).  Portrait of Bernard Bolzano 1781–1848 (right).
}}

Galileo Galilei presented what was later coined Galileo's paradox in his book *Two New Sciences* (1638), where he presents a seeming paradox in infinite sequences of numbers. It goes as follows: for each perfect square  1, 4, 9, 16, and so on, there is a unique square root {{tmath|1= \textstyle(\sqrt{n^2} = n) }} 1, 2, 3, 4, and so on. Therefore, there are as many perfect squares as there are square roots. However, every number is a square root, since it can be squared, but not every number is a perfect square. Moreover, the proportion of perfect squares as one passes larger values diminishes, and is eventually smaller than any given fraction. Galileo denied that this was fundamentally contradictory, however he concluded that this meant we could not compare the sizes of infinite sets, missing the opportunity to discover cardinality.

In *A Treatise of Human Nature* (1739), David Hume is quoted for saying *"When two numbers are so combined, as that the one has always a unit answering to every unit of the other, we pronounce them equal",* now called *Hume's principle*, which was used extensively by Gottlob Frege later during the rise of set theory.

Bernard Bolzano's *Paradoxes of the Infinite* (*Paradoxien des Unendlichen*, 1851) is often considered the first systematic attempt to introduce the concept of sets into mathematical analysis. In this work, Bolzano defended the notion of actual infinity, presented an early formulation of what would later be recognized as one-to-one correspondence between infinite sets. He discussed examples such as the pairing between the intervals  and  by the relation , and revisited Galileo's paradox. However, he too resisted saying that these sets were, in that sense, the same size. While *Paradoxes of the Infinite* anticipated several ideas central to later set theory, the work had little influence on contemporary mathematics, in part due to its posthumous publication and limited circulation.

### Early set theory
#### Georg Cantor

The concept of cardinality emerged nearly fully formed in the work of Georg Cantor during the 1870s and 1880s, in the context of mathematical analysis. In a series of papers beginning with *On a Property of the Collection of All Real Algebraic Numbers* (1874), Cantor introduced the idea of comparing the sizes of infinite sets, through the notion of one-to-one correspondence. He showed that the set of real numbers was, in this sense, strictly larger than the set of natural numbers using a nested intervals argument. This result was later refined into the more widely known diagonal argument of 1891, published in *Über eine elementare Frage der Mannigfaltigkeitslehre,* where he also proved the more general result (now called Cantor's Theorem) that the power set of any set is strictly larger than the set itself.

Cantor introduced the notion cardinal numbers together with ordinal numbers, which he viewed as abstractions of sets. For a given set , he wrote {{tmath|1= \textstyle\overline{M} }} to mean the abstraction of that set from its elements, while maintaining their order, and cardinal numbers were a double abstraction, written . Specifically, his definition was "the general concept which, with the aid of our intelligence, results from a set when we abstract from the nature of its various elements and from the order of their being given." This definition was considered to be imprecise, unclear, and purely psychological, and it would be some time before the concept was put on more rigorous footing. He also introduced the Aleph sequence for infinite cardinal numbers. These notations appeared in correspondence and were formalized in his later writings, particularly the series *Beiträge zur Begründung der transfiniten Mengenlehre* (18951897). In these works, Cantor developed an arithmetic of cardinal numbers, defining addition, multiplication, and exponentiation of cardinal numbers. This led to the formulation of the Continuum Hypothesis: whether {{tmath|1=2^{\aleph_0} = \aleph_1 }}. Cantor was unable to resolve CH and left it as an open problem.

#### Other contributors
Parallel to Cantor’s development, Richard Dedekind independently formulated many advanced theorems of set theory, and helped establish set-theoretic foundations of algebra and arithmetic. Dedekind's ** (1888) emphasized structural properties over extensional definitions, and supported the bijective formulation of size and number. Dedekind was in correspondence with Cantor during the development of set theory; he supplied Cantor with a proof of the countability of the algebraic numbers, and gave feedback and modifications on Cantor's proofs before publishing.

After Cantor's 1883 proof that all finite-dimensional spaces  have the same cardinality, in 1890, Giuseppe Peano introduced the Peano curve, which was a more visual proof that the unit interval  has the same cardinality as the unit square on . This created a new area of mathematical analysis studying what is now called space-filling curves. In 1894, attempting  to formalize Cantor's "cardinal number", Peano introduced "definition by abstraction": if one can define an equivalence relation, then one may define "that property" which an equivalence class describes. However, this was received harshly by Russell for not considering that this property may not be unique.

German logician Gottlob Frege attempted to ground the concepts of number and arithmetic in logic using Cantor's theory of cardinality and Hume's principle in *Die Grundlagen der Arithmetik* (1884) and the subsequent *Grundgesetze der Arithmetik* (1893, 1903). Frege defined cardinal numbers as equivalence classes of sets under equinumerosity. However, Frege's approach to set theory was later shown to be flawed. His approach was eventually reformalized by Bertrand Russell and Alfred Whitehead in *Principia Mathematica* (19101913, vol. II) using a theory of types. Though Russell initially had difficulties accepting Cantor's and Frege’s intuitions of cardinality. This definition of cardinal numbers is now referred to as the *FregeRussell* definition. This definition was eventually superseded by the convention established by John von Neumann in 1928 which uses representatives to define cardinal numbers.

At the Paris conference of the International Congress of Mathematicians in 1900, David Hilbert, one of the most influential mathematicians of the time, gave a speech wherein he presented ten unsolved problems (of a total of 23, later published, now called *Hilbert's problems*). Of these, he placed "Cantor's problem" (now called the Continuum Hypothesis) as the first on the list. This list of problems would prove to be very influential in 20th century mathematics, and attracted a lot of attention from other mathematicians toward Cantor's theory of cardinality.

### Axiomatic set theory


In 1908, Ernst Zermelo proposed the first axiomatization of set theory, now called Zermelo set theory, primarily to support his earlier (1904) proof of the Well-ordering theorem, which showed that all cardinal numbers could be represented as Alephs, though the proof required a controversial principle now known as the Axiom of Choice (AC). Zermelo's system would later be extended by Abraham Fraenkel and Thoralf Skolem in the 1920s to create the standard foundation of set theory, called Zermelo–Fraenkel set theory (ZFC, "C" for the Axiom of Choice). ZFC provided a rigorous foundation through which infinite cardinals could be systematically studied while avoiding the paradoxes of naive set theory.

Ignoring possible foundational issues, during the early 1900s, Felix Hausdorff would begin studying "exorbitant numbers": roughly, very large cardinal numbers, or what are now called inaccessible cardinals. This work would be continued and popularized by several other influential set theorists such as Paul Mahlowho introduced Mahlo cardinalsas well as Wacław Sierpiński, and Alfred Tarski. Their work would eventually be known collectively as the study of large cardinals.

In 1940, Kurt Gödel showed that the Continuum Hypothesis (CH) cannot be *disproved* from the axioms of ZFC by showing that both CH and AC hold in his constructible universe: an inner model of ZFC. The existence of a model of ZFC in which additional axioms hold shows that the additional axioms are (relatively) consistent with ZFC. In 1963, Paul Cohen showed that CH cannot be *proven* from the ZFC axioms, which showed that CH was independent from ZFC. To prove his result, Cohen developed the method of forcing, which has become a standard tool in set theory. Cohen was awarded the Fields Medal in 1966 for his proof.

## See also


* *Cardinal and Ordinal Numbers*
*
* Infinitary combinatorics
* Multiplicity (mathematics)
* Natural density
* Numerosity (mathematics)
* Subcountability

## References
### Notes


### Citations


### Primary sources

*
*
*
*
*
*
*
*
*
*
*
*
*
*
*
*
*
*  Alt URL
*
*  
*
*
*
*
*


### Secondary and tertiary sources

*
*  Alt URL
*
*
*
*
*
*
*
*
*
*
*
*
*
*
*
*  Alt URL
*  Archived on 2016-01-06
*
*
*  Alt URL
*
*
*
*
*
*  Alt URL
*
*
*  Alt URL
*
*
*
*
*
*
* Jaina mathematics
 |url=https://mathshistory.st-andrews.ac.uk/HistTopics/Jaina_mathematics/
 |access-date=2025-06-09
 |website=MacTutor History of Mathematics Archive
 |via=University of St Andrews, Scotland
}}
*
*  Alt URL
*
*
*
*
*  Alt URL
*  Alt URL
*
*
*
*

