import { API_BASE } from './apiBase'
import { getAuthHeaders, authedFetch } from './auth'
import { z } from 'zod'

export { API_BASE }

export const PrereqInfoSchema = z.object({
  id: z.string(),
  label: z.string(),
  status: z.string(),
  mastery_pct: z.number(),
})

const AuthResSchema = z.object({
  token: z.string(),
  student_id: z.string(),
  name: z.string(),
  diagnostic_completed: z.boolean(),
})

const ScoresSchema = z.object({
  lifetime_points: z.number(),
  weekly_score: z.number(),
  speed_bonus: z.number(),
  concepts_mastered: z.number(),
  current_streak: z.number(),
  level: z.string(),
  xp_total: z.number(),
  xp_today: z.number(),
  daily_xp_goal: z.number(),
  spaced_reps: z.record(z.string(), z.number()).optional(),
  avg_learning_speed: z.number().optional(),
  xp_since_quiz: z.number().optional(),
  quiz_due: z.boolean().optional(),
  quiz_active: z.boolean().optional(),
  // Server policy, not a client constant. This is the scheduling gate, not a
  // learner-facing number: the interface says when a check is ready.
  quiz_gate_xp: z.number().optional(),
})

// oxlint-disable-next-line anti-slop/no-unknown-parameters -- zod parse boundary: input is untrusted by design
function validateResponse<T>(schema: z.ZodType<T>, data: unknown, name: string): T {
  const result = schema.safeParse(data)
  if (!result.success) {
    console.error(`API validation error (${name}):`, result.error.issues)
    throw new Error(`API response validation failed for ${name}`)
  }
  return result.data as T
}

const GraphNodeSchema = z.object({
  id: z.string(),
  label: z.string(),
  domain: z.string(),
  prerequisites: z.array(z.string()),
  grading_type: z.string().optional(),
})

const GraphResSchema = z.object({
  nodes: z.array(GraphNodeSchema),
  count: z.number().optional(),
})

export const ConceptProgressSchema = z.object({
  concept_id: z.string().optional(),
  status: z.string(),
  streak: z.number(),
  best_streak: z.number().optional(),
  avg_response_time: z.number().optional(),
  attempts: z.number().optional(),
  completed: z.boolean().optional(),
  /**
   * Attainment, 0..1, derived server-side from the status tier and progress within it.
   *
   * The graph used to recompute this as `streak / mastery_threshold.streak`, from a
   * *bundled* copy of the corpus. That was wrong twice over: the streak resets on every
   * tier advance, so the bar fell from 100% to ~10% three times on the way up, and a
   * bundled threshold silently disagrees with the server's whenever the corpus is
   * rebuilt. Optional so a payload predating the field still parses.
   */
  mastery_pct: z.number().optional(),
})

const ProgressMapSchema = z.record(z.string(), ConceptProgressSchema)

const WeaknessEntrySchema = z.object({
  id: z.string(),
  label: z.string(),
  weakness: z.number().optional(),
})

const WeaknessResSchema = z.object({
  by_domain: z.record(z.string(), z.array(WeaknessEntrySchema)),
})

const LessonInfoSchema = z.object({
  title: z.string(),
  body: z.string().optional(),
  concepts: z.array(z.string()),
  progress: z.record(z.string(), ConceptProgressSchema).optional(),
  prerequisites: z.array(PrereqInfoSchema).optional(),
  dependents: z.array(PrereqInfoSchema).optional(),
})

const LessonsResSchema = z.object({
  lessons: z.record(z.string(), z.array(LessonInfoSchema)),
})

export interface StartSessionRes {
  student_id: string
  session_id: string
  question: Question | null
}

export interface Question {
  concept_id: string
  concept_name: string
  question: string
  is_review: boolean
  attempt_id?: string
  lesson?: {
    Title: string
    Body: string
    Concepts: string[]
  }
  diagram?: string
  grading_type?: string
  scaffold_active?: boolean
  scaffold_step?: number
  scaffold_total?: number
}

export interface AnswerRes {
  result: AnswerResult | null
  next_question: Question | null
  done: boolean
}

export interface AnswerResult {
  correct: boolean
  feedback: string
  new_status: string
  /** Authoritative mastery evidence score 0..1, from internal/mastery.BuildEvidence. */
  evidence_score: number
  /** The same score in plain language. Presentational; claims nothing about retention. */
  evidence_band: string
  explanation: string
  streak: number
  required_streak: number
  xp: number
  expected_answer: string
  ungraded?: boolean
  scaffold_active?: boolean
  scaffold_step?: number
  scaffold_total?: number
  /** One sentence naming the mistake, when it is certain. Empty otherwise. */
  diagnosis?: string
}

export interface GraphRes {
  nodes: {
    id: string
    label: string
    domain: string
    prerequisites: string[]
    grading_type: string
  }[]
  count: number
}

export type LeaderboardView = 'weekly' | 'all_time'

// No student id reaches the browser. `id` is a username (or a one-way
// derivation of the account id), which is what the avatar route and React keys
// use; a guest id would have been a credential on a public page.
export interface LeaderboardEntry {
  rank: number
  id: string
  name: string
  avatar?: LeaderboardAvatar
  xp: number
  xp_total: number
  mastered: number
  level: string
  is_me?: boolean
}

export interface LeaderboardAvatar {
  kind: 'custom' | 'generated'
  url?: string
  style?: string
  version?: number
}

export interface LeaderboardPosition {
  rank: number
  xp: number
  eligible: boolean
  guest?: boolean
  in_page: boolean
}

export interface LeaderboardBoard {
  view: LeaderboardView
  page: number
  page_size: number
  total: number
  total_pages: number
  entries: LeaderboardEntry[]
  me?: LeaderboardPosition
  week_start: string
  week_end: string
}

export interface Scores {
  lifetime_points: number
  weekly_score: number
  speed_bonus: number
  concepts_mastered: number
  current_streak: number
  level: string
  xp_total: number
  xp_today: number
  daily_xp_goal: number
  spaced_reps?: Record<string, number>
  avg_learning_speed?: number
  paused_until?: string
  // Quiz gate signal (backend); absent on old mocks.
  xp_since_quiz?: number
  quiz_due?: boolean
  quiz_active?: boolean
  quiz_gate_xp?: number
}

export interface ConceptProgress {
  concept_id?: string
  status: string
  streak: number
  best_streak?: number
  avg_response_time?: number
  attempts?: number
  // Set by the lesson-completion endpoint (PR1, server-owned). Optional so
  // pre-PR1 payloads keep validating; unlock rules treat absent as false.
  completed?: boolean
  placement_seeded?: boolean
  is_provisional?: boolean
}

// oxlint-disable-next-line anti-slop/no-unknown-parameters -- catch value: unknown is the correct error contract
export function isConflictError(err: unknown): boolean {
  return err instanceof Error && (err as { status?: number }).status === 409
}

// oxlint-disable-next-line anti-slop/no-unknown-parameters -- catch value: unknown is the correct error contract
export function isTooQuickError(err: unknown): boolean {
  return err instanceof Error && (err as { status?: number }).status === 400
}

// oxlint-disable-next-line anti-slop/no-unknown-parameters -- catch value: unknown is the correct error contract
export function isNotFoundError(err: unknown): boolean {
  return err instanceof Error && (err as { status?: number }).status === 404
}

// oxlint-disable-next-line anti-slop/no-unknown-parameters -- catch value: unknown is the correct error contract
export function getErrorStatus(err: unknown): number | undefined {
  if (err instanceof Error) {
    const s = (err as { status?: number }).status
    if (typeof s === 'number') return s
  }
  return undefined
}

// oxlint-disable-next-line anti-slop/no-unknown-parameters -- catch value: unknown is the correct error contract
export function getErrorMessage(err: unknown): string {
  if (err instanceof Error) {
    const e = err as { serverMessage?: string }
    if (typeof e.serverMessage === 'string' && e.serverMessage) return e.serverMessage
    return err.message
  }
  return String(err)
}

async function throwWithResponse(res: Response, fallbackMsg: string): Promise<never> {
  let serverMessage = ''
  try {
    const data = await res.json() as { error?: string }
    if (data && typeof data.error === 'string') serverMessage = data.error
  } catch { /* ignore json parse */ }
  const msg = serverMessage ? `${fallbackMsg}: ${serverMessage}` : fallbackMsg
  throw Object.assign(new Error(msg), { status: res.status, serverMessage })
}

export async function getProgress(studentID: string): Promise<Record<string, ConceptProgress>> {
  const headers = { ...getAuthHeaders() } satisfies Record<string, string>
  const res = await authedFetch(`${API_BASE}/api/progress/${studentID}`, { headers })
  if (!res.ok) return {}
  return validateResponse(ProgressMapSchema, await res.json(), 'getProgress') as Record<string, ConceptProgress>
}

export async function getScores(studentID: string): Promise<Scores> {
  const headers = { ...getAuthHeaders() } satisfies Record<string, string>
  const res = await authedFetch(`${API_BASE}/api/scores/${studentID}`, { headers })
  if (!res.ok) throw new Error(`Scores fetch failed: ${res.status}`)
  return validateResponse(ScoresSchema, await res.json(), 'getScores') as Scores
}

export async function getGraph(): Promise<GraphRes> {
  const res = await authedFetch(`${API_BASE}/api/graph`)
  if (!res.ok) throw new Error(`Graph fetch failed: ${res.status}`)
  return validateResponse(GraphResSchema, await res.json(), 'getGraph') as GraphRes
}

export interface DomainSummary {
  id: string
  conceptCount: number
  masteredCount: number
  unlockedCount: number
  inProgressCount: number
}

export interface TopicSummary {
  conceptId: string
  title: string
  state: string
  depth: number
}

export interface DomainDetail {
  domain: { id: string }
  topics: TopicSummary[]
}

/**
 * The curriculum browser's data. The server decides each domain's counts and
 * each topic's state; the client only renders them (ADR-042). Fails soft to an
 * empty list so a network blip degrades to the empty state, never to a crash.
 */
export async function getCurriculumDomains(): Promise<DomainSummary[]> {
  const res = await authedFetch(`${API_BASE}/api/curriculum/domains`, { cache: 'no-store' })
  if (!res.ok) return []
  const data = (await res.json()) as { domains?: DomainSummary[] }
  return data.domains ?? []
}

export async function getCurriculumDomain(domainId: string): Promise<DomainDetail | null> {
  const res = await authedFetch(`${API_BASE}/api/curriculum/domains/${encodeURIComponent(domainId)}`, { cache: 'no-store' })
  if (!res.ok) return null
  return res.json() as Promise<DomainDetail>
}

export async function getLeaderboard(
  view: LeaderboardView,
  page: number,
  pageSize = 25,
): Promise<LeaderboardBoard> {
  const params = new URLSearchParams({ view, page: String(page), page_size: String(pageSize) })
  const res = await authedFetch(`${API_BASE}/api/leaderboard?${params.toString()}`, { cache: 'no-store' })
  if (!res.ok) await throwWithResponse(res, 'Could not load the leaderboard')
  const data = (await res.json()) as Partial<LeaderboardBoard>
  // A deploy in progress can pair this bundle with the previous API, whose
  // response was a bare array. Render the error state rather than crash on
  // `entries.length` underneath the error boundary.
  if (!data || !Array.isArray(data.entries)) {
    throw new Error('The leaderboard response was not a page')
  }
  return data as LeaderboardBoard
}

export interface LeagueMemberLine {
  rank: number
  id: string
  name: string
  avatar?: LeaderboardAvatar
  xp: number
  is_me?: boolean
}

export interface MyLeague {
  week: string
  week_start: string
  week_end: string
  tier: string
  league_no: number
  size: number
  promote: number
  demote: number
  members: LeagueMemberLine[]
  me?: { rank: number; xp: number }
}

