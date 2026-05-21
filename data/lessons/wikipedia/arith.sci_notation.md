> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Scientific_notation) — CC BY-SA 4.0

# Scientific notation

E series|the food additive codes|E number|the letter|E

**Scientific notation** is a way of expressing numbers that are too large or too small to be conveniently written in decimal form; since to do so would require writing out an inconveniently long string of digits. It may be referred to as **scientific form** or **standard index form**, or **standard form** in the United Kingdom. This base ten notation is commonly used by scientists, mathematicians, and engineers, in part because it can simplify certain arithmetic operations. On scientific calculators, it is usually known as "SCI" display mode.

{| class="wikitable" style="float:right; margin:5px;"
!Decimal notation
!Scientific notation
|-
|
|
|-
|
|
|-
|
|
|-
|
|
|-
|
|
|-
|
|
|-
|
|
|-
|
|
|}

In scientific notation, nonzero numbers are written in the form

or *m* times ten raised to the power of *n*, where *n* is an integer, and the coefficient *m* is a nonzero real number (usually between 1 and 10 in absolute value, and nearly always written as a terminating decimal). The integer *n* is called the exponent and the real number *m* is called the *significand* or *mantissa*. The term "mantissa" can be ambiguous where logarithms are involved, because it is also the traditional name of the fractional part of the common logarithm. If the number is negative then a minus sign precedes *m*, as in ordinary decimal notation. In normalized notation, the exponent is chosen so that the absolute value (modulus) of the significand *m* is at least 1 but less than 10.

Decimal floating point is a computer arithmetic system closely related to scientific notation.

## History

For performing calculations with a slide rule, standard form expression is required. Thus, the use of scientific notation increased as engineers and educators used that tool. See Slide rule#History.

## Styles
### Normalized notation

Any real number can be written in the form in many ways: for example, 350 can be written as or or. In *normalized* scientific notation (called "standard form" in the United Kingdom), the exponent *n* is chosen so that the absolute value of *m* remains at least one and less than ten (1 ≤ |*m*| < 10). Thus 350 is written as. This form allows easy comparison of numbers: numbers with bigger exponents are (due to the normalization) larger than those with smaller exponents, and subtraction of exponents gives an estimate of the number of orders of magnitude separating the numbers. It is also the form that is required when using tables of common logarithms. In normalized notation, the exponent *n* is negative for a number with absolute value between 0 and 1 (e.g. 0.5 is written as ). The 10 and exponent are often omitted when the exponent is 0. For a series of numbers that are to be added or subtracted (or otherwise compared), it can be convenient to use the same value of *n* for all elements of the series.

Normalized scientific form is the typical form of expression of large numbers in many fields, unless an unnormalized or differently normalized form, such as engineering notation, is desired. Normalized scientific notation is often called **exponential notation** – although the latter term is more general and also applies when *m* is not restricted to the range 1 to 10 (as in engineering notation for instance) and to bases other than 10 (for example, ).

### Engineering notation

Engineering notation (often named "ENG" on scientific calculators) differs from normalized scientific notation in that the exponent *n* is restricted to multiples of 3. Consequently, the absolute value of *m* is in the range 1 ≤ |*m*| < 1000, rather than 1 ≤ |*m*| < 10. Though similar in concept, engineering notation is rarely called *scientific notation*. Engineering notation allows the numbers to explicitly match their corresponding SI prefixes, which facilitates reading and oral communication. For example, can be read as "twelve-point-five nanometres" and written as , while its scientific notation equivalent would likely be read out as "one-point-two-five times ten-to-the-negative-eight metres".

### E notation
{| class="wikitable" style="float:right; margin: 0.5em 0 1.3em 1.4em"
!Explicit notation
!E notation
|-
|
|
|-
|
|
|-
|
|
|-
|
|
|-
|
|
|-
|
|
|-
|
|
|-
|
|
|}

