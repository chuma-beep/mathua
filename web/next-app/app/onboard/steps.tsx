'use client'

import type { RefObject } from 'react'
import Link from 'next/link'
import KatexContent from '../../components/KatexContent'
import SectionHeader from '../../components/SectionHeader'
import ProgressBar from '../../components/ProgressBar'
import ReportButton from '../../components/ReportButton'
import BriefingCard, { DIAGNOSTIC_BRIEFING } from '../../components/BriefingCard'
import SubmitErrorBlock, { type SubmitError } from '../../components/SubmitErrorBlock'
import DiagnosticResults from '../../components/DiagnosticResults'
import Loading from '../../components/Loading'
import type { DiagnosticProgress, GoalPlanRes } from '../../lib/api'
import type { AnswerFormat } from '../../lib/answerFormat'
import ChoiceOptions from '../../components/ChoiceOptions'
import { domainLabels, type DomainInfo } from './domains'
import { MathAnswerInput, type MathFocusHandle } from '@/components/math/MathInput'
import { focusClearOfTabs } from '../../lib/focusClearOfTabs'

export function WelcomeStep({
  domains,
  loading,
  hasPaused,
  selectedCount,
  confirming,
  onToggle,
  onSelectAll,
  onStart,
  onResume,
  onBegin,
  onCancel,
  needsName,
  demoName,
  onDemoName,
  onStartDemo,
  onCancelDemo,
}: {
  domains: DomainInfo[]
  loading: boolean
  hasPaused: boolean
  selectedCount: number
  confirming: boolean
  onToggle: (name: string) => void
  onSelectAll: () => void
  onStart: () => void
  onResume: () => void
  onBegin: () => void
  onCancel: () => void
  needsName: boolean
  demoName: string
  onDemoName: (name: string) => void
  onStartDemo: () => void
  onCancelDemo: () => void
}) {
  if (confirming) {
    const selected = domains.filter(d => d.selected)
    const shown = selected.slice(0, 3).map(d => domainLabels(d.name))
    const extra = selected.length > 3 ? ` +${selected.length - 3} more` : ''
    return (
      <div className="max-w-2xl mx-auto px-2">
        <div className="text-center mb-4 mt-6 sm:mt-8 min-w-0">
          <div className="font-mono text-[11px] uppercase text-mathua-muted mb-3">One quick check</div>
          <h1 className="font-serif text-2xl sm:text-3xl font-medium text-mathua-primary px-2">
            Before you begin
          </h1>
          <p className="text-mathua-secondary text-sm mt-3 px-2">
            Starting test on: {shown.join(', ')}{extra} · {selectedCount} concepts in scope
          </p>
          <p className="text-mathua-muted text-xs mt-2 px-2">
            The test adapts as it goes and usually takes about 30–45 minutes.
          </p>
        </div>
        <BriefingCard eyebrow="What to expect" items={DIAGNOSTIC_BRIEFING} />
        <div className="mt-6 flex flex-col sm:flex-row gap-3 justify-center px-2">
          <button
            type="button"
            onClick={onBegin}
            disabled={loading}
            className="border border-mathua-blue bg-mathua-blue text-white hover:opacity-90 rounded-none h-12 min-h-[44px] px-6 sm:px-10 font-medium text-sm disabled:opacity-50"
          >
            {loading ? (<><Loading inline size={13} /> Loading…</>) : 'Begin diagnostic →'}
          </button>
          <button
            type="button"
            onClick={onCancel}
            disabled={loading}
            className="font-mono text-xs text-mathua-muted hover:text-mathua-primary min-h-[44px] px-4 disabled:opacity-50"
          >
            ← Back
          </button>
        </div>
      </div>
    )
  }
  // A guest who has no token yet. The name is not an account: the server makes
  // a student row so a session has an owner, and this screen says plainly that
  // the row cannot be signed back into, rather than letting a guest believe
  // their run was saved.
  if (needsName) {
    return (
      <div className="max-w-2xl mx-auto px-2">
        <div className="text-center mb-4 mt-6 sm:mt-8 min-w-0">
          <div className="font-mono text-[11px] uppercase text-mathua-muted mb-3">Guest run</div>
          <h1 className="font-serif text-2xl sm:text-3xl font-medium text-mathua-primary px-2">
            What should we call you?
          </h1>
          <p className="text-mathua-secondary text-sm mt-3 px-2">
            You are not signed in, so this run is a demo. Mathua tests you against a temporary
            profile with this name. Nothing here can be signed into later, and the result is not
            attached to an account.
          </p>
        </div>
        <div className="border border-mathua-border bg-mathua-surface p-4 sm:p-6">
          <label htmlFor="demo-name" className="block font-mono text-[11px] uppercase text-mathua-muted mb-2">
            Your name
          </label>
          <input
            id="demo-name"
            type="text"
            value={demoName}
            onChange={e => onDemoName(e.target.value)}
            maxLength={40}
            autoComplete="nickname"
            placeholder="e.g. Sam"
            className="w-full h-11 px-3 bg-mathua-bg border border-mathua-border text-mathua-primary text-sm rounded-none focus:outline-none focus:border-mathua-blue"
          />
          <p className="text-mathua-muted text-xs mt-2">
            It only labels the demo. You can sign in later and take the diagnostic on a real account.
          </p>
        </div>
        <div className="mt-6 flex flex-col sm:flex-row gap-3 justify-center px-2">
          <button
            type="button"
            onClick={onStartDemo}
            disabled={loading || demoName.trim() === ''}
            className="border border-mathua-blue bg-mathua-blue text-white hover:opacity-90 rounded-none h-12 min-h-[44px] px-6 sm:px-10 font-medium text-sm disabled:opacity-50"
          >
            {loading ? (<><Loading inline size={13} /> Loading…</>) : 'Run the demo →'}
          </button>
          <button
            type="button"
            onClick={onCancelDemo}
            disabled={loading}
            className="font-mono text-xs text-mathua-muted hover:text-mathua-primary min-h-[44px] px-4 disabled:opacity-50"
          >
            ← Back
          </button>
        </div>
        <p className="text-center mt-4">
          <Link href="/login" className="font-mono text-xs text-mathua-blue hover:underline">
            Sign in for a saved run →
          </Link>
        </p>
      </div>
    )
  }
  return (
    <div className="max-w-4xl mx-auto">
      <div className="text-center mb-4 mt-6 sm:mt-8 px-2 min-w-0">
        <div className="font-mono text-[11px] uppercase text-mathua-muted mb-3">Welcome to Mathua</div>
        <h1 className="font-serif text-2xl sm:text-4xl font-medium text-mathua-primary px-2">
          Let&apos;s find your starting point
        </h1>
        <p className="text-mathua-secondary text-sm mt-3 max-w-[500px] mx-auto px-2">
          Select what you want to learn. We&apos;ll test your current knowledge and build a personalized plan.
        </p>
      </div>

      <div className="flex gap-3 justify-center mb-6 px-2">
        <button
          type="button"
          onClick={onSelectAll}
          aria-pressed={domains.length > 0 && domains.every(d => d.selected)}
          className="bg-mathua-surface border border-mathua-border rounded-none h-10 min-h-[36px] px-6 text-sm text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
        >
          {domains.length > 0 && domains.every(d => d.selected) ? 'Deselect everything' : 'Select everything'}
        </button>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2 sm:gap-3 mb-8 min-w-0">
        {domains.map(d => {
          const label = domainLabels(d.name)
          return (
            <button
              type="button"
              key={d.name}
              onClick={() => onToggle(d.name)}
              className={`rounded-none p-3 sm:p-4 text-left transition-all text-sm min-h-[60px] min-w-0 overflow-hidden ${
                d.selected
                  ? 'bg-mathua-blue text-white'
                  : 'bg-mathua-surface border border-mathua-border text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue'
              }`}
            >
              <div className="font-medium truncate sm:whitespace-normal sm:line-clamp-2 break-words text-[13px] sm:text-sm">{label}</div>
              <div className={`font-mono text-[10px] mt-1 truncate ${d.selected ? 'text-white/70' : 'text-mathua-muted'}`}>
                {d.count} concepts
              </div>
            </button>
          )
        })}
      </div>

      <div className="text-center px-4">
        {hasPaused && (
          <button
            type="button"
            onClick={onResume}
            disabled={loading}
            className="border border-mathua-blue bg-mathua-blue text-white hover:opacity-90 rounded-none h-12 min-h-[44px] px-6 sm:px-10 font-medium text-sm disabled:opacity-50 max-w-full mb-3"
          >
            {loading ? (<><Loading inline size={13} /> Loading…</>) : 'Continue diagnostic →'}
          </button>
        )}
        <button
          type="button"
          onClick={onStart}
          disabled={selectedCount === 0 || loading}
          className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 min-h-[44px] px-6 sm:px-10 font-medium text-sm disabled:opacity-50 max-w-full"
        >
          {loading ? (<><Loading inline size={13} /> Loading…</>) : `Start diagnostic test (${selectedCount} concepts in scope)`}
        </button>
      </div>

      {/* The way out, and the reason it is a link rather than a flag.
          A learner can start any topic whose prerequisites are met without ever
          taking the diagnostic: /learn?concept=… is gated server-side, not by
          this flag, and /domains is already browsable. Nothing here records a
          result, seeds placement, or grants credit — it only makes an existing
          path findable from the screen a learner is actually looking at. */}
      <div className="text-center px-4 mt-6">
        <Link
          href="/domains"
          className="font-mono text-xs text-mathua-muted hover:text-mathua-blue underline underline-offset-2 min-h-[44px] inline-flex items-center"
        >
          Browse topics and start learning without the test →
        </Link>
      </div>
    </div>
  )
}

