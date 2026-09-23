import { describe, it, expect } from 'vitest'
import { subdomainLeaf, formatConceptPath } from '../lib/conceptDisplay'

describe('subdomainLeaf', () => {
  it('strips the domain namespace prefix', () => {
    expect(subdomainLeaf('abstract_algebra.rings')).toBe('rings')
    expect(subdomainLeaf('arithmetic.order_of_operations')).toBe('order_of_operations')
  })

  it('leaves bare subdomains alone', () => {
    expect(subdomainLeaf('rings')).toBe('rings')
    expect(subdomainLeaf('')).toBe('')
  })
})

describe('formatConceptPath', () => {
  it('renders domain.leaf without doubling', () => {
    expect(formatConceptPath('abstract_algebra', 'abstract_algebra.rings', 'abstract.rings.polynomial'))
      .toBe('abstract_algebra.rings · abstract.rings.polynomial')
  })

  it('falls back to domain when subdomain is empty or identical', () => {
    expect(formatConceptPath('topology', '', 'topo.basic.open')).toBe('topology · topo.basic.open')
    expect(formatConceptPath('topology', 'topology', 'topo.basic.open')).toBe('topology · topo.basic.open')
  })
})
