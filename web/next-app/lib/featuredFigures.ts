import type { VisualSpec } from '@/components/visual/types'

// The three figures the announcement shows. They are the ones that read without
// a caption: the unit circle driving a sine wave, a vortex whose circulation is
// the curl times the area, and a washer slice that is the volume formula.
//
// This is launch content, not a gallery. Change these three specs to change what
// the announcement page shows.
export const FEATURED_FIGURES: VisualSpec[] = [
  {
    type: 'trig.unitcircle',
    title: 'The unit circle and the sine wave',
    caption: 'Turn the angle and the height of the point is the sine. The wave beside it is the same height, unrolled.',
    params: { angle: 45, showSin: true, showCos: true },
  },
  {
    type: 'calc.vectorfield',
    title: 'Circulation is curl times area',
    caption: 'The red loop is the path. The circulation around it equals the curl at the centre times the area inside, which is Green\'s theorem drawn rather than asserted.',
    params: { preset: 'vortex', radius: 1.2, showStreamlines: true },
  },
  {
    type: 'calc.solidrev',
    title: 'The slice is the formula',
    caption: 'A washer is the outer disk minus the hole. Move the slice and the volume it contributes is the number in the readout.',
    params: { mode: 'washer', slice: 1.2 },
  },
]
