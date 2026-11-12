# Federated learning

Federated learning trains a shared model across clients that keep their data local. Clients compute updates, and a server aggregates them.

## The round

Each round selects a subset of clients, sends the model and collects updates. The number of rounds sets the communication cost.

## Averaging

The server averages client updates, often weighted by sample count. A weighted average reduces to the pooled gradient when clients are honest.

## The heterogeneity

Client distributions differ, so their gradients point in different directions. That disagreement slows convergence relative to independent and identical data.

To check your understanding, select clients, weight by samples and measure the disagreement.

Federated learning moves computation to the data instead of the data to the model.

Each section above supplies the reasoning for one practice family, so work a numeric instance of each before moving on. Write the first step of every instance down before you compute the final value, because the grader checks the number but the method is what the next concept assumes. If an instance refuses to settle, redraw it with smaller numbers and check the sign of every term; the arithmetic in these families is exact, so a mismatch means a sign or an order of operations, not an approximation.