export function DiagnosticStep({
  question,
  conceptName,
  progress,
  questionCount,
  answerInput,
  lastResult,
  accuracy,
  loading,
  inputRef,
  conceptId,
  sessionId,
  onInputChange,
  onSubmit,
  onDontKnow,
  onNext,
  done,
  answerFormat,
  submitError,
  planError,
  onSkip,
  skipsLeft,
  onRestart,
  onRetryPlan,
  retryAvailable,
  onRetryQuestion,
}: {
  question: string
  conceptName: string
  progress: DiagnosticProgress | null
  questionCount: number
  answerInput: string
  lastResult: { correct: boolean; feedback: string } | null
  accuracy: { correct: number; total: number }
  loading: boolean
  inputRef: RefObject<MathFocusHandle | null>
  conceptId: string
  sessionId: string
  onInputChange: (value: string) => void
  onSubmit: () => void
  onDontKnow: () => void
  onNext: () => void
  done: boolean
  answerFormat: AnswerFormat
  submitError: SubmitError | null
  planError: string
  onSkip: () => void
  skipsLeft: number
  onRestart: () => void
  onRetryPlan: () => void
  retryAvailable: boolean
  onRetryQuestion: () => void
}) {
  // The number of the question on screen. `progress.answered` already counts it once
  // its feedback is showing, so adding one then labelled the result of question 1 as
  // "Question 2" before Next had been pressed.
  const shownNumber = progress ? progress.answered + (lastResult ? 0 : 1) : questionCount
  return (
    <>
      <SectionHeader label={`Question ${shownNumber}`} title={conceptName} />
      <div className="max-w-2xl mx-auto px-2 sm:px-0 min-w-0 overflow-hidden">
        <ProgressBar
          answered={shownNumber}
          coverDone={progress ? progress.cover_done : 0}
          coverSize={progress ? progress.cover_size : 0}
        />
        <div className="bg-mathua-surface border border-mathua-border rounded-none p-4 sm:p-6 mb-6 w-full max-w-full min-w-0 overflow-hidden">
          <div className="bg-mathua-code border border-mathua-border rounded-none p-4 sm:p-6 text-center mb-4 w-full max-w-full min-w-0 overflow-hidden">
            <div className="w-full max-w-full min-w-0 overflow-hidden">
              <KatexContent className="text-mathua-primary text-lg font-mono font-light whitespace-pre-wrap break-words">{question}</KatexContent>
            </div>
          </div>
          <ChoiceOptions question={question} value={answerInput} onPick={onInputChange} disabled={loading || lastResult !== null} />
          <form
            onSubmit={e => { e.preventDefault(); onSubmit() }}
            className="flex flex-col sm:flex-row gap-3 min-w-0"
          >
            <label htmlFor="onboard-answer" className="sr-only">Your answer</label>
            <MathAnswerInput
              id="onboard-answer"
              value={answerInput}
              onChange={onInputChange}
              onSubmit={onSubmit}
              conceptId={conceptId}
              disabled={loading || lastResult !== null}
              status={lastResult ? (lastResult.correct ? 'correct' : 'incorrect') : 'default'}
              inputMode={answerFormat.inputMode}
              placeholder="Your answer..."
              className="sm:flex-1"
              focusRef={inputRef}
            />
            <button
              type="submit"
              disabled={!answerInput.trim() || loading || lastResult !== null}
              className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 min-h-[36px] px-6 font-medium text-sm disabled:opacity-50 shrink-0 w-auto self-end sm:self-auto"
            >
              Check Answer
            </button>
          </form>
          <div className="mt-2 flex items-center justify-between gap-2">
            <p className="font-mono text-[11px] text-mathua-muted">{answerFormat.hint}</p>
            {!lastResult && !loading && (
              <button
                type="button"
                onClick={onDontKnow}
                className="shrink-0 font-mono text-[11px] text-mathua-muted hover:text-mathua-red underline underline-offset-2"
              >
                I don&apos;t know
              </button>
            )}
          </div>
          <div className="mt-2 flex justify-end">
            <ReportButton
              key={question}
              conceptId={conceptId}
              kind="question"
              question={question}
              source="diagnostic"
              sessionId={sessionId}
              answerGiven={answerInput}
              answerHint={answerFormat.hint}
              answerInputMode={answerFormat.inputMode}
            />
          </div>
        </div>

        {lastResult && (
          <div className={`bg-mathua-surface border rounded-none p-4 mb-4 text-center ${lastResult.correct ? 'border-mathua-green' : 'border-mathua-red'}`}>
            <KatexContent className={lastResult.correct ? 'text-mathua-green' : 'text-mathua-red'}>{lastResult.feedback}</KatexContent>
            {!lastResult.correct && !done && retryAvailable && (
              <button
                type="button"
                onClick={onRetryQuestion}
                disabled={loading}
                className="mt-2 font-mono text-[11px] text-mathua-muted hover:text-mathua-blue underline underline-offset-2 disabled:opacity-50"
              >
                I made a silly mistake — retry
              </button>
            )}
          </div>
        )}

        {submitError ? (
          <div className="mb-4">
            <SubmitErrorBlock
              error={submitError}
              onRetry={onSubmit}
              onSkip={onSkip}
              skipsLeft={skipsLeft}
              onRestart={onRestart}
              restartLabel="Restart diagnostic"
              retrying={loading}
            />
          </div>
        ) : lastResult && !done ? (
          <div className="mb-4 text-center">
            <button
              type="button"
              ref={focusClearOfTabs}
              onClick={onNext}
              className="clear-of-tabs border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint rounded-none h-12 px-8 font-medium text-sm shrink-0"
            >
              Next →
            </button>
          </div>
        ) : lastResult && done && planError ? (
          <div className="mb-4 border border-mathua-red bg-mathua-surface p-4 text-left" role="alert">
            <p className="font-mono text-xs text-mathua-red">Couldn&apos;t load your results.</p>
            <p className="mt-1 font-mono text-[11px] text-mathua-muted break-words [overflow-wrap:anywhere]">{planError}</p>
            <div className="mt-3">
              <button
                type="button"
                onClick={onRetryPlan}
                disabled={loading}
                className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue-faint px-4 h-10 font-mono text-xs disabled:opacity-50"
              >
                {loading ? 'Loading…' : 'Load results'}
              </button>
            </div>
          </div>
        ) : lastResult && done ? (
          <div className="mb-4 text-center text-mathua-muted text-xs font-mono">
            Preparing results…
          </div>
        ) : null}

        <div className="text-center text-mathua-muted text-xs font-mono">
          {accuracy.total > 0 && `${accuracy.correct}/${accuracy.total} correct`}
        </div>
      </div>
    </>
  )
}

export function ResultsStep({ plan, onStartPractice, demo }: { plan: GoalPlanRes; onStartPractice: () => void; demo?: boolean }) {
  return (
    <div className="max-w-3xl mx-auto min-w-0 overflow-hidden px-2 sm:px-0">
      <SectionHeader label="Your results" title="Here's what we found" />
      {demo && (
        <div className="mt-4 border border-mathua-border bg-mathua-surface px-4 py-3">
          <p className="text-mathua-secondary text-sm">
            This was a demo run. The plan below is real, but it is not attached to an account, so
            it will not be here when this browser session ends.
          </p>
          <p className="mt-2">
            <Link href="/login" className="font-mono text-xs text-mathua-blue hover:underline">
              Sign in to keep a plan →
            </Link>
          </p>
        </div>
      )}
      <div className="mt-6 min-w-0 overflow-hidden">
        <DiagnosticResults plan={plan} onStartPractice={onStartPractice} />
      </div>
    </div>
  )
}
