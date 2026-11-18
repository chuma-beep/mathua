import type { ComponentType } from 'react'
import type { ZodTypeAny } from 'zod'

// The visual subsystem's contract. A shard names one of the templates below and
// hands it parameters; the registry validates and the dispatcher renders.
//
// `params` is unknown on purpose. Each template owns its schema and parses its
// own input, so adding a template never widens a shared type and a malformed
// payload is a quiet fallback rather than a broken lesson.

export interface VisualSpec {
  type: string
  title?: string
  caption?: string
  params?: unknown
}

export interface TemplateProps {
  params: unknown
  reducedMotion: boolean
  /** Authored chrome from the shard's VisualSpec, not from the parameter object. */
  title?: string
  caption?: string
}

export interface TemplateModule {
  default: ComponentType<TemplateProps>
}

export interface VisualTemplateEntry {
  schema: ZodTypeAny
  load: () => Promise<TemplateModule>
}

export interface ReadoutItem {
  label: string
  value: string
}
