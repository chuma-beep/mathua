'use client'

import { useId, useRef, useState, type ChangeEvent, type FormEvent } from 'react'
import { submitReport } from '../lib/api'

const MAX = 2000
const MIN = 10

// The contact box.
//
// It is one textarea and one button because the destination already exists: a
// message is stored as a report whose kind is "message", so it appears in the same
// audited queue as a wrong question, with the same close button and the same record
// of who decided what. A separate inbox would have been a second page to remember;
// this one is the page the administrator already opens whenever a question is
// reported.
//
// It works signed out. A locked-out learner is the main reason the box exists, and
// they cannot file a question report that names a concept they cannot reach. The
// server's per-reporter limiter and the per-IP write limiter cover the abuse that
// follows from accepting anonymous text.
export default function ContactBox() {
  const fieldId = useId()
  const [text, setText] = useState('')
  const [image, setImage] = useState<File | null>(null)
  const [imageError, setImageError] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [sent, setSent] = useState(false)
  const fileRef = useRef<HTMLInputElement>(null)

  const trimmed = text.trim()
  const canSend = !busy && !imageError

  function pickImage(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0] ?? null
    setImageError('')
    if (!file) {
      setImage(null)
      return
    }
    if (file.size > 4 * 1024 * 1024) {
      setImage(null)
      setImageError('Screenshot must be 4MB or smaller.')
      event.target.value = ''
      return
    }
    if (!['image/png', 'image/jpeg', 'image/webp', 'image/gif'].includes(file.type)) {
      setImage(null)
      setImageError('Screenshot must be PNG, JPEG, GIF, or WebP.')
      event.target.value = ''
      return
    }
    setImage(file)
  }

  async function send(event: FormEvent) {
    event.preventDefault()
    if (!canSend) return
    if (trimmed.length < MIN) {
      setError(`Add at least ${MIN} characters so the message is actionable.`)
      return
    }
    if (imageError) return
    setBusy(true)
    setError('')
    try {
      await submitReport({
        kind: 'message',
        source: 'contact',
        reason: 'other',
        detail: trimmed,
        image: image ?? undefined,
      })
      setSent(true)
      setText('')
      setImage(null)
      setImageError('')
      if (fileRef.current) fileRef.current.value = ''
    } catch (err) {
      const message = err instanceof Error ? err.message : ''
      setError(
        message.includes('429')
          ? 'Too many messages in the last hour. Try again later.'
          : 'The message did not send. Try again, or use the X link below.',
      )
    } finally {
      setBusy(false)
    }
  }

  if (sent) {
    return (
      <div className="border border-mathua-border bg-mathua-surface p-4">
        <p role="status" className="font-mono text-xs text-mathua-primary">
          Message received. It is in the queue and will be read.
        </p>
        <button
          type="button"
          onClick={() => setSent(false)}
          className="mt-3 border border-mathua-border px-3 h-9 font-mono text-xs text-mathua-muted hover:text-mathua-blue transition-colors"
        >
          Write another
        </button>
      </div>
    )
  }

  return (
    <form onSubmit={send} className="border border-mathua-border bg-mathua-surface p-4">
      <label htmlFor={fieldId} className="font-mono text-[10px] uppercase tracking-wider text-mathua-secondary">
        Your message
      </label>
      <p id={`${fieldId}-hint`} className="mt-1 font-mono text-[10px] leading-relaxed text-mathua-muted">
        Say what happened and where. Do not include a password, yours or anyone else&rsquo;s.
        If you need an answer, message me on X instead.
      </p>
      <textarea
        id={fieldId}
        value={text}
        onChange={(event) => {
          setText(event.target.value)
          if (error) setError('')
        }}
        maxLength={MAX}
        rows={6}
        required
        aria-describedby={`${fieldId}-hint ${fieldId}-count`}
        className="mt-2 w-full resize-y border border-mathua-border bg-mathua-bg px-2 py-2 font-mono text-xs text-mathua-primary outline-none placeholder:text-mathua-muted focus:border-mathua-blue"
        placeholder="What went wrong, and on which page."
      />
      <div className="mt-1 flex items-baseline justify-between gap-3 font-mono text-[10px] text-mathua-muted">
        <span id={`${fieldId}-count`}>
          {text.length}/{MAX}
        </span>
        <span>{trimmed.length > 0 && trimmed.length < MIN ? `at least ${MIN} characters` : ''}</span>
      </div>
      <div className="mt-3">
        <label htmlFor={`${fieldId}-shot`} className="font-mono text-[10px] uppercase tracking-wider text-mathua-secondary">
          Screenshot (optional)
        </label>
        <input
          ref={fileRef}
          id={`${fieldId}-shot`}
          type="file"
          accept="image/png,image/jpeg,image/webp,image/gif"
          onChange={pickImage}
          className="mt-1 block w-full font-mono text-[10px] text-mathua-muted file:mr-3 file:border file:border-mathua-border file:bg-mathua-bg file:px-2 file:py-1 file:font-mono file:text-[10px] file:text-mathua-secondary"
        />
        <p className="mt-1 font-mono text-[10px] text-mathua-muted">
          PNG, JPEG, GIF or WebP, up to 4MB. It is visible only to the moderator.
        </p>
        {image && (
          <p className="mt-1 font-mono text-[10px] text-mathua-secondary">
            Attached: {image.name}
          </p>
        )}
      </div>
      {(error || imageError) && (
        <p role="alert" className="mt-2 font-mono text-[11px] text-mathua-red">
          {error || imageError}
        </p>
      )}
      <button
        type="submit"
        disabled={!canSend}
        className="mt-3 border border-mathua-blue px-3 h-9 font-mono text-xs text-mathua-blue hover:bg-mathua-blue-faint transition-colors disabled:border-mathua-border disabled:text-mathua-muted disabled:hover:bg-transparent"
      >
        {busy ? 'Sending...' : 'Send message'}
      </button>
    </form>
  )
}
