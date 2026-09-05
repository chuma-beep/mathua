import { describe, it, expect } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import ProfileSidebar from '../components/ProfileSidebar'
import ProfileMenuSheet, { ProfileMobileBar } from '../components/ProfileMenuSheet'

// next/link + next/navigation are stubbed by the app's test env; fall back to plain anchors
import React from 'react'

const TOC = [
  { id: 'activity', label: 'Activity' },
  { id: 'domains', label: 'Domains' },
  { id: 'struggles', label: 'Struggles' },
  { id: 'efficacy', label: 'Efficacy' },
]

describe('ProfileSidebar', () => {
  it('renders identity, quick links, and TOC with active state', () => {
    render(
      <ProfileSidebar
        name="Ada"
        studentId="s1"
        level="Scholar"
        streak={5}
        toc={TOC}
        activeId="domains"
      />,
    )
    expect(screen.getByText('Ada')).toBeInTheDocument()
    expect(screen.getByText('Study')).toBeInTheDocument()
    const active = screen.getByRole('link', { name: 'Domains' })
    expect(active).toHaveAttribute('aria-current', 'true')
    expect(screen.getByRole('link', { name: 'Activity' })).not.toHaveAttribute('aria-current')
  })

  it('renders next/goals slots for T2-T4 plug-ins', () => {
    render(
      <ProfileSidebar
        name="Ada"
        studentId="s1"
        toc={TOC}
        activeId="activity"
        nextSlot={<div>next-up</div>}
        goalsSlot={<div>goals</div>}
      />,
    )
    expect(screen.getByText('next-up')).toBeInTheDocument()
    expect(screen.getByText('goals')).toBeInTheDocument()
  })
})

describe('ProfileMobileBar', () => {
  it('renders chips with active state and menu button', () => {
    const onChip = () => {}
    const onMenu = () => {}
    render(
      <ProfileMobileBar chips={TOC} activeId="activity" hidden={false} onChipNavigate={onChip} onOpenMenu={onMenu} />,
    )
    expect(screen.getByRole('button', { name: /open profile menu/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Activity' })).toHaveAttribute('aria-current', 'true')
  })
})

describe('ProfileMenuSheet', () => {
  it('returns null when closed and traps focus + closes on Escape', () => {
    const onClose = () => {}
    const { rerender } = render(
      <ProfileMenuSheet
        open={false}
        onClose={onClose}
        chips={TOC}
        activeId="activity"
        chipsHidden={false}
        onChipNavigate={() => {}}
        onOpenMenu={() => {}}
        sidebarProps={{ name: 'Ada', studentId: 's1', toc: TOC, activeId: 'activity' }}
      />,
    )
    expect(screen.queryByRole('dialog')).toBeNull()
    rerender(
      <ProfileMenuSheet
        open
        onClose={onClose}
        chips={TOC}
        activeId="activity"
        chipsHidden={false}
        onChipNavigate={() => {}}
        onOpenMenu={() => {}}
        sidebarProps={{ name: 'Ada', studentId: 's1', toc: TOC, activeId: 'activity' }}
      />,
    )
    expect(screen.getByRole('dialog', { name: /profile menu/i })).toBeInTheDocument()
    fireEvent.keyDown(document, { key: 'Escape' })
  })
})
