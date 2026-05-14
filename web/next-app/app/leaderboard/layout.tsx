import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Leaderboard \u2014 Mathua',
  description: 'Weekly leaderboard for Mathua — compete, improve, and rise through the ranks.',
}

export default function LeaderboardLayout({ children }: { children: React.ReactNode }) {
  return children
}
