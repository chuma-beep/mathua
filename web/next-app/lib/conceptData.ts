import conceptsData from '../data/concepts.json'

// Typed view over data/concepts.json. The raw JSON import infers a large
// literal type that is awkward to iterate, which previously forced `as any[]`
// casts at every call site. One named contract here keeps callers typed.
export interface ConceptRecord {
  id: string
  label: string
  domain: string
  subdomain?: string
  grading_type?: string
  prerequisites?: string[]
  mastery_threshold?: { streak: number; avg_time_seconds: number }
}

export const concepts = conceptsData as ConceptRecord[]
