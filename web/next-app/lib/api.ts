const API_BASE = process.env.NEXT_PUBLIC_API_URL || ''

import { getAuthHeaders } from './auth'
import { z } from 'zod'

const QuestionSchema = z.object({
  concept_id: z.string(),
  concept_name: z.string(),
  question: z.string(),
  is_review: z.boolean(),
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
})

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function validateResponse(_schema: z.ZodTypeAny, data: unknown, _name: string) {
  const result = _schema.safeParse(data)
  if (!result.success) {
    console.error(`API validation error (${_name}):`, result.error.issues)
  }
  return data
}

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
}

export interface ConceptProgress {
  concept_id: string
  status: string
  streak: number
  best_streak: number
  avg_response_time: number
  attempts: number
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

export async function submitAnswer(
  sessionID: string,
  answer: string,
  elapsed: number,
): Promise<AnswerRes> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...getAuthHeaders(),
  }
  const res = await fetch(`${API_BASE}/api/answer`, {
    method: 'POST',
    headers,
    body: JSON.stringify({ session_id: sessionID, answer, elapsed }),
  })
  if (!res.ok) throw new Error(`Answer submit failed: ${res.status}`)
  return validateResponse(AnswerResSchema, await res.json(), 'submitAnswer') as AnswerRes
}

export async function getProgress(studentID: string): Promise<Record<string, ConceptProgress>> {
  const res = await fetch(`${API_BASE}/api/progress/${studentID}`)
  if (!res.ok) return {}
  return res.json()
}

export async function getScores(studentID: string): Promise<Scores> {
  const res = await fetch(`${API_BASE}/api/scores/${studentID}`)
  if (!res.ok) throw new Error(`Scores fetch failed: ${res.status}`)
  return validateResponse(ScoresSchema, await res.json(), 'getScores') as Scores
}

export async function getGraph(): Promise<GraphRes> {
  const res = await fetch(`${API_BASE}/api/graph`)
  if (!res.ok) throw new Error(`Graph fetch failed: ${res.status}`)
  return res.json()
}

export async function getLeaderboard(): Promise<LeaderboardEntry[]> {
  const res = await fetch(`${API_BASE}/api/leaderboard`)
  if (!res.ok) return []
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
	return res.json()
}

export interface UserSettings {
	show_timer?: boolean
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

export interface ConceptProgress {
    status: string
    mastery: number
    streak: number
}

export interface PrereqInfo {
    id: string
    label: string
    status: string
    mastery_pct: number
}

export interface LessonInfo {
    title: string
    body: string
    concepts: string[]
    progress?: Record<string, ConceptProgress>
    prerequisites?: PrereqInfo[]
    dependents?: PrereqInfo[]
}

export interface LessonsRes {
	lessons: Record<string, LessonInfo[]>
}

export async function getLessons(studentId?: string): Promise<LessonsRes> {
	const url = studentId ? `${API_BASE}/api/lessons?student_id=${encodeURIComponent(studentId)}` : `${API_BASE}/api/lessons`
	const res = await fetch(url)
	if (!res.ok) return { lessons: {} }
	return res.json()
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
	const res = await fetch(`${API_BASE}/api/concepts/${encodeURIComponent(conceptId)}`)
	if (!res.ok) throw new Error(`Concept detail fetch failed: ${res.status}`)
	return res.json()
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

export async function getLessonPractice(conceptId: string, count = 5): Promise<LessonPracticeRes> {
	const res = await fetch(`${API_BASE}/api/lessons/${encodeURIComponent(conceptId)}/practice?count=${count}`)
	if (!res.ok) throw new Error(`Lesson practice fetch failed: ${res.status}`)
	return res.json()
}

export interface DueReviewsRes {
	count: number
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
): Promise<AnswerRes> {
	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...getAuthHeaders(),
	}
	const res = await fetch(`${API_BASE}/api/reviews/answer`, {
		method: 'POST',
		headers,
		body: JSON.stringify({ session_id: sessionID, answer, elapsed }),
	})
	if (!res.ok) throw new Error(`Review answer submit failed: ${res.status}`)
	return res.json()
}
