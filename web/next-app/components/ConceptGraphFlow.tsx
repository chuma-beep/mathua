'use client'

import { useMemo, useState, useEffect, useCallback, useRef } from 'react'
import {
  ReactFlow,
  ReactFlowProvider,
  Background,
  BackgroundVariant,
  Controls,
  MiniMap,
  Handle,
  Position,
  MarkerType,
  BaseEdge,
  getSmoothStepPath,
  useReactFlow,
  type Node,
  type Edge,
  type NodeProps,
  type EdgeProps,
} from '@xyflow/react'
import '@xyflow/react/dist/base.css'
import Fuse from 'fuse.js'
import dagre from 'dagre'
import type { MasteryStatus } from '../lib/graphStatus'

export type { MasteryStatus }

export interface GraphConcept {
  id: string
  label: string
  domain: string
  prerequisites: string[]
}

interface ConceptGraphFlowProps {
  concepts: GraphConcept[]
  conceptStatuses?: Record<string, MasteryStatus>
  conceptProgress?: Record<string, number>
  theme?: 'dark' | 'light'
  onPathNodes?: string[]
  onNodeSelect?: (id: string) => void
  selectedId?: string | null
  onSelectionChange?: (id: string | null) => void
}

const STATUS_COLORS: Record<MasteryStatus, string> = {
  mastered: '#10b981',
  practicing: '#60a5fa',
  learning: '#f59e0b',
  unseen: '#71717a',
  locked: '#52525b',
}

const STATUS_LABELS: Record<MasteryStatus, string> = {
  mastered: 'Mastered',
  practicing: 'Practicing',
  learning: 'Learning',
  unseen: 'Unseen',
  locked: 'Locked',
}

const AMBIENT_KEY = 'mathua_graph_flow'

function domainColor(domain: string): string {
  let hash = 0
  for (let i = 0; i < domain.length; i++) hash = (hash * 31 + domain.charCodeAt(i)) | 0
  const hue = Math.abs(hash) % 360
  return `hsl(${hue} 55% 55%)`
}

type ConceptNodeData = {
  label: string
  domain: string
  status: MasteryStatus
  progress?: number
  onPath: boolean
  dimmed: boolean
  selected: boolean
  onSelect?: (id: string) => void
}

type ConceptFlowNode = Node<ConceptNodeData, 'concept'>

function ConceptNode({ id, data }: NodeProps<ConceptFlowNode>) {
  const statusColor = STATUS_COLORS[data.status] ?? STATUS_COLORS.unseen
  const pulseClass =
    data.status === 'mastered'
      ? ' node-pulse-mastered'
      : data.status === 'learning'
        ? ' node-breathe-learning'
        : ''
  const progressPct =
    typeof data.progress === 'number' ? Math.round(data.progress * 100) : null
  const ariaLabel = `${data.label}, ${data.domain.replace(/_/g, ' ')} domain, ${STATUS_LABELS[data.status] ?? data.status}${progressPct !== null ? `, ${progressPct}% toward mastery` : ''}`
  return (
    <div
      className={`concept-node${pulseClass}`}
      role="button"
      tabIndex={0}
      aria-label={ariaLabel}
      aria-pressed={data.selected}
      onKeyDown={e => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault()
          e.stopPropagation()
          data.onSelect?.(id)
        }
      }}
      style={{
        width: 180,
        height: 40,
        position: 'relative' as const,
        opacity: data.dimmed ? 0.22 : 1,
        transition: 'opacity 150ms ease',
        fontFamily: "'IBM Plex Mono', monospace",
        fontSize: 11,
        color: 'var(--text-primary)',
        cursor: 'pointer',
      }}
    >
      <div
        style={{
          width: '100%',
          height: '100%',
          position: 'relative' as const,
          display: 'flex',
          alignItems: 'center',
          gap: 8,
          padding: '0 10px',
          background: 'var(--surface-elevated)',
          border: data.selected
            ? `1.5px solid ${statusColor}`
            : data.onPath
              ? '1px solid var(--accent-teal)'
              : '1px solid var(--border)',
          borderRadius: 4,
          overflow: 'hidden',
        }}
      >
        <Handle type="target" position={Position.Left} style={{ opacity: 0, pointerEvents: 'none' }} />
        <span style={{ width: 3, height: 22, borderRadius: 2, flexShrink: 0, background: domainColor(data.domain) }} />
        <span style={{ width: 8, height: 8, borderRadius: '50%', flexShrink: 0, background: statusColor }} />
        <span style={{ whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{data.label}</span>
        <Handle type="source" position={Position.Right} style={{ opacity: 0, pointerEvents: 'none' }} />
        <div
          data-progress={progressPct ?? 0}
          aria-hidden
          style={{
            position: 'absolute',
            left: 8,
            right: 8,
            bottom: 3,
            height: 4,
            borderRadius: 2,
            background: 'var(--border)',
          }}
        >
          <div
            style={{
              height: '100%',
              width: `${progressPct ?? 0}%`,
              background: statusColor,
              borderRadius: 2,
              transition: 'width 300ms ease',
            }}
          />
        </div>
      </div>
      <div
        className="node-progress-tip"
        aria-hidden
        style={{
          position: 'absolute',
          top: '100%',
          left: '50%',
          transform: 'translateX(-50%)',
          marginTop: 4,
          zIndex: 10,
          whiteSpace: 'nowrap' as const,
          pointerEvents: 'none' as const,
          color: 'var(--text-secondary)',
          background: 'var(--surface-elevated)',
          border: '0.5px solid var(--border-strong)',
          borderRadius: 3,
          padding: '3px 8px',
        }}
      >
        {progressPct !== null ? `${progressPct}% toward mastery` : 'Not started'}
      </div>
    </div>
  )
}

