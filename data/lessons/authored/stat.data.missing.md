# Missing-data mechanisms and imputation

Missing data can distort an analysis. The mechanism determines whether a correction is possible.

## Identify the missingness mechanism

MCAR means missingness is independent of the data. MAR means it depends on observed variables. MNAR means it depends on the missing value itself.

## Impute a value

Mean imputation replaces a missing value by the observed mean. It preserves the sample size but shrinks variance. Multiple imputation adds the uncertainty back.

## Measure the bias

Compare the complete-case estimate with the imputed estimate. A large difference signals that the missingness is informative. The direction of the bias depends on the mechanism.

A complete answer states the governing relation, substitutes the given values and reads the result in the language of the problem. Where a second method is available, compare the two results before moving to the next question. The three sections are one chain, and the practice sequence asks one family from each section so the learner keeps the definition, the computation and the interpretation connected.