Calculators and computer programs typically present very large or small numbers using scientific notation, and some can be configured to uniformly present all numbers that way. Because superscript exponents, like the ⟨\(^{7}\)⟩ in 10\(^{7}\), can be inconvenient to display or type, the letter "E" or "e" (for "exponent") is often used to represent "times ten raised to the power of", so that the notation *m* E *n* for a decimal significand *m* and integer exponent *n* means the same as *m* × 10\(^{*n*}\). For example is written as or , and is written as or. Sometimes the positive power is explicitly shown ( or ). While common in computer output, this abbreviated version of scientific notation is discouraged for published documents by some style guides.

Most popular programming languages – including Fortran, C/C++, Python, and JavaScript – use this "E" notation, which comes from Fortran and was present in the first version released for the IBM 704 in 1956. The E notation was already used by the developers of SHARE Operating System (SOS) for the IBM 709 in 1958. Later versions of Fortran (at least since FORTRAN IV as of 1961) also use "D" to signify double precision numbers in scientific notation, and newer Fortran compilers use "Q" to signify quadruple precision. The MATLAB programming language supports the use of either "E" or "D".

The ALGOL 60 (1960) programming language uses a subscript ten "\(_{10}\)" character instead of the letter "E", for example: 6.022\(_{10}\)23. This presented a challenge for computer systems which did not provide such a character, so ALGOL W (1966) replaced the symbol by a single quote, e.g. 6.022'+23, and some Soviet ALGOL variants allowed the use of the Cyrillic letter "ю", e.g.. Subsequently, the ALGOL 68 programming language provided a choice of characters: , or \(_{10}\). The ALGOL "\(_{10}\)" character was included in the Soviet GOST 10859 text encoding (1964), and was added to Unicode 5.2 (2009) as. Some programming languages use other symbols. For instance, Simula uses (or for long), as in. Mathematica supports the shorthand notation (reserving the letter for the mathematical constant *e*).

The first pocket calculators supporting scientific notation appeared in 1972. To enter numbers in scientific notation calculators include a button labeled "EXP" or "×10\(^{*x*}\)", among other variants. The displays of pocket calculators of the 1970s did not display an explicit symbol between significand and exponent; instead, one or more digits were left blank (e.g. 6.022 23, as seen in the HP-25), or a pair of smaller and slightly raised digits were reserved for the exponent (e.g. 6.022 \(^{23}\), as seen in the Commodore PR100). In 1976, Hewlett-Packard calculator user Jim Davidson coined the term *decapower* for the scientific-notation exponent to distinguish it from "normal" exponents, and suggested the letter "D" as a separator between significand and exponent in typewritten numbers (for example, ); these gained some currency in the programmable calculator user community. The letters "E" or "D" were used as a scientific-notation separator by Sharp pocket computers released between 1987 and 1995, "E" used for 10-digit numbers and "D" used for 20-digit double-precision numbers. The Texas Instruments TI-83 and TI-84 series of calculators (1996–present) use a small capital E for the separator.

In 1962, Ronald O. Whitaker of Rowco Engineering Co. proposed a power-of-ten system nomenclature where the exponent would be circled, e.g. 6.022 × 10\(^{3}\) would be written as "6.022③".

## Significant figures

A significant figure is a digit in a number that adds to its precision. This includes all nonzero numbers, zeroes between significant digits, and zeroes indicated to be significant. Leading and trailing zeroes are not significant digits, because they exist only to show the scale of the number. Unfortunately, this leads to ambiguity. The number is usually read to have five significant figures: 1, 2, 3, 0, and 4, the final two zeroes serving only as placeholders and adding no precision. The same number, however, would be used if the last two digits were also measured precisely and found to equal 0 – seven significant figures.

When a number is converted into normalized scientific notation, it is scaled down to a number between 1 and 10. All of the significant digits remain, but the placeholding zeroes are no longer required. Thus would become if it had five significant digits. If the number were known to six or seven significant figures, it would be shown as or. Thus, an additional advantage of scientific notation is that the number of significant figures is unambiguous.

