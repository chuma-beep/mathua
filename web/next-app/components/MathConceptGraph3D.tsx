'use client'

/* eslint-disable react/no-unknown-property -- R3F/Drei JSX elements use non-HTML attributes */

import React, { useRef, useState, useCallback, useMemo, useEffect } from 'react'
import { Canvas, useFrame } from '@react-three/fiber'
import { OrbitControls, Html } from '@react-three/drei'
import * as THREE from 'three'
import { layoutDAG3D, type ConceptLayoutInput } from '../lib/layoutDAG3D'

export interface ConceptDef {
  id: string
  label: string
  domain: string
  prerequisites: string[]
}

export type MasteryStatus = 'mastered' | 'practicing' | 'learning' | 'unseen' | 'locked'

interface MathConceptGraph3DProps {
  theme?: 'dark' | 'light'
  concepts: ConceptDef[]
  conceptStatuses?: Record<string, MasteryStatus>
  onPathNodes?: string[]
  onNodeSelect?: (nodeId: string) => void
}

const DOMAIN_COLORS: Record<string, string> = {
  counting:             '#c8a96e',
  arithmetic:           '#4db8a0',
  fractions:            '#a8a0f0',
  prealgebra:           '#7dd3fc',
  algebra:              '#e8a849',
  geometry:             '#86efac',
  trigonometry:         '#fda4af',
  complex_numbers:      '#e879f9',
  precalculus:          '#fda4af',
  calculus:             '#e879f9',
  linear_algebra:       '#86efac',
  statistics:           '#67e8f9',
  discrete_math:        '#f0abab',
  number_theory:        '#a8e6cf',
  differential_equations:'#fdba74',
  abstract_algebra:     '#c4b5fd',
  topology:             '#f9a8d4',
}

const DOMAIN_COLORS_LIGHT: Record<string, string> = {
  counting:             '#a0814a',
  arithmetic:           '#3a9a8a',
  fractions:            '#8a80d8',
  prealgebra:           '#5ab8dc',
  algebra:              '#c08a30',
  geometry:             '#68c88c',
  trigonometry:         '#e88a99',
  complex_numbers:      '#c868e8',
  precalculus:          '#e88a99',
  calculus:             '#c868e8',
  linear_algebra:       '#68c88c',
  statistics:           '#50c8d8',
  discrete_math:        '#d08a8a',
  number_theory:        '#80c8a8',
  differential_equations:'#d09050',
  abstract_algebra:     '#a090d0',
  topology:             '#d080b0',
}

const FALLBACK_COLOR = '#5a6577'
const FALLBACK_COLOR_LIGHT = '#888'

function domainColor(domain: string, theme: 'dark' | 'light'): string {
  const map = theme === 'dark' ? DOMAIN_COLORS : DOMAIN_COLORS_LIGHT
  return map[domain] ?? (theme === 'dark' ? FALLBACK_COLOR : FALLBACK_COLOR_LIGHT)
}

const LINK_COLOR = '#8ed8d0'
const LINK_COLOR_LIGHT = '#207068'

const ACTIVE_LINK_COLOR = '#60a5fa'
const ACTIVE_LINK_COLOR_LIGHT = '#2563eb'

const NODE_RADIUS = 0.20

interface RenderNode {
  id: string
  name: string
  domain: string
  position: [number, number, number]
  status: MasteryStatus | null
  onPath: boolean
}

interface Link {
  source: string
  target: string
}

interface GraphSceneProps {
  nodes: RenderNode[]
  links: Link[]
  activeId: string
  onSelect: (id: string) => void
  theme: 'dark' | 'light'
}

const STATUS_COLORS_DARK: Record<string, string> = {
  mastered:   '#4db8a0',
  practicing: '#60a5fa',
  learning:   '#e8a849',
  unseen:     '#5a6577',
  locked:     '#2a2d35',
}

const STATUS_COLORS_LIGHT: Record<string, string> = {
  mastered:   '#3a9a8a',
  practicing: '#3b82f6',
  learning:   '#c08a30',
  unseen:     '#9ca3af',
  locked:     '#d0d0d0',
}

