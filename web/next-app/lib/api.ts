const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

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

export async function healthCheck(): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/api/health`)
    return res.ok
  } catch {
    return false
  }
}

export async function signup(name: string, username: string, password: string): Promise<{ token: string; student_id: string; name: string }> {
  const res = await fetch(`${API_BASE}/api/auth/signup`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, username, password }),
  })
  if (!res.ok) throw new Error('Signup failed')
  return res.json()
}

export async function login(username: string, password: string): Promise<{ token: string; student_id: string; name: string }> {
  const res = await fetch(`${API_BASE}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })
  if (!res.ok) throw new Error('Invalid credentials')
  return res.json()
}
