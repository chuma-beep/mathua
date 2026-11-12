# Minimum-cost flow

A minimum-cost flow sends a required amount from source to sink at the least total cost, with capacities and per-unit edge prices.

## The cost

For a flow \(f\) and edge costs \(c_e\), the total cost is \(\sum_e c_e f_e\). The problem minimises this sum subject to capacity constraints and conservation at every vertex except the source and sink.

## Residual costs

Sending flow along an edge creates a residual edge in the opposite direction with cost \(-c_e\). Negative residual costs signal that rerouting can lower the total, which is why successive shortest paths need reduced costs or potentials.

## Optimality

A feasible flow of the required value is optimal exactly when the residual network contains no negative-cost cycle. The cycle condition is the flow analogue of complementary slackness.

To check your understanding, multiply flow by unit cost, reverse an edge and build the residual network.

Capacity and conservation fix the feasible flows, and the cost function picks the cheapest one.