function nodeDisplayColor(node: RenderNode, theme: 'dark' | 'light'): string {
  const showStatus = node.status !== null

  if (showStatus) {
    if (node.onPath) return theme === 'dark' ? '#ffdd88' : '#b08020'
    if (node.status === 'locked' || node.status === 'unseen') return theme === 'dark' ? '#2a2d35' : '#d0d0d0'
    const statusMap = theme === 'dark' ? STATUS_COLORS_DARK : STATUS_COLORS_LIGHT
    return statusMap[node.status] ?? domainColor(node.domain, theme)
  }

  if (node.onPath) return theme === 'dark' ? '#ffdd88' : '#b08020'
  return domainColor(node.domain, theme)
}

const monoFont = "'IBM Plex Mono', monospace"
const serifFont = "'IBM Plex Serif', serif"
const bodyFont = "'IBM Plex Serif', serif"

const tooltipStyle: React.CSSProperties = {
  background: 'var(--bg)',
  border: '0.5px solid var(--border)',
  borderRadius: 0,
  whiteSpace: 'nowrap',
  minWidth: 'max-content',
  boxShadow: '0 4px 20px rgba(0,0,0,0.3)',
}

const NodeMesh = React.memo(function NodeMesh({
  node,
  isActive,
  isHovered,
  onHover,
  onSelect,
  theme,
}: {
  node: RenderNode
  isActive: boolean
  isHovered: boolean
  onHover: (id: string | null) => void
  onSelect: (id: string) => void
  theme: 'dark' | 'light'
}) {
  const meshRef = useRef<THREE.Mesh>(null)
  const ringRef = useRef<THREE.Mesh>(null)
  const color = nodeDisplayColor(node, theme)
  const showStatus = node.status !== null && node.status !== 'locked' && node.status !== 'unseen'
  const statusLabel = showStatus ? node.status : node.domain
  const [isMobile, setIsMobile] = useState(false)

  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 640)
    checkMobile()
    window.addEventListener('resize', checkMobile)
    return () => window.removeEventListener('resize', checkMobile)
  }, [])

  useFrame(() => {
    if (meshRef.current) {
      const mat = meshRef.current.material as THREE.MeshStandardMaterial
      const targetIntensity = isActive ? 2.0 : isHovered ? 1.2 : node.onPath ? 0.8 : node.status === 'locked' ? 0.1 : 0.3
      mat.emissiveIntensity = THREE.MathUtils.lerp(mat.emissiveIntensity, targetIntensity, 0.25)
    }
    if (ringRef.current) {
      const mat = ringRef.current.material as THREE.MeshBasicMaterial
      const targetOpacity = node.onPath ? 0.18 : 0.04
      mat.opacity = THREE.MathUtils.lerp(mat.opacity, targetOpacity, 0.2)
    }
  })

  return (
    <group position={node.position as [number, number, number]}>
      <mesh
        ref={ringRef}
        onPointerOver={() => onHover(node.id)}
        onPointerOut={() => onHover(null)}
        onClick={() => onSelect(node.id)}
      >
        <sphereGeometry args={[NODE_RADIUS * 2.2, 12, 12]} />
        <meshBasicMaterial color={node.onPath ? '#c8a96e' : color} transparent opacity={0.04} depthWrite={false} />
      </mesh>
      <mesh ref={meshRef}>
        <sphereGeometry args={[NODE_RADIUS, 16, 16]} />
        <meshStandardMaterial
          color={color}
          emissive={color}
          emissiveIntensity={node.status === 'locked' ? 0.05 : 0.3}
          roughness={node.status === 'locked' ? 0.7 : 0.2}
          metalness={0.1}
        />
      </mesh>
      {isHovered && (
          <Html center distanceFactor={isMobile ? 18 : 12} style={{ pointerEvents: 'none', zIndex: 20 }}>
          <div style={{ ...tooltipStyle, padding: isMobile ? '10px 14px' : '8px 12px', maxWidth: isMobile ? '220px' : '240px' }}>
            <div style={{ color: 'var(--text-primary)', fontSize: isMobile ? '14px' : '12px', fontFamily: monoFont }}>{node.name}</div>
            <div style={{ color, fontSize: isMobile ? '13px' : '12px', textTransform: 'uppercase', marginTop: '3px', fontFamily: monoFont }}>{statusLabel}</div>
            {showStatus && (
              <div style={{ color: 'var(--text-muted)', fontSize: isMobile ? '12px' : '11px', marginTop: '2px', textTransform: 'uppercase', fontFamily: monoFont }}>
                {node.status} {node.onPath ? '· on path' : ''}
              </div>
            )}
            {!showStatus && (
              <div style={{ color: 'var(--text-muted)', fontSize: isMobile ? '12px' : '11px', marginTop: '2px', fontFamily: monoFont }}>
                Not started
              </div>
            )}
          </div>
        </Html>
      )}
    </group>
  )
})