// A signed-in learner's own league. Unauthenticated callers get null rather
// than an error: "sign in to compete" is an answer, not a failure.
export async function getMyLeague(): Promise<MyLeague | null> {
  const res = await authedFetch(`${API_BASE}/api/leagues`, { cache: 'no-store' })
  if (res.status === 401) return null
  if (!res.ok) await throwWithResponse(res, 'Could not load your league')
  return res.json() as Promise<MyLeague>
}

export interface ShareReport {
  student_id: string
  name: string
  scores: Scores
  activity: DailyActivity[]
  progress: Record<string, ConceptProgress>
  weakness: Record<string, number>
  attempts?: AttemptRecord[]
}

export interface AttemptRecord {
  session_id: string
  student_id: string
  concept_id: string
  concept_name?: string
  answer: string
  expected: string
  correct: boolean
  elapsed_seconds: number
  timestamp: string
  question: string
  source: string
  explanation: string
}

export interface AttemptsRes {
  attempts: AttemptRecord[]
  total: number
}

export async function getAttempts(params?: {
  source?: string
  concept_id?: string
  incorrect_only?: boolean
  limit?: number
  offset?: number
}): Promise<AttemptsRes> {
  const q = new URLSearchParams()
  if (params?.source) q.set('source', params.source)
  if (params?.concept_id) q.set('concept_id', params.concept_id)
  if (params?.incorrect_only) q.set('incorrect_only', '1')
  if (params?.limit != null) q.set('limit', String(params.limit))
  if (params?.offset) q.set('offset', String(params.offset))
  const qs = q.toString()
  const res = await authedFetch(`${API_BASE}/api/attempts${qs ? `?${qs}` : ''}`, {
    headers: { ...getAuthHeaders() },
    cache: 'no-store',
  })
  if (!res.ok) throw new Error(`Attempts fetch failed: ${res.status}`)
  return res.json()
}

export async function enableShare(): Promise<{ enabled: boolean; token: string }> {
  const res = await authedFetch(`${API_BASE}/api/share`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ enabled: true }),
  })
  if (!res.ok) throw new Error(`Enable share failed: ${res.status}`)
  return res.json()
}

export async function disableShare(): Promise<void> {
  await authedFetch(`${API_BASE}/api/share`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ enabled: false }),
  })
}

export async function getShareReport(token: string): Promise<ShareReport | null> {
  const res = await authedFetch(`${API_BASE}/api/share/${encodeURIComponent(token)}`, { cache: 'no-store' })
  if (!res.ok) return null
  return res.json()
}

export interface CourseProgress {
  total: number
  mastered: number
  pct: number
  days_remaining: number
}

export interface CourseStatus {
  id: string
  name: string
  grade: string
  description: string
  targets: string[]
  progress?: CourseProgress
}

export async function getTranscript(): Promise<CourseStatus[]> {
  const res = await authedFetch(`${API_BASE}/api/transcript`, { headers: { ...getAuthHeaders() }, cache: 'no-store' })
  if (!res.ok) return []
  const data = await res.json()
  return data.courses ?? []
}

export interface EfficacyReport {
  concepts_touched: number
  first_pass_rate: number
  second_pass_rate: number
  avg_attempts_per_concept: number
  total_attempts: number
}

export async function getEfficacy(): Promise<EfficacyReport | null> {
  const res = await authedFetch(`${API_BASE}/api/efficacy`, { headers: { ...getAuthHeaders() } })
  if (!res.ok) return null
  return res.json()
}

export interface EfficacyWeek {
  week_start: string
  attempts: number
  concepts_touched: number
  first_pass_rate: number
  second_pass_rate: number
  active_students: number
}

export interface EfficacyTrend {
  weeks: EfficacyWeek[]
  total_students: number
  returning_students: number
  retention_rate: number
  first_pass_trend: number
}

// Product-wide longitudinal efficacy (weekly first-pass rate + retention).
export async function getEfficacyTrend(): Promise<EfficacyTrend | null> {
  const res = await authedFetch(`${API_BASE}/api/efficacy/trend`, { headers: { ...getAuthHeaders() } })
  if (!res.ok) return null
  return res.json()
}

export interface ConfigRes {
	auth_enabled: boolean
	google_client_id?: string
	providers?: string[]
	/**
	 * Whether optional AI augmentation is configured on the server. A boolean
	 * only: the provider chain, model ids and per-identity quota are operational
	 * detail the client has no use for. When this is false the AI calls below
	 * are not attempted and the deterministic paths stand alone.
	 */
	ai_enabled?: boolean
	/**
	 * Whether the friends surface is exposed. The storage layer ships dark, so
	 * a deployment without MATHUA_FRIENDS=on answers 404 on every friends route
	 * and the navigation must not offer a page that cannot load.
	 */
	friends_enabled?: boolean
}

export async function getConfig(): Promise<ConfigRes | null> {
	const res = await authedFetch(`${API_BASE}/api/config`)
	if (!res.ok) return null
	return res.json()
}

export async function googleOneTap(idToken: string): Promise<AuthRes> {
	const res = await authedFetch(`${API_BASE}/api/auth/google`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ id_token: idToken }) })
	if (!res.ok) throw new Error('Google sign-in failed')
	return validateResponse(AuthResSchema, await res.json(), 'googleOneTap') as AuthRes
}

export async function healthCheck(): Promise<boolean> {
  try {
    const res = await authedFetch(`${API_BASE}/api/health`)
    return res.ok
  } catch {
    return false
  }
}

export interface AuthRes {
  token: string
  student_id: string
  name: string
  diagnostic_completed: boolean
}

export async function signup(name: string, username: string, password: string, email?: string, guestId?: string): Promise<AuthRes> {
	const res = await authedFetch(`${API_BASE}/api/auth/signup`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ name, username, password, email: email ?? '', student_id: guestId ?? '' }),
	})
	if (!res.ok) await throwWithResponse(res, `Signup failed`)
	return validateResponse(AuthResSchema, await res.json(), 'signup') as AuthRes
}

export async function login(username: string, password: string): Promise<AuthRes> {
	const res = await authedFetch(`${API_BASE}/api/auth/login`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ username, password }),
	})
	if (res.status === 429) throw new Error('Too many attempts, try again shortly')
	if (!res.ok) throw new Error('Invalid credentials')
	return validateResponse(AuthResSchema, await res.json(), 'login') as AuthRes
}

// Goal-based diagnostic API

export interface GoalPathRes {
	concepts: { id: string; label: string; domain: string; prerequisites: string[] }[]
	count: number
}

export interface DiagnosticProgress {
  answered: number
  estimated_total: number
  min_total: number
  max_total: number
  probe_budget: number
  done: boolean
}

export interface GoalDiagStartRes {
  session_id: string
  student_id?: string
  concept_id?: string
  concept_name?: string
  question?: string
  grading_type?: string
  answer_hint?: string
  answer_input_mode?: 'numeric' | 'text'
  done?: boolean
  progress?: DiagnosticProgress
}

export interface GoalDiagAnswerRes {
  done: boolean
  correct?: boolean
  feedback?: string
  concept_id?: string
  concept_name?: string
  question?: string
  grading_type?: string
  answer_hint?: string
  answer_input_mode?: 'numeric' | 'text'
  retry_available?: boolean
  progress?: DiagnosticProgress
}

export interface GoalPlanRes {
  readiness: number
  total_tested: number
  correct_count: number
  weak_areas: Record<string, { id: string; label: string }[]>
  strong_areas: Record<string, string[]>
  frontier_label?: string
  frontier_idx?: number
  frontier_conditional?: boolean
  conditionally_completed?: string[]
  placement_course_id?: string
  completion_estimates?: Record<string, string>
}

export interface WeaknessRes {
	by_domain: Record<string, { id: string; label: string; weakness: number }[]>
}

export async function getGoalPath(conceptIds: string[]): Promise<GoalPathRes> {
	const res = await authedFetch(`${API_BASE}/api/goal`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify({ concept_ids: conceptIds }),
	})
	if (!res.ok) throw new Error(`Goal path fetch failed: ${res.status}`)
	return res.json()
}

export async function startGoalDiagnostic(conceptIds: string[]): Promise<GoalDiagStartRes> {
	const res = await authedFetch(`${API_BASE}/api/goal/diagnostic`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify({ concept_ids: conceptIds }),
	})
	if (!res.ok) await throwWithResponse(res, `Goal diagnostic start failed: ${res.status}`)
	return res.json()
}

export async function startGoalDiagnosticName(name: string, conceptIds: string[]): Promise<GoalDiagStartRes> {
	const res = await authedFetch(`${API_BASE}/api/goal/diagnostic`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ name, concept_ids: conceptIds }),
	})
	if (!res.ok) await throwWithResponse(res, `Goal diagnostic start failed: ${res.status}`)
	return res.json()
}

export async function submitGoalAnswer(
	sessionId: string,
	conceptId: string,
	answer: string,
	elapsed: number,
	dontKnow?: boolean,
): Promise<GoalDiagAnswerRes> {
	const res = await authedFetch(`${API_BASE}/api/goal/diagnostic/answer`, {
		// The server records each verdict in diagnostic_answers, keyed by session
		// and concept, and replays a duplicate instead of grading twice. That
		// makes one automatic retry safe in both cases: if the first request
		// committed, the retry receives the stored verdict; if it never arrived,
		// the retry is graded fresh. Without it, a single dropped packet on
		// mobile data ended the whole run.
		//
		// Twelve seconds, not twenty-five: with the retry the worst case is half
		// a minute, and learners reported the screen sitting there for "a good
		// minute" before it gave up. A tutor endpoint that has not answered in
		// twelve seconds is better retried than waited on.
		timeoutMs: 12000,
		retryAttempts: 1,
		method: 'POST',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify({ session_id: sessionId, concept_id: conceptId, answer, elapsed, dont_know: dontKnow ?? false }),
	})
	if (!res.ok) await throwWithResponse(res, `Goal answer failed: ${res.status}`)
	return res.json()
}

export async function resumeGoalDiagnostic(sessionId: string): Promise<GoalDiagAnswerRes & { session_id: string }> {
	const res = await authedFetch(`${API_BASE}/api/goal/diagnostic/resume?session_id=${encodeURIComponent(sessionId)}`, {
		headers: { ...getAuthHeaders() },
	})
	if (!res.ok) await throwWithResponse(res, `Goal resume failed: ${res.status}`)
	return res.json()
}

// Skip the pending diagnostic question: settles the concept with no evidence
// recorded and serves the next question. Escape hatch for unanswerable items.
export async function skipGoalQuestion(sessionId: string): Promise<GoalDiagAnswerRes> {
	const res = await authedFetch(`${API_BASE}/api/goal/diagnostic/skip`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify({ session_id: sessionId }),
	})
	if (!res.ok) await throwWithResponse(res, `Goal skip failed: ${res.status}`)
	return res.json()
}

// Void the superseded question's miss and restore it ("silly mistake"
// retry). The concept_id names the question being retried and must match
// the voidable attempt, a stale client that already moved on gets a 400
// instead of voiding the wrong answer.
export async function retryGoalQuestion(sessionId: string, conceptId: string): Promise<GoalDiagAnswerRes> {
	const res = await authedFetch(`${API_BASE}/api/goal/diagnostic/retry`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify({ session_id: sessionId, concept_id: conceptId }),
	})
	if (!res.ok) await throwWithResponse(res, `Goal retry failed: ${res.status}`)
	return res.json()
}

export async function getGoalPlan(sessionId: string): Promise<GoalPlanRes> {
	const res = await authedFetch(`${API_BASE}/api/goal/plan`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify({ session_id: sessionId }),
	})
	if (!res.ok) await throwWithResponse(res, `Goal plan failed: ${res.status}`)
	return res.json()
}

export async function setDailyXPGoal(goal: number): Promise<void> {
	const res = await authedFetch(`${API_BASE}/api/goals/xp`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify({ goal }),
	})
	if (!res.ok) throw new Error(`Set goal failed: ${res.status}`)
}

// ── Destinations & study plans ──