const nodeTypes = { concept: ConceptNode }

type FlowEdgeData = {
  variant: 'plain' | 'flow' | 'ambient'
  highlighted: boolean
}

type ConceptFlowEdge = Edge<FlowEdgeData, 'flowedge'>

function FlowEdge({ sourceX, sourceY, targetX, targetY, sourcePosition, targetPosition, data, markerEnd }: EdgeProps<ConceptFlowEdge>) {
  const [path] = getSmoothStepPath({
    sourceX,
    sourceY,
    targetX,
    targetY,
    sourcePosition,
    targetPosition,
    borderRadius: 12,
  })
  const variantClass = data?.variant === 'flow' ? ' edge-flow' : data?.variant === 'ambient' ? ' edge-ambient' : ''
  return (
    <BaseEdge
      path={path}
      markerEnd={markerEnd}
      className={variantClass}
      style={{
        stroke: data?.highlighted ? 'var(--accent-blue)' : 'var(--border)',
        strokeWidth: data?.highlighted ? 1.6 : 1,
        opacity: data?.variant === 'ambient' && !data?.highlighted ? 0.5 : 1,
      }}
    />
  )
}

const edgeTypes = { flowedge: FlowEdge }

function layoutWithDagre(concepts: GraphConcept[]): Map<string, { x: number; y: number }> {
  const g = new dagre.graphlib.Graph()
  g.setGraph({ rankdir: 'LR', nodesep: 16, ranksep: 70, marginx: 20, marginy: 20 })
  g.setDefaultEdgeLabel(() => ({}))
  for (const c of concepts) g.setNode(c.id, { width: 180, height: 40 })
  const known = new Set(concepts.map(c => c.id))
  for (const c of concepts) {
    for (const p of c.prerequisites) {
      if (known.has(p)) g.setEdge(p, c.id)
    }
  }
  dagre.layout(g)

  const domainOrder: string[] = []
  const seenDomains = new Set<string>()
  for (const c of concepts) {
    if (!seenDomains.has(c.domain)) {
      seenDomains.add(c.domain)
      domainOrder.push(c.domain)
    }
  }

  const byRank = new Map<number, string[]>()
  for (const c of concepts) {
    const n = g.node(c.id)
    if (!n) continue
    const rank = Math.round(n.x)
    if (!byRank.has(rank)) byRank.set(rank, [])
    byRank.get(rank)!.push(c.id)
  }

  const positions = new Map<string, { x: number; y: number }>()
  byRank.forEach((ids) => {
    ids.sort((a, b) => {
      const ca = concepts.find(x => x.id === a)!
      const cb = concepts.find(x => x.id === b)!
      const d = domainOrder.indexOf(ca.domain) - domainOrder.indexOf(cb.domain)
      return d !== 0 ? d : ca.label.localeCompare(cb.label)
    })
    let prevDomain: string | null = null
    let y = 0
    for (const id of ids) {
      const c = concepts.find(x => x.id === id)!
      if (prevDomain !== null && c.domain !== prevDomain) y += 28
      positions.set(id, { x: g.node(id).x - 90, y })
      prevDomain = c.domain
      y += 40 + 16
    }
  })
  return positions
}

