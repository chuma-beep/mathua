'use client'

import { useEffect, useRef } from 'react'

/**
 * A confirmation step for a privileged action.
 *
 * Promotion and removal of an administrator access the database, take effect immediately, and
 * are recorded in an immutable audit log — none of which a mis-click can be undone by. Two
 * details in here are deliberate rather than decorative:
 *
 *  - **Nothing happens on the first click.** The dialog is the action. A row of one-tap
 *    privileged buttons is a list of things a learner or an administrator can change by reaching
 *    for the wrong one, and the audit log means the mistake is permanent even when the role is
 *    restored.
 *  - **Focus starts on Cancel**, and Escape cancels. The safe choice is the one under the
 *    thumb and the one the keyboard reaches first. Auto-focusing the confirming button would make
 *    a stray Enter key with the dialog already open enough to grant administrator access.
 *
 * This is a courtesy, not a control. The server authorizes and validates regardless of what the
 * dialog was asked, which is why nothing here is described as protection.
 */

interface ConfirmDialogProps {
  open: boolean
  title: string
  /** Names the account and the consequence. Plain text, so it reads aloud and prints. */
  children: React.ReactNode
  confirmLabel: string
  /** `destructive` tints the confirming button red, for access removal. */
  tone?: 'default' | 'destructive'
  busy?: boolean
  onConfirm: () => void
  onCancel: () => void
}

export default function ConfirmDialog({
  open,
  title,
  children,
  confirmLabel,
  tone = 'default',
  busy = false,
  onConfirm,
  onCancel,
}: ConfirmDialogProps) {
  const cancelRef = useRef<HTMLButtonElement>(null)

  useEffect(() => {
    if (!open) return
    cancelRef.current?.focus()
    const onKey = (e: KeyboardEvent) => {
      // Escape cancels. Tab is left alone: the dialog has two controls, so the browser's own
      // wrap-around is already correct, and trapping focus here would be more machinery than
      // two buttons can justify.
      if (e.key === 'Escape' && !busy) onCancel()
    }
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  }, [open, busy, onCancel])

  if (!open) return null

  const confirmClass = tone === 'destructive'
    ? 'border border-mathua-red text-mathua-red hover:bg-mathua-red-faint'
    : 'border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint'

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4"
      style={{ background: 'color-mix(in oklab, var(--bg) 75%, transparent)' }}
    >
      {/* The backdrop is a button so it is reachable and announced, and so a click outside is an
          explicit cancel rather than a dead zone that silently swallows the click. */}
      <button
        type="button"
        aria-label="Cancel"
        tabIndex={-1}
        onClick={() => { if (!busy) onCancel() }}
        className="absolute inset-0 cursor-default"
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="confirm-dialog-title"
        className="relative w-full max-w-[440px] border border-mathua-border bg-mathua-surface p-5"
      >
        <h2 id="confirm-dialog-title" className="font-serif text-lg text-mathua-primary">
          {title}
        </h2>
        <div className="mt-2 font-mono text-xs text-mathua-secondary space-y-1">
          {children}
        </div>
        <div className="mt-5 flex flex-wrap justify-end gap-2">
          <button
            ref={cancelRef}
            type="button"
            onClick={onCancel}
            disabled={busy}
            className="border border-mathua-border px-4 h-11 font-mono text-xs text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue disabled:opacity-40"
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={onConfirm}
            disabled={busy}
            className={`px-4 h-11 font-mono text-xs disabled:opacity-40 ${confirmClass}`}
          >
            {busy ? 'Working…' : confirmLabel}
          </button>
        </div>
      </div>
    </div>
  )
}