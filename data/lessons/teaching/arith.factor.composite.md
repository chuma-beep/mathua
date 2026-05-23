> Content sourced from [ORCCA](https://pcc.edu/orcca) (Open Resources for Community College Algebra) — CC BY 4.0
# Set Notation and Types of Numbers

When we talk about *how many* or *how much* of something we have, it often makes
      sense to use different types of numbers. For example, if we are counting dogs in a shelter,
      the possibilities are only $0,1,2,\ldots$ . (It would be difficult to have $\frac{1}{2}$ of a dog.) On the other hand if you were weighing a dog in pounds, it
      doesn't make sense to only allow yourself to work with whole numbers. The dog might weigh
      something like $28.35$ pounds. These examples highlight how certain kinds of numbers are
      appropriate for certain situations. We'll classify various types of numbers in this section.

## Set Notation

What is the mathematical difference between these three "lists" ? \[
        28, 31, 30\qquad\{28, 31, 30\}\qquad(28, 31, 30)
      \] To a mathematician, the last one, $(28, 31, 30)$ is an *ordered* triple. What
      matters is not merely the three numbers, but *also* the order in which they come. The
      ordered triple $(28, 31, 30)$ is not the same as $(30, 31, 28)$ ; they have the same
      numbers in them, but the order has changed. For some context, February has $28$ days; *then* March has $31$ days; *then* April has $30$ days. The order of
      the three numbers is meaningful in that context.

With curly braces and $\{28, 31, 30\}$ , a mathematician sees a collection of numbers and
      does not particularly care in which order they are written. Such a collection is called a *set* . All that matters is that these numbers are part of a collection. They've been *written* in some particular order because that's necessary to write them down.
      But you might as well have put the three numbers in a bag and shaken up the bag. For some
      context, maybe your favorite three NBA players have jersey numbers $30$ , $31$ , and $28$ , and you like them all equally well. It doesn't really matter what order you use to
      list them.

So we can say: \{28, 31, 30\}\amp=\{30, 31, 28\}\amp(28, 31, 30)\amp\neq(30, 31, 28)

What about just writing $28, 31, 30$ ? This list of three numbers is ambiguous. Without
      the curly braces or parentheses, it's unclear to a reader if the order is important. *Set notation* is the use of curly braces to surround a list/collection of numbers, and we will use set notation frequently in this section.

**Set Notation**
Practice using (and not using) set notation.
- Without knowing which error code is most common, express this set mathematically.
  *Solution*: Since we only have to describe a collection of three numbers and their order doesn't
              matter, we can write .
- Error code `500` is the most common. Error code `403` is the least common of
              these three. And that leaves `404` in the middle. Express the error codes in a
              mathematical way that appreciates how frequently they happen, from most often to least
              often.
  *Solution*: Now we must describe the same three numbers and we want readers to know that the order
              we are writing the numbers matters. We can write .

## Different Number Sets

In the introduction, we mentioned how different sets of numbers are appropriate for different
      situations. Here are the basic sets of numbers that are used in basic algebra.

Natural Numbers When we count, we begin: $1, 2, 3, \dots$ and continue on in that pattern. These
            numbers are known as *natural numbers* . $\mathbb{N}=\{1,2,3,\dots \}$ Whole Numbers If we include zero, then we have the set of *whole numbers* . $\{0,1,2,3,\dots \}$ has no standard symbol, but some options are $\mathbb{N}_0$ , $\mathbb{N}\cup\{0\}$ , and $\mathbb{Z}_{\geq0}$ . Integers If we include the negatives of whole numbers, then we have the set of *integers* . $\mathbb{Z}=\{\dots,-3,-2,-1,0,1,2,3,\dots \}$ . A $\mathbb{Z}$ is used because one word in German for "numbers" is "Zahlen" . Rational Numbers A *rational number* is any number that *can* be written as a fraction
            of integers, where the denominator is nonzero. Alternatively, a *rational number* is any number that *can* be written with a decimal that terminates or that repeats. $\mathbb{Q}=\left\{0,1,-1,2,\frac{1}{2},-\frac{1}{2},-2,3,\frac{1}{3},-\frac{1}{3},-3,\frac{3}{2},\frac{2}{3}\ldots\right\}$ $\mathbb{Q}=\left\{0,1,-1,2,0.5,-0.5,-2,3,0.\overline{3},-0.\overline{3},-3,1.5,0.\overline{6}\ldots\right\}$ A $\mathbb{Q}$ is used because fractions are *q* uotients of integers. Irrational Numbers Any number that *cannot* be written as a fraction of integers belongs to the set
            of *irrational numbers* . Another way to say this is that any number whose
            decimal places goes on forever without repeating is an *irrational number* .
            Some examples include $\pi\approx3.1415926\ldots$ , $\sqrt{15}\approx3.87298\ldots$ , $e\approx2.71828\ldots$ There is no standard symbol for the set of irrational numbers. Real Numbers Any number that can be marked somewhere on a number line is a *real number* .
            Real numbers might be the only numbers you are familiar with. For a number to *not* be real, you have to start considering things called *complex numbers* , which are not our concern right now. The set of real numbers can be denoted with $\mathbb{R}$ for short.

> **Rational Numbers in Other Forms**
> Any number that *can* be written as a ratio of integers is rational, even if it's not
        written that way at first. For example, these numbers might not look rational to you at
        first glance: $-4$ , $\sqrt{9}$ , $0\pi$ , and $\sqrt[3]{\sqrt{5}+2}-\sqrt[3]{\sqrt{5}-2}$ . But they are all rational, because they
        can respectively be written as $\frac{-4}{1}$ , $\frac{3}{1}$ , $\frac{0}{1}$ ,
        and $\frac{1}{1}$ .

**Determine If Numbers Are This Type or That Type**
Determine which numbers from the set $\left\{-102, -7.25, 0, \frac{\pi}{4}, 2, \frac{10}{3}, \sqrt{19}, \sqrt{25}, 10.\overline{7} \right\}$ are natural numbers, whole numbers, integers, rational numbers, irrational numbers, and
          real numbers.

*Solution*
All of these numbers are real numbers, because all of these numbers can be positioned on
          the real number line.
Each real number is either rational or irrational, and not both. $-102$ , $-7.25$ , $0$ , and $2$ are rational because we can see directly that their
          decimal expressions terminate. $10.\overline{7}$ is also rational, because its
          decimal expression repeats. $\frac{10}{3}$ is rational because it is a ratio of
          integers. And last but not least, $\sqrt{25}$ is rational, because that's the same
          thing as $5$ .
This leaves only $\frac{\pi}{4}$ and $\sqrt{19}$ as irrational numbers. Their
          decimal expressions go on forever without entering a repetitive cycle.
Only $-102$ , $0$ , $2$ , and $\sqrt{25}$ (which is really $5$ ) are
          integers.
Of these, only $0$ , $2$ , and $\sqrt{25}$ are whole numbers, because whole
          numbers exclude the negative integers.
Of these, only $2$ and $\sqrt{25}$ are natural numbers, because the natural
          numbers exclude $0$ .


In the introduction, we mentioned that the different types of numbers are appropriate in
          different situation. Which number set do you think is most appropriate in each of the
          following situations?
- The number of people in a math class that play the ukulele. This number is best considered as a .
  *Solution*: The number of people who play the ukulele could be $0,1,2,\dots$ , so the whole
              numbers are the appropriate set.
- The hypotenuse's length in a given right triangle. This number is best considered as a .
  *Solution*: A hypotenuse's length could be $1$ , $1.2$ , $\sqrt{2}$ (which is
              irrational), or any other positive number. So the real numbers are the appropriate
              set.