export interface DestinationStatus {
	id: string
	name: string
	description: string
	total: number
	mastered: number
	pct: number
}

export async function getDestinations(): Promise<DestinationStatus[]> {
	const res = await authedFetch(`${API_BASE}/api/destinations`, { headers: { ...getAuthHeaders() } })
	if (!res.ok) throw new Error(`Destinations fetch failed: ${res.status}`)
	const data = await res.json()
	return data.destinations ?? []
}

export interface WorkloadEstimate {
	total: number
	mastered: number
	remaining: string[]
	xp_remaining: number
	time_min_remaining: number
	reviews_due: number
	quizzes_ahead: number
	assessment_min: number
	diagnostic_min: number
	days: number
	finish_date: string
	required_per_day: number
	feasible: boolean
	lines: { learning: number; assessment: number; diagnostic: number }
}

export interface EstimateRes {
	destination: string
	estimate: WorkloadEstimate
	pace: { rate: number; source: 'measured' | 'goal'; adherence: number; active_days: number; trailing_days: number }
	probes: string[]
	plan_delta_days?: number
	fresh_start?: boolean
}

export interface EstimateParams {
	daily_goal?: number
	deadline_days?: number
	rest_days?: number
	diagnostic_min?: number
}

export async function getEstimate(destId: string, params: EstimateParams = {}): Promise<EstimateRes> {
	const q = new URLSearchParams()
	if (params.daily_goal) q.set('daily_goal', String(params.daily_goal))
	if (params.deadline_days) q.set('deadline_days', String(params.deadline_days))
	if (params.rest_days !== undefined) q.set('rest_days', String(params.rest_days))
	if (params.diagnostic_min) q.set('diagnostic_min', String(params.diagnostic_min))
	const res = await authedFetch(`${API_BASE}/api/destinations/${encodeURIComponent(destId)}/estimate?${q.toString()}`, {
		headers: { ...getAuthHeaders() },
	})
	if (!res.ok) throw new Error(`Estimate fetch failed: ${res.status}`)
	return res.json()
}

export interface StudyPlan {
	destination: string
	daily_goal: number
	deadline_days: number
	rest_days: number
	xp_remaining: number
	created: string
}

export async function savePlan(body: { destination: string; daily_goal: number; deadline_days: number; rest_days: number }): Promise<StudyPlan> {
	const res = await authedFetch(`${API_BASE}/api/plans`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify(body),
	})
	if (!res.ok) throw new Error(`Save plan failed: ${res.status}`)
	return res.json()
}

export async function getCurrentPlan(): Promise<{ plan: StudyPlan | null; delta_days?: number; xp_remaining?: number }> {
	const res = await authedFetch(`${API_BASE}/api/plans/current`, { headers: { ...getAuthHeaders() } })
	if (!res.ok) return { plan: null }
	return res.json()
}

export async function resetAccount(phrase: string): Promise<void> {
	const res = await authedFetch(`${API_BASE}/api/account/reset`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify({ phrase }),
	})
	// Read the {error} envelope rather than res.text(): throwing the raw body
	// put a literal `{"error":"…"}` in the danger zone, which reads as a
	// broken app instead of an explanation.
	if (!res.ok) await throwWithResponse(res, 'Reset failed')
}

// The server's own account payload. `role` was missing here and zod strips unknown keys, so the
// field was being silently discarded on the way in: `/api/auth/me` sent it, `getMe()` dropped
// it, and the Admin navigation entry could never appear. A schema is a filter, so a field the
// type does not name is a field the application does not have, which is the second half of
// ADR-043's duplicate-schema lesson, arriving from the other direction.
const MeSchema = z.object({
	student_id: z.string(),
	name: z.string(),
	has_password: z.boolean().optional(),
	role: z.enum(['student', 'moderator', 'admin', 'owner']).optional(),
	/**
	 * The server's own record that placement happened. The client must read it
	 * here rather than trusting `mathua_user` or inferring it from having shown
	 * the last diagnostic question: local storage can be stale on another device,
	 * and "the final question was displayed" is not a completion.
	 */
	diagnostic_completed: z.boolean().optional(),
	diagnostic_completed_at: z.string().optional(),
})

export type MeInfo = z.infer<typeof MeSchema>

export async function getMe(): Promise<MeInfo> {
	const res = await authedFetch(`${API_BASE}/api/auth/me`, { headers: { ...getAuthHeaders() } })
	if (!res.ok) throw new Error(`Me fetch failed: ${res.status}`)
	return validateResponse(MeSchema, await res.json(), 'getMe')
}

export async function deleteAccount(body: { phrase: string; password?: string }): Promise<void> {
	const res = await authedFetch(`${API_BASE}/api/account`, {
		method: 'DELETE',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify(body),
	})
	// Same reason as resetAccount: never surface the raw JSON envelope.
	if (!res.ok) await throwWithResponse(res, 'Delete failed')
}

export async function getWeaknesses(): Promise<WeaknessRes> {
  const res = await authedFetch(`${API_BASE}/api/weaknesses`, {
    headers: { ...getAuthHeaders() },
  })
  if (!res.ok) return { by_domain: {} }
  return validateResponse(WeaknessResSchema, await res.json(), 'getWeaknesses') as WeaknessRes
}

export interface UserSettings {
	show_timer?: boolean
	pause_until?: string | null
	avatar_preset?: number | null
	avatar_dicebear?: { style: string; seed: string } | null
	avatar_custom?: boolean
	avatar_version?: number | null
	// How the leaderboard and league rows name this learner. The server reads
	// the same key from the opaque settings blob; a missing key means "named",
	// which is the behaviour that existed before the setting.
	leaderboard_privacy?: 'named' | 'anonymous' | 'hidden'
	// Whether accepted friends can see this learner's name and shared graph.
	// A missing key means "friends", the RFC default; `nobody` hides the
	// learner from lists and revokes profile and graph access immediately.
	friends_visibility?: 'friends' | 'nobody'
}

export async function getSettings(): Promise<UserSettings> {
	const res = await authedFetch(`${API_BASE}/api/settings`, {
		headers: { ...getAuthHeaders() },
	})
	if (!res.ok) return {}
	return res.json()
}

export async function updateSettings(settings: UserSettings): Promise<void> {
	const res = await authedFetch(`${API_BASE}/api/settings`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify(settings),
	})
	if (!res.ok) throw new Error(`Update settings failed: ${res.status}`)
}

export interface ProfileUpdateRes {
	student_id: string
	name: string
	email?: string
}

export async function updateProfileName(name: string, email?: string): Promise<ProfileUpdateRes> {
	const res = await authedFetch(`${API_BASE}/api/profile`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify(email === undefined ? { name } : { name, email }),
	})
	if (!res.ok) await throwWithResponse(res, `Update profile failed: ${res.status}`)
	return res.json()
}

export async function requestPasswordReset(identifier: string): Promise<void> {
	const res = await authedFetch(`${API_BASE}/api/auth/reset/request`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ identifier }),
	})
	if (!res.ok) await throwWithResponse(res, `Reset request failed: ${res.status}`)
}

export async function completePasswordReset(token: string, password: string): Promise<AuthRes> {
	const res = await authedFetch(`${API_BASE}/api/auth/reset/complete`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ token, password }),
	})
	if (!res.ok) await throwWithResponse(res, `Reset failed: ${res.status}`)
	return validateResponse(AuthResSchema, await res.json(), 'resetComplete') as AuthRes
}

export async function changePassword(currentPassword: string, newPassword: string): Promise<void> {
	const res = await authedFetch(`${API_BASE}/api/auth/password`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify({ current_password: currentPassword, new_password: newPassword }),
	})
	if (!res.ok) await throwWithResponse(res, `Change password failed: ${res.status}`)
}

export const OAUTH_PROVIDERS = ['google', 'github', 'facebook', 'microsoft', 'apple'] as const
export type OAuthProvider = (typeof OAUTH_PROVIDERS)[number]

export const OAUTH_LABELS: Record<OAuthProvider, string> = {
	google: 'Google',
	github: 'GitHub',
	facebook: 'Facebook',
	microsoft: 'Microsoft',
	apple: 'Apple',
}

// startOAuthLogin begins a redirect OAuth flow. intent=link binds the dance
// to the current account (needs a link token); otherwise it's a login.
export function startOAuthLogin(provider: string, opts: { intent?: 'login' | 'link'; linkToken?: string; ret?: string } = {}) {
	const q = new URLSearchParams()
	if (opts.intent === 'link') q.set('intent', 'link')
	if (opts.linkToken) q.set('link_token', opts.linkToken)
	q.set('return', opts.ret || '/profile')
	window.location.href = `${API_BASE}/api/auth/${provider}/login?${q.toString()}`
}

export interface IdentityInfo {
	provider: string
	email: string
}

export async function getIdentities(): Promise<IdentityInfo[]> {
	const res = await authedFetch(`${API_BASE}/api/auth/identities`, {
		headers: { ...getAuthHeaders() },
	})
	if (!res.ok) return []
	return res.json()
}

export async function deleteIdentity(provider: string): Promise<void> {
	const res = await authedFetch(`${API_BASE}/api/auth/identities/${encodeURIComponent(provider)}`, {
		method: 'DELETE',
		headers: { ...getAuthHeaders() },
	})
	if (!res.ok) await throwWithResponse(res, `Disconnect failed: ${res.status}`)
}

export async function createLinkToken(): Promise<string> {
	const res = await authedFetch(`${API_BASE}/api/auth/link-token`, {
		method: 'POST',
		headers: { ...getAuthHeaders() },
	})
	if (!res.ok) await throwWithResponse(res, `Link failed: ${res.status}`)
	const data = await res.json()
	return data.link_token as string
}

export async function requestEmailVerification(): Promise<void> {
	const res = await authedFetch(`${API_BASE}/api/auth/email/request`, {
		method: 'POST',
		headers: { ...getAuthHeaders() },
	})
	if (!res.ok) await throwWithResponse(res, `Verification request failed: ${res.status}`)
}

export async function verifyEmail(token: string): Promise<void> {
	const res = await authedFetch(`${API_BASE}/api/auth/email/verify`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ token }),
	})
	if (!res.ok) await throwWithResponse(res, `Verification failed: ${res.status}`)
}

export async function uploadAvatar(file: File): Promise<void> {
	const form = new FormData()
	form.append('avatar', file)
	const res = await authedFetch(`${API_BASE}/api/avatar`, {
		method: 'POST',
		headers: { ...getAuthHeaders() },
		body: form,
	})
	if (!res.ok) throw new Error(`Avatar upload failed: ${res.status}`)
}

export async function deleteAvatar(): Promise<void> {
	const res = await authedFetch(`${API_BASE}/api/avatar`, {
		method: 'DELETE',
		headers: { ...getAuthHeaders() },
	})
	if (!res.ok) throw new Error(`Avatar delete failed: ${res.status}`)
}

export function avatarImageUrl(v?: string | number): string {
	const suffix = v !== undefined ? `?v=${encodeURIComponent(String(v))}` : ''
	return `${API_BASE}/api/avatar/me${suffix}`
}

export interface PrereqInfo {
    id: string
    label: string
    status: string
    mastery_pct: number
}

export interface LessonInfo {
    title: string
    body?: string
    concepts: string[]
    progress?: Record<string, ConceptProgress>
    prerequisites?: PrereqInfo[]
    dependents?: PrereqInfo[]
}

// Lesson bodies ship separately (the list payload is metadata-only) and are
// fetched on selection. No client-side memoization: a stale body surviving a
// backend restart is worse than an extra request. no-store also bypasses the
// Next.js data cache in dev.
// One figure, table or other instructional block referenced by a lesson.
//
// `src` is the reference exactly as the lesson body writes it; resolution to a served
// URL happens in KatexContent, which is the only place that knows the serving base.
// `line` is the 1-based line in the body, so a renderer can interleave figures with the
// prose rather than appending them.
export interface LessonAsset {
	id: string
	kind: 'image' | 'table'
	src?: string
	alt?: string
	caption?: string
	line: number
	section?: string
	concept_ids?: string[]
	source?: string
}

