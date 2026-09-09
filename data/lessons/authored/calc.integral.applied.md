# Applications of the Definite Integral

## Work by a variable force

When a force \\(F(x)\\) acts along the \\(x\\)-axis from \\(x = a\\) to \\(x = b\\), the work done is

\\[W = \int_{a}^{b} F(x)\,dx\\]

For a constant force this reduces to \\(W = F\cdot d\\). A spring with Hooke's-law force \\(F(x) = kx\\) stretched from \\(0\\) to \\(L\\) stores work \\(W = kL^{2}/2\\).

### Example

\\(F(x) = 2x\\) from \\(0\\) to \\(3\\): \\(W = \int_{0}^{3} 2x\,dx = 9\\) joules.

## Center of mass

A thin rod on \\([a,b]\\) with density \\(\rho(x)\\) has total mass \\(M = \int_{a}^{b}\rho(x)\,dx\\) and center of mass

\\[\bar{x} = \frac{1}{M}\int_{a}^{b} x\,\rho(x)\,dx\\]

For uniform density the center is the midpoint. Denser regions pull \\(\bar{x}\\) toward themselves: for \\(\rho(x) = x\\) on \\([0,1]\\), \\(M = 1/2\\) and \\(\bar{x} = 2/3\\).

## Surface of revolution

Revolving \\(y = f(x) \geq 0\\) about the \\(x\\)-axis sweeps bands of radius \\(f(x)\\) and slant width \\(\sqrt{1+(f')^{2}}\,dx\\):

\\[S = 2\pi\int_{a}^{b} f(x)\sqrt{1+(f'(x))^{2}}\,dx\\]

The \\(2\pi f(x)\\) factor is the circumference; the square root corrects for slant. Revolving a semicircle of radius \\(r\\) recovers \\(S = 4\pi r^{2}\\).
