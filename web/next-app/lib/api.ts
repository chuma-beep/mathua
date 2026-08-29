export const API_BASE = process.env.NEXT_PUBLIC_API_URL || ''

import { getAuthHeaders } from './auth'
import { z } from 'zod'

const QuestionSchema = z.object({
  concept_id: z.string(),
  concept_name: z.string(),
  question: z.string(),
  is_review: z.boolean(),
  attempt_id: z.string().optional(),
  lesson: z.object({ Title: z.string(), Body: z.string(), Concepts: z.array(z.string()) }).optional(),
  diagram: z.string().optional(),
})

const StartSessionResSchema = z.object({
  student_id: z.string(),
  session_id: z.string(),
  question: QuestionSchema.nullable(),
})

const AnswerResultSchema = z.object({
  correct: z.boolean(),
  feedback: z.string(),
  new_status: z.string(),
  explanation: z.string(),
  streak: z.number(),
  required_streak: z.number(),
  xp: z.number(),
  expected_answer: z.string(),
})

const AnswerResSchema = z.object({
  result: AnswerResultSchema.nullable(),
  next_question: QuestionSchema.nullable(),
  done: z.boolean(),
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
})

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

const PrereqInfoSchema = z.object({
  id: z.string(),
  label: z.string(),
  status: z.string(),
  mastery_pct: z.number(),
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

const ConceptDetailResSchema = z.object({
  concept: z.object({
    id: z.string(),
    label: z.string(),
    domain: z.string(),
    subdomain: z.string().optional(),
  }),
  lesson: z.object({ title: z.string(), body: z.string(), concepts: z.array(z.string()) }).optional(),
  prerequisites: z.array(PrereqInfoSchema),
  dependents: z.array(PrereqInfoSchema).optional(),
  unlocked: z.boolean(),
  progress: z.object({
    status: z.string(),
    streak: z.number(),
    required_streak: z.number(),
    mastery_pct: z.number(),
  }).optional(),
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
  explanation: string
  streak: number
  required_streak: number
  xp: number
  expected_answer: string
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

export interface LeaderboardEntry {
  rank: number
  name: string
  mastered: number
  streak: number
  level: string
  score: number
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
}

export interface ConceptProgress {
  concept_id?: string
  status: string
  streak: number
  best_streak?: number
  avg_response_time?: number
  attempts?: number
}

export async function startSession(): Promise<StartSessionRes> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...getAuthHeaders(),
  }
  const res = await fetch(`${API_BASE}/api/session`, {
    method: 'POST',
    headers,
  })
  if (!res.ok) throw new Error(`Session start failed: ${res.status}`)
  return validateResponse(StartSessionResSchema, await res.json(), 'startSession') as StartSessionRes
}

export async function startSessionName(name: string): Promise<StartSessionRes> {
  const res = await fetch(`${API_BASE}/api/session`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name }),
  })
  if (!res.ok) throw new Error(`Session start failed: ${res.status}`)
  return res.json()
}

export function isConflictError(err: unknown): boolean {
  return err instanceof Error && (err as unknown as Record<string, unknown>).status === 409
}

export function isTooQuickError(err: unknown): boolean {
  return err instanceof Error && (err as unknown as Record<string, unknown>).status === 400
}

export function getErrorMessage(err: unknown): string {
  if (err instanceof Error) {
    const e = err as unknown as Record<string, unknown>
    if (typeof e.serverMessage === 'string' && e.serverMessage) return e.serverMessage as string
    return err.message
  }
  return String(err)
}

function throwWithStatus(msg: string, status: number): never {
  throw Object.assign(new Error(msg), { status })
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

export async function submitAnswer(
  sessionID: string,
  answer: string,
  elapsed: number,
  attemptID: string,
): Promise<AnswerRes> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...getAuthHeaders(),
  }
  const res = await fetch(`${API_BASE}/api/answer`, {
    method: 'POST',
    headers,
    body: JSON.stringify({ session_id: sessionID, attempt_id: attemptID, answer, elapsed }),
  })
  if (!res.ok) await throwWithResponse(res, `Answer submit failed: ${res.status}`)
  return validateResponse(AnswerResSchema, await res.json(), 'submitAnswer') as AnswerRes
}

export async function getCurrentQuestion(sessionID: string): Promise<Question | null> {
  const headers: Record<string, string> = { ...getAuthHeaders() }
  const res = await fetch(`${API_BASE}/api/session/current?session_id=${encodeURIComponent(sessionID)}`, { headers })
  if (!res.ok) return null
  const data = await res.json()
  if (!data.question) return null
  return validateResponse(QuestionSchema, data.question, 'getCurrentQuestion') as Question
}