// A lesson plus its figure inventory. `/study` renders the body verbatim and needs
// nothing from `assets`; the list exists so a figure can be addressed, counted and
// repaired rather than being trapped inside the markdown.
export interface LessonDocument {
	title: string
	body: string
	source?: string
	assets: LessonAsset[]
}

// Lesson bodies ship separately (the list payload is metadata-only) and are
// fetched on selection. No client-side memoization: a stale body surviving a
// backend restart is worse than an extra request. no-store also bypasses the
// Next.js data cache in dev.
export async function getLessonDocument(title: string): Promise<LessonDocument> {
	const res = await authedFetch(`${API_BASE}/api/lessons/body?title=${encodeURIComponent(title)}`, {
		cache: 'no-store',
	})
	if (!res.ok) throw new Error(`Lesson body fetch failed: ${res.status}`)
	const data = (await res.json()) as Partial<LessonDocument>
	// The server sends [] rather than null, but a proxy or an older build could omit the
	// field, and every caller iterates it.
	return {
		title: data.title ?? title,
		body: data.body ?? '',
		source: data.source,
		assets: data.assets ?? [],
	}
}

export async function getLessonBody(title: string): Promise<string> {
	return (await getLessonDocument(title)).body
}

export interface LessonsRes {
	lessons: Record<string, LessonInfo[]>
}

export async function getLessons(studentId?: string): Promise<LessonsRes> {
	const url = studentId ? `${API_BASE}/api/lessons?student_id=${encodeURIComponent(studentId)}` : `${API_BASE}/api/lessons`
	const res = await authedFetch(url, { cache: 'no-store' })
	if (!res.ok) return { lessons: {} }
	return validateResponse(LessonsResSchema, await res.json(), 'getLessons') as LessonsRes
}

export interface PracticeQuestion {
	question: string
	// The answer and explanation are deliberately absent. They are returned by
	// POST /api/study/answer after the learner submits, so the served question
	// cannot be inspected for its verdict before evidence is produced.
	// Declared by the generator that produced the question, because the
	// concept's grading type describes the grader and not the answer syntax.
	answer_hint?: string
	answer_input_mode?: 'numeric' | 'text'
}

export interface LessonPracticeRes {
	questions: PracticeQuestion[]
	concept_id: string
}

export interface KpInfo {
	label: string
	section?: string
	subgoals: string[]
	worked_example: string
	// Figures belonging to this learning step, selected from the lesson by section (or by
	// an authored asset_ids list on the shard). `/learn` shows these; `/study` shows the
	// whole lesson. Optional so a response from a build predating the field still renders.
	assets?: LessonAsset[]
}

export interface LessonKpsRes {
	concept_id: string
	kps: KpInfo[]
	diagram?: string
}

export async function getLessonPractice(conceptId: string, count = 5, opts?: { seed?: number; exclude?: string[]; difficulty?: number; kp_ordinal?: number }): Promise<LessonPracticeRes> {
	const params = new URLSearchParams({ count: String(count) })
	if (opts?.seed) params.set('seed', String(opts.seed))
	if (opts?.difficulty !== undefined) params.set('difficulty', String(opts.difficulty))
	if (opts?.kp_ordinal) params.set('kp_ordinal', String(opts.kp_ordinal))
	for (const ex of opts?.exclude ?? []) params.append('exclude', ex)
	const res = await authedFetch(`${API_BASE}/api/lessons/${encodeURIComponent(conceptId)}/practice?${params.toString()}`)
	if (!res.ok) throw new Error(`Lesson practice fetch failed: ${res.status}`)
	return res.json()
}

// The knowledge-point shard, shared between concurrent callers for the length of the request.
//
// `/learn` needs it twice on one page: LearnStepper renders the worked example for the current
// section, and the reference panel renders every block for the same concept. Two callers, one
// payload, and the shard is authored content that does not change while the tab is open, so
// this is the in-flight cache rather than the session cache, because a later visit should
// re-read it (a shard can be corrected, and a stale copy of a teaching section is a wrong
// lesson rather than a stale progress bar).
const lessonKPsInFlight = new Map<string, Promise<LessonKpsRes>>()

export function getLessonKPs(conceptId: string): Promise<LessonKpsRes> {
	const existing = lessonKPsInFlight.get(conceptId)
	if (existing) return existing
	const p = (async () => {
		const res = await authedFetch(`${API_BASE}/api/lessons/${encodeURIComponent(conceptId)}/kp`, { cache: 'no-store' })
		if (!res.ok) return { concept_id: conceptId, kps: [] }
		return res.json() as Promise<LessonKpsRes>
	})().finally(() => {
		lessonKPsInFlight.delete(conceptId)
	})
	lessonKPsInFlight.set(conceptId, p)
	return p
}

/** What kind of work a recommendation is. Mirrors scheduler.Kind in Go. */
export type RecommendationKind = 'learn' | 'practice' | 'review' | 'mastery_check'

/** Why the task was chosen, in language a learner can act on. Mirrors scheduler.Reason. */
export type RecommendationReason = 'new' | 'review_due' | 'needs_practice' | 'prerequisite' | 'mastery_ready'

/**
 * One worthwhile task, decided by the engine.
 *
 * The wording is the server's, and that is the point: `/profile` and `/learn` used to format
 * one head two ways, "Currently working towards X / Answer it" beside "Next up / Continue: X",
 * sometimes with different labels for the same concept, because each derived its own copy.
 */
export interface Recommendation {
  id: string
  conceptId: string
  conceptTitle: string
  kind: RecommendationKind
  reason: RecommendationReason
  priority: number
  estimatedMinutes?: number
  xp?: number
  action: { type: RecommendationKind; href: string }
  badge: string
  detail: string
  cta: string
}

/**
 * The engine's answer: one dominant task plus the alternatives the learner may pick.
 *
 * "The engine chooses; the learner disposes", so `alternatives` is part of the contract rather
 * than an afterthought, and a learner who disagrees with the head can act on it.
 */
export interface RecommendationRes {
  primary: Recommendation | null
  alternatives: Recommendation[]
  generatedAt: string
}

/**
 * GET /api/next, the single answer to "what should I do now?".
 *
 * A successful response with no head degrades to the same orientation card the
 * server uses, so an older build or a proxy that drops the field cannot turn
 * this page into a render error.
 *
 * A *failed* request does not. It used to return that card too, and the comment
 * here used to defend it as graceful degradation. It is not: with the API
 * blocked, the learn page told the learner "All learned, nothing is due right
 * now" while their queue was untouched and simply unreachable. A learning app
 * that reports an outage as an empty queue is lying about the one fact it
 * exists to state, and the honest failure state cannot be reached at all
 * because the error never leaves this function.
 */
export async function getNext(exclude?: string[]): Promise<RecommendationRes> {
  const fallback: RecommendationRes = {
    primary: {
      id: 'graph',
      conceptId: '',
      conceptTitle: 'Concept graph',
      kind: 'learn',
      reason: 'new',
      priority: 0,
      action: { type: 'learn', href: '/graph' },
      badge: 'All learned',
      detail: 'Nothing is due right now, the graph shows what builds on what you have',
      cta: 'See the graph →',
    },
    alternatives: [],
    generatedAt: new Date().toISOString(),
  }
  const params = new URLSearchParams()
  if (exclude?.length) params.set('exclude', exclude.join(','))
  const qs = params.toString()
  const res = await authedFetch(`${API_BASE}/api/next${qs ? `?${qs}` : ''}`, { cache: 'no-store', timeoutMs: 10000 })
  if (!res.ok) await throwWithResponse(res, 'Could not load what to study next')
  const data = (await res.json()) as Partial<RecommendationRes>
  if (!data.primary) return fallback
  return { primary: data.primary, alternatives: data.alternatives ?? [], generatedAt: data.generatedAt ?? '' }
}

export interface EligibilityPrereq {
	conceptId: string
	title: string
	state: string
	met: boolean
}

export interface ReadinessRes {
	concept_id: string
	ready: boolean
	weak: { id: string; label: string; status: string; mastery_pct: number }[]
	missing: { id: string; label: string; status: string; mastery_pct: number }[]
	/**
	 * The server's eligibility verdict, present on current servers and absent on
	 * older payloads. The client renders it; it never computes it. `eligible`
	 * false means the topic's prerequisites are unmet and /learn will refuse it.
	 */
	eligible?: boolean
	state?: string
	reason?: string
	prerequisites?: EligibilityPrereq[]
}

export async function getLessonReadiness(conceptId: string): Promise<ReadinessRes> {
	const res = await authedFetch(`${API_BASE}/api/lessons/${encodeURIComponent(conceptId)}/readiness`, { cache: 'no-store' })
	// Fail soft, in the direction that matters: a blip must not block a learner
	// whose prerequisites are satisfied. The practice endpoint still refuses a
	// genuinely locked topic, so the server remains the authority.
	if (!res.ok) return { concept_id: conceptId, ready: true, weak: [], missing: [], eligible: true }
	return res.json()
}

/** The canonical eligibility read. Thin alias so call sites name the intent. */
export async function getLessonEligibility(conceptId: string): Promise<ReadinessRes> {
	return getLessonReadiness(conceptId)
}

export interface StudyAnswerRes {
	correct: boolean
	feedback: string
	explanation?: string
	new_status?: string
	streak?: number
	required_streak?: number
	xp: number
	expected_answer?: string
	ungraded?: boolean
	scaffold_active?: boolean
	scaffold_step?: number
	scaffold_total?: number
	/** One sentence naming the mistake, when it is certain. Empty otherwise. */
	diagnosis?: string
	/**
	 * Authoritative mastery evidence score for this concept, 0..1, from the same
	 * `BuildEvidence` the ladder decides on.
	 *
	 * Absent on an ungraded answer: nothing was recorded, so there is no new evidence.
	 */
	evidence_score?: number
	/**
	 * `evidence_score` in plain language, from `mastery.EvidenceBand`. Presentational only,
	 * and deliberately claiming nothing about retention, the evidence model has no spacing
	 * term, so it cannot support such a claim.
	 *
	 * Empty before the evidence floor is met, which is what stops a single correct answer
	 * from reading as a strong result.
	 */
	evidence_band?: string
}

interface StudyAnswerBody {
	concept_id: string
	answer: string
	elapsed: number
	student_id?: string
	question?: string
	submission_id: string
}

function newSubmissionId(): string {
	const c = globalThis.crypto
	if (c && typeof c.randomUUID === 'function') return c.randomUUID()
	return `study-${Date.now()}-${Math.random().toString(36).slice(2)}`
}

// No `expected` is sent: the server anchors every served question against its
// own answer, so a client-supplied one is either redundant or an attempt to
// grade itself. A 409 means the server has no record of this question (the
// anchor's TTL expired, or the tab was left open too long) and the caller
// should re-serve rather than resubmit.
export async function submitStudyAnswer(conceptId: string, answer: string, elapsed: number, question?: string): Promise<StudyAnswerRes> {
	const { getGuestId } = await import('./auth')
	const headers = { 'Content-Type': 'application/json', ...getAuthHeaders() } satisfies Record<string, string>
	const body: StudyAnswerBody = {
		concept_id: conceptId,
		answer,
		elapsed,
		question,
		submission_id: newSubmissionId(),
	}
	const guestId = getGuestId()
	if (guestId && !('Authorization' in headers)) {
		body.student_id = guestId
	}
	const res = await authedFetch(`${API_BASE}/api/study/answer`, {
		// The same submission_id goes out on every retry, so the server
		// replays the stored verdict instead of writing a second attempt.
		retryAttempts: 1,
		timeoutMs: 25000,
		method: 'POST',
		headers,
		body: JSON.stringify(body),
	})
	if (!res.ok) await throwWithResponse(res, `Study answer failed: ${res.status}`)
	return res.json() as Promise<StudyAnswerRes>
}