function EdgeLines({ links, positionMap, activeId, theme }: { links: Link[], positionMap: Map<string, [number, number, number]>, activeId: string, theme: 'dark' | 'light' }) {
  const activeLinks = useMemo(() => links.filter(l => l.source === activeId || l.target === activeId), [links, activeId])
  const edgeColor = theme === 'dark' ? ACTIVE_LINK_COLOR : ACTIVE_LINK_COLOR_LIGHT

  const geometry = useMemo(() => {
    const positions = new Float32Array(activeLinks.length * 6)
    const colors = new Float32Array(activeLinks.length * 6)

    activeLinks.forEach((link, i) => {
      const sp = positionMap.get(link.source)
      const tp = positionMap.get(link.target)
      if (!sp || !tp) return

      const idx = i * 6
      positions[idx]   = sp[0]; positions[idx+1] = sp[1]; positions[idx+2] = sp[2]
      positions[idx+3] = tp[0]; positions[idx+4] = tp[1]; positions[idx+5] = tp[2]

      const c = new THREE.Color(edgeColor)
      colors[idx] = c.r; colors[idx+1] = c.g; colors[idx+2] = c.b
      colors[idx+3] = c.r; colors[idx+4] = c.g; colors[idx+5] = c.b
    })

    const geo = new THREE.BufferGeometry()
    geo.setAttribute('position', new THREE.BufferAttribute(positions, 3))
    geo.setAttribute('color', new THREE.BufferAttribute(colors, 3))
    return geo
  }, [activeLinks, edgeColor])

  const lineRef = useRef<THREE.LineSegments>(null)
  const clockRef = useRef(0)

  useFrame((_, delta) => {
    clockRef.current += delta
    if (lineRef.current) {
      const mat = lineRef.current.material as THREE.LineBasicMaterial
      mat.opacity = 0.6 + 0.4 * Math.sin(clockRef.current * 4)
    }
  })

  if (activeLinks.length === 0) return null

  return (
    <lineSegments ref={lineRef} geometry={geometry}>
      <lineBasicMaterial vertexColors transparent opacity={theme === 'dark' ? 0.9 : 1.0} />
    </lineSegments>
  )
}

function AllEdges({ links, positionMap, theme }: { links: Link[], positionMap: Map<string, [number, number, number]>, theme: 'dark' | 'light' }) {
  const edgeColor = theme === 'dark' ? LINK_COLOR : LINK_COLOR_LIGHT

  const geometry = useMemo(() => {
    const positions = new Float32Array(links.length * 6)
    const colors = new Float32Array(links.length * 6)

    links.forEach((link, i) => {
      const sp = positionMap.get(link.source)
      const tp = positionMap.get(link.target)
      if (!sp || !tp) return

      const idx = i * 6
      positions[idx]   = sp[0]; positions[idx+1] = sp[1]; positions[idx+2] = sp[2]
      positions[idx+3] = tp[0]; positions[idx+4] = tp[1]; positions[idx+5] = tp[2]

      const c = new THREE.Color(edgeColor)
      colors[idx] = c.r; colors[idx+1] = c.g; colors[idx+2] = c.b
      colors[idx+3] = c.r; colors[idx+4] = c.g; colors[idx+5] = c.b
    })

    const geo = new THREE.BufferGeometry()
    geo.setAttribute('position', new THREE.BufferAttribute(positions, 3))
    geo.setAttribute('color', new THREE.BufferAttribute(colors, 3))
    return geo
  }, [links, edgeColor])

  return (
    <lineSegments geometry={geometry}>
      <lineBasicMaterial vertexColors transparent opacity={theme === 'dark' ? 0.15 : 0.25} />
    </lineSegments>
  )
}

