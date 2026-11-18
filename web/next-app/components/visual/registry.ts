import type { VisualTemplateEntry } from './types'
import {
  complexMultiplySchema,
  matrix2dSchema,
  normalSchema,
  riemannSchema,
  tangentSchema,
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
}

export const VISUAL_TYPES = Object.keys(VISUAL_TEMPLATES)

export function isVisualType(type: string): boolean {
  return Object.prototype.hasOwnProperty.call(VISUAL_TEMPLATES, type)
}