export interface DueReviewsRes {
	count: number
}

export interface DailyActivity {
	date: string
	questions: number
	correct: number
	concepts: string[]
}

export async function getActivity(days: number = 365): Promise<DailyActivity[]> {
	const headers = { ...getAuthHeaders() } satisfies Record<string, string>
	const res = await authedFetch(`${API_BASE}/api/activity?days=${days}`, { headers })
	if (!res.ok) return []
	return res.json()
}

export async function getDueReviews(): Promise<DueReviewsRes> {
	const headers = { ...getAuthHeaders() } satisfies Record<string, string>
	const res = await authedFetch(`${API_BASE}/api/reviews/due`, { headers })
	if (!res.ok) return { count: 0 }
	return res.json()
}

export async function startReviewSession(): Promise<StartSessionRes> {
	const headers = {
		'Content-Type': 'application/json',
		...getAuthHeaders(),
	} satisfies Record<string, string>
	const res = await authedFetch(`${API_BASE}/api/reviews/session`, {
		method: 'POST',
		headers,
	})
	if (!res.ok) throw new Error(`Review session start failed: ${res.status}`)
	return res.json()
}

export async function submitReviewAnswer(
	sessionID: string,
	answer: string,
	elapsed: number,
	attemptID: string,
): Promise<AnswerRes> {
	const headers = {
		'Content-Type': 'application/json',
		...getAuthHeaders(),
	} satisfies Record<string, string>
	const res = await authedFetch(`${API_BASE}/api/reviews/answer`, {
		method: 'POST',
		retryAttempts: 1,
		timeoutMs: 25000,
		headers,
		body: JSON.stringify({ session_id: sessionID, attempt_id: attemptID, answer, elapsed }),
	})
	if (!res.ok) await throwWithResponse(res, `Review answer submit failed: ${res.status}`)
	return res.json()
}

// Quiz (actionable on the mastery-check schedule, 80% difficulty target, own
// grading path, anonymous demo ungated)
export interface QuizStartRes {
	session_id: string
	student_id?: string
	concept_id?: string
  concept_name?: string
  question?: string
  done?: boolean
  grading_type?: string
  answer_hint?: string
  answer_input_mode?: 'numeric' | 'text'
  // Closed-book response metadata; the time limit classifies slow answers, not a hard deadline.
  closed_book?: boolean
  time_limit_seconds?: number
  questions_total?: number
}
export interface QuizAnswerRes {
  done: boolean
  correct?: boolean
  // A verdict, or one of the grader's bare status tokens ("Not equivalent"). Never the
  // solution: a test measures knowledge and does not teach, so the response carries no
  // explanation and no answer-revealing diagnosis. Instruction is reached through /learn.
  feedback?: string
  xp?: number
  new_status?: string
  /**
   * Authoritative mastery evidence score, 0..1.
   *
   * Absent on an ungraded answer (nothing was recorded, so there is no new evidence).
   * This is a progress reading, not a mastery verdict, `new_status` is the verdict, and
   * the score alone never means anything without the evidence floor the ladder also
   * requires.
   */
  evidence_score?: number
  concept_id?: string
  concept_name?: string
  question?: string
  grading_type?: string
  answer_hint?: string
  answer_input_mode?: 'numeric' | 'text'
	// Immediate remedial concepts and whether the server will allow a retake.
	remedial?: string[]
	retake_available?: boolean
	closed_book?: boolean
	time_limit_seconds?: number
	questions_total?: number
}

interface QuizSessionBody {
	student_id?: string
}

export async function startQuizSession(): Promise<QuizStartRes> {
	const { getGuestId } = await import('./auth')
	const headers = { 'Content-Type': 'application/json', ...getAuthHeaders() } satisfies Record<string, string>
	const body: QuizSessionBody = {}
	const guestId = getGuestId()
	if (guestId && !('Authorization' in headers)) body.student_id = guestId
	const res = await authedFetch(`${API_BASE}/api/quiz/session`, {
		method: 'POST',
		headers,
		body: JSON.stringify(body),
	})
	if (!res.ok) await throwWithResponse(res, `Quiz session start failed: ${res.status}`)
	return res.json()
}
interface QuizAnswerBody {
	session_id: string
	concept_id: string
	answer: string
	elapsed: number
	student_id?: string
}

export async function submitQuizAnswer(sessionId: string, conceptId: string, answer: string, elapsed: number, dontKnow?: boolean): Promise<QuizAnswerRes> {
	const { getGuestId } = await import('./auth')
	const headers = { 'Content-Type': 'application/json', ...getAuthHeaders() } satisfies Record<string, string>
	const body: QuizAnswerBody & { dont_know: boolean } = { session_id: sessionId, concept_id: conceptId, answer, elapsed, dont_know: dontKnow ?? false }
	const guestId = getGuestId()
	if (guestId && !('Authorization' in headers)) body.student_id = guestId
	const res = await authedFetch(`${API_BASE}/api/quiz/answer`, {
    // The engine replays an attempt id it has already graded, so repeating
    // this exact body cannot grade or award twice.
    retryAttempts: 1,
		timeoutMs: 25000,
		method: 'POST',
		headers,
		body: JSON.stringify(body),
	})
	if (!res.ok) await throwWithResponse(res, `Quiz answer failed: ${res.status}`)
	return res.json()
}

interface QuizSkipBody {
	session_id: string
	student_id?: string
}

// Skip the pending quiz question: advances with no grade, no XP, no remedial.
export async function skipQuizQuestion(sessionId: string): Promise<QuizAnswerRes> {
	const { getGuestId } = await import('./auth')
	const headers = { 'Content-Type': 'application/json', ...getAuthHeaders() } satisfies Record<string, string>
	const body: QuizSkipBody = { session_id: sessionId }
	const guestId = getGuestId()
	if (guestId && !('Authorization' in headers)) body.student_id = guestId
	const res = await authedFetch(`${API_BASE}/api/quiz/skip`, {
		method: 'POST',
		headers,
		body: JSON.stringify(body),
	})
	if (!res.ok) await throwWithResponse(res, `Quiz skip failed: ${res.status}`)
	return res.json()
}

// Timed quiz. This is a separate assessment from the mastery check: its gate
// counts positive XP earned after the durable rollout baseline, and the run is
// server-timed.
export interface TimedQuizRunSummary {
	run_id: string
	status: string
	started_at: number
	deadline_at: number
	completed_at?: number
	question_count: number
	correct_count: number
	unanswered_count: number
	xp_earned: number
}

export interface TimedQuizSummaryRes {
	eligible_xp: number
	gate_xp: number
	due: boolean
	server_now: number
	active?: TimedQuizRunSummary
	latest?: TimedQuizRunSummary
}

export interface TimedQuizQuestion {
	id: string
	item_id: string
	concept_id: string
	concept_name: string
	position: number
	question: string
	grading_type?: string
	answer_hint?: string
	answer_input_mode?: 'numeric' | 'text'
	served_at: number
	status: string
}

export interface TimedQuizSessionRes {
	run_id: string
	status: string
	started_at: number
	deadline_at: number
	server_now: number
	duration_seconds: number
	questions_total: number
	closed_book: boolean
	items: TimedQuizQuestion[]
}

export interface TimedQuizResultItem {
	item_id: string
	concept_id: string
	concept_name: string
	position: number
	question: string
	expected_answer: string
	answer: string
	correct: boolean
	status: string
	explanation: string
	xp_earned: number
	elapsed_seconds: number
}

export interface TimedQuizRemediation {
	concept_id: string
	concept_name: string
	href: string
}

export interface TimedQuizResultRes extends TimedQuizRunSummary {
	score_percent: number
	items: TimedQuizResultItem[]
	remediation: TimedQuizRemediation[]
}

export interface TimedQuizAnswerRes {
	run_id: string
	item_id: string
	concept_id?: string
	concept_name?: string
	correct?: boolean
	feedback?: string
	explanation?: string
	xp?: number
	new_status?: string
	completed?: boolean
	result?: TimedQuizResultRes
	expired?: boolean
	skipped?: boolean
	ungraded?: boolean
}

export async function getTimedQuiz(): Promise<TimedQuizSummaryRes> {
	const headers = { ...getAuthHeaders() } satisfies Record<string, string>
	const res = await authedFetch(`${API_BASE}/api/timed-quiz`, { headers })
	if (!res.ok) await throwWithResponse(res, `Timed quiz summary failed: ${res.status}`)
	return res.json()
}

export async function startTimedQuizSession(): Promise<TimedQuizSessionRes | { expired: true; result: TimedQuizResultRes }> {
	const headers = { 'Content-Type': 'application/json', ...getAuthHeaders() } satisfies Record<string, string>
	const res = await authedFetch(`${API_BASE}/api/timed-quiz/session`, { method: 'POST', headers })
	if (!res.ok) await throwWithResponse(res, `Timed quiz session start failed: ${res.status}`)
	return res.json()
}

export async function submitTimedQuizAnswer(
	runId: string,
	itemId: string,
	answer: string,
	elapsed: number,
	submissionId: string,
	dontKnow = false,
): Promise<TimedQuizAnswerRes> {
	const headers = { 'Content-Type': 'application/json', ...getAuthHeaders() } satisfies Record<string, string>
	const res = await authedFetch(`${API_BASE}/api/timed-quiz/answer`, {
		method: 'POST',
		retryAttempts: 1,
		timeoutMs: 25000,
		headers,
		body: JSON.stringify({ run_id: runId, item_id: itemId, answer, elapsed, submission_id: submissionId, dont_know: dontKnow }),
	})
	if (!res.ok) await throwWithResponse(res, `Timed quiz answer failed: ${res.status}`)
	return res.json()
}

export async function skipTimedQuizItem(runId: string, itemId: string): Promise<TimedQuizAnswerRes> {
	const headers = { 'Content-Type': 'application/json', ...getAuthHeaders() } satisfies Record<string, string>
	const res = await authedFetch(`${API_BASE}/api/timed-quiz/skip`, {
		method: 'POST',
		headers,
		body: JSON.stringify({ run_id: runId, item_id: itemId }),
	})
	if (!res.ok) await throwWithResponse(res, `Timed quiz skip failed: ${res.status}`)
	return res.json()
}

export async function finalizeTimedQuiz(runId: string): Promise<TimedQuizResultRes> {
	const headers = { 'Content-Type': 'application/json', ...getAuthHeaders() } satisfies Record<string, string>
	const res = await authedFetch(`${API_BASE}/api/timed-quiz/finalize`, {
		method: 'POST',
		headers,
		body: JSON.stringify({ run_id: runId }),
	})
	if (!res.ok) await throwWithResponse(res, `Timed quiz finalize failed: ${res.status}`)
	return res.json()
}

export async function getTimedQuizResult(runId: string): Promise<TimedQuizResultRes> {
	const headers = { ...getAuthHeaders() } satisfies Record<string, string>
	const res = await authedFetch(`${API_BASE}/api/timed-quiz/result?id=${encodeURIComponent(runId)}`, { headers })
	if (!res.ok) await throwWithResponse(res, `Timed quiz result failed: ${res.status}`)
	return res.json()
}

export async function validateToken(): Promise<{ valid: boolean; student_id: string }> {
	const headers = getAuthHeaders()
	if (!headers.Authorization) return { valid: false, student_id: '' }
	// Canonical route is /api/auth/me (the server has no /api/me).
	// Network failures reject so callers never mistake "offline" for "logged out".
	const res = await authedFetch(`${API_BASE}/api/auth/me`, { headers })
	if (!res.ok) return { valid: false, student_id: '' }
	try {
		const me = await res.json()
		return { valid: true, student_id: typeof me?.student_id === 'string' ? me.student_id : '' }
	} catch {
		// Unparseable body is NOT a valid session, never report valid:true.
		return { valid: false, student_id: '' }
	}
}

