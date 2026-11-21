import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import AdminAudiencePage from '../app/admin/audience/page'
import { fetchAdminLocations } from '../lib/api'
import type { AdminLocations } from '../lib/api'

vi.mock('next/navigation', () => ({
  usePathname: () => '/admin/audience',
}))

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  fetchAdminLocations: vi.fn(),
}))

vi.mock('../hooks/useAdminMe', () => ({
  useAdminMe: () => ({
    me: { id: 'root', role: 'owner', permissions: ['admin.access'], assignable_roles: ['admin'] },
    loading: false,
    can: () => true,
    isStaff: true,
  }),
}))

vi.mock('../components/admin/AdminShell', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}))

// The globe is WebGL; jsdom has no context. The page's numbers are the part
// that can be wrong, and the table carries them independently of the canvas.
vi.mock('../components/admin/CountryGlobe', () => ({
  default: ({ countries }: { countries: unknown[] }) => (
    <div data-testid="globe">{countries.length} markers</div>
  ),
}))

vi.mock('../components/admin/DeviceGallery', () => ({
  default: ({ devices }: { devices: unknown[] }) => (
    <div data-testid="devices-3d">{devices.length} devices</div>
  ),
}))

function locations(over: Partial<AdminLocations> = {}): AdminLocations {
  const daily = Array.from({ length: 30 }, (_, i) => ({
    day: `2026-11-${String(i + 1).padStart(2, '0')}`,
    active: i % 4,
  }))
  const deviceDaily = daily.map((d, i) => ({
    day: d.day,
    mobile: i % 3,
    desktop: i % 2,
    tablet: i % 5 === 0 ? 1 : 0,
  }))
  return {
    countries: [
      { code: 'NG', known: 5, active_7d: 4, active_prev_7d: 2, new_7d: 2, new_prev_7d: 0 },
      { code: 'GB', known: 2, active_7d: 2, active_prev_7d: 1, new_7d: 0, new_prev_7d: 0 },
    ],
    located: 7,
    total: 10,
    active_7d: 6,
    active_prev_7d: 3,
    new_7d: 2,
    new_prev_7d: 0,
    visits_7d: 9,
    visits_prev_7d: 4,
    visits_30d: 42,
    devices: [
      { device: 'mobile', active_7d: 4, active_prev_7d: 2 },
      { device: 'desktop', active_7d: 2, active_prev_7d: 1 },
    ],
    daily,
    device_daily: deviceDaily,
    generated_at: '2026-11-21T10:00:00Z',
    ...over,
  }
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(fetchAdminLocations).mockResolvedValue(locations())
})

describe('the admin audience page', () => {
  it('separates where accounts are, where activity grew, and where signups came from', async () => {
    render(<AdminAudiencePage />)
    await screen.findByText('7 of 10')

    expect(screen.getByText('Active this week')).toBeTruthy()
    expect(screen.getByText('+100% against the previous 7 days')).toBeTruthy()
    expect(screen.getByText('First recorded week, no prior week to compare')).toBeTruthy()
    expect(screen.getByText('Visits this week')).toBeTruthy()
    expect(screen.getByText('+125% against the previous 7 days')).toBeTruthy()
    expect(screen.getByText('42 visits')).toBeTruthy()

    // The globe is decorative; the same two countries must be readable in
    // plain text, with the prior window and the change beside the count.
    // It loads through next/dynamic, so the first paint has no canvas.
    expect((await screen.findByTestId('globe')).textContent).toBe('2 markers')
    expect(screen.getAllByText('Nigeria').length).toBeGreaterThan(0)
    expect(screen.getAllByText('United Kingdom').length).toBeGreaterThan(0)
    expect(screen.getByText('Previous 7d')).toBeTruthy()
    expect(screen.getByText('Peak 3 learners in a day')).toBeTruthy()

    // A visit is a learner-day, so it is reported beside the head count and
    // the device split, never instead of it.
    expect(screen.getByText('Devices this week')).toBeTruthy()
    expect(screen.getByText('Visits per active learner')).toBeTruthy()
    expect(screen.getByText('1.5')).toBeTruthy()
    // The gallery loads through next/dynamic, like the globe.
    expect((await screen.findByTestId('devices-3d')).textContent).toBe('2 devices')
    expect(screen.getByText('Last 30 days')).toBeTruthy()
    expect(screen.getByText('Distinct learners per day')).toBeTruthy()
    expect(screen.getByText('One class per learner per day')).toBeTruthy()
  })

  it('says recording starts now rather than drawing an empty globe', async () => {
    vi.mocked(fetchAdminLocations).mockResolvedValue(locations({ countries: [], located: 0, total: 4, active_7d: 0, active_prev_7d: 0, new_7d: 0, visits_7d: 0, visits_prev_7d: 0, visits_30d: 0, devices: [] }))
    render(<AdminAudiencePage />)

    await screen.findByText('No countries recorded yet.')
    expect(screen.queryByTestId('globe')).toBeNull()
    expect(screen.getByText(/The first country is stored the next time/)).toBeTruthy()
  })

  it('reports a failed read instead of rendering zeroes as fact', async () => {
    vi.mocked(fetchAdminLocations).mockRejectedValue(new Error('Could not load locations'))
    render(<AdminAudiencePage />)

    const alert = await screen.findByRole('alert')
    expect(alert.textContent).toContain('Could not load locations')
  })
})