### Estimated final digits
It is customary in scientific measurement to record all the definitely known digits from the measurement and to estimate at least one additional digit if there is any information at all available on its value. The resulting number contains more information than it would without the extra digit, which may be considered a significant digit because it conveys some information leading to greater precision in measurements and in aggregations of measurements (adding them or multiplying them together).

More detailed information about the precision of a value written in scientific notation can be conveyed through additional notation. For instance, the accepted value of the mass of the proton can be expressed as , which is shorthand for. However, it is unclear whether an error expressed in this way ( in this case) is the maximum possible error, standard error, or some other confidence interval.

## Use of spaces
In normalized scientific notation, in E notation, and in engineering notation, the space (which in typesetting may be represented by a normal width space or a thin space) that is allowed *only* before and after "×" or in front of "E" is sometimes omitted, though it is less common to do so before the alphabetical character.

## Further examples of scientific notation
* An electron's mass is about. In scientific notation, this is written. * The Earth's mass is about. In scientific notation, this is written. * The Earth's circumference is approximately. In scientific notation, this In engineering notation, this is written. In SI writing style, this may be written (**).
* An inch is defined as *exactly*. Using scientific notation, this value can be uniformly expressed to any desired precision, from the nearest tenth of a millimeter to the nearest nanometer , or beyond.
* Hyperinflation means that too much money is put into circulation, perhaps by printing banknotes, chasing too few goods. It is sometimes defined as inflation of 50% or more in a single month. In such conditions, money rapidly loses its value. Some countries have had events of inflation of 1 million percent or more in a single month, which usually results in the rapid abandonment of the currency. For example, in November 2008 the monthly inflation rate of the Zimbabwean dollar reached 79.6 billion percent (470% per day); the approximate value with three significant figures would be %, or more simply a rate of. ## Converting numbers
Converting a number in these cases means to either convert the number into scientific notation form, convert it back into decimal form or to change the exponent part of the equation. None of these alter the actual number, only how it's expressed.

### Decimal to scientific
First, move the decimal separator point sufficient places, *n*, to put the number's value within a desired range, between 1 and 10 for normalized notation. If the decimal was moved to the left, append × 10*\(^{n}\)*; to the right, × 10*\(^{−n}\)*. To represent the number in normalized scientific notation, the decimal separator would be moved 6 digits to the left and × 10\(^{6}\) appended, resulting in. The number would have its decimal separator shifted 3 digits to the right instead of the left and yield as a result.

### Scientific to decimal
Converting a number from scientific notation to decimal notation, first remove the × 10*\(^{n}\)* on the end, then shift the decimal separator *n* digits to the right (positive *n*) or left (negative *n*). The number would have its decimal separator shifted 6 digits to the right and become , while would have its decimal separator moved 3 digits to the left and be. ### Exponential
Conversion between different scientific notation representations of the same number with different exponential values is achieved by performing opposite operations of multiplication or division by a power of ten on the significand and an subtraction or addition of one on the exponent part. The decimal separator in the significand is shifted *x* places to the left (or right) and *x* is added to (or subtracted from) the exponent, as shown below.

 = = = 1234

## Basic operations

Given two numbers in scientific notation,

\[
x_0=m_0\times10^{n_0}
\]

and

\[
x_1=m_1\times10^{n_1}
\]

Multiplication and division are performed using the rules for operation with exponentiation:

\[
x_0 x_1=m_0 m_1\times10^{n_0+n_1}
\]

and

\[
\frac{x_0}{x_1}=\frac{m_0}{m_1}\times10^{n_0-n_1}
\]

Some examples are:

\[
5.67\times10^{-5} \times 2.34\times10^2 \approx 13.3\times10^{-5+2} = 13.3\times10^{-3} = 1.33\times10^{-2}
\]

and

\[
\frac{2.34\times10^2}{5.67\times10^{-5 \approx 0.413\times10^{2-(-5)} = 0.413\times10^{7} = 4.13\times10^6
\]

Addition and subtraction require the numbers to be represented using the same exponential part, so that the significand can be simply added or subtracted:
block indent
