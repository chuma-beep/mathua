import type { VisualTemplateEntry } from './types'
import {
  basisChangeSchema,
  complexFieldSchema,
  coordinateSchema,
  complexMultiplySchema,
  complexMapSchema,
  contourIntegralSchema,
  conicSchema,
  contourSchema,
  eigenSchema,
  familySchema,
  limitSchema,
  matrix2dSchema,
  normalSchema,
  regressionSchema,
  riemannSchema,
  samplingSchema,
  solid3dSchema,
  spreadSchema,
  tangentSchema,
  taylorSchema,
  vectorCombinationSchema,
  transform2dSchema,
} from './schemas'

// The closed registry. Adding a template is one entry here, one schema above,
// one component under templates/, and one id in data/visualizations/types.json.
// Nothing else in the lesson pipeline changes.
//
// The loaders are dynamic imports: a lesson with no visual, or with a different
// visual, never downloads this template's code.
export const VISUAL_TEMPLATES: Record<string, VisualTemplateEntry> = {
  'linalg.matrix2d': {
    schema: matrix2dSchema,
    load: () => import('./templates/MatrixTransform2D'),
  },
  'calc.tangent': {
    schema: tangentSchema,
    load: () => import('./templates/TangentLine'),
  },
  'calc.riemann': {
    schema: riemannSchema,
    load: () => import('./templates/RiemannSum'),
  },
  'complex.multiply': {
    schema: complexMultiplySchema,
    load: () => import('./templates/ComplexMultiply'),
  },
  'stat.normal': {
    schema: normalSchema,
    load: () => import('./templates/NormalCurve'),
  },
  'geo.transform2d': {
    schema: transform2dSchema,
    load: () => import('./templates/Transform2D'),
  },
  'geo.conic': {
    schema: conicSchema,
    load: () => import('./templates/ConicSection'),
  },
  'calc.contour': {
    schema: contourSchema,
    load: () => import('./templates/ContourMap'),
  },
  'complex.levels': {
    schema: complexFieldSchema,
    load: () => import('./templates/ComplexLevels'),
  },
  'complex.contour': {
    schema: contourIntegralSchema,
    load: () => import('./templates/ContourIntegral'),
  },
  'stat.family': {
    schema: familySchema,
    load: () => import('./templates/DistributionFamilies'),
  },
  'geo.coord': {
    schema: coordinateSchema,
    load: () => import('./templates/CoordinateWorkbench'),
  },
  'complex.map': {
    schema: complexMapSchema,
    load: () => import('./templates/ComplexMapping'),
  },
  'linalg.eigen': {
    schema: eigenSchema,
    load: () => import('./templates/EigenDirections'),
  },
  'calc.taylor': {
    schema: taylorSchema,
    load: () => import('./templates/TaylorSeries'),
  },
  'calc.limit': {
    schema: limitSchema,
    load: () => import('./templates/LimitApproach'),
  },
  'stat.sampling': {
    schema: samplingSchema,
    load: () => import('./templates/SamplingDistribution'),
  },
  'stat.regression': {
    schema: regressionSchema,
    load: () => import('./templates/LeastSquares'),
  },
  'solid3d': {
    schema: solid3dSchema,
    load: () => import('./templates/Solid3D'),
  },
  'stat.spread': {
    schema: spreadSchema,
    load: () => import('./templates/SpreadMeasure'),
  },
  'linalg.vector': {
    schema: vectorCombinationSchema,
    load: () => import('./templates/VectorCombination'),
  },
  'linalg.basis': {
    schema: basisChangeSchema,
    load: () => import('./templates/BasisChange'),
  },
}

export const VISUAL_TYPES = Object.keys(VISUAL_TEMPLATES)

export function isVisualType(type: string): boolean {
  return Object.prototype.hasOwnProperty.call(VISUAL_TEMPLATES, type)
}
