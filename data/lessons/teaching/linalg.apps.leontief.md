# Leontief Input-Output Models

An economy can be divided into sectors that consume one another's output. A Leontief model records those internal demands and asks how much each sector must produce to meet a final demand.

## Consumption and Technology Matrices

Let \(A\) be the technology matrix. Its entry \(a_{ij}\) is the amount of sector \(i\)'s output needed to produce one unit of sector \(j\)'s output. Each column therefore describes the inputs consumed by one sector.

For a two-sector economy the matrix has four entries: \(a_{11}\), \(a_{12}\), \(a_{21}\) and \(a_{22}\). The first column lists what sector 1 consumes, and the second column lists what sector 2 consumes.

The row sums of \(A\) measure how much of a sector's output is consumed by all sectors per unit of total production. A row sum below 1 leaves output for final demand.

## Solving (I - A)x = d

Let \(x\) be total production and \(d\) the final demand. Internal consumption is \(Ax\), so the balance equation is

\(x = Ax + d\).

Rearranging gives

\((I-A)x = d\).

When \(I-A\) is invertible, the production plan is

\(x = (I-A)^{-1}d\).

For a one-sector economy with technology coefficient \(a\), this reduces to \(x=d/(1-a)\). A sector that consumes half of its own output must produce twice the final demand.

## Productive Economies and the Inverse

The matrix \(I-A\) is invertible when its determinant is nonzero. For a two-sector model,

\(\det(I-A)=(1-a_{11})(1-a_{22})-a_{12}a_{21}\).

A positive determinant with nonnegative entries indicates that the economy can meet final demand without negative production. Such an economy is called productive.

The inverse \((I-A)^{-1}\) is the total requirements matrix. Its entries show how much extra output each sector must generate, directly and indirectly, for one additional unit of final demand. Leontief models therefore turn a system of balance equations into a practical planning tool.
