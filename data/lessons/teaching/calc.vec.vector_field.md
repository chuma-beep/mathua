# Vector Fields

## A Vector at Every Point

A vector field assigns a vector to each point of a region. In two dimensions it can be written \(\mathbf F(x,y)=\langle P(x,y),Q(x,y)\rangle\), and in three dimensions it has a third component \(R(x,y,z)\). The components are ordinary functions, so a field can be evaluated at a point by substituting the coordinates into each component.

A vector field can be visualised by drawing short arrows at many points. The arrows show direction and relative magnitude, but they do not show the absolute magnitude unless a scale is given. Fields that point away from a source, circulate around an axis, or stay parallel are the basic patterns to recognise.

## Divergence and Curl

Divergence measures the local expansion of a field. In three dimensions it is the sum \(P_x+Q_y+R_z\) of the partial derivatives of the components. A positive divergence means the field spreads out, while negative divergence means it converges. For a linear field such as \(\langle ax,by,cz\rangle\), the divergence is the constant \(a+b+c\).

Curl measures local rotation. It is the vector of differences of partial derivatives, and for \(\mathbf F=\langle ay,bz,cx\rangle\) it equals \(\langle -b,-c,-a\rangle\). A field with zero curl is called irrotational. Divergence and curl are independent pieces of information: a field can expand without rotating, rotate without expanding, or do both.

## Evaluating a Field

To evaluate a field, substitute the point into every component. For \(\mathbf F=\langle ay,bz,cx\rangle\) at \((x_0,y_0,z_0)\), the first component is \(ay_0\), the second is \(bz_0\) and the third is \(cx_0\). The result is a single vector, not a scalar, so keep the three components separate in the answer.

The same field can be used to ask for a component, a magnitude or a directional derivative. Read the question carefully: a component is one number, a magnitude is a square root, and a divergence or curl is a different combination of derivatives. Checking which operation is requested is the main defence against answering the wrong question.