interface SearchResult {
  id: string
  label: string
  domain: string
}

function SearchOverlay({
  concepts,
  onSelect,
}: {
  concepts: GraphConcept[]
  onSelect: (id: string) => void
}) {
  const [query, setQuery] = useState('')
  const [open, setOpen] = useState(false)
  const boxRef = useRef<HTMLDivElement>(null)

  const fuse = useMemo(
    () =>
      new Fuse(concepts, {
        keys: [
          { name: 'label', weight: 2 },
          { name: 'domain', weight: 1 },
        ],
        threshold: 0.4,
        minMatchCharLength: 2,
      }),
    [concepts]
  )

  const results = useMemo<SearchResult[]>(() => {
    if (!query.trim()) return []
    return fuse
      .search(query.trim())
      .slice(0, 8)
      .map(r => ({ id: r.item.id, label: r.item.label, domain: r.item.domain }))
  }, [query, fuse])

  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (boxRef.current && !boxRef.current.contains(e.target as globalThis.Node)) setOpen(false)
    }
    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [])

  return (
    <div
      ref={boxRef}
      className="absolute top-2.5 left-[62px] sm:left-16 z-[5] w-[min(calc(100%-76px),240px)] sm:w-[240px] font-mono"
    >
      <input
        value={query}
        onChange={e => {
          setQuery(e.target.value)
          setOpen(true)
        }}
        onFocus={() => setOpen(true)}
        aria-label="Search concepts"
        onKeyDown={e => {
          if (e.key === 'Enter' && results.length > 0) {
            onSelect(results[0].id)
            setOpen(false)
          }
          if (e.key === 'Escape') setOpen(false)
        }}
        placeholder="Search concepts…"
        className="w-full min-h-[36px] h-9 px-2.5 text-[11px] font-mono text-mathua-primary bg-mathua-surface-elevated border-[0.5px] border-mathua-border-strong rounded-[4px] outline-none"
      />
      {open && results.length > 0 && (
        <div className="mt-1 bg-mathua-surface-elevated border-[0.5px] border-mathua-border rounded-[4px] overflow-hidden max-h-[min(50dvh,240px)] overflow-y-auto">
          {results.map(r => (
            <button
              key={r.id}
              type="button"
              onClick={() => {
                onSelect(r.id)
                setOpen(false)
                setQuery('')
              }}
              className="flex items-center gap-1.5 w-full text-left px-2.5 py-2 min-h-[36px] text-[11px] font-mono text-mathua-secondary bg-transparent border-none border-b-[0.5px] border-mathua-border last:border-b-0 cursor-pointer hover:bg-mathua-code"
            >
              <span className="w-[3px] h-3.5 rounded-[2px] shrink-0" style={{ background: domainColor(r.domain) }} />
              <span className="truncate flex-1 min-w-0">{r.label}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  )
}

function AmbientToggle({
  enabled,
  onToggle,
}: {
  enabled: boolean
  onToggle: () => void
}) {
  return (
    <button
      type="button"
      onClick={onToggle}
      title="Toggle edge flow animation"
      aria-label="Toggle edge flow animation"
      style={{
        position: 'absolute',
        top: 10,
        right: 10,
        zIndex: 5,
        height: 30,
        padding: '0 12px',
        fontSize: 10,
        textTransform: 'uppercase',
        letterSpacing: '0.05em',
        fontFamily: "'IBM Plex Mono', monospace",
        color: enabled ? '#fff' : 'var(--text-muted)',
        background: enabled ? 'var(--accent-blue)' : 'var(--surface-elevated)',
        border: '0.5px solid var(--border-strong)',
        borderRadius: 4,
        cursor: 'pointer',
      }}
    >
      ⏵ Flow
    </button>
  )
}

const LIST_TOGGLE_LABEL = 'List'

function ListToggleButton({
  open,
  onToggle,
}: {
  open: boolean
  onToggle: () => void
}) {
  return (
    <button
      type="button"
      onClick={onToggle}
      title="Browse concepts as a keyboard-accessible list"
      aria-expanded={open}
      aria-label={open ? 'Close concept list' : 'Open concept list'}
      className={`absolute top-2.5 left-2.5 z-[6] min-h-[36px] h-9 px-2.5 text-[10px] uppercase tracking-[0.05em] font-mono border-[0.5px] border-mathua-border-strong rounded-[4px] cursor-pointer transition-colors ${open ? 'bg-mathua-teal text-white' : 'bg-mathua-surface-elevated text-mathua-muted hover:text-mathua-primary'}`}
    >
      {LIST_TOGGLE_LABEL}
    </button>
  )
}

function ListView({
  concepts,
  conceptStatuses,
  selectedId,
  onSelect,
  onClose,
}: {
  concepts: GraphConcept[]
  conceptStatuses?: Record<string, MasteryStatus>
  selectedId: string | null
  onSelect: (id: string) => void
  onClose: () => void
}) {
  const [query, setQuery] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    inputRef.current?.focus()
  }, [])

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase()
    if (!q) return concepts
    return concepts.filter(
      c =>
        c.label.toLowerCase().includes(q) ||
        c.id.toLowerCase().includes(q) ||
        c.domain.toLowerCase().includes(q)
    )
  }, [concepts, query])

  const groups = useMemo(() => {
    const order: string[] = []
    const byDomain = new Map<string, GraphConcept[]>()
    for (const c of filtered) {
      if (!byDomain.has(c.domain)) {
        byDomain.set(c.domain, [])
        order.push(c.domain)
      }
      byDomain.get(c.domain)!.push(c)
    }
    return order.map(d => ({ domain: d, items: byDomain.get(d)! }))
  }, [filtered])

  return (
    <div
      role="dialog"
      aria-label="Concept list"
      style={{
        position: 'absolute',
        inset: '48px 8px 8px 8px',
        zIndex: 6,
        display: 'flex',
        flexDirection: 'column',
        background: 'var(--surface-elevated)',
        border: '0.5px solid var(--border)',
        borderRadius: 4,
        overflow: 'hidden',
      }}
    >
      <div style={{ padding: 8 }}>
        <input
          ref={inputRef}
          value={query}
          onChange={e => setQuery(e.target.value)}
          onKeyDown={e => {
            if (e.key === 'Escape') {
              e.preventDefault()
              onClose()
            }
          }}
          placeholder="Filter concepts…"
          aria-label="Filter concepts"
          style={{
            width: '100%',
            height: 28,
            padding: '0 8px',
            fontSize: 11,
            fontFamily: 'inherit',
            color: 'var(--text-primary)',
            background: 'var(--surface-elevated)',
            border: '0.5px solid var(--border-strong)',
            borderRadius: 3,
            outline: 'none',
          }}
        />
      </div>
      <div style={{ flex: 1, overflowY: 'auto', padding: '0 8px 8px' }}>
        {groups.map(g => (
          <div key={g.domain}>
            <div
              style={{
                margin: '8px 0 4px',
                fontSize: 9,
                textTransform: 'uppercase',
                letterSpacing: '0.08em',
                color: 'var(--text-muted)',
                fontFamily: "'IBM Plex Mono', monospace",
              }}
            >
              {g.domain.replace(/_/g, ' ')} · {g.items.length}
            </div>
            {g.items.map(c => {
              const status = conceptStatuses?.[c.id]
              return (
                <button
                  key={c.id}
                  type="button"
                  onClick={() => {
                    onSelect(c.id)
                    onClose()
                  }}
                  aria-current={c.id === selectedId ? 'true' : undefined}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 8,
                    width: '100%',
                    textAlign: 'left',
                    padding: '5px 8px',
                    fontSize: 11,
                    fontFamily: 'inherit',
                    color:
                      c.id === selectedId
                        ? 'var(--text-primary)'
                        : 'var(--text-secondary)',
                    background:
                      c.id === selectedId
                        ? 'var(--surface-elevated)'
                        : 'transparent',
                    border: 'none',
                    borderBottom: '0.5px solid var(--border)',
                    borderRadius: 3,
                    cursor: 'pointer',
                  }}
                >
                  <span
                    style={{
                      width: 7,
                      height: 7,
                      borderRadius: '50%',
                      flexShrink: 0,
                      background:
                        status != null ? STATUS_COLORS[status] : STATUS_COLORS.unseen,
                    }}
                    aria-hidden
                  />
                  <span
                    style={{
                      whiteSpace: 'nowrap',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      flex: 1,
                    }}
                  >
                    {c.label}
                  </span>
                  {status && (
                    <span style={{ color: 'var(--text-muted)', fontSize: 9 }}>
                      {STATUS_LABELS[status]}
                    </span>
                  )}
                </button>
              )
            })}
          </div>
        ))}
        {groups.length === 0 && (
          <div
            style={{
              padding: 12,
              fontSize: 11,
              color: 'var(--text-muted)',
              fontFamily: "'IBM Plex Mono', monospace",
            }}
          >
            No matching concepts
          </div>
        )}
      </div>
    </div>
  )
}

