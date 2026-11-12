# Damped sinusoidal graphs

A damped sinusoid multiplies a wave by a decaying factor, so its amplitude shrinks while its period stays fixed.

## The decaying envelope

The factor is an exponential that bounds the wave from above and below. For
\(y = 5e^{-0.2t}\cos t\) the envelope is \(y = \pm 5e^{-0.2t}\). At \(t = 0\) the amplitude is 5
and at \(t = 10\) it is \(5e^{-2} \approx 0.677\), so the wave has lost most of its height in ten
units of time. The cosine still runs at one radian per unit of time, so the wave keeps its rhythm
while it loses size.

## Half-lives of amplitude

The amplitude halves when the exponential reaches one half, at \(t = \ln 2/0.2 \approx 3.466\),
where \(5e^{-0.693} = 2.5\). A second half-life takes the amplitude to 1.25 at about
\(t = 6.931\). A third takes it to 0.625 at about \(t = 10.397\). Equal time steps divide the
amplitude by equal factors, which is what a half-life means.

## Reading a damped graph

The period is still \(2\pi \approx 6.283\), so the peaks sit at \(t = 0\), \(t = 6.283\) and
\(t = 12.566\). Their heights are 5, then \(5e^{-1.257} \approx 1.423\), then
\(5e^{-2.513} \approx 0.405\). Each peak is 0.285 of the height of the one before, because
\(e^{-0.2(2\pi)} \approx 0.2846\). A graph whose successive peaks shrink by a constant ratio is
damped; the ratio gives the decay rate.

Measure the ratio between successive peaks, not the difference.
