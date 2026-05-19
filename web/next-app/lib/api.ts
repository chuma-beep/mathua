const API_BASE = process.env.NEXT_PUBLIC_API_URL || ''

import { getAuthHeaders } from './auth'

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
  return res.json()
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
  return res.json()
}

export async function getProgress(studentID: string): Promise<Record<string, ConceptProgress>> {
  const res = await fetch(`${API_BASE}/api/progress/${studentID}`)
  if (!res.ok) return {}
  return res.json()
}

export async function getScores(studentID: string): Promise<Scores> {
  const res = await fetch(`${API_BASE}/api/scores/${studentID}`)
  if (!res.ok) throw new Error(`Scores fetch failed: ${res.status}`)
  return res.json()
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
  return res.json()
}

export async function login(username: string, password: string): Promise<AuthRes> {
	const res = await fetch(`${API_BASE}/api/auth/login`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ username, password }),
	})
	if (!res.ok) throw new Error('Invalid credentials')
	return res.json()
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
	await fetch(`${API_BASE}/api/settings`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
		body: JSON.stringify(settings),
	})
}

export interface LessonInfo {
	title: string
	body: string
	concepts: string[]
}

export interface LessonsRes {
	lessons: Record<string, LessonInfo[]>
}

export async function getLessons(): Promise<LessonsRes> {
	const res = await fetch(`${API_BASE}/api/lessons`)
	if (!res.ok) return { lessons: {} }
	return res.json()
}
