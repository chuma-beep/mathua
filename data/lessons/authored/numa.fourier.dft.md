# The discrete Fourier transform

The discrete Fourier transform writes a finite signal as a sum of complex exponentials. Each coefficient measures one frequency.

## The transform

For a length-\(N\) signal, \(X_k=\sum_{n=0}^{N-1}x_n e^{-2\pi i kn/N}\). The index \(k\) labels the frequency.

The zero-frequency coefficient \(X_0\) is the sum of the samples. It measures the average or DC level.

The transform is linear, so the DFT of a sum is the sum of the DFTs.

## Special coefficients

At \(k=N/2\), the exponential alternates sign. The coefficient is the alternating sum of the samples.

For a real signal, the remaining coefficients come in conjugate pairs. The transform still has \(N\) complex outputs.

The Nyquist coefficient is real for a real signal, which is why it can be read as an alternating sum.

## Computational cost

A direct evaluation of all \(N\) outputs costs \(N^2\) complex multiplications. The FFT reduces that to order \(N\log N\).

The DFT is the mathematical object. The FFT is an algorithm for computing it quickly.

The saving matters most when \(N\) is large, which is why real-time signal processing depends on the FFT.
