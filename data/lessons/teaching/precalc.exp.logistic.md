# Logistic Growth Models

## The logistic formula

The model is P(t) = L / (1 + A e^(-kt)).
L is the carrying capacity of the environment.
A and k set the starting point and the speed of growth.

The logistic model starts like exponential growth and then bends toward a ceiling. The bend appears because the factor 1 - P/L shrinks as the population approaches L. In a restricted habitat that ceiling is usually set by food, space or waste.

## Reading the parameters

The initial population is L / (1 + A).
As t grows the exponential term vanishes and P approaches L.
A large k makes the population approach the capacity sooner.

When t is zero the term e^0 equals one, so the denominator is 1 + A and the starting value follows immediately. When t becomes large the exponential term tends to zero and the denominator tends to one, which is why the model levels off at L.

## Where growth is fastest

The steepest point of the curve is at half the carrying capacity.
Below L/2 the population grows faster each period.
Above L/2 the growth slows as the environment fills.

The growth rate is the product of k, the current population and the unused fraction of capacity. That product is largest when the population is exactly L/2. Harvesting or conservation decisions often use this point as a reference.