// Content reports, users complain about questions, explanations,
// lesson bodies, worked examples, or diagrams.
export type ReportKind = 'question' | 'explanation' | 'lesson_body' | 'worked_example' | 'diagram' | 'message'
export type ReportReason = 'wrong_answer' | 'bad_explanation' | 'unclear' | 'formatting' | 'other'
// The server's vocabulary, not the old confirmed/fixed pair, which the
// moderation endpoint stopped accepting. A type that names a status the
// server rejects is a bug waiting for a caller.
export type ReportStatus = 'open' | 'reviewing' | 'resolved' | 'dismissed'

export interface SubmitReportInput {
	concept_id?: string
	kind: ReportKind
	question?: string
	expected?: string
	explanation?: string
	lesson_id?: string
	source?: string
	session_id?: string
	attempt_id?: string
	reason: ReportReason
	detail?: string
	// What the learner typed and the format the field offered. Sent from the
	// payload already in hand so a moderator can rebuild the exact attempt.
	answer_given?: string
	grading_type?: string
	answer_hint?: string
	answer_input_mode?: string
	/** Optional screenshot. When present the request becomes multipart, so the
	 * text and the image are one report rather than two retryable operations. */
	image?: File
}

export interface SubmitReportRes {
	id: number
	status: ReportStatus
	has_image?: boolean
}

export async function submitReport(rawInput: SubmitReportInput): Promise<SubmitReportRes> {
	const { getGuestId } = await import('./auth')
	const { image, ...input } = rawInput
	const guestId = getGuestId()
	const authHeaders = getAuthHeaders()
	const reporterID = guestId && !('Authorization' in authHeaders) ? guestId : ''

	if (image) {
		const form = new FormData()
		for (const [key, value] of Object.entries(input)) {
			if (value !== undefined && value !== null && value !== '') form.append(key, String(value))
		}
		if (reporterID) form.append('reporter_id', reporterID)
		form.append('image', image)
		const res = await authedFetch(`${API_BASE}/api/reports`, {
			method: 'POST',
			headers: authHeaders,
			body: form,
		})
		if (!res.ok) throw new Error(`Report submit failed: ${res.status}`)
		return res.json() as Promise<SubmitReportRes>
	}

	const headers = { 'Content-Type': 'application/json', ...authHeaders } satisfies Record<string, string>
	const body = { ...input, ...(reporterID ? { reporter_id: reporterID } : {}) }
	const res = await authedFetch(`${API_BASE}/api/reports`, {
		method: 'POST',
		headers,
		body: JSON.stringify(body),
	})
	if (!res.ok) throw new Error(`Report submit failed: ${res.status}`)
	return res.json() as Promise<SubmitReportRes>
}

// ── Role-based administration ──────────────────────────────────────────────
//
// A different mechanism from the triage session above, on purpose. Triage is gated by a shared
// password from the deployment environment and confers no account identity; these calls carry
// the learner's own token and are authorized by the role on that account. The server decides, 
// a 403 here means the backend said no, and hiding the pages in the UI is a convenience on top
// of that, never the control.
//
// Every one of these uses authedFetch. A 401 means the caller's own session is not signed in,
// a 403 means the account lacks the permission, and authedFetch's handling of a 401 is correct
// for the first case. The server decides; hiding a page in the UI is a convenience on top.

const RoleSchema = z.enum(['student', 'moderator', 'admin', 'owner'])

// AdminUser mirrors the server's AdminUser projection. It has no field for a password hash, a
// share token or learning state, and the schema is the reason: a response carrying a field this
// type does not name cannot reach a page, whatever the server sends.
const AdminUserSchema = z.object({
	id: z.string(),
	name: z.string(),
	username: z.string(),
	email: z.string(),
	role: RoleSchema,
	created_at: z.string(),
	has_password: z.boolean(),
})

export type AdminUser = z.infer<typeof AdminUserSchema>
export type Role = z.infer<typeof RoleSchema>

// One definition, used by the overview's "recent activity" and by the audit page. Two schemas
// for one row type is how the two surfaces come to disagree about what an event looks like.
const AdminAuditEventSchema = z.object({
	id: z.number(),
	actor_id: z.string(),
	action: z.string(),
	entity_type: z.string(),
	entity_id: z.string(),
	before_json: z.string(),
	after_json: z.string(),
	created_at: z.string(),
})

const AdminOverviewSchema = z.object({
	users: z.number(),
	// Live learner counts: a session seen in the last fifteen minutes, and an
	// attempt since midnight UTC. Guests are separate because a page full of
	// guest sessions is not the same as people who can come back tomorrow.
	active_now: z.number(),
	active_now_guests: z.number(),
	active_today: z.number(),
	active_today_guests: z.number(),
	staff: z.number(),
	owners: z.number(),
	concepts: z.number(),
	domains: z.number(),
	questions: z.number(),
	reports: z.object({
		open: z.number(),
		reviewing: z.number(),
		resolved: z.number(),
		dismissed: z.number(),
		outstanding: z.number(),
	}),
	recent_activity: z.array(AdminAuditEventSchema),
})

export type AdminOverview = z.infer<typeof AdminOverviewSchema>
export type AdminAuditEvent = z.infer<typeof AdminAuditEventSchema>

export interface AdminUsersRes {
	users: AdminUser[]
}

export interface AdminAuditRes {
	events: AdminAuditEvent[]
}

