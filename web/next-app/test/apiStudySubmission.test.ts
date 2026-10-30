import { describe, expect, it, vi, beforeEach } from 'vitest'

vi.mock('../lib/auth', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/auth')>()
  return {
    ...mod,
    authedFetch: vi.fn(),
    getAuthHeaders: () => ({}),
    getGuestId: () => null,
  }
})

import { authedFetch } from '../lib/auth'
import { getErrorMessage, getTopicLockedPrerequisites, submitStudyAnswer } from '../lib/api'

const ok = () =>
  new Response(JSON.stringify({ correct: true, feedback: 'Correct', xp: 1 }), {
    status: 200,
    headers: { 'Content-Type': 'application/json' },
  })

describe('submitStudyAnswer submission id', () => {
  beforeEach(() => {
    vi.mocked(authedFetch).mockReset()
    vi.mocked(authedFetch).mockImplementation(async () => ok())
  })

  it('sends one non-empty submission id per logical answer', async () => {
    await submitStudyAnswer('a', '4', 5, '2+2')
    await submitStudyAnswer('a', '4', 5, '2+2')

    expect(authedFetch).toHaveBeenCalledTimes(2)
    const bodies = vi.mocked(authedFetch).mock.calls.map((call) => {
      const init = call[1] as RequestInit
      return JSON.parse(String(init.body)) as { submission_id?: string }
    })
    expect(bodies[0].submission_id).toBeTruthy()
    expect(bodies[1].submission_id).toBeTruthy()
    expect(bodies[0].submission_id).not.toBe(bodies[1].submission_id)
    expect(vi.mocked(authedFetch).mock.calls[0][1]).toMatchObject({ retryAttempts: 1 })
  })

  it('preserves the structured topic-lock reason and prerequisite ids', async () => {
    vi.mocked(authedFetch).mockResolvedValueOnce(new Response(JSON.stringify({
      error: { code: 'TOPIC_LOCKED', message: 'This topic has unmet prerequisites.', details: { prerequisiteConceptIds: ['arith.mult.2_5_10'] } },
    }), { status: 403 }))

    const error = await submitStudyAnswer('arith.exp.concept', '256', 5, '4^5').catch(e => e)
    expect(getErrorMessage(error)).toBe('This topic has unmet prerequisites.')
    expect(getTopicLockedPrerequisites(error)).toEqual(['arith.mult.2_5_10'])
  })

  it('keeps ordinary forbidden errors distinct from topic locks', async () => {
    vi.mocked(authedFetch).mockResolvedValueOnce(new Response(JSON.stringify({ error: 'student_id not permitted' }), { status: 403 }))
    const error = await submitStudyAnswer('a', '4', 5, '2+2').catch(e => e)
    expect(getErrorMessage(error)).toBe('student_id not permitted')
    expect(getTopicLockedPrerequisites(error)).toBeNull()
  })
})
