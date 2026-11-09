# Nondeterministic finite automata

A nondeterministic finite automaton may have several transitions for one symbol and may change state without reading input. Determinisation turns it into an equivalent DFA.

## States and transitions

An NFA has a finite state set, an alphabet, a transition relation and a set of accepting states. A transition relation may offer zero, one or several next states. Epsilon edges allow a state change without consuming a letter.

## Subset construction

The subset construction starts from the epsilon closure of the start state. Each letter moves every state in the current subset and then closes the result under epsilon edges. The reachable subsets become the states of an equivalent DFA.

## Membership and epsilon closures

A string is accepted when some run ends in an accepting state. Simulating all runs at once means carrying a set of states and closing it after every letter. The epsilon closure is the set of states reachable by epsilon edges alone.

The three sections form one chain. The first names the object, the second performs the computation and the third reads the result in context. Practice asks one family from each section so the three moves stay connected.