export async function adminOverview(): Promise<AdminOverview> {
	const res = await authedFetch(`${API_BASE}/api/admin/overview`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load the overview')
	return validateResponse(AdminOverviewSchema, await res.json(), 'adminOverview')
}

// The account page reads one response: profile, role history, counts and the reports that
// account filed. Splitting it into four requests would give the page four ways to half-load,
// and the counts would disagree with the list while the slower one was still arriving.
const AdminActivitySchema = z.object({
	xp_total: z.number(),
	xp_today: z.number(),
	concepts_touched: z.number(),
	concepts_mastered: z.number(),
	attempts: z.number(),
	last_attempt_at: z.string(),
	last_active_day: z.string(),
	diagnostic_completed_at: z.string(),
	reports_filed: z.number(),
})

// Only the fields the page renders. A report carries the learner's answer and the expected
// answer; the account page does not need either to say who sent it, and a schema that omits
// them is why they cannot appear here by accident.
const AdminUserReportSchema = z.object({
	id: z.number(),
	kind: z.string(),
	concept_id: z.string(),
	reason: z.string(),
	status: z.string(),
	detail: z.string(),
	created_at: z.string(),
})

export type AdminUserReport = z.infer<typeof AdminUserReportSchema>
export type AdminActivity = z.infer<typeof AdminActivitySchema>

// The account detail page shows the same picture the learner sees in their
// profile. The server resolves the precedence once: a custom upload becomes a
// relative by-id route, a Google photo stays absolute, and DiceBear or a
// legacy preset travels as data for the client to draw.
const AdminUserAvatarSchema = z.object({
	avatar_url: z.string().optional(),
	avatar_dicebear: z.object({ style: z.string(), seed: z.string() }).nullable().optional(),
	avatar_preset: z.number().nullable().optional(),
})

export type AdminUserAvatar = z.infer<typeof AdminUserAvatarSchema>

export interface AdminUserDetail {
	user: AdminUser
	avatar: AdminUserAvatar
	history: AdminAuditEvent[]
	activity: AdminActivity
	reports: AdminUserReport[]
}

export async function getAdminUser(id: string): Promise<AdminUserDetail> {
	const res = await authedFetch(`${API_BASE}/api/admin/users/${encodeURIComponent(id)}`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load that account')
	const data = await res.json()
	return validateResponse(
		z.object({
			user: AdminUserSchema,
			avatar: AdminUserAvatarSchema,
			history: z.array(AdminAuditEventSchema),
			activity: AdminActivitySchema,
			reports: z.array(AdminUserReportSchema),
		}),
		data,
		'adminUserDetail',
	)
}

// The administrator's view of the weekly leagues.
//
// The public league endpoint answers only "where am I", so an administrator cannot see a cohort
// at all. This payload is the pyramid: every league that exists in a week, and one roster when a
// tier and league number are given. Member rows carry the account id, which the public board
// deliberately withholds, because this is the surface where a row has to link to an account.
const AdminLeagueMemberSchema = z.object({
	rank: z.number(),
	name: z.string(),
	username: z.string(),
	student_id: z.string(),
	xp: z.number(),
	avatar_url: z.string().optional(),
	avatar_dicebear: z.object({ style: z.string(), seed: z.string() }).nullable().optional(),
	avatar_preset: z.number().nullable().optional(),
})

const AdminLeagueSchema = z.object({
	tier: z.string(),
	league_no: z.number(),
	size: z.number(),
	movers: z.number(),
	members: z.array(AdminLeagueMemberSchema).optional(),
})

const AdminLeaguesSchema = z.object({
	week: z.string(),
	week_start: z.string(),
	week_end: z.string(),
	promote: z.number(),
	demote: z.number(),
	seats: z.number(),
	leagues: z.array(AdminLeagueSchema),
})

export type AdminLeague = z.infer<typeof AdminLeagueSchema>
export type AdminLeagueMember = z.infer<typeof AdminLeagueMemberSchema>
export type AdminLeagues = z.infer<typeof AdminLeaguesSchema>

/**
 * Read the week's leagues, or one league's standings.
 *
 * All three parameters travel together or not at all: a tier without a league number is refused
 * by the server rather than ignored, because a page that asked for one tier and then rendered the
 * whole pyramid would be showing something other than what it requested.
 */
export async function getAdminLeagues(
	query: { week?: string; tier?: string; league?: number } = {},
): Promise<AdminLeagues> {
	const params = new URLSearchParams()
	if (query.week) params.set('week', query.week)
	if (query.tier && query.league) {
		params.set('tier', query.tier)
		params.set('league', String(query.league))
	}
	const qs = params.toString()
	const res = await authedFetch(`${API_BASE}/api/admin/leagues${qs ? `?${qs}` : ''}`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load leagues')
	return validateResponse(AdminLeaguesSchema, await res.json(), 'adminLeagues')
}

export async function searchAdminUsers(q: string): Promise<AdminUsersRes> {
	const params = new URLSearchParams()
	if (q.trim()) params.set('q', q.trim())
	const qs = params.toString()
	const res = await authedFetch(`${API_BASE}/api/admin/users${qs ? `?${qs}` : ''}`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load users')
	const data = (await res.json()) as { users?: unknown }
	return { users: validateResponse(z.array(AdminUserSchema), data.users ?? [], 'adminUsers') }
}

/**
 * Set an account's role.
 *
 * `role` is a closed union here as well as on the server. Both ends reject anything else, and
 * the reason they both bother is that coercing would let the response say one thing and the row
 * say another, which is exactly the sort of quiet disagreement the rest of this client has been
 * bitten by.
 */
export async function setUserRole(userId: string, role: Role): Promise<{ user: AdminUser; changed: boolean }> {
	const res = await authedFetch(`${API_BASE}/api/admin/users/${encodeURIComponent(userId)}`, {
		method: 'PATCH',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ role }),
	})
	if (!res.ok) await throwWithResponse(res, 'Could not change the role')
	const data = (await res.json()) as { user?: unknown; changed?: unknown }
	return {
		user: validateResponse(AdminUserSchema, data.user, 'adminUserRole'),
		changed: data.changed === true,
	}
}

export async function listAdminAudit(): Promise<AdminAuditRes> {
	const res = await authedFetch(`${API_BASE}/api/admin/audit`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load the audit log')
	const data = (await res.json()) as { events?: unknown }
	return { events: validateResponse(z.array(AdminAuditEventSchema), data.events ?? [], 'adminAudit') }
}

// ── Admin V2 ──────────────────────────────────────────────────────────────
//
// The moderation queue, contributor management and content inspection. The permission list from
// /api/admin/me drives which sections render; every call is still authorized server-side.

export interface AdminMe {
	id: string
	role: Role
	permissions: string[]
	assignable_roles: Role[]
}

export async function adminMe(): Promise<AdminMe> {
	const res = await authedFetch(`${API_BASE}/api/admin/me`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load your admin access')
	const data = (await res.json()) as AdminMe
	return data
}

// A moderation record. It is never deleted; resolved/dismissed stamp who decided and why.
const QuestionReportSchema = z.object({
	id: z.number(),
	reporter_id: z.string(),
	reporter_username: z.string().optional(),
	reporter_name: z.string().optional(),
	concept_id: z.string(),
	kind: z.string(),
	question: z.string(),
	expected: z.string(),
	explanation: z.string(),
	lesson_id: z.string(),
	source: z.string(),
	session_id: z.string(),
	attempt_id: z.string(),
	// Captured with the report so a moderator can replay the attempt rather
	// than guess it: what the learner typed, and the field they typed it in.
	// Older rows predate the capture, hence optional.
	answer_given: z.string().optional(),
	grading_type: z.string().optional(),
	answer_hint: z.string().optional(),
	// The Go field has no omitempty, so a report filed before the capture
	// existed carries "" rather than omitting the key. Without the preprocess
	// that empty string fails the enum and one old row hides the whole queue.
	answer_input_mode: z.preprocess(
		v => (v === '' ? undefined : v),
		z.enum(['numeric', 'text']).optional(),
	),
	reason: z.string(),
	detail: z.string(),
	status: z.string(),
	resolution: z.string(),
	resolved_by: z.string(),
	resolved_at: z.string(),
	created_at: z.string(),
	// Older rows and non-image reports carry false/absent. The moderation page
	// fetches the bytes separately because an <img> tag cannot send the token.
	has_image: z.boolean().optional(),
	unread: z.boolean().optional(),
})
export type QuestionReport = z.infer<typeof QuestionReportSchema>

export type ModerationStatus = 'open' | 'reviewing' | 'resolved' | 'dismissed'

export interface MyReportsResponse {
	reports: QuestionReport[]
	unread: number
}

// The learner-facing view of the reports this account filed. Moderation is
// separate; this endpoint returns only the caller's rows.
export async function listMyReports(): Promise<MyReportsResponse> {
	const res = await authedFetch(`${API_BASE}/api/reports/mine`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load your reports')
	const data = (await res.json()) as { reports?: unknown; unread?: unknown }
	return {
		reports: validateResponse(z.array(QuestionReportSchema), data.reports ?? [], 'myReports'),
		unread: typeof data.unread === 'number' ? data.unread : 0,
	}
}

export async function markMyReportsRead(): Promise<void> {
	const res = await authedFetch(`${API_BASE}/api/reports/mine/read`, { method: 'POST' })
	if (!res.ok) await throwWithResponse(res, 'Could not mark your reports as read')
}

export async function clearMyReports(): Promise<void> {
	const res = await authedFetch(`${API_BASE}/api/reports/mine/clear`, { method: 'POST' })
	if (!res.ok) await throwWithResponse(res, 'Could not clear your reports')
}

const AnnouncementSchema = z.object({
	id: z.string(),
	title: z.string(),
	body: z.string(),
	published: z.boolean(),
	created_at: z.string(),
	published_at: z.string(),
	read: z.boolean().optional(),
})
export type Announcement = z.infer<typeof AnnouncementSchema>

export interface MyAnnouncementsResponse {
	announcements: Announcement[]
	unread: number
}

export async function listAnnouncements(): Promise<MyAnnouncementsResponse> {
	const res = await authedFetch(`${API_BASE}/api/announcements`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load announcements')
	const data = (await res.json()) as { announcements?: unknown; unread?: unknown }
	return {
		announcements: validateResponse(z.array(AnnouncementSchema), data.announcements ?? [], 'announcements'),
		unread: typeof data.unread === 'number' ? data.unread : 0,
	}
}

export async function markAnnouncementsRead(): Promise<void> {
	const res = await authedFetch(`${API_BASE}/api/announcements/read`, { method: 'POST' })
	if (!res.ok) await throwWithResponse(res, 'Could not mark announcements as read')
}

export async function clearAnnouncements(): Promise<void> {
	const res = await authedFetch(`${API_BASE}/api/announcements/clear`, { method: 'POST' })
	if (!res.ok) await throwWithResponse(res, 'Could not clear announcements')
}

// Friends. The flag is server-owned; a 404 here means the deployment has not
// cleared the surface, not that the page is broken, so callers should treat it
// as "not available" rather than an error to retry.
export interface FriendAvatar {
	avatar_url?: string
	avatar_dicebear?: { style: string; seed: string } | null
	avatar_preset?: number | null
}

export interface FriendIdentity {
	profile_id: string
	friend_code: string
}

export interface FriendSummary {
	profile_id: string
	name: string
	username?: string
	avatar: FriendAvatar
	xp_total: number
	mastered: number
	friends_since?: string
}

export interface FriendRequest {
	profile_id: string
	name: string
	username?: string
	avatar: FriendAvatar
	direction: 'incoming' | 'outgoing'
	created_at: string
}

export interface FriendsFeed {
	identity: FriendIdentity
	friends: FriendSummary[]
	incoming: FriendRequest[]
	outgoing: FriendRequest[]
}

export interface FriendGraphNode {
	concept_id: string
	status: string
}

export interface FriendLeaderboardRow {
	profile_id: string
	name: string
	username?: string
	avatar: FriendAvatar
	xp: number
	mastered: number
	is_me: boolean
}

export async function getFriends(): Promise<FriendsFeed> {
	const res = await authedFetch(`${API_BASE}/api/friends`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load friends')
	return (await res.json()) as FriendsFeed
}

async function friendPost(path: string, body: Record<string, string>): Promise<Response> {
	const res = await authedFetch(`${API_BASE}/api/friends/${path}`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body),
	})
	if (!res.ok) await throwWithResponse(res, 'Friend action failed')
	return res
}

export async function sendFriendRequest(query: string): Promise<'pending' | 'accepted'> {
	const res = await friendPost('requests', { query })
	const data = (await res.json()) as { status?: string }
	return data.status === 'accepted' ? 'accepted' : 'pending'
}

export async function acceptFriendRequest(profileId: string): Promise<void> {
	await friendPost('requests/accept', { profile_id: profileId })
}

export async function declineFriendRequest(profileId: string): Promise<void> {
	await friendPost('requests/decline', { profile_id: profileId })
}

export async function cancelFriendRequest(profileId: string): Promise<void> {
	await friendPost('requests/cancel', { profile_id: profileId })
}

export async function removeFriend(profileId: string): Promise<void> {
	await friendPost('remove', { profile_id: profileId })
}

export async function blockFriend(profileId: string): Promise<void> {
	await friendPost('block', { profile_id: profileId })
}

export async function unblockFriend(profileId: string): Promise<void> {
	await friendPost('unblock', { profile_id: profileId })
}

export async function rotateFriendCode(): Promise<string> {
	const res = await friendPost('code/rotate', {})
	const data = (await res.json()) as { friend_code?: string }
	return data.friend_code ?? ''
}

export async function getFriendProfile(profileId: string): Promise<FriendSummary> {
	const res = await authedFetch(`${API_BASE}/api/friends/profile?profile_id=${encodeURIComponent(profileId)}`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load this friend')
	return (await res.json()) as FriendSummary
}

export async function getFriendGraph(profileId: string): Promise<FriendGraphNode[]> {
	const res = await authedFetch(`${API_BASE}/api/friends/graph?profile_id=${encodeURIComponent(profileId)}`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load this graph')
	const data = (await res.json()) as { nodes?: FriendGraphNode[] }
	return data.nodes ?? []
}

export async function getFriendsLeaderboard(view: 'weekly' | 'all_time'): Promise<FriendLeaderboardRow[]> {
	const res = await authedFetch(`${API_BASE}/api/friends/leaderboard?view=${view}`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load the friends leaderboard')
	const data = (await res.json()) as { rows?: FriendLeaderboardRow[] }
	return data.rows ?? []
}

export async function listAdminAnnouncements(): Promise<{ announcements: Announcement[] }> {
	const res = await authedFetch(`${API_BASE}/api/admin/announcements`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load announcements')
	const data = (await res.json()) as { announcements?: unknown }
	return { announcements: validateResponse(z.array(AnnouncementSchema), data.announcements ?? [], 'adminAnnouncements') }
}

export interface AnnouncementInput {
	title: string
	body: string
	published: boolean
}

export async function createAdminAnnouncement(input: AnnouncementInput): Promise<Announcement> {
	const res = await authedFetch(`${API_BASE}/api/admin/announcements`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(input),
	})
	if (!res.ok) await throwWithResponse(res, 'Could not create the announcement')
	const data = (await res.json()) as { announcement?: unknown }
	return validateResponse(AnnouncementSchema, data.announcement, 'createAnnouncement')
}

export async function updateAdminAnnouncement(
	id: string,
	input: Partial<AnnouncementInput>,
): Promise<Announcement> {
	const res = await authedFetch(`${API_BASE}/api/admin/announcements/${encodeURIComponent(id)}`, {
		method: 'PATCH',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(input),
	})
	if (!res.ok) await throwWithResponse(res, 'Could not update the announcement')
	const data = (await res.json()) as { announcement?: unknown }
	return validateResponse(AnnouncementSchema, data.announcement, 'updateAnnouncement')
}

export async function listAdminReports(status?: string): Promise<{ reports: QuestionReport[] }> {
	const qs = status && status !== 'all' ? `?status=${encodeURIComponent(status)}` : ''
	const res = await authedFetch(`${API_BASE}/api/admin/reports${qs}`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load reports')
	const data = (await res.json()) as { reports?: unknown }
	return { reports: validateResponse(z.array(QuestionReportSchema), data.reports ?? [], 'adminReports') }
}

export async function getAdminReportImage(id: number): Promise<Blob> {
	const res = await authedFetch(`${API_BASE}/api/admin/reports/${id}/image`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load the screenshot')
	return res.blob()
}

export async function getAdminReport(id: number): Promise<{ report: QuestionReport; history: AdminAuditEvent[] }> {
	const res = await authedFetch(`${API_BASE}/api/admin/reports/${id}`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load the report')
	const data = (await res.json()) as { report?: unknown; history?: unknown }
	return {
		report: validateResponse(QuestionReportSchema, data.report, 'adminReport'),
		history: validateResponse(z.array(AdminAuditEventSchema), data.history ?? [], 'adminReportHistory'),
	}
}

// Replays a stored attempt through the same grader the learner faced. The
// endpoint is pure: it never enters SubmitAnswer, so no XP, streak, attempt or
// mastery changes. `writes: false` in the response is the server saying so.
export interface AdminGradeRes {
	correct: boolean
	score: number
	feedback?: string
	diagnosis?: string
	unavailable?: boolean
	graded_answer: string
	expected: string
	answer_hint?: string
	answer_input_mode?: 'numeric' | 'text'
	writes: boolean
}

export async function gradeAdminReport(id: number, answer: string): Promise<AdminGradeRes> {
	const res = await authedFetch(`${API_BASE}/api/admin/reports/${id}/grade`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ answer }),
	})
	if (!res.ok) await throwWithResponse(res, 'Could not grade the reported answer')
	return (await res.json()) as AdminGradeRes
}

export async function moderateReport(id: number, status: ModerationStatus, resolution: string): Promise<QuestionReport> {
	const res = await authedFetch(`${API_BASE}/api/admin/reports/${id}`, {
		method: 'PATCH',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ status, resolution }),
	})
	if (!res.ok) await throwWithResponse(res, 'Could not update the report')
	const data = (await res.json()) as { report?: unknown }
	return validateResponse(QuestionReportSchema, data.report, 'adminReportUpdate')
}

export interface BulkModerationRes {
	updated: number
	skipped: number
	failed: { id: number; error: string }[]
}

// One decision across a concept's unsettled reports: a fix is usually one bug,
// not one report. The server still writes one transaction, stamp and audit row
// per report, this call only saves the operator from retyping the same reason.
export async function bulkModerateReports(
	selector: { ids?: number[]; concept_id?: string },
	status: ModerationStatus,
	resolution: string,
): Promise<BulkModerationRes> {
	const res = await authedFetch(`${API_BASE}/api/admin/reports/bulk`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ ...selector, status, resolution }),
	})
	if (!res.ok) await throwWithResponse(res, 'Could not close the reports')
	const data = (await res.json()) as Partial<BulkModerationRes>
	return {
		updated: data.updated ?? 0,
		skipped: data.skipped ?? 0,
		failed: Array.isArray(data.failed) ? data.failed : [],
	}
}

export async function listAdmins(): Promise<{ admins: AdminUser[] }> {
	const res = await authedFetch(`${API_BASE}/api/admin/admins`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load staff')
	const data = (await res.json()) as { admins?: unknown }
	return { admins: validateResponse(z.array(AdminUserSchema), data.admins ?? [], 'adminAdmins') }
}

const AdminInvitationSchema = z.object({
	id: z.number(),
	email: z.string(),
	role: RoleSchema,
	invited_by: z.string(),
	created_at: z.string(),
	expires_at: z.string(),
	accepted_at: z.string(),
	accepted_by: z.string(),
	revoked_at: z.string(),
	revoked_by: z.string(),
})
export type AdminInvitation = z.infer<typeof AdminInvitationSchema>

export async function listAdminInvitations(): Promise<{ invitations: AdminInvitation[] }> {
	const res = await authedFetch(`${API_BASE}/api/admin/invitations`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load invitations')
	const data = (await res.json()) as { invitations?: unknown }
	return { invitations: validateResponse(z.array(AdminInvitationSchema), data.invitations ?? [], 'adminInvitations') }
}

export async function createAdminInvitation(email: string, role: Role): Promise<{ id: number; token: string; email: string; role: Role; expires_at: string }> {
	const res = await authedFetch(`${API_BASE}/api/admin/invitations`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ email, role }),
	})
	if (!res.ok) await throwWithResponse(res, 'Could not create the invitation')
	return res.json() as Promise<{ id: number; token: string; email: string; role: Role; expires_at: string }>
}

export async function acceptAdminInvitation(token: string): Promise<{ role: Role }> {
	const res = await authedFetch(`${API_BASE}/api/admin/invitations/accept`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ token }),
	})
	if (!res.ok) await throwWithResponse(res, 'Could not accept the invitation')
	return res.json() as Promise<{ role: Role }>
}