- The proportion of people in a math class that have a cat. This number is best considered as a .
  *Solution*: This proportion will be a ratio of integers, as both the total number of people in the
              class and the number of people who have a cat are integers. So the rational numbers
              are the appropriate set.
- The number of people in the room with you who have the same birthday as you. This number is best considered as a .
  *Solution*: We know that the number of people must be a counting number, and since *you* are in the room with yourself, there is at least one person in that room with your
              birthday. So the natural numbers are the appropriate set.
- The total revenue (in dollars) generated for ticket sales at a Timbers soccer game. This number is best considered as a .
  *Solution*: The total revenue will be some number of dollars and cents, such as $\$631{,}897.15$ , which is a terminating decimal and thus a rational number. So
              the rational numbers are the appropriate set.

## Converting Repeating Decimals to Fractions

We have learned that a terminating decimal number is a rational number. It's easy to convert a
      terminating decimal number into a fraction of integers: you just need to multiply and divide
      by one of the numbers in the set $\{10,100,1000,\ldots\}$ . For example, when we say the
      number $0.123$ out loud, we say "one hundred and twenty-three thousandths" . While
      that's a lot to say, it makes it obvious that this number can be written as a ratio: \[
        0.123=\frac{123}{1000}
      \] .
      Similarly, \[
        21.28=\frac{2128}{100}=\frac{532\cdot4}{25\cdot4}=\frac{532}{25}
      \] ,
      demonstrating how *any* terminating decimal can be written as a fraction.

Repeating decimals can also be written as a fraction. To understand how, use a calculator to
      find the decimal for, say, $\frac{73}{99}$ and $\frac{189}{999}$ You will find that \[
        \frac{73}{99}=0.73737373\ldots=0.\overline{73}\qquad\frac{189}{999}=0.189189189\ldots=0.\overline{189}
      \] .
      The pattern is that dividing a number by a number from $\{9,99,999,\ldots\}$ with the
      same number of digits will create a repeating decimal that starts as "0." and then
      repeats the numerator. We can use this observation to reverse engineer some fractions from
      repeating decimals.

- Write the rational number $0.772772772\ldots$ as a fraction.
  *Solution*: The *three* -digit number $772$ repeats after the decimal. So we will make
              use of the *three* -digit denominator $999$ . And we have $\frac{772}{999}$ .
- Write the rational number $0.69696969\ldots$ as a fraction.
  *Solution*: The *two* -digit number $69$ repeats after the decimal. So we will make use
              of the *two* -digit denominator $99$ . And we have $\frac{69}{99}$ .
              But this fraction can be reduced to $\frac{23}{33}$ .

Converting a repeating decimal to a fraction is not always quite this straightforward. There
      are complications if the number takes a few digits before it begins repeating. For your
      interest, here is one example on how to do that.

**Example**
