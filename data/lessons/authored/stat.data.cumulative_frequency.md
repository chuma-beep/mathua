# Cumulative frequency tables

A cumulative frequency table adds the counts of a frequency table as it goes, so each entry gives the number of observations at or below that class boundary.

## Running totals

A cumulative frequency is a running total, not a separate count. Take classes with frequencies 4, 7, 9 and 5. The cumulative entries are 4, 11, 20 and 25. Each entry repeats the previous one plus the next frequency, so the last entry equals the number of observations.

## Finding a total

The final cumulative entry is the sample size, which makes the table a check on every earlier step. If a table of 25 observations ends at 23, two observations have been dropped or misread. The total also locates the median: for 25 observations the median is the value at position 13, which falls in the third class with cumulative entry 20.

## Recovering a class

Subtract consecutive cumulative entries to recover a frequency that was lost. If the entries at the second and third boundaries are 11 and 20, the third class holds \(20 - 11 = 9\) observations. The same subtraction recovers every class except the first, whose frequency is its own cumulative entry. A gap of zero is a legitimate empty class.

Cumulative counts are cheap to build and cheap to audit, which is why they appear wherever a distribution has to be read quickly.
