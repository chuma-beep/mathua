# Column generation

Column generation solves a program with too many columns to write down. A restricted master carries a few columns and a pricing problem finds the next useful one.

## Too many columns

The full program has one variable for every member of a huge set. Most of those variables are zero at the optimum, so only a few need to exist.

## The pricing problem

The dual of the restricted master supplies prices. Pricing looks for the column with the most negative reduced cost under those prices.

## The restricted master

The master holds the columns found so far and is solved to optimality each round. It grows by one column per improving price until none is negative.

To check your understanding, take the column's cost, compare the reduced costs and count the columns added.

Pricing is where the structure of the application enters. A fast pricing oracle is what makes the method practical.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller coefficients and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