export async function getProgress(studentID: string): Promise<Record<string, ConceptProgress>> {
  const res = await fetch(`${API_BASE}/api/progress/${studentID}`)
  if (!res.ok) return {}
  return validateResponse(ProgressMapSchema, await res.json(), 'getProgress') as Record<string, ConceptProgress>
}

export async function getScores(studentID: string): Promise<Scores> {
  const res = await fetch(`${API_BASE}/api/scores/${studentID}`)
  if (!res.ok) throw new Error(`Scores fetch failed: ${res.status}`)
  return validateResponse(ScoresSchema, await res.json(), 'getScores') as Scores
}

export async function getGraph(): Promise<GraphRes> {
  const res = await fetch(`${API_BASE}/api/graph`)
  if (!res.ok) throw new Error(`Graph fetch failed: ${res.status}`)
  return validateResponse(GraphResSchema, await res.json(), 'getGraph') as GraphRes
}

export async function getLeaderboard(): Promise<LeaderboardEntry[]> {
  const res = await fetch(`${API_BASE}/api/leaderboard`)
  if (!res.ok) return []
  return res.json()
}

export interface LeagueMember {
  student_id: string
  name: string
  tier: string
  total_mastered: number
  weekly_mastered: number
  moved: number
}

export interface League {
  tier: string
  members: LeagueMember[]
}

export interface LeagueBoard {
  week: string
  leagues: League[]
}

export async function getLeagues(): Promise<LeagueBoard> {
  const res = await fetch(`${API_BASE}/api/leagues`, { headers: { ...getAuthHeaders() } })
  if (!res.ok) return { week: '', leagues: [] }
  return res.json()
}

export interface ShareReport {
  student_id: string
  name: string
  scores: Scores
  activity: DailyActivity[]
  progress: Record<string, ConceptProgress>
  weakness: Record<string, number>
}

export async function enableShare(): Promise<{ enabled: boolean; token: string }> {
  const res = await fetch(`${API_BASE}/api/share`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ enabled: true }),
  })
  if (!res.ok) throw new Error(`Enable share failed: ${res.status}`)
  return res.json()
}

export async function disableShare(): Promise<void> {
  await fetch(`${API_BASE}/api/share`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ enabled: false }),
  })
}

export async function getShareReport(token: string): Promise<ShareReport | null> {
  const res = await fetch(`${API_BASE}/api/share/${encodeURIComponent(token)}`, { cache: 'no-store' })
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
  const res = await fetch(`${API_BASE}/api/transcript`, { headers: { ...getAuthHeaders() }, cache: 'no-store' })
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
  const res = await fetch(`${API_BASE}/api/efficacy`, { headers: { ...getAuthHeaders() } })
  if (!res.ok) return null
  return res.json()
}

export interface ConfigRes {
	auth_enabled: boolean
}

export async function getConfig(): Promise<ConfigRes> {
	const res = await fetch(`${API_BASE}/api/config`)
	if (!res.ok) return { auth_enabled: false }
	return res.json()
}

