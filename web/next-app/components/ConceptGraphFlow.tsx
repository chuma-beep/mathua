'use client'

import { useMemo, useState, useEffect, useCallback } from 'react'
import {
  ReactFlow,
  Background,
  BackgroundVariant,
  Controls,
  MiniMap,
  Handle,
  Position,
  MarkerType,
  type Node,
  type Edge,
  type NodeProps,
} from '@xyflow/react'
import '@xyflow/react/dist/base.css'
import dagre from 'dagre'

export type MasteryStatus = 'mastered' | 'practicing' | 'learning' | 'unseen' | 'locked'

export interface GraphConcept {
  id: string
  label: string
  domain: string
  prerequisites: string[]
}

interface ConceptGraphFlowProps {
  concepts: GraphConcept[]
  conceptStatuses?: Record<string, MasteryStatus>
  theme?: 'dark' | 'light'
  onPathNodes?: string[]
  onNodeSelect?: (id: string) => void
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
  onPath: boolean
  dimmed: boolean
  selected: boolean
}

type ConceptFlowNode = Node<ConceptNodeData, 'concept'>

function ConceptNode({ data }: NodeProps<ConceptFlowNode>) {
  const statusColor = STATUS_COLORS[data.status] ?? STATUS_COLORS.unseen
  return (
    <div
      style={{
        width: 180,
        height: 40,
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
        opacity: data.dimmed ? 0.22 : 1,
        transition: 'opacity 150ms ease',
        fontFamily: "'IBM Plex Mono', monospace",
        fontSize: 11,
        color: 'var(--text-primary)',
        overflow: 'hidden',
      }}
    >
      <Handle type="target" position={Position.Left} style={{ opacity: 0, pointerEvents: 'none' }} />
      <span
        style={{
          width: 3,
          height: 22,
          borderRadius: 2,
          flexShrink: 0,
          background: domainColor(data.domain),
        }}
      />
      <span
        style={{
          width: 8,
          height: 8,
          borderRadius: '50%',
          flexShrink: 0,
          background: statusColor,
        }}
      />
      <span style={{ whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
        {data.label}
      </span>
      <Handle type="source" position={Position.Right} style={{ opacity: 0, pointerEvents: 'none' }} />
    </div>
  )
}

const nodeTypes = { concept: ConceptNode }

function layoutWithDagre(concepts: GraphConcept[]): Map<string, { x: number; y: number }> {
  const g = new dagre.graphlib.Graph()
  g.setGraph({ rankdir: 'LR', nodesep: 16, ranksep: 70, marginx: 20, marginy: 20 })
  g.setDefaultEdgeLabel(() => ({}))
  for (const c of concepts) g.setNode(c.id, { width: 180, height: 40 })
  for (const c of concepts) {
    for (const p of c.prerequisites) {
      if (concepts.some(x => x.id === p)) g.setEdge(p, c.id)
    }
  }
  dagre.layout(g)
  const positions = new Map<string, { x: number; y: number }>()
  for (const c of concepts) {
    const n = g.node(c.id)
    if (n) positions.set(c.id, { x: n.x - 90, y: n.y - 20 })
  }
  return positions
}

export default function ConceptGraphFlow({
  concepts,
  conceptStatuses,
  theme = 'dark',
  onPathNodes,
  onNodeSelect,
}: ConceptGraphFlowProps) {
  const [isMobile, setIsMobile] = useState(false)
  const [selectedId, setSelectedId] = useState<string | null>(null)

  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 640)
    checkMobile()
    window.addEventListener('resize', checkMobile)
    return () => window.removeEventListener('resize', checkMobile)
  }, [])

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
    if (!selectedId) return null
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
    walkUp(selectedId)
    walkDown(selectedId)
    return { upstream, downstream }
  }, [selectedId, conceptById, dependentsMap])

  const flowNodes = useMemo<ConceptFlowNode[]>(() => {
    const positions = layoutWithDagre(concepts)
    return concepts.map(c => {
      const inNeighborhood =
        !neighborhood ||
        c.id === selectedId ||
        neighborhood.upstream.has(c.id) ||
        neighborhood.downstream.has(c.id)
      return {
        id: c.id,
        type: 'concept' as const,
        position: positions.get(c.id) ?? { x: 0, y: 0 },
        data: {
          label: c.label,
          domain: c.domain,
          status: conceptStatuses?.[c.id] ?? 'unseen',
          onPath: onPathNodes?.includes(c.id) ?? false,
          dimmed: !inNeighborhood,
          selected: c.id === selectedId,
        },
      }
    })
  }, [concepts, conceptStatuses, onPathNodes, neighborhood, selectedId])

  const flowEdges = useMemo<Edge[]>(() => {
    const known = new Set(concepts.map(c => c.id))
    const edges: Edge[] = []
    for (const c of concepts) {
      for (const p of c.prerequisites) {
        if (!known.has(p)) continue
        const highlighted =
          !!selectedId && (p === selectedId || c.id === selectedId)
        const inNeighborhood =
          !neighborhood ||
          (neighborhood.upstream.has(p) &&
            (neighborhood.upstream.has(c.id) || c.id === selectedId)) ||
          (p === selectedId && neighborhood.downstream.has(c.id)) ||
          (neighborhood.downstream.has(c.id) &&
            (neighborhood.downstream.has(p) || p === selectedId))
        edges.push({
          id: `${p}->${c.id}`,
          source: p,
          target: c.id,
          style: {
            stroke: highlighted
              ? 'var(--accent-blue)'
              : inNeighborhood
                ? 'var(--border-strong)'
                : 'var(--border)',
            opacity: neighborhood && !inNeighborhood ? 0.15 : 1,
            strokeWidth: highlighted ? 1.5 : 1,
          },
          markerEnd: { type: MarkerType.ArrowClosed, width: 12, height: 12 },
        })
      }
    }
    return edges
  }, [concepts, selectedId, neighborhood])

  const handleNodeClick = useCallback(
    (_: React.MouseEvent, node: Node) => {
      setSelectedId(node.id)
    },
    []
  )

  const selected = selectedId ? conceptById.get(selectedId) : null
  const selectedStatus = selectedId ? conceptStatuses?.[selectedId] ?? 'unseen' : null

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
        style={{
          height: isMobile ? 320 : 520,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          border: '0.5px solid var(--border)',
          background: 'var(--graph-surface)',
          color: 'var(--text-muted)',
          fontSize: 14,
          fontFamily: "'IBM Plex Mono', monospace",
        }}
      >
        No concepts to display
      </div>
    )
  }

  return (
    <div>
      <div
        style={{
          height: isMobile ? 320 : 520,
          width: '100%',
          overflow: 'hidden',
          border: '0.5px solid var(--border)',
          background: 'var(--graph-surface)',
        }}
      >
        <ReactFlow
          nodes={flowNodes}
          edges={flowEdges}
          nodeTypes={nodeTypes}
          onNodeClick={handleNodeClick}
          onPaneClick={() => setSelectedId(null)}
          fitView
          fitViewOptions={{ padding: 0.15, maxZoom: 1 }}
          minZoom={0.03}
          maxZoom={2.5}
          onlyRenderVisibleElements
          nodesDraggable={false}
          nodesConnectable={false}
          elementsSelectable
          proOptions={{ hideAttribution: false }}
        >
          <Background
            variant={BackgroundVariant.Dots}
            gap={24}
            size={1}
            color={theme === 'dark' ? '#2a2a30' : '#d4d4d8'}
          />
          <Controls
            showInteractive={false}
            style={{
              background: 'var(--surface-elevated)',
              borderColor: 'var(--border)',
              color: 'var(--text-primary)',
            }}
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
      </div>

      {selected && (
        <div
          style={{
            marginTop: 10,
            padding: '12px 14px',
            border: '0.5px solid var(--border)',
            background: 'var(--surface)',
            fontFamily: "'IBM Plex Mono', monospace",
          }}
        >
          <div style={{ display: 'flex', alignItems: 'center', gap: 10, flexWrap: 'wrap' }}>
            <span style={{ fontSize: 14, fontWeight: 600, color: 'var(--text-primary)' }}>
              {selected.label}
            </span>
            <span
              style={{
                fontSize: 10,
                textTransform: 'uppercase',
                letterSpacing: '0.05em',
                padding: '2px 8px',
                borderRadius: 3,
                background: domainColor(selected.domain),
                color: '#fff',
              }}
            >
              {selected.domain.replace(/_/g, ' ')}
            </span>
            {selectedStatus && (
              <span style={{ fontSize: 11, color: STATUS_COLORS[selectedStatus] }}>
                ● {STATUS_LABELS[selectedStatus]}
              </span>
            )}
            <button
              onClick={() => onNodeSelect?.(selected.id)}
              style={{
                marginLeft: 'auto',
                fontSize: 12,
                fontFamily: 'inherit',
                color: '#60a5fa',
                background: 'none',
                border: 'none',
                cursor: 'pointer',
              }}
            >
              Open concept →
            </button>
          </div>
          {(prereqList.length > 0 || unlocksList.length > 0) && (
            <div style={{ display: 'flex', gap: 24, marginTop: 10, flexWrap: 'wrap' }}>
              {prereqList.length > 0 && (
                <div>
                  <div
                    style={{
                      color: 'var(--text-muted)',
                      fontSize: 11,
                      textTransform: 'uppercase',
                      marginBottom: 4,
                    }}
                  >
                    requires
                  </div>
                  <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap' }}>
                    {prereqList.map(p => (
                      <button
                        key={p.id}
                        onClick={() => setSelectedId(p.id)}
                        style={{
                          fontSize: 11,
                          fontFamily: 'inherit',
                          color: 'var(--text-secondary)',
                          background: 'var(--surface-elevated)',
                          border: '0.5px solid var(--border)',
                          borderRadius: 3,
                          padding: '2px 8px',
                          cursor: 'pointer',
                        }}
                      >
                        {p.label}
                      </button>
                    ))}
                  </div>
                </div>
              )}
              {unlocksList.length > 0 && (
                <div>
                  <div
                    style={{
                      color: 'var(--text-muted)',
                      fontSize: 11,
                      textTransform: 'uppercase',
                      marginBottom: 4,
                    }}
                  >
                    unlocks
                  </div>
                  <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap' }}>
                    {unlocksList.map(u => (
                      <button
                        key={u.id}
                        onClick={() => setSelectedId(u.id)}
                        style={{
                          fontSize: 11,
                          fontFamily: 'inherit',
                          color: 'var(--text-secondary)',
                          background: 'var(--surface-elevated)',
                          border: '0.5px solid var(--border)',
                          borderRadius: 3,
                          padding: '2px 8px',
                          cursor: 'pointer',
                        }}
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