function Particles({ theme }: { theme: 'dark' | 'light' }) {
  const positions = useMemo(() => {
    const pos = new Float32Array(800 * 3)
    for (let i = 0; i < 800; i++) {
      const theta = Math.random() * Math.PI * 2
      const phi = Math.acos(2 * Math.random() - 1)
      const r = 60 * Math.cbrt(Math.random())
      pos[i*3] = r * Math.sin(phi) * Math.cos(theta)
      pos[i*3+1] = r * Math.sin(phi) * Math.sin(theta)
      pos[i*3+2] = r * Math.cos(phi)
    }
    return pos
  }, [])

  const color = theme === 'dark' ? '#4a5568' : '#aaa'

  return (
    <points>
      <bufferGeometry>
        <bufferAttribute attach="attributes-position" count={800} array={positions} itemSize={3} />
      </bufferGeometry>
      <pointsMaterial size={0.06} color={color} transparent opacity={0.8} sizeAttenuation />
    </points>
  )
}

function GraphScene({ nodes, links, activeId, positionMap, onSelect, theme }: GraphSceneProps & { positionMap: Map<string, [number, number, number]> }) {
  const [hovered, setHovered] = useState<string | null>(null)
  const controlsRef = useRef<any>(null)

  useEffect(() => {
    if (controlsRef.current) {
      controlsRef.current.autoRotate = hovered === null
    }
  }, [hovered])

  return (
    <>
      <ambientLight intensity={0.4} />
      <pointLight position={[10, 10, 10]} intensity={0.8} />
      <pointLight position={[-10, -10, -10]} intensity={0.4} color={theme === 'dark' ? '#c8a96e' : '#a0814a'} />

      <AllEdges links={links} positionMap={positionMap} theme={theme} />
      <EdgeLines links={links} positionMap={positionMap} activeId={activeId} theme={theme} />

      {nodes.map(node => (
        <NodeMesh
          key={node.id}
          node={node}
          isActive={node.id === activeId}
          isHovered={node.id === hovered}
          onHover={setHovered}
          onSelect={onSelect}
          theme={theme}
        />
      ))}

      <Particles theme={theme} />

      <OrbitControls
        ref={controlsRef}
        enablePan={false}
        enableZoom={true}
        minDistance={12}
        maxDistance={45}
        autoRotate={true}
        autoRotateSpeed={1.2}
        dampingFactor={0.05}
        enableDamping={true}
      />
    </>
  )
}

function InfoPanel({ activeId, concepts, conceptStatuses, onPathNodes, theme }: {
  activeId: string
  concepts: ConceptDef[]
  conceptStatuses?: Record<string, MasteryStatus>
  onPathNodes?: string[]
  theme: 'dark' | 'light'
}) {
  const [isMobile, setIsMobile] = useState(false)
  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 640)
    checkMobile()
    window.addEventListener('resize', checkMobile)
    return () => window.removeEventListener('resize', checkMobile)
  }, [])

  const concept = useMemo(() => concepts.find(c => c.id === activeId), [concepts, activeId])
  if (!concept) return null

  const status = conceptStatuses?.[concept.id] ?? null
  const onPath = onPathNodes?.includes(concept.id) ?? false
  const color = domainColor(concept.domain, theme)

  const prereqConcepts = concept.prerequisites
    .flatMap(id => {
      const c = concepts.find(c => c.id === id)
      return c ? [c] : []
    }) as ConceptDef[]

  const unlockedBy = concepts.filter(c => c.prerequisites.includes(concept.id))

  return (
    <div style={{
      background: 'transparent',
      borderTop: '0.5px solid var(--border)',
      padding: isMobile ? '0.5rem' : '0.75rem 1rem',
      borderRadius: 0,
    }}>
      <div style={{ color: 'var(--text-muted)', fontSize: '12px', textTransform: 'uppercase', fontFamily: monoFont }}>selected concept</div>
      <div style={{ color: 'var(--text-primary)', fontSize: isMobile ? '16px' : '18px', marginTop: '4px', fontFamily: serifFont, fontWeight: 400 }}>{concept.label}</div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginTop: '8px', flexWrap: 'wrap' as const }}>
        <span style={{ color, fontSize: isMobile ? '12px' : '12px', fontFamily: monoFont }}>● {concept.domain}</span>
        {status && (
          <span style={{
            color: status === 'mastered' ? 'var(--accent-teal)' : status === 'learning' ? '#e8a849' : 'var(--text-muted)',
            fontSize: isMobile ? '12px' : '12px',
            textTransform: 'uppercase',
            fontFamily: monoFont,
          }}>
            {status}
          </span>
        )}
        {onPath && (
          <span style={{ color: 'var(--accent-blue)', fontSize: isMobile ? '12px' : '12px', fontFamily: monoFont }}>on path</span>
        )}
      </div>
      {prereqConcepts.length > 0 && (
        <div style={{ marginTop: '8px' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '12px', textTransform: 'uppercase', marginBottom: '4px', fontFamily: monoFont }}>prerequisites</div>
          <div style={{ display: 'flex', gap: '6px', flexWrap: 'wrap' as const }}>
            {prereqConcepts.map(p => (
              <span key={p.id} style={{ color: 'var(--text-secondary)', fontSize: '12px', fontFamily: monoFont }}>
                {p.label}
              </span>
            ))}
          </div>
        </div>
      )}
      {unlockedBy.length > 0 && (
        <div style={{ marginTop: '6px' }}>
          <div style={{ color: 'var(--text-muted)', fontSize: '12px', textTransform: 'uppercase', marginBottom: '4px', fontFamily: monoFont }}>unlocks</div>
          <div style={{ display: 'flex', gap: '6px', flexWrap: 'wrap' as const }}>
            {unlockedBy.map(u => (
              <span key={u.id} style={{ color: 'var(--text-secondary)', fontSize: '12px', fontFamily: monoFont }}>
                {u.label}
              </span>
            ))}
          </div>
        </div>
      )}
      <div style={{ marginTop: '10px', display: 'flex', gap: '8px' }}>
        <a
          href={`/study?concept=${encodeURIComponent(concept.id)}`}
          style={{ color: '#60a5fa', fontSize: '12px', fontFamily: monoFont, textDecoration: 'none' }}
        >
          Study →
        </a>
        <a
          href={`/session`}
          style={{ color: 'var(--accent-teal)', fontSize: '12px', fontFamily: monoFont, textDecoration: 'none' }}
        >
          Practice →
        </a>
      </div>
    </div>
  )
}