export async function healthCheck(): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/api/health`)
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

export async function signup(name: string, username: string, password: string): Promise<AuthRes> {
  const res = await fetch(`${API_BASE}/api/auth/signup`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, username, password }),
  })
  if (!res.ok) throw new Error('Signup failed')
  return validateResponse(AuthResSchema, await res.json(), 'signup') as AuthRes
}

export async function login(username: string, password: string): Promise<AuthRes> {
	const res = await fetch(`${API_BASE}/api/auth/login`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ username, password }),
	})
	if (!res.ok) throw new Error('Invalid credentials')
	return validateResponse(AuthResSchema, await res.json(), 'login') as AuthRes
}

// Goal-based diagnostic API

export interface GoalPathRes {
	concepts: { id: string; label: string; domain: string; prerequisites: string[] }[]
	count: number
}

export interface GoalDiagStartRes {
	session_id: string
	student_id?: string
	concept_id?: string
	concept_name?: string
	question?: string
	done?: boolean
}

export interface GoalDiagAnswerRes {
	done: boolean
	correct?: boolean
	feedback?: string
	concept_id?: string
	concept_name?: string
	question?: string
}

export interface GoalPlanRes {
	readiness: number
	total_tested: number
	correct_count: number
	weak_areas: Record<string, { id: string; label: string }[]>
	strong_areas: Record<string, string[]>
}

export interface WeaknessRes {
	by_domain: Record<string, { id: string; label: string; weakness: number }[]>
}

export async function getGoalPath(conceptIds: string[]): Promise<GoalPathRes> {
	const res = await fetch(`${API_BASE}/api/goal`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify({ concept_ids: conceptIds }),
	})
	if (!res.ok) throw new Error(`Goal path fetch failed: ${res.status}`)
	return res.json()
}

export async function startGoalDiagnostic(conceptIds: string[]): Promise<GoalDiagStartRes> {
	const res = await fetch(`${API_BASE}/api/goal/diagnostic`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify({ concept_ids: conceptIds }),
	})
	if (!res.ok) throw new Error(`Goal diagnostic start failed: ${res.status}`)
	return res.json()
}

export async function startGoalDiagnosticName(name: string, conceptIds: string[]): Promise<GoalDiagStartRes> {
	const res = await fetch(`${API_BASE}/api/goal/diagnostic`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ name, concept_ids: conceptIds }),
	})
	if (!res.ok) throw new Error(`Goal diagnostic start failed: ${res.status}`)
	return res.json()
}

export async function submitGoalAnswer(
	sessionId: string,
	conceptId: string,
	answer: string,
	elapsed: number,
): Promise<GoalDiagAnswerRes> {
	const res = await fetch(`${API_BASE}/api/goal/diagnostic/answer`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify({ session_id: sessionId, concept_id: conceptId, answer, elapsed }),
	})
	if (!res.ok) throw new Error(`Goal answer failed: ${res.status}`)
	return res.json()
}

export async function getGoalPlan(sessionId: string): Promise<GoalPlanRes> {
	const res = await fetch(`${API_BASE}/api/goal/plan`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify({ session_id: sessionId }),
	})
	if (!res.ok) throw new Error(`Goal plan failed: ${res.status}`)
	return res.json()
}

export async function setDailyXPGoal(goal: number): Promise<void> {
	const res = await fetch(`${API_BASE}/api/goals/xp`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify({ goal }),
	})
	if (!res.ok) throw new Error(`Set goal failed: ${res.status}`)
}

export async function getWeaknesses(): Promise<WeaknessRes> {
  const res = await fetch(`${API_BASE}/api/weaknesses`, {
    headers: { ...getAuthHeaders() },
  })
  if (!res.ok) return { by_domain: {} }
  return validateResponse(WeaknessResSchema, await res.json(), 'getWeaknesses') as WeaknessRes
}

export interface UserSettings {
	show_timer?: boolean
	pause_until?: string | null
}

export async function getSettings(): Promise<UserSettings> {
	const res = await fetch(`${API_BASE}/api/settings`, {
		headers: { ...getAuthHeaders() },
	})
	if (!res.ok) return {}
	return res.json()
}

export async function updateSettings(settings: UserSettings): Promise<void> {
	const res = await fetch(`${API_BASE}/api/settings`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify(settings),
	})
	if (!res.ok) throw new Error(`Update settings failed: ${res.status}`)
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
export async function getLessonBody(title: string): Promise<string> {
	const res = await fetch(`${API_BASE}/api/lessons/body?title=${encodeURIComponent(title)}`, {
		cache: 'no-store',
	})
	if (!res.ok) throw new Error(`Lesson body fetch failed: ${res.status}`)
	const data = (await res.json()) as { title: string; body: string }
	return data.body
}

export interface LessonsRes {
	lessons: Record<string, LessonInfo[]>
}

export async function getLessons(studentId?: string): Promise<LessonsRes> {
	const url = studentId ? `${API_BASE}/api/lessons?student_id=${encodeURIComponent(studentId)}` : `${API_BASE}/api/lessons`
	const res = await fetch(url)
	if (!res.ok) return { lessons: {} }
	return validateResponse(LessonsResSchema, await res.json(), 'getLessons') as LessonsRes
}

export interface ConceptDetailRes {
	concept: {
		id: string
		label: string
		domain: string
		subdomain: string
	}
	lesson?: {
		title: string
		body: string
		concepts: string[]
	}
	prerequisites: {
		id: string
		label: string
		status: string
		mastery_pct: number
	}[]
	dependents?: {
		id: string
		label: string
		status: string
		mastery_pct: number
	}[]
	unlocked: boolean
	progress?: {
		status: string
		streak: number
		required_streak: number
		mastery_pct: number
	}
}

export async function getConceptDetail(conceptId: string): Promise<ConceptDetailRes> {
	const res = await fetch(`${API_BASE}/api/concepts/${encodeURIComponent(conceptId)}`, {
		cache: 'no-store',
	})
	if (!res.ok) throw new Error(`Concept detail fetch failed: ${res.status}`)
	return validateResponse(ConceptDetailResSchema, await res.json(), 'getConceptDetail') as ConceptDetailRes
}

export interface PracticeQuestion {
	question: string
	answer: string
	explanation: string
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
}

export interface LessonKpsRes {
	concept_id: string
	kps: KpInfo[]
	diagram?: string
}

export async function getLessonPractice(conceptId: string, count = 5): Promise<LessonPracticeRes> {
	const res = await fetch(`${API_BASE}/api/lessons/${encodeURIComponent(conceptId)}/practice?count=${count}`)
	if (!res.ok) throw new Error(`Lesson practice fetch failed: ${res.status}`)
	return res.json()
}

export async function getLessonKPs(conceptId: string): Promise<LessonKpsRes> {
	const res = await fetch(`${API_BASE}/api/lessons/${encodeURIComponent(conceptId)}/kp`, { cache: 'no-store' })
	if (!res.ok) return { concept_id: conceptId, kps: [] }
	return res.json()
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
}

export async function submitStudyAnswer(conceptId: string, answer: string, expected: string, elapsed: number): Promise<StudyAnswerRes> {
	const { getGuestId } = await import('./auth')
	const headers: Record<string, string> = { 'Content-Type': 'application/json', ...getAuthHeaders() }
	const body: Record<string, unknown> = { concept_id: conceptId, answer, expected, elapsed }
	const guestId = getGuestId()
	if (guestId && !headers.Authorization) {
		body.student_id = guestId
	}
	const res = await fetch(`${API_BASE}/api/study/answer`, {
		method: 'POST',
		headers,
		body: JSON.stringify(body),
	})
	if (!res.ok) throw new Error(`Study answer failed: ${res.status}`)
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
	const headers: Record<string, string> = { ...getAuthHeaders() }
	const res = await fetch(`${API_BASE}/api/activity?days=${days}`, { headers })
	if (!res.ok) return []
	return res.json()
}

export async function getDueReviews(): Promise<DueReviewsRes> {
	const headers: Record<string, string> = { ...getAuthHeaders() }
	const res = await fetch(`${API_BASE}/api/reviews/due`, { headers })
	if (!res.ok) return { count: 0 }
	return res.json()
}

export async function startReviewSession(): Promise<StartSessionRes> {
	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...getAuthHeaders(),
	}
	const res = await fetch(`${API_BASE}/api/reviews/session`, {
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
	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...getAuthHeaders(),
	}
	const res = await fetch(`${API_BASE}/api/reviews/answer`, {
		method: 'POST',
		headers,
		body: JSON.stringify({ session_id: sessionID, attempt_id: attemptID, answer, elapsed }),
	})
	if (!res.ok) await throwWithResponse(res, `Review answer submit failed: ${res.status}`)
	return res.json()
}

// Quiz (actionable every 150 XP, 80% difficulty, guest unlimited retake, own grading path)
export interface QuizStartRes {
	session_id: string
	student_id?: string
	concept_id?: string
	concept_name?: string
	question?: string
	done?: boolean
}
export interface QuizAnswerRes {
	done: boolean
	correct?: boolean
	feedback?: string
	xp?: number
	new_status?: string
	concept_id?: string
	concept_name?: string
	question?: string
}

export async function startQuizSession(): Promise<QuizStartRes> {
	const { getGuestId } = await import('./auth')
	const headers: Record<string, string> = { 'Content-Type': 'application/json', ...getAuthHeaders() }
	const body: Record<string, unknown> = {}
	const guestId = getGuestId()
	if (guestId && !headers.Authorization) body.student_id = guestId
	const res = await fetch(`${API_BASE}/api/quiz/session`, {
		method: 'POST',
		headers,
		body: JSON.stringify(body),
	})
	if (!res.ok) throw new Error(`Quiz session start failed: ${res.status}`)
	return res.json()
}
export async function submitQuizAnswer(sessionId: string, conceptId: string, answer: string, elapsed: number): Promise<QuizAnswerRes> {
	const { getGuestId } = await import('./auth')
	const headers: Record<string, string> = { 'Content-Type': 'application/json', ...getAuthHeaders() }
	const body: Record<string, unknown> = { session_id: sessionId, concept_id: conceptId, answer, elapsed }
	const guestId = getGuestId()
	if (guestId && !headers.Authorization) body.student_id = guestId
	const res = await fetch(`${API_BASE}/api/quiz/answer`, {
		method: 'POST',
		headers,
		body: JSON.stringify(body),
	})
	if (!res.ok) await throwWithResponse(res, `Quiz answer failed: ${res.status}`)
	return res.json()
}

export async function validateToken(): Promise<{ valid: boolean; student_id: string }> {
	const headers = getAuthHeaders()
	if (!headers.Authorization) return { valid: false, student_id: '' }
	const res = await fetch(`${API_BASE}/api/me`, { headers })
	if (!res.ok) return { valid: false, student_id: '' }
	try {
		const me = await res.json()
		return { valid: true, student_id: typeof me?.student_id === 'string' ? me.student_id : '' }
	} catch {
		return { valid: true, student_id: '' }
	}
}
