import type { VisualTemplateEntry } from './types'
import {
  analysisSchema,
  basisChangeSchema,
  complexFieldSchema,
  coordinateSchema,
  complexMultiplySchema,
  complexMapSchema,
  contourIntegralSchema,
  conicSchema,
  contourSchema,
  discreteSchema,
  eigenSchema,
  familySchema,
  fractionLineSchema,
  functionGraphSchema,
  gaussianSchema,
  graphSchema,
  groupSchema,
  accumulationSchema,
  conditionalAreaSchema,
  inferenceSchema,
  limitSchema,
  matrix2dSchema,
  modelFitSchema,
  modularSchema,
  normalSchema,
  numberLineSchema,
  odeSlopeSchema,
  optimizationSchema,
  phaseSchema,
  probabilityTreeSchema,
  regressionSchema,
  riemannSchema,
  rootsSchema,
  samplingSchema,
  seriesSumSchema,
  sinusoidSchema,
  solidRevSchema,
  solid3dSchema,
  spreadSchema,
  surfaceSchema,
  tangentSchema,
  taylorSchema,
  vectorCombinationSchema,
  transform2dSchema,
  unitCircleSchema,
  vectorFieldSchema,
  regionSchema,
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
  'trig.unitcircle': {
    schema: unitCircleSchema,
    load: () => import('./templates/UnitCircle'),
  },
  'trig.sinusoid': {
    schema: sinusoidSchema,
    load: () => import('./templates/SinusoidExplorer'),
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
  'linalg.elimination': {
    schema: gaussianSchema,
    load: () => import('./templates/GaussianElimination'),
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
  'arith.line': {
    schema: numberLineSchema,
    load: () => import('./templates/NumberLineLab'),
  },
  'ode.slopefield': {
    schema: odeSlopeSchema,
    load: () => import('./templates/SlopeField'),
  },
  'numa.roots': {
    schema: rootsSchema,
    load: () => import('./templates/RootIteration'),
  },
  'graph.network': {
    schema: graphSchema,
    load: () => import('./templates/GraphNetwork'),
  },
  'prob.tree': {
    schema: probabilityTreeSchema,
    load: () => import('./templates/ProbabilityTree'),
  },
  'frac.line': {
    schema: fractionLineSchema,
    load: () => import('./templates/FractionLine'),
  },
  'function.graph': {
    schema: functionGraphSchema,
    load: () => import('./templates/FunctionGraph'),
  },
  'nt.modular': {
    schema: modularSchema,
    load: () => import('./templates/ModularLab'),
  },
  'abstract.group': {
    schema: groupSchema,
    load: () => import('./templates/GroupTable'),
  },
  'discrete.lab': {
    schema: discreteSchema,
    load: () => import('./templates/DiscreteLab'),
  },
  'optim.path': {
    schema: optimizationSchema,
    load: () => import('./templates/OptimizationLab'),
  },
  'analysis.lab': {
    schema: analysisSchema,
    load: () => import('./templates/AnalysisLab'),
  },
  'ml.fit': {
    schema: modelFitSchema,
    load: () => import('./templates/ModelFit'),
  },
  'stat.inference': {
    schema: inferenceSchema,
    load: () => import('./templates/InferenceLab'),
  },
  'ode.phase': {
    schema: phaseSchema,
    load: () => import('./templates/PhasePortrait'),
  },
  'calc.vectorfield': {
    schema: vectorFieldSchema,
    load: () => import('./templates/VectorFieldLab'),
  },
  'calc.surface': {
    schema: surfaceSchema,
    load: () => import('./templates/SurfaceLab'),
  },
  'calc.region': {
    schema: regionSchema,
    load: () => import('./templates/RegionIntegral'),
  },
  'calc.seriessum': {
    schema: seriesSumSchema,
    load: () => import('./templates/SeriesLab'),
  },
  'calc.accumulation': {
    schema: accumulationSchema,
    load: () => import('./templates/AccumulationLab'),
  },
  'calc.solidrev': {
    schema: solidRevSchema,
    load: () => import('./templates/SolidRevolution'),
  },
  'prob.area': {
    schema: conditionalAreaSchema,
    load: () => import('./templates/ConditionalArea'),
  },
}

export const VISUAL_TYPES = Object.keys(VISUAL_TEMPLATES)

export function isVisualType(type: string): boolean {
  return Object.prototype.hasOwnProperty.call(VISUAL_TEMPLATES, type)
}
