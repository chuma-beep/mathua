# Finite State Machines

## States and Transitions

A finite state machine has a finite set of states, an input alphabet and a transition function that maps a state and an input symbol to a new state. Some states are designated as accepting.

The machine reads an input string one symbol at a time, starting in the start state. The string is accepted when the machine finishes in an accepting state.

The entire behaviour is determined by the transition table, so two machines are equivalent when their tables produce the same accept or reject answer for every string.

## Running a Machine

To run a machine, follow the arrow for each input symbol and record the state after every step. A common design counts positions modulo a fixed number, so the state after k steps is the starting state plus k modulo the modulus.

The number of states is fixed by the design, not by the length of the input. This is what separates finite automata from machines with unbounded memory.

Diagrams and tables are two views of the same transition function.

## Counting States and Strings

The product of two machines has one state for each pair of states, so its state count is the product of the two state counts. The number of strings of length n over an alphabet of size a is a to the power n.

Accepting counts depend on the transition structure. For a modulo counter with a single accepting state, exactly one residue class is accepted in each length.

These counts explain why finite automata are useful for pattern matching and protocol checking.
