'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import dynamic from 'next/dynamic'
import { askAssistant, getConfig, type AssistantPart, type AssistantTurn } from '../lib/api'
import { splitMath } from '../lib/assistantMath'
import { ArrowUp, Bell, Close, MessagesSquare } from './icons'
import { useAnnouncementsFeed } from './announcements/useAnnouncementsFeed'

// KaTeX is a large dependency and the dock is mounted in the root layout, so it
// would otherwise be pulled into every route's bundle. Loading it through a
// dynamic import means a reply with no mathematics, an orientation answer, which
// is most of them, never downloads it.
const AssistantMath = dynamic(() => import('./AssistantMath'), { ssr: false })

// AssistantText renders one prose segment of a reply. Mathematics inside it is
// split out and rendered; everything else is the characters the model wrote,
// including any Markdown it used despite being asked not to.
function AssistantText({ text }: { text: string }) {
	const parts = splitMath(text)
	if (parts.length === 1 && parts[0].kind === 'text') return <>{text}</>
	return (
		<>
			{parts.map((p, i) =>
				p.kind === 'math' ? (
					<AssistantMath key={i} tex={p.tex} display={p.display} />
				) : (
					<span key={i}>{p.text}</span>
				),
			)}
		</>
	)
}

// The product assistant, as a dock in the corner.
//
// It is mounted once in the root layout rather than per page, so the conversation
// survives navigation within the app: a visitor who asks "where do I start" and
// then follows a link can come back to the answer.
//
// Three constraints are deliberate.
//
// **It always renders, and the server decides what it can do.** `ai_enabled` is
// the one fact `/api/config` exposes (ADR-055). Hiding the button when AI is off
// made a switched-off feature indistinguishable from a missing one: the same
// page looked identical whether the assistant was unimplemented, unconfigured or
// broken. With a provider the panel answers; without one it says so once and
// offers the same destinations as links, which is what a visitor wanted anyway.
// It never pretends to answer, and it never leaves an empty box behind.
//
// **Replies are text, never HTML.** The server returns plain prose and a list of
// links whose hrefs it chose itself, so there is no markup to sanitise and no URL
// the model could have written. The links are rendered as Next `Link`s from that
// list, never parsed out of the reply.
//
// **A failed request is not an empty conversation.** A spent quota and an outage
// produce different copy, because "try again later" is wrong advice for a limit
// that resets tomorrow.

type Bubble = AssistantTurn & { parts?: AssistantPart[] }

const SUGGESTIONS = ['What is Mathua?', 'Where should I start?']

// AREAS is what the panel offers when the assistant cannot answer. These are
// routes the app owns regardless of AI, so the dock is useful in a deployment
// with no provider, which is every fresh checkout until a key is set.
const AREAS = [
	{ label: 'Start learning', href: '/learn' },
	{ label: 'Take the diagnostic', href: '/diagnose' },
	{ label: 'Plan a goal', href: '/plan' },
	{ label: 'Concept graph', href: '/graph' },
	{ label: 'Glossary', href: '/glossary' },
]

type Status = 'loading' | 'on' | 'off'

// Consent before the first message.
//
// The dock is the one surface where a visitor's own words leave the server, so it is the one
// place that has to say so before they type. The gate is deliberately small: a notice, one
// checkbox and a Start button, inside the panel they have already opened. Nothing is sent
// until it is accepted, and an unaccepted panel has no input at all rather than a dead one.
//
// It is an attestation, not a check. No account is required, because registration is the
// easiest thing for a minor to complete and would have made signing up the loophole. Email
// ownership shows that someone can receive mail, not how old they are.
//
// The key is versioned: if what leaves the server ever changes, bump the version and every
// visitor is asked again, instead of being covered by an acceptance of the old wording.
export const ASSISTANT_CONSENT_KEY = 'mathua.assistant-consent'
export const ASSISTANT_CONSENT_VERSION = 'v1'

// 'unknown' is the first client render, before storage has been read. It is not 'needed',
// because that would flash the notice at every visitor who had already accepted it, and it
// is not 'given', because that would render an input nobody has agreed to yet.
type Consent = 'unknown' | 'given' | 'needed'

function failureCopy(e: unknown): string {
	const message = e instanceof Error ? e.message : ''
	if (message.includes('429')) {
		return 'That is today\u2019s limit for the assistant. The rest of Mathua works as usual.'
	}
	return 'The assistant is unavailable right now. Everything else still works.'
}

