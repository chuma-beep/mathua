'use client'

export interface BriefingItem {
  title: string
  detail: string
}

// Pre-test expectations, MA-grounded: diagnostics run 25–45 adaptive
// questions (~30–45 min per engine bounds), quizzes are timed + closed-book.
// Stated once up front so the test itself stays distraction-free.
export const DIAGNOSTIC_BRIEFING: BriefingItem[] = [
  {
    title: '25–45 questions · about 30–45 minutes',
    detail: 'The test adapts to you. Clear answers mean fewer questions.',
  },
  {
    title: 'Answer from what you know — no searching',
    detail: 'Honest answers put you at the right starting point. Looked-up answers put you in the wrong place.',
  },
  {
    title: 'Speed counts',
    detail: 'Fast correct answers show mastery and shorten the test.',
  },
  {
    title: 'Pause anytime',
    detail: 'Leave and pick up where you stopped.',
  },
]

export const QUIZ_BRIEFING: BriefingItem[] = [
  {
    title: 'Up to 5 questions · timed · closed book',
    detail: 'Each question has its own timer. Keep books and notes shut.',
  },
  {
    title: 'Covers what you just learned',
    detail: 'Questions come from your recent lessons.',
  },
  {
    title: 'Miss one and review it on the spot',
    detail: 'A wrong answer opens a short review right away.',
  },
  {
    title: '20 XP for each correct answer · retake anytime',
    detail: 'No penalty for retaking.',
  },
]

export default function BriefingCard({ eyebrow, items }: { eyebrow: string; items: BriefingItem[] }) {
  return (
    <div className="border border-mathua-border bg-mathua-surface p-4 text-left">
      <p className="font-mono text-[10px] uppercase tracking-wider text-mathua-muted mb-3">
        {eyebrow}
      </p>
      <ul className="space-y-3">
        {items.map(item => (
          <li key={item.title} className="min-w-0">
            <p className="font-mono text-xs text-mathua-primary">{item.title}</p>
            <p className="font-mono text-[11px] text-mathua-muted mt-0.5 break-words [overflow-wrap:anywhere]">
              {item.detail}
            </p>
          </li>
        ))}
      </ul>
    </div>
  )
}