export async function revokeAdminInvitation(id: number): Promise<AdminInvitation> {
	const res = await authedFetch(`${API_BASE}/api/admin/invitations/${id}/revoke`, { method: 'POST' })
	if (!res.ok) await throwWithResponse(res, 'Could not revoke the invitation')
	const data = (await res.json()) as { invitation?: unknown }
	return validateResponse(AdminInvitationSchema, data.invitation, 'adminInvitationRevoke')
}

// ── Optional AI augmentation ──
//
// Both calls are best-effort enhancements, never a dependency. The server
// refuses them when AI is disabled (503) or the caller's daily budget is spent
// (429); getGoalScope folds either into used_ai=false so a caller falls back to
// its deterministic match with no special case. No `expected_answer` is ever a
// field here: the graded answer stays on the server.

export interface AssistantLink {
	/** The id the assistant named, e.g. "glossary" or a destination id. */
	id: string
	label: string
	/** An internal route the server chose; never built from the model's text. */
	href: string
}

export interface AssistantTurn {
	role: 'user' | 'assistant'
	content: string
}

/**
 * One run of a reply: either text or a link. Links arrive as parts rather than
 * as a list beside the text, so the label appears where the assistant wrote it
 * and is not printed a second time underneath.
 */
export interface AssistantPart {
	text?: string
	link?: AssistantLink
}

export interface AssistantReply {
	reply: string
	segments: AssistantPart[]
}

/**
 * Ask the product assistant. It answers "what is Mathua" and "where do I
 * start"; it is not a tutor and holds no concept or mastery context. `history`
 * is replayed so the assistant can follow up on its own clarifying question; the
 * server bounds it and keeps only user and assistant turns. Throws when AI is
 * unavailable or the quota is spent, so the caller shows "unavailable" rather
 * than a wrong answer.
 */
export async function askAssistant(body: { message: string; history?: AssistantTurn[] }): Promise<AssistantReply> {
	const headers = { 'Content-Type': 'application/json', ...getAuthHeaders() } satisfies Record<string, string>
	const res = await authedFetch(`${API_BASE}/api/ai/assistant`, {
		method: 'POST',
		headers,
		body: JSON.stringify(body),
	})
	if (!res.ok) await throwWithResponse(res, `Assistant unavailable: ${res.status}`)
	const data = (await res.json()) as AssistantReply
	return { reply: data.reply, segments: data.segments ?? [] }
}

export interface GoalScopeRes {
	/** Destination ids the model chose, validated against the server's own list. */
	destinations: string[]
	/** False when AI is off, failed, or chose nothing, use the deterministic match. */
	used_ai: boolean
	/**
	 * Why AI did not contribute, set only when used_ai is false. The server
	 * distinguishes "the provider failed" from "the model chose nothing", so the
	 * UI can say "suggestions unavailable" instead of silently showing an empty
	 * list for both.
	 */
	reason?: 'disabled' | 'no_candidates' | 'failed' | 'no_match'
}

/**
 * Ask the model to map a free-text goal to destinations. Never throws: a
 * failure is reported as used_ai=false, because the caller already has a
 * deterministic answer on screen and must not lose it to a provider outage.
 */
export async function getGoalScope(goal: string): Promise<GoalScopeRes> {
	const headers = { 'Content-Type': 'application/json', ...getAuthHeaders() } satisfies Record<string, string>
	const res = await authedFetch(`${API_BASE}/api/ai/goal-scope`, {
		method: 'POST',
		headers,
		body: JSON.stringify({ goal }),
	})
	if (!res.ok) return { destinations: [], used_ai: false }
	return res.json() as Promise<GoalScopeRes>
}

export interface AdminContentConcept {
	id: string
	label: string
	domain: string
	subdomain: string
	grading_type: string
	prerequisites: string[]
	// The format the question will actually be served with, and whether it was
	// declared by the generator rather than derived from the grading type. The
	// concept record alone cannot answer this: grading_type names the grader.
	answer_hint: string
	answer_input_mode: 'numeric' | 'text'
	answer_overridden: boolean
}

export async function listAdminContent(
	q = '',
	opts: { overridden?: boolean } = {},
): Promise<{ concepts: AdminContentConcept[]; total: number }> {
	const params = new URLSearchParams()
	if (q.trim()) params.set('q', q.trim())
	if (opts.overridden) params.set('overridden', '1')
	const qs = params.toString() ? `?${params.toString()}` : ''
	const res = await authedFetch(`${API_BASE}/api/admin/content${qs}`, { cache: 'no-store' })
	if (!res.ok) await throwWithResponse(res, 'Could not load content')
	const data = (await res.json()) as { concepts?: AdminContentConcept[]; total?: number }
	return { concepts: data.concepts ?? [], total: data.total ?? 0 }
}

// The lab. Both routes are read-only: generate serves what the generator would
// produce for a concept, and grade runs the same pure GradeAnswer the report
// replay uses. `writes: false` in each response is the server stating that.
export interface LabQuestion {
	question: string
	answer: string
	explanation: string
	answer_hint: string
	answer_input_mode: 'numeric' | 'text'
}

export interface LabGenerateRes {
	concept_id: string
	concept_name: string
	grading_type: string
	answer_hint: string
	answer_input_mode: 'numeric' | 'text'
	difficulty: number
	questions: LabQuestion[]
	writes: boolean
}

export async function generateLabQuestions(
	conceptId: string,
	count = 3,
	difficulty = 0.5,
): Promise<LabGenerateRes> {
	const res = await authedFetch(`${API_BASE}/api/admin/lab/generate`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ concept_id: conceptId, count, difficulty }),
	})
	if (!res.ok) await throwWithResponse(res, 'Could not generate questions')
	return (await res.json()) as LabGenerateRes
}

export async function gradeLabAnswer(
	conceptId: string,
	expected: string,
	answer: string,
): Promise<AdminGradeRes> {
	const res = await authedFetch(`${API_BASE}/api/admin/lab/grade`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ concept_id: conceptId, expected, answer }),
	})
	if (!res.ok) await throwWithResponse(res, 'Could not grade the answer')
	const data = (await res.json()) as AdminGradeRes
	return data
}

// Learner reports read only what the server can substantiate. The response
// deliberately carries its limitations, so the page can state what it does not
// know instead of inferring a history from current totals.
export interface LearnerReportPeriod {
  days: number
  since?: string
  until: string
  label: string
}

export interface LearnerReportOverview {
  total: number
  mastered: number
  in_progress: number
  unlocked: number
  locked: number
  xp_lifetime: number
  xp_period: number
  activity_days: number
  questions: number
  correct: number
}

export interface LearnerDomainSummary {
  id: string
  conceptCount: number
  masteredCount: number
  unlockedCount: number
  inProgressCount: number
}

export interface LearnerReportGap {
  concept_id: string
  label: string
  state?: string
}

export interface LearnerMasteryRow {
  concept_id: string
  label: string
  domain: string
  state: string
  status: string
  decaying: boolean
  mastery_pct: number
  mastered_at?: string
  attempts: number
  streak: number
  last_attempted?: string
  next_review?: string
  placement_seeded: boolean
  provisional: boolean
  prerequisite_gaps: LearnerReportGap[] | null
  href: string
}

export interface LearnerAssessment {
  id: string
  type: string
  title: string
  status: string
  started_at?: string
  completed_at?: string
  questions: number
  correct: number
  accuracy?: number
  concepts: string[] | null
  findings: LearnerReportGap[] | null
  href?: string
  note?: string
}

export interface LearnerHistoryEvent {
  id: string
  kind: string
  at: string
  title: string
  concept_id?: string
  source?: string
  detail?: string
  href?: string
  correct?: boolean
  xp?: number
}

export interface LearnerReport {
  generated_at: string
  learner: Record<string, string>
  period: LearnerReportPeriod
  overview: LearnerReportOverview
  domains: LearnerDomainSummary[] | null
  mastery: LearnerMasteryRow[] | null
  assessments: LearnerAssessment[] | null
  activity: { date: string; questions: number; correct: number; concepts: string[] | null }[] | null
  history: LearnerHistoryEvent[] | null
  history_total: number
  limitations: string[] | null
}

export interface LearnerReportHistory {
  events: LearnerHistoryEvent[]
  total: number
  limit: number
  offset: number
}

export async function getLearnerReport(params: { days?: number; domain?: string } = {}): Promise<LearnerReport> {
  const qs = new URLSearchParams()
  if (params.days !== undefined) qs.set('days', String(params.days))
  if (params.domain) qs.set('domain', params.domain)
  const query = qs.toString()
  const res = await authedFetch(`${API_BASE}/api/learner/reports${query ? `?${query}` : ''}`, {
    headers: { ...getAuthHeaders() },
    cache: 'no-store',
  })
  if (!res.ok) throw new Error(`Reports fetch failed: ${res.status}`)
  return res.json()
}

export async function getLearnerReportHistory(params: { limit?: number; offset?: number } = {}): Promise<LearnerReportHistory> {
  const qs = new URLSearchParams()
  if (params.limit !== undefined) qs.set('limit', String(params.limit))
  if (params.offset !== undefined) qs.set('offset', String(params.offset))
  const query = qs.toString()
  const res = await authedFetch(`${API_BASE}/api/learner/reports/history${query ? `?${query}` : ''}`, {
    headers: { ...getAuthHeaders() },
    cache: 'no-store',
  })
  if (!res.ok) throw new Error(`Report history fetch failed: ${res.status}`)
  return res.json()
}
