'use client'

import type { RefObject } from 'react'
import KatexContent from '../../components/KatexContent'
import SectionHeader from '../../components/SectionHeader'
import ProgressBar from '../../components/ProgressBar'
import SymbolPalette from '../../components/SymbolPalette'
import ReportButton from '../../components/ReportButton'
import BriefingCard, { DIAGNOSTIC_BRIEFING } from '../../components/BriefingCard'
import SubmitErrorBlock, { type SubmitError } from '../../components/SubmitErrorBlock'
import DiagnosticResults from '../../components/DiagnosticResults'
import Loading from '../../components/Loading'
import type { DiagnosticProgress, GoalPlanRes } from '../../lib/api'
import type { AnswerFormat } from '../../lib/answerFormat'
import { domainLabels, type DomainInfo } from './domains'
import { Input } from '@/components/ui/input'

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
}) {
  if (confirming) {
    const selected = domains.filter(d => d.selected)
    const shown = selected.slice(0, 3).map(d => domainLabels[d.name] || d.name)
    const extra = selected.length > 3 ? ` +${selected.length - 3} more` : ''
    return (
      <div className="max-w-2xl mx-auto px-2">
        <div className="text-center mb-4 mt-6 sm:mt-8 min-w-0">
          <div className="font-mono text-[11px] uppercase text-mathua-muted mb-3">One quick check</div>
          <h1 className="font-serif text-2xl sm:text-3xl font-medium text-mathua-primary px-2">
            Before you begin
          </h1>
          <p className="text-mathua-secondary text-sm mt-3 px-2">
            Starting test on: {shown.join(', ')}{extra} · {selectedCount} concepts
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
          className="bg-mathua-surface border border-mathua-border rounded-none h-10 min-h-[36px] px-6 text-sm text-mathua-secondary hover:border-mathua-blue hover:text-mathua-blue"
        >
          Select everything
        </button>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2 sm:gap-3 mb-8 min-w-0">
        {domains.map(d => {
          const label = domainLabels[d.name] || d.name
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
          className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 min-h-[44px] px-6 sm:px-10 font-medium text-sm disabled:opacity-50 max-w-full"
        >
          {loading ? (<><Loading inline size={13} /> Loading…</>) : `Start diagnostic test (${selectedCount} concepts)`}
        </button>
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
  onNext,
  done,
  answerFormat,
  submitError,
  planError,
  onSkip,
  skipsLeft,
  onRestart,
  onRetryPlan,
}: {
  question: string
  conceptName: string
  progress: DiagnosticProgress | null
  questionCount: number
  answerInput: string
  lastResult: { correct: boolean; feedback: string } | null
  accuracy: { correct: number; total: number }
  loading: boolean
  inputRef: RefObject<HTMLInputElement | null>
  conceptId: string
  sessionId: string
  onInputChange: (value: string) => void
  onSubmit: () => void
  onNext: () => void
  done: boolean
  answerFormat: AnswerFormat
  submitError: SubmitError | null
  planError: string
  onSkip: () => void
  skipsLeft: number
  onRestart: () => void
  onRetryPlan: () => void
}) {
  return (
    <>
      <SectionHeader label={`Question ${progress ? progress.answered + 1 : questionCount}`} title={conceptName} />
      <div className="max-w-2xl mx-auto px-2 sm:px-0 min-w-0 overflow-hidden">
        <ProgressBar
          answered={progress ? progress.answered + 1 : questionCount}
          coverDone={progress ? progress.cover_done : 0}
          coverSize={progress ? progress.cover_size : 0}
        />
        <div className="bg-mathua-surface border border-mathua-border rounded-none p-4 sm:p-6 mb-6 w-full max-w-full min-w-0 overflow-hidden">
          <div className="bg-mathua-code border border-mathua-border rounded-none p-4 sm:p-6 text-center mb-4 w-full max-w-full min-w-0 overflow-hidden">
            <div className="w-full max-w-full min-w-0 overflow-hidden">
              <KatexContent className="text-mathua-primary text-lg font-mono font-light whitespace-pre-wrap break-words">{question}</KatexContent>
            </div>
          </div>
          <form
            onSubmit={e => { e.preventDefault(); onSubmit() }}
            className="flex flex-col sm:flex-row gap-3 min-w-0"
          >
            <label htmlFor="onboard-answer" className="sr-only">Your answer</label>
            <Input
              ref={inputRef}
              id="onboard-answer"
              type="text"
              value={answerInput}
              onChange={(e) => onInputChange(e.target.value)}
              placeholder="Your answer..."
              enterKeyHint="go"
              inputMode={answerFormat.inputMode}
              disabled={loading || lastResult !== null}
              className="sm:flex-1"
            />
            <button
              type="submit"
              disabled={!answerInput.trim() || loading || lastResult !== null}
              className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 min-h-[36px] px-6 font-medium text-sm disabled:opacity-50 shrink-0 w-auto self-end sm:self-auto"
            >
              Check Answer
            </button>
          </form>
          <p className="mt-2 font-mono text-[11px] text-mathua-muted">{answerFormat.hint}</p>
          <SymbolPalette targetRef={inputRef} onInsert={onInputChange} />
          <div className="mt-2 flex justify-end">
            <ReportButton
              key={question}
              conceptId={conceptId}
              kind="question"
              question={question}
              source="diagnostic"
              sessionId={sessionId}
            />
          </div>
        </div>

        {lastResult && (
          <div className={`bg-mathua-surface border rounded-none p-4 mb-4 text-center ${lastResult.correct ? 'border-mathua-green' : 'border-mathua-red'}`}>
            <KatexContent className={lastResult.correct ? 'text-mathua-green' : 'text-mathua-red'}>{lastResult.feedback}</KatexContent>
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
              autoFocus
              onClick={onNext}
              className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white rounded-none h-12 px-8 font-medium text-sm shrink-0"
            >
              Next →
            </button>
          </div>
        ) : lastResult && done && planError ? (
          <div className="mb-4 border border-mathua-red/60 bg-mathua-surface p-4 text-left" role="alert">
            <p className="font-mono text-xs text-mathua-red">Couldn&apos;t load your results.</p>
            <p className="mt-1 font-mono text-[11px] text-mathua-muted break-words [overflow-wrap:anywhere]">{planError}</p>
            <div className="mt-3">
              <button
                type="button"
                onClick={onRetryPlan}
                disabled={loading}
                className="border border-mathua-blue text-mathua-blue hover:bg-mathua-blue hover:text-white px-4 h-10 font-mono text-xs disabled:opacity-50"
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

export function ResultsStep({ plan, onStartPractice }: { plan: GoalPlanRes; onStartPractice: () => void }) {
  return (
    <div className="max-w-3xl mx-auto min-w-0 overflow-hidden px-2 sm:px-0">
      <SectionHeader label="Your results" title="Here's what we found" />
      <div className="mt-6 min-w-0 overflow-hidden">
        <DiagnosticResults plan={plan} onStartPractice={onStartPractice} />
      </div>
    </div>
  )
}