function GraphInner({
  concepts,
  conceptStatuses,
  conceptProgress,
  theme = 'dark',
  onPathNodes,
  onNodeSelect,
  selectedId,
  onSelectionChange,
}: ConceptGraphFlowProps) {
  const [isMobile, setIsMobile] = useState(false)
  const [internalSelected, setInternalSelected] = useState<string | null>(selectedId ?? null)
  const [ambientOn, setAmbientOn] = useState(false)
  const [listOpen, setListOpen] = useState(false)
  const [inView, setInView] = useState(true)
  const [tabVisible, setTabVisible] = useState(true)
  const wrapperRef = useRef<HTMLDivElement>(null)
  const hoverRef = useRef(false)
  const { setCenter, zoomIn, zoomOut, fitView, getViewport, setViewport } = useReactFlow()

  const effectiveSelected = selectedId !== undefined ? selectedId : internalSelected

  useEffect(() => {
    try {
      setAmbientOn(window.localStorage.getItem(AMBIENT_KEY) === '1')
    } catch {}
  }, [])

  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 640)
    checkMobile()
    window.addEventListener('resize', checkMobile)
    return () => window.removeEventListener('resize', checkMobile)
  }, [])

  useEffect(() => {
    const el = wrapperRef.current
    if (!el || typeof IntersectionObserver === 'undefined') return
    const obs = new IntersectionObserver(entries => setInView(entries[0]?.isIntersecting ?? true), {
      threshold: 0.02,
    })
    obs.observe(el)
    return () => obs.disconnect()
  }, [])

  useEffect(() => {
    const onVis = () => setTabVisible(document.visibilityState === 'visible')
    document.addEventListener('visibilitychange', onVis)
    return () => document.removeEventListener('visibilitychange', onVis)
  }, [])

  const ambientActive = ambientOn && !isMobile && inView && tabVisible

  const select = useCallback(
    (id: string | null) => {
      setInternalSelected(id)
      onSelectionChange?.(id)
    },
    [onSelectionChange]
  )

  const listOpenRef = useRef(listOpen)
  useEffect(() => {
    listOpenRef.current = listOpen
  }, [listOpen])

  // Keyboard viewport controls. Active while focus is inside the graph or the
  // pointer hovers it, and never while typing in an input.
  useEffect(() => {
    function isEditable(t: EventTarget | null): boolean {
      if (!(t instanceof HTMLElement)) return false
      return (
        t.isContentEditable ||
        t.tagName === 'INPUT' ||
        t.tagName === 'TEXTAREA' ||
        t.tagName === 'SELECT'
      )
    }

    function onKeyDown(e: KeyboardEvent) {
      if (e.metaKey || e.ctrlKey || e.altKey) return
      if (isEditable(e.target)) return
      const wrapper = wrapperRef.current
      const engaged =
        (wrapper !== null && wrapper.contains(document.activeElement)) || hoverRef.current
      if (!engaged) return

      switch (e.key) {
        case '+':
        case '=':
          e.preventDefault()
          zoomIn({ duration: 200 })
          break
        case '-':
        case '_':
          e.preventDefault()
          zoomOut({ duration: 200 })
          break
        case '0':
          e.preventDefault()
          fitView({ padding: 0.15, maxZoom: 1, duration: 300 })
          break
        case 'ArrowLeft':
        case 'ArrowRight':
        case 'ArrowUp':
        case 'ArrowDown': {
          e.preventDefault()
          const vp = getViewport()
          const step = 90
          const next = { ...vp }
          if (e.key === 'ArrowLeft') next.x += step
          if (e.key === 'ArrowRight') next.x -= step
          if (e.key === 'ArrowUp') next.y += step
          if (e.key === 'ArrowDown') next.y -= step
          setViewport(next, { duration: 200 })
          break
        }
        case 'Escape':
          e.preventDefault()
          if (listOpenRef.current) {
            setListOpen(false)
          } else {
            select(null)
          }
          break
      }
    }

    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [zoomIn, zoomOut, fitView, getViewport, setViewport, select])

  // React Flow re-emits selection changes whenever this callback's identity
  // changes (every render, if inline) and once on mount with an empty list —
  // both would clear the selection and wipe a ?concept= deep link. Keep the
  // identity stable and treat this stream as select-only: deselection goes
  // through onPaneClick and the Escape handler instead.
  const handleRFSelectionChange = useCallback(
    ({ nodes }: { nodes: ConceptFlowNode[] }) => {
      const first = nodes[0]?.id
      if (first) select(first)
    },
    [select]
  )

  const conceptById = useMemo(() => new Map(concepts.map(c => [c.id, c])), [concepts])

  const dependentsMap = useMemo(() => {
    const m = new Map<string, string[]>()
    for (const c of concepts) {
      for (const p of c.prerequisites) {
        if (!m.has(p)) m.set(p, [])
        m.get(p)!.push(c.id)
      }
    }
    return m
  }, [concepts])

  const neighborhood = useMemo(() => {
    if (!effectiveSelected) return null
    const upstream = new Set<string>()
    const downstream = new Set<string>()
    const walkUp = (id: string) => {
      for (const p of conceptById.get(id)?.prerequisites ?? []) {
        if (!upstream.has(p)) {
          upstream.add(p)
          walkUp(p)
        }
      }
    }
    const walkDown = (id: string) => {
      for (const d of dependentsMap.get(id) ?? []) {
        if (!downstream.has(d)) {
          downstream.add(d)
          walkDown(d)
        }
      }
    }
    walkUp(effectiveSelected)
    walkDown(effectiveSelected)
    return { upstream, downstream }
  }, [effectiveSelected, conceptById, dependentsMap])

  const layout = useMemo(() => layoutWithDagre(concepts), [concepts])

  const flowNodes = useMemo<ConceptFlowNode[]>(() => {
    return concepts.map(c => {
      const inNeighborhood =
        !neighborhood ||
        c.id === effectiveSelected ||
        neighborhood.upstream.has(c.id) ||
        neighborhood.downstream.has(c.id)
      return {
        id: c.id,
        type: 'concept' as const,
        position: layout.get(c.id) ?? { x: 0, y: 0 },
        data: {
          label: c.label,
          domain: c.domain,
          status: conceptStatuses?.[c.id] ?? 'unseen',
          progress: conceptProgress?.[c.id],
          onPath: onPathNodes?.includes(c.id) ?? false,
          dimmed: !inNeighborhood,
          selected: c.id === effectiveSelected,
          onSelect: select,
        },
      }
    })
  }, [concepts, conceptStatuses, conceptProgress, onPathNodes, neighborhood, effectiveSelected, layout, select])

  const knownIds = useMemo(() => new Set(concepts.map(c => c.id)), [concepts])

  const allEdgesList = useMemo(() => {
    const list: { source: string; target: string }[] = []
    for (const c of concepts) {
      for (const p of c.prerequisites) {
        if (knownIds.has(p)) list.push({ source: p, target: c.id })
      }
    }
    return list
  }, [concepts, knownIds])

  const flowEdges = useMemo<ConceptFlowEdge[]>(() => {
    return allEdgesList.map(({ source, target }) => {
      const highlighted = !!effectiveSelected && (source === effectiveSelected || target === effectiveSelected)
      const variant: FlowEdgeData['variant'] = highlighted
        ? 'flow'
        : ambientActive
          ? 'ambient'
          : 'plain'
      return {
        id: `${source}->${target}`,
        source,
        target,
        type: 'flowedge' as const,
        data: { variant, highlighted },
        zIndex: highlighted ? 1 : 0,
        markerEnd: highlighted ? { type: MarkerType.ArrowClosed, width: 12, height: 12 } : undefined,
      }
    })
  }, [allEdgesList, effectiveSelected, ambientActive])

  useEffect(() => {
    if (!effectiveSelected || isMobile) return
    const pos = layout.get(effectiveSelected)
    if (!pos) return
    setCenter(pos.x + 90, pos.y + 20, { zoom: 0.85, duration: 500 })
  }, [effectiveSelected, layout, setCenter, isMobile])

  const handlePaneClick = useCallback(() => {
    select(null)
  }, [select])

  const toggleAmbient = useCallback(() => {
    setAmbientOn(prev => {
      const next = !prev
      try {
        window.localStorage.setItem(AMBIENT_KEY, next ? '1' : '0')
      } catch {}
      return next
    })
  }, [])

  const selected = effectiveSelected ? conceptById.get(effectiveSelected) ?? null : null
  const selectedStatus = effectiveSelected ? conceptStatuses?.[effectiveSelected] ?? 'unseen' : null

  const prereqList = useMemo(() => {
    if (!selected) return []
    return selected.prerequisites
      .filter(p => conceptById.has(p))
      .map(p => ({ id: p, label: conceptById.get(p)!.label }))
  }, [selected, conceptById])

  const unlocksList = useMemo(() => {
    if (!selected) return []
    return (dependentsMap.get(selected.id) ?? []).map(d => ({
      id: d,
      label: conceptById.get(d)?.label ?? d,
    }))
  }, [selected, dependentsMap, conceptById])

  if (concepts.length === 0) {
    return (
      <div
        className="flex items-center justify-center border-[0.5px] border-mathua-border bg-mathua-graph-surface text-mathua-muted font-mono text-sm w-full max-w-full"
        style={{
          height: isMobile ? 'clamp(320px,60dvh,520px)' : 'clamp(320px,50dvh,520px)',
        }}
      >
        No concepts to display
      </div>
    )
  }

  return (
    <div className="w-full max-w-full min-w-0 overflow-hidden">
      <div
        ref={wrapperRef}
        data-testid="graph-wrapper"
        className={`${!ambientActive && ambientOn && !isMobile ? 'graph-flow-paused' : ''} relative w-full max-w-full min-w-0 overflow-hidden border-[0.5px] border-mathua-border bg-mathua-graph-surface`.trim()}
        onMouseEnter={() => {
          hoverRef.current = true
        }}
        onMouseLeave={() => {
          hoverRef.current = false
        }}
        aria-keyshortcuts="plus minus 0 ArrowLeft ArrowRight ArrowUp ArrowDown Escape"
        style={{
          height: isMobile ? 'clamp(320px,60dvh,520px)' : 'clamp(320px,50dvh,520px)',
          paddingBottom: 'env(safe-area-inset-bottom)',
        }}
      >
        <style>{`.react-flow__controls-button{width:44px!important;height:44px!important;min-height:44px!important;min-width:44px!important} .react-flow__minimap{width:80px!important;height:50px!important} @media(min-width:640px){.react-flow__controls-button{width:28px!important;height:28px!important;min-height:28px!important;min-width:28px!important} .react-flow__minimap{width:160px!important;height:100px!important}}`}</style>
        {/* MIT-licensed library; hiding the attribution badge is permitted.
            The LICENSE copyright notice remains in node_modules untouched. */}
        <ReactFlow
          nodes={flowNodes}
          edges={flowEdges}
          nodeTypes={nodeTypes}
          edgeTypes={edgeTypes}
          onNodeClick={(_, node) => select(node.id)}
          onPaneClick={handlePaneClick}
          onSelectionChange={handleRFSelectionChange}
          defaultEdgeOptions={{ type: 'flowedge' }}
          fitView
          fitViewOptions={{ padding: 0.15, maxZoom: 1 }}
          minZoom={0.03}
          maxZoom={2.5}
          onlyRenderVisibleElements
          nodesDraggable={false}
          nodesConnectable={false}
          elementsSelectable
          proOptions={{ hideAttribution: true }}
        >
          <Background
            variant={BackgroundVariant.Dots}
            gap={24}
            size={1}
            color={theme === 'dark' ? '#2a2a30' : '#d4d4d8'}
          />
          <Controls
            showInteractive={false}
            style={{ background: 'var(--surface-elevated)', borderColor: 'var(--border)', color: 'var(--text-primary)' }}
          />
          <MiniMap
            pannable
            zoomable
            maskColor={theme === 'dark' ? 'rgba(0,0,0,0.6)' : 'rgba(255,255,255,0.6)'}
            style={{ background: 'var(--surface-elevated)' }}
            nodeColor={node => {
              const n = node as ConceptFlowNode
              return STATUS_COLORS[n.data?.status] ?? STATUS_COLORS.unseen
            }}
          />
        </ReactFlow>
        <SearchOverlay
          concepts={concepts}
          onSelect={id => {
            select(id)
          }}
        />
        <ListToggleButton open={listOpen} onToggle={() => setListOpen(o => !o)} />
        {listOpen && (
          <ListView
            concepts={concepts}
            conceptStatuses={conceptStatuses}
            selectedId={effectiveSelected}
            onSelect={select}
            onClose={() => setListOpen(false)}
          />
        )}
        {!isMobile && <AmbientToggle enabled={ambientOn} onToggle={toggleAmbient} />}
      </div>

      {!isMobile && (
        <div
          style={{
            marginTop: 4,
            textAlign: 'right',
            fontSize: 10,
            letterSpacing: '0.05em',
            color: 'var(--text-muted)',
            fontFamily: "'IBM Plex Mono', monospace",
            userSelect: 'none' as const,
          }}
        >
          + − zoom · 0 fit · ← → ↑ ↓ pan · Esc clear
        </div>
      )}

      {selected && (
        <div className="mt-2.5 p-3 sm:p-3.5 border-[0.5px] border-mathua-border bg-mathua-surface font-mono w-full max-w-full min-w-0 overflow-hidden">
          <div className="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-2.5 flex-wrap min-w-0">
            <span className="text-[14px] font-semibold text-mathua-primary truncate min-w-0">{selected.label}</span>
            <span
              className="text-[10px] uppercase tracking-[0.05em] px-2 py-0.5 rounded-[3px] text-white shrink-0"
              style={{ background: domainColor(selected.domain) }}
            >
              {selected.domain.replace(/_/g, ' ')}
            </span>
            {selectedStatus && (
              <span className="text-[11px] shrink-0" style={{ color: STATUS_COLORS[selectedStatus] }}>● {STATUS_LABELS[selectedStatus]}</span>
            )}
            <button
              type="button"
              onClick={() => onNodeSelect?.(selected.id)}
              className="text-[12px] font-mono text-[#60a5fa] bg-transparent border border-mathua-blue/30 rounded-none px-3 py-1.5 min-h-[36px] sm:min-h-[36px] w-full sm:w-auto sm:ml-auto shrink-0 inline-flex items-center justify-center hover:bg-mathua-blue hover:text-white transition-colors"
            >
              Open concept →
            </button>
          </div>
          {(prereqList.length > 0 || unlocksList.length > 0) && (
            <div className="flex gap-4 sm:gap-6 mt-2.5 flex-wrap min-w-0">
              {prereqList.length > 0 && (
                <div className="min-w-0 flex-1 sm:flex-none">
                  <div className="text-mathua-muted text-[11px] uppercase mb-1">
                    requires
                  </div>
                  <div className="flex gap-1.5 flex-wrap">
                    {prereqList.map(p => (
                      <button
                        key={p.id}
                        type="button"
                        onClick={() => select(p.id)}
                        className="text-[11px] font-mono text-mathua-secondary bg-mathua-surface-elevated border-[0.5px] border-mathua-border rounded-[3px] px-2.5 py-1.5 min-h-[36px] cursor-pointer hover:border-mathua-blue hover:text-mathua-blue transition-colors"
                      >
                        {p.label}
                      </button>
                    ))}
                  </div>
                </div>
              )}
              {unlocksList.length > 0 && (
                <div className="min-w-0 flex-1 sm:flex-none">
                  <div className="text-mathua-muted text-[11px] uppercase mb-1">
                    unlocks
                  </div>
                  <div className="flex gap-1.5 flex-wrap">
                    {unlocksList.map(u => (
                      <button
                        key={u.id}
                        type="button"
                        onClick={() => select(u.id)}
                        className="text-[11px] font-mono text-mathua-secondary bg-mathua-surface-elevated border-[0.5px] border-mathua-border rounded-[3px] px-2.5 py-1.5 min-h-[36px] cursor-pointer hover:border-mathua-blue hover:text-mathua-blue transition-colors"
                      >
                        {u.label}
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  )
}

export default function ConceptGraphFlow(props: ConceptGraphFlowProps) {
  return (
    <ReactFlowProvider>
      <GraphInner {...props} />
    </ReactFlowProvider>
  )
}
