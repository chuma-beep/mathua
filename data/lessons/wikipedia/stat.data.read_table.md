> Content sourced from [Wikipedia](https://en.wikipedia.org/wiki/Table_%28information%29) — CC BY-SA 4.0

# Read data from tables and charts

A **table** is an arrangement of information or data, typically in rows and columns, or possibly in a more complex structure. Tables are widely used in communication, research, and data analysis. Tables appear in print media, handwritten notes, computer software, architectural ornamentation, traffic signs, and many other places. The precise conventions and terminology for describing tables vary depending on the context. Further, tables differ significantly in variety, structure, flexibility, notation, representation and use. Information or data conveyed in table form is said to be in **tabular** format (adjective). In books and technical articles, tables are typically presented apart from the main text in numbered and captioned floating blocks.

A **table cell** is one grouping within a chart table used for storing information or data. Cells are grouped horizontally (rows of cells) and vertically (columns of cells). Each cell contains information relating to the combination of the row and column headings it is collinear with.

## Basic description
A table consists of an ordered arrangement of **rows** and **columns**. This is a simplified description of the most basic kind of table. Certain considerations follow from this simplified description:

* the term **row** has several common synonyms (e.g., record, k-tuple, n-tuple, vector);
* the term **column** has several common synonyms (e.g., field, parameter, property, attribute, stanchion);
* a column is usually identified by a name;
* a column name can consist of a word, phrase or a numerical index;
* the intersection of a row and a column is called a cell.

The elements of a table may be grouped, segmented, or arranged in many different ways, and even nested recursively. Additionally, a table may include metadata, annotations, a header, a footer or other ancillary features.

### Simple table
The following illustrates a simple table with four columns and nine rows. The first row is not counted, because it is only used to display the column names. This is called a "header row".

{| class="wikitable col1left col2left" style="text-align:center"
|+Age table
|-
! First name !! Last name !! Age !! Gender
|-
| Tinu || Elejogun || 14 || F
|-
| Javier || Zapata || 28 || M
|-
| Lily || McGarrett || 18 || F
|-
| Olatunkbo || Chijiaku || 22 || M
|-
| Adrienne || Anthoula || 22 || M
|-
| Axelia || Athanasios || 22 || M
|-
| Jon-Kabat || Zinn || 22 || M
|-
| Thabang || Mosoa || 15 || F
|-
| Rhian || Ellis || 12 || M
|-
|}

### Multi-dimensional table


The concept of **dimension** is also a part of basic terminology. Any "simple" table can be represented as a "multi-dimensional"
table by normalizing the data values into ordered hierarchies. A common example of such a table is a multiplication table.

{| class="wikitable" style="text-align:center"
|+ style="white-space:nowrap" |Multiplication table
|-
!×!! 1 !! 2 !! 3
|-
! 1
| 1 || 2 || 3
|-
! 2
| 2 || 4 || 6
|-
! 3
| 3 || 6 || 9
|}

In multi-dimensional tables, each cell in the body of the table (and the value of that cell) relates to the values at the beginnings of the column (i.e. the header), the row, and other structures in more complex tables. This is an injective relation: each combination of the values of the headers row (row 0, for lack of a better term) and the headers column (column 0 for lack of a better term) is related to a unique cell in
the table:

* Column 1 and row 1 will only correspond to cell (1,1);
* Column 1 and row 2 will only correspond to cell (2,1) etc.

The first column often presents information dimension description by which the rest of the table is navigated. This column is called "stub column". Tables may contain three or multiple dimensions and can be classified by the number of dimensions. Multi-dimensional tables may have super-rows - rows that describe additional dimensions for the rows that are presented below that row and are usually grouped in a tree-like structure.  This structure is typically visually presented with an appropriate number of white spaces in front of each stub's label.

In literature tables often present numerical values, cumulative statistics, categorical values, and at times parallel descriptions in form of text. They can condense large amount of information to a limited space and therefore they are popular in scientific literature in many fields of study.


## Generic representation
As a communication tool, a table allows a form of generalization of information from an unlimited number of different social or scientific contexts. It provides a familiar way to convey information that might otherwise not be obvious or readily understood.

For example, in the following diagram, two alternate representations of the same information are presented side by side. On the left is the NFPA 704 standard "fire diamond" with example values indicated and on the right is a simple table displaying the same values, along with additional information. Both representations convey essentially the same information, but the tabular representation is arguably more comprehensible to someone who is not familiar with the NFPA 704 standard. The tabular representation may not, however, be ideal for every circumstance (for example because of space limitations, or safety reasons).
{| class="wikitable" style="margin-left:auto; margin-right:auto; text-align:center;"
|+Fire diamond
|-
! Standard Representation
! Tabular Representation
|-
| style="text-align:center;" |
|
{| class="wikitable" style="text-align: center"
|-
|+ style="white-space:nowrap" |Risk levels of hazardous materials in this facility
|-
! Health Risk !! Flammability !! Reactivity !! Special
|-
| Level 3 || Level 2 || Level 1 ||
|}
|}

## Information technology
### Software applications
Modern software applications give users the ability to generate, format, and edit tables and tabular data for a wide variety of uses, such as in: word processing and spreadsheet applications, presentation software, and in HTML or another markup language.

#### HTML
Table cells are a key component in HTML and webpage building. It is part of the  component. A programmer may specify dimensions for a table cell, and use them to hold sections of webpages. A table cell in HTML is a non-empty element and is supposed to always be closed. There are two different kinds of table cell in HTML, namely: normal table cell and header cell. **<td>** denotes a table cell, the name implying 'data', while **<th>** denotes a table 'header'. The two can be used interchangeably, but it is recommended that header cell be only used for the top and side headers of a table. Furthermore, a table cell must be nested within a **<table>** tag *and* a **<tr>** (table row) tag. If there are more table cell tags in any given row than in any other, the particular **<tr>** must be given a  attribute declaring how many columns of cells wide it should be. By using the **rowspan** and **colspan** attributes, developers can combine multiple rows or columns, allowing them to design more complex and visually structured tables.

The following table illustrates usage of colspan and rowpan:
{| width="80%" align="center" style="font-variant:small-caps"
| style="background:var(--background-color-inverted, black)" width="5%" |
| style="background:var(--background-color-inverted, black)" width="5%" |
| style="background:var(--background-color-inverted, black)" width="5%" |
|<-- *This row has three* **table data cells**
|-
| colspan="2" style="background:var(--background-color-progressive, blue);" |
| style="background:var(--background-color-inverted, black)" |
|<-- *This row has two. The first uses* **colspan="2"**
|-
| rowspan="2" style="background:var(--background-color-destructive, red)" |
| style="background:var(--background-color-inverted, black)" |
| style="background:var(--background-color-inverted, black)" |
|<-- *This row has three table data cells, but one spans two rows because it uses* **rowspan="2"**
|-
| style="background:var(--background-color-inverted, black)" |
| style="background:var(--background-color-inverted, black)" |
|<-- *This row has only two table data cells, because its first is being taken up*
|}

The following is an example of an HTML table containing 4 cells:
{| border="1"
|Cell 1
|Cell 2
|-
|Cell 3
|Cell 4
|}
HTML source:


      Cell 1


      Cell 2




      Cell 3


      Cell 4




### Software development
Tables have uses in software development for both high-level specification and low-level implementation.
Usage in software specification can encompass ad hoc inclusion of simple decision tables in textual documents through to the use of tabular specification methodologies, examples of which include Software Cost Reduction and Statestep.
Proponents of tabular techniques, among whom David Parnas is prominent, emphasize their understandability, as well as the quality and cost advantages of a format allowing systematic inspection, while corresponding shortcomings experienced with a graphical notation were cited in motivating the development of at least two tabular approaches.

At a programming level, software may be implemented using constructs generally represented or understood as tabular, whether to store data (perhaps to memoize earlier results), for example, in arrays or hash tables, or control tables determining the flow of program execution in response to various events or inputs.

### Databases
Database systems often store data in structures called tables; in which columns are data fields and rows represent data records.

## Other uses
There are several specific situations in which tables are routinely used as a matter of custom or formal convention.

### Publishing


### Mathematics

Examples include:
* Arithmetic (Multiplication table)
* Logic (Truth table)

### Natural sciences
In natural sciences, uses of tables include the periodic table in chemistry, and tide tables in oceanography.

#### Periodic table


#### Tide table


## Historical relationship to furniture
In medieval counting houses, the tables were covered with a piece of checkered cloth, to count money.*Exchequer* is an archaic term for the English institution which accounted for money owed to the monarch. Thus the checkerboard tables of stacks of coins are a concrete realization of this information.

## See also
* Chart
* Diagram
* Abstract data type
* Column (database)
* Information graphics
* Periodic table
* Reference table
* Row (database)
* Table (database)
* Table (HTML)
* Tensor
* Dependent and independent variables
* Zebra striping

## References


## External links

*

