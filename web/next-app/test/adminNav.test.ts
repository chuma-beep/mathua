import { describe, expect, it } from 'vitest'
import { adminBreadcrumb, isAdminNavActive } from '../components/admin/nav'

describe('admin navigation state', () => {
  it('treats Overview as the section root rather than a prefix', () => {
    expect(isAdminNavActive('/admin', '/admin')).toBe(true)
    expect(isAdminNavActive('/admin/users', '/admin')).toBe(false)
    expect(isAdminNavActive('/admin/reports', '/admin')).toBe(false)
  })

  it('keeps a section active on its subpages', () => {
    expect(isAdminNavActive('/admin/users/detail', '/admin/users')).toBe(true)
    expect(isAdminNavActive('/admin/lab/play', '/admin/lab')).toBe(true)
  })

  it('builds a breadcrumb only below a section root', () => {
    expect(adminBreadcrumb('/admin/users')).toEqual([])
    expect(adminBreadcrumb('/admin/users/detail')).toEqual([
      { label: 'Users', href: '/admin/users' },
      { label: 'Account', href: '/admin/users/detail' },
    ])
    expect(adminBreadcrumb('/admin/lab/play')[1].label).toBe('Playground')
  })
})
