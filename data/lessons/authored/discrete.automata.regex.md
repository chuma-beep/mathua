# Regular expressions and finite automata

Regular expressions and finite automata describe the same languages. A DFA can be converted to a regular expression by eliminating states, and a regular expression can be converted to an NFA.

## DFAs and regular languages

A deterministic finite automaton has one transition per state and letter. Its accepted language is regular. Counting accepted words of a fixed length is a path-counting problem on the transition graph.

## Counting words by length

Start with one path at the start state and propagate counts through the transitions. After the required number of letters, add the counts at accepting states. Summing over lengths gives the number of accepted words up to a bound.

## Membership tests

Run the DFA from the start state, following one transition per letter. The word is accepted exactly when the final state is accepting. This test is deterministic and takes time proportional to the word length.

The three sections form one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section so the three moves stay connected.