export default function AssistantDock() {
	const [status, setStatus] = useState<Status>('loading')
	const [open, setOpen] = useState(false)
	const [consent, setConsent] = useState<Consent>('unknown')
	const [consentChecked, setConsentChecked] = useState(false)
	const [turns, setTurns] = useState<Bubble[]>([])
	const [draft, setDraft] = useState('')
	const [busy, setBusy] = useState(false)
	const [error, setError] = useState('')
	const inputRef = useRef<HTMLInputElement>(null)
	const listRef = useRef<HTMLDivElement>(null)

	useEffect(() => {
		let alive = true
		getConfig()
			.then(c => {
				if (alive) setStatus(c?.ai_enabled ? 'on' : 'off')
			})
			.catch(() => {
				// A config read that fails is not a reason to show an assistant that
				// cannot answer, but it is also not a reason to hide a control the
				// page now always offers; the panel says it is unavailable instead.
				if (alive) setStatus('off')
			})
		return () => {
			alive = false
		}
	}, [])

	// Read the acceptance once, on mount. Storage can be unavailable, a private window,
	// blocked cookies and the safe failure is to show the notice again, not to assume it
	// was given.
	useEffect(() => {
		let stored = ''
		try {
			stored = window.localStorage.getItem(ASSISTANT_CONSENT_KEY) ?? ''
		} catch {
			stored = ''
		}
		setConsent(stored === ASSISTANT_CONSENT_VERSION ? 'given' : 'needed')
	}, [])

	const acceptConsent = useCallback(() => {
		try {
			window.localStorage.setItem(ASSISTANT_CONSENT_KEY, ASSISTANT_CONSENT_VERSION)
		} catch {
			// Not persisting only means this visitor is asked again next time. The answer
			// they just gave is still the answer, so the panel opens either way.
		}
		setConsent('given')
	}, [])

	useEffect(() => {
		if (open) inputRef.current?.focus()
	}, [open])

	useEffect(() => {
		const list = listRef.current
		if (list) list.scrollTop = list.scrollHeight
	}, [turns, busy, error])

	useEffect(() => {
		if (!open) return
		const onKey = (e: KeyboardEvent) => {
			if (e.key === 'Escape') setOpen(false)
		}
		window.addEventListener('keydown', onKey)
		return () => window.removeEventListener('keydown', onKey)
	}, [open])

	const send = useCallback(
		async (text: string) => {
			const message = text.trim()
			if (!message || busy || status !== 'on' || consent !== 'given') return
			setError('')
			setDraft('')
			// Replay what both sides said, so a clarifying question can be answered.
			// The server keeps only these two roles and bounds the count.
			const history: AssistantTurn[] = turns.map(t => ({ role: t.role, content: t.content }))
			setTurns(t => [...t, { role: 'user', content: message }])
			setBusy(true)
			try {
				const res = await askAssistant({ message, history })
				setTurns(t => [...t, { role: 'assistant', content: res.reply, parts: res.segments }])
			} catch (e) {
				setError(failureCopy(e))
			} finally {
				setBusy(false)
			}
		},
		[busy, status, turns, consent],
	)

	const announcements = useAnnouncementsFeed()
	const unreadAnnouncement = announcements.newestUnread

	return (
		// The mobile offset carries the safe-area inset for the same reason the tab
		// bar does: on a phone with a home indicator, a plain 5rem lands the trigger
		// on top of the tabs. On desktop the tab bar is gone, so the inset is too.
		<div className="fixed right-4 bottom-[calc(80px+env(safe-area-inset-bottom))] z-40 flex flex-col items-end gap-2 print:hidden lg:bottom-6">
			{open && (
				<section
					role="dialog"
					aria-label="Mathua assistant"
					className="flex h-[min(28rem,70dvh)] w-[min(22rem,calc(100vw-2rem))] flex-col overflow-hidden rounded-none border border-mathua-border bg-mathua-bg shadow-xl"
				>
					<header className="flex items-center justify-between border-b border-mathua-border px-3 py-2">
						<p className="font-mono text-[11px] tracking-wider text-mathua-muted uppercase">Assistant</p>
						<button
							type="button"
							onClick={() => setOpen(false)}
							className="rounded-none p-1 text-mathua-muted hover:bg-mathua-surface-elevated hover:text-mathua-primary"
							aria-label="Close assistant"
						>
							<Close size={14} />
						</button>
					</header>

					{status !== 'on' && (
						<div className="flex-1 space-y-3 overflow-y-auto px-3 py-3">
							<p className="text-[13px] leading-relaxed text-mathua-muted">
								{status === 'loading'
									? 'Connecting to the assistant…'
									: 'The assistant is unavailable right now. These go to the same places:'}
							</p>
							{status === 'off' && (
								<div className="flex flex-wrap gap-2">
									{AREAS.map(a => (
										<Link
											key={a.href}
											href={a.href}
											onClick={() => setOpen(false)}
											className="rounded-none border border-mathua-border px-2 py-1 font-mono text-[11px] text-mathua-blue hover:border-mathua-blue hover:bg-mathua-surface-elevated"
										>
											{a.label}
										</Link>
									))}
								</div>
							)}
						</div>
					)}

					{status === 'on' && consent === 'needed' && (
						<div className="flex-1 space-y-3 overflow-y-auto px-3 py-3">
							<p className="font-mono text-[11px] tracking-wider text-mathua-muted uppercase">
								Before you send a message
							</p>
							<p className="text-[13px] leading-relaxed text-mathua-primary">
								Your text goes to an AI provider so it can answer. Do not type names, email
								addresses, school details or anything about another person.
							</p>
							<p className="text-[12px] leading-relaxed text-mathua-muted">
								Your answers, progress and account details are not sent. The message is not
								stored after it is answered. Replies can be wrong.
							</p>
							<label className="flex cursor-pointer items-start gap-2 text-[12px] leading-relaxed text-mathua-primary">
								<input
									type="checkbox"
									checked={consentChecked}
									onChange={e => setConsentChecked(e.target.checked)}
									className="mt-0.5 accent-mathua-blue"
								/>
								<span>I am 13 or older, and I understand this.</span>
							</label>
							<div className="flex items-center gap-2">
								<button
									type="button"
									onClick={acceptConsent}
									disabled={!consentChecked}
									className="rounded-none border border-mathua-blue bg-mathua-blue-faint px-3 py-1.5 font-mono text-[11px] text-mathua-blue disabled:border-mathua-border disabled:bg-transparent disabled:text-mathua-muted"
								>
									Start
								</button>
								<Link
									href="/privacy"
									onClick={() => setOpen(false)}
									className="font-mono text-[11px] text-mathua-muted underline decoration-mathua-border underline-offset-2 hover:text-mathua-blue"
								>
									Privacy notice
								</Link>
							</div>
						</div>
					)}

					{status === 'on' && consent === 'given' && (
						<>
					<div ref={listRef} aria-live="polite" className="flex-1 space-y-3 overflow-y-auto px-3 py-3">
						{turns.length === 0 && (
							<div className="space-y-2">
								<p className="text-[13px] leading-relaxed text-mathua-muted">
									Ask what Mathua is or where to begin. Answers are short, and you can follow the links.
								</p>
								<div className="flex flex-wrap gap-2">
									{SUGGESTIONS.map(s => (
										<button
											key={s}
											type="button"
											onClick={() => void send(s)}
											className="rounded-none border border-mathua-border px-2 py-1 font-mono text-[11px] text-mathua-blue hover:border-mathua-blue hover:bg-mathua-surface-elevated"
										>
											{s}
										</button>
									))}
								</div>
							</div>
						)}

						{turns.map((t, i) => (
							<div key={i} className={t.role === 'user' ? 'flex justify-end' : 'flex justify-start'}>
								<div
									className={[
										'max-w-[85%] rounded-none px-3 py-2 text-[13px] leading-relaxed whitespace-pre-wrap',
										t.role === 'user'
											? 'bg-mathua-surface-elevated text-mathua-primary'
											: 'border border-mathua-border text-mathua-primary',
									].join(' ')}
								>
									{t.parts && t.parts.length > 0
										? t.parts.map((p, j) =>
												p.link ? (
													<Link
														key={j}
														href={p.link.href}
														onClick={() => setOpen(false)}
														className="text-mathua-blue underline decoration-mathua-border underline-offset-2 hover:decoration-mathua-blue"
													>
														{p.link.label}
													</Link>
												) : (
													<AssistantText key={j} text={p.text ?? ''} />
												),
											)
										: t.content}
								</div>
							</div>
						))}

						{busy && <p className="font-mono text-[11px] text-mathua-muted">Thinking…</p>}
						{error && <p className="text-[12px] text-mathua-red">{error}</p>}
					</div>

					<form
						onSubmit={e => {
							e.preventDefault()
							void send(draft)
						}}
						className="flex items-center gap-2 border-t border-mathua-border px-2 py-2"
					>
						<input
							ref={inputRef}
							value={draft}
							onChange={e => setDraft(e.target.value)}
							maxLength={2000}
							placeholder="Ask about Mathua"
							aria-label="Your question"
							className="min-w-0 flex-1 bg-transparent px-1 text-[13px] text-mathua-primary outline-none placeholder:text-mathua-muted"
						/>
						<button
							type="submit"
							disabled={busy || !draft.trim()}
							aria-label="Send"
							className="rounded-none p-1.5 text-mathua-blue disabled:text-mathua-muted"
						>
							<ArrowUp size={16} />
						</button>
					</form>
						</>
					)}
				</section>
			)}

			<button
				type="button"
				onClick={() => setOpen(v => !v)}
				aria-expanded={open}
				aria-label={open ? 'Close assistant' : 'Open assistant'}
				className="flex h-11 w-11 items-center justify-center rounded-none text-mathua-primary hover:text-mathua-blue"
			>
				{open ? <Close size={18} /> : <MessagesSquare size={18} />}
			</button>

			{unreadAnnouncement ? (
				<Link
					href="/announcements"
					onClick={() => void announcements.markRead(unreadAnnouncement.id)}
					aria-label={`Unread announcement: ${unreadAnnouncement.title}`}
					title={unreadAnnouncement.title}
					className="relative flex h-11 w-11 items-center justify-center rounded-none text-mathua-muted hover:text-mathua-blue"
				>
					<Bell size={18} />
					<span className="absolute right-2 top-2 h-2 w-2 rounded-full bg-mathua-red" aria-hidden="true" />
				</Link>
			) : null}
		</div>
	)
}