export default function MathConceptGraph3D({
  theme = 'dark',
  concepts,
  conceptStatuses,
  onPathNodes,
  onNodeSelect,
}: MathConceptGraph3DProps) {
  const [isMobile, setIsMobile] = useState(false)

  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < 640)
    checkMobile()
    window.addEventListener('resize', checkMobile)
    return () => window.removeEventListener('resize', checkMobile)
  }, [])

  const { nodes, links, positionMap } = useMemo(() => {
    const input: ConceptLayoutInput[] = concepts.map(c => ({ id: c.id, prerequisites: c.prerequisites }))
    const layout = layoutDAG3D(input)

    const allLinks: Link[] = []
    for (const c of concepts) {
      for (const prereq of c.prerequisites) {
        allLinks.push({ source: prereq, target: c.id })
      }
    }

    const nodeList: RenderNode[] = concepts.map(c => ({
      id: c.id,
      name: c.label,
      domain: c.domain,
      position: layout.positions[c.id] ?? [0, 0, 0],
      status: conceptStatuses?.[c.id] ?? null,
      onPath: onPathNodes?.includes(c.id) ?? false,
    }))

    const pm = new Map<string, [number, number, number]>()
    for (const n of nodeList) {
      pm.set(n.id, n.position)
    }

    return { nodes: nodeList, links: allLinks, positionMap: pm }
  }, [concepts, conceptStatuses, onPathNodes])

  const [activeId, setActiveId] = useState(() => nodes[0]?.id ?? '')

  useEffect(() => {
    if (nodes.length > 0 && !nodes.find(n => n.id === activeId)) {
      setActiveId(nodes[0].id)
    }
  }, [nodes, activeId])

  return (
    <div>
      <div style={{
        height: isMobile ? '320px' : '520px',
        width: '100%',
        borderRadius: 0,
        overflow: 'hidden',
        border: '0.5px solid var(--border)',
      }}>
        <Canvas
          camera={{ position: [0, 0, 28], fov: 60 }}
          gl={{ alpha: true, premultipliedAlpha: true }}
          dpr={[1, 2]}
        >
          <GraphScene
            nodes={nodes}
            links={links}
            activeId={activeId}
            positionMap={positionMap}
            onSelect={(id) => {
              setActiveId(id)
              onNodeSelect?.(id)
            }}
            theme={theme}
          />
        </Canvas>
      </div>
      <InfoPanel
        activeId={activeId}
        concepts={concepts}
        conceptStatuses={conceptStatuses}
        onPathNodes={onPathNodes}
        theme={theme}
      />
    </div>
  )
}
