# Aliasing

Aliasing occurs when a signal is sampled too slowly. High frequencies then appear as lower frequencies in the sampled data.

## The Nyquist limit

The Nyquist frequency is half the sampling rate. Frequencies above it cannot be represented unambiguously.

To avoid aliasing, sample faster than twice the highest frequency in the signal.

The limit is a property of the sampling process, not of the reconstruction algorithm.

## Folding frequencies

A frequency \(f\) sampled at \(f_s\) appears at \(|f-kf_s|\) for an integer \(k\) that brings it below Nyquist.

The alias is not a physical lower frequency. It is the same samples interpreted at a different frequency.

Two different frequencies can produce identical samples, so the original cannot be recovered without more information.

## Preventing aliasing

An anti-aliasing filter removes frequencies above Nyquist before sampling. Increasing the sampling rate also raises the limit.

Once aliasing has occurred, no post-processing can separate the folded frequencies in general.

A higher sampling rate costs storage and power, so the filter and rate are chosen together.
