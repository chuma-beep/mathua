'use client'

/* eslint-disable react/no-unknown-property -- R3F/Drei JSX elements use non-HTML attributes */

import React, { useRef, useState, useCallback, useMemo, useEffect, useLayoutEffect } from 'react'
import { Canvas, useFrame, useThree, type ThreeEvent } from '@react-three/fiber'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
import * as THREE from 'three'
import { layoutDAG3D, type ConceptLayoutInput } from '../lib/layoutDAG3D'
import { loadGraphPayload, loadPositionEntries } from '../lib/graphPositions'
import type { GraphPayload } from '../lib/graphPayload'
import { domainColor, domainLabel, DOMAIN_ORDER } from '../lib/graphDomains'
import { useIsMobile } from '../hooks/useIsMobile'
import { type MasteryStatus } from '../lib/graphStatus'

export type { MasteryStatus }

export interface ConceptDef {
  id: string
  label: string
  domain: string
  prerequisites: string[]
}

interface MathConceptGraph3DProps {
  theme?: 'dark' | 'light'
  conceptStatuses?: Record<string, MasteryStatus>
  onPathNodes?: string[]
  onNodeSelect?: (nodeId: string) => void
}

type Vec3 = [number, number, number]
type LodTier = 'far' | 'mid' | 'close'

const LINK_COLOR = '#8ed8d0'
const LINK_COLOR_LIGHT = '#207068'
const ACTIVE_LINK_COLOR = '#60a5fa'
const ACTIVE_LINK_COLOR_LIGHT = '#2563eb'

const NODE_RADIUS = 0.2
const LABEL_POOL = 72
const HIGHLIGHT = new THREE.Color('#ffffff')
const GOLD = new THREE.Color('#c8a96e')
const DIM_DARK = new THREE.Color('#11151f')
const DIM_LIGHT = new THREE.Color('#eae8e0')

const EDGE_ALPHA: Record<LodTier, number> = { far: 0.07, mid: 0.16, close: 0.26 }
const EDGE_ALPHA_LIGHT: Record<LodTier, number> = { far: 0.1, mid: 0.2, close: 0.32 }
const FOCUS_EDGE_ALPHA = 0.6
const UNFOCUSED_EDGE_FACTOR = 0.12

const STATUS_COLORS_DARK = {
  mastered:   '#4db8a0',
  practicing: '#60a5fa',
  learning:   '#e8a849',
  unseen:     '#5a6577',
  locked:     '#2a2d35',
} satisfies Record<string, string>

const STATUS_COLORS_LIGHT = {
  mastered:   '#3a9a8a',
  practicing: '#3b82f6',
  learning:   '#c08a30',
  unseen:     '#9ca3af',
  locked:     '#d0d0d0',
} satisfies Record<string, string>

const monoFont = "var(--font-ibm-plex-mono), 'IBM Plex Mono', monospace"
const serifFont = "var(--font-ibm-plex-serif), 'IBM Plex Serif', serif"

interface RenderNode {
  id: string
  name: string
  domain: string
  position: Vec3
  status: MasteryStatus | null
  onPath: boolean
  importance: number
}

interface Link {
  source: string
  target: string
}

interface BuiltGraph {
  nodes: RenderNode[]
  links: Link[]
  positionMap: Map<string, Vec3>
  nodeById: Map<string, RenderNode>
  prereqsById: Map<string, string[]>
  unlocksById: Map<string, string[]>
  domainCentroids: Record<string, Vec3>
  byImportance: RenderNode[]
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

function usePrefersReducedMotion(): boolean {
  const [reduced, setReduced] = useState(false)
  useEffect(() => {
    if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') return
    const mql = window.matchMedia('(prefers-reduced-motion: reduce)')
    setReduced(mql.matches)
    const onChange = (e: MediaQueryListEvent) => setReduced(e.matches)
    mql.addEventListener('change', onChange)
    return () => mql.removeEventListener('change', onChange)
  }, [])
  return reduced
}

function useCoarsePointer(): boolean {
  const [coarse, setCoarse] = useState(false)
  useEffect(() => {
    if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') return
    const mql = window.matchMedia('(pointer: coarse)')
    setCoarse(mql.matches)
    const onChange = (e: MediaQueryListEvent) => setCoarse(e.matches)
    mql.addEventListener('change', onChange)
    return () => mql.removeEventListener('change', onChange)
  }, [])
  return coarse
}

// ── Data ────────────────────────────────────────────────

function buildGraph(
  payload: GraphPayload,
  entries: Array<[string, Vec3]>,
  statuses: Record<string, MasteryStatus> | undefined,
  onPathNodes: string[] | undefined
): BuiltGraph {
  const positionMap = new Map<string, Vec3>(entries)
  const idAt = (i: number) => payload.nodes[i][0]

  const prereqIds: string[][] = payload.nodes.map(() => [])
  for (const [source, target] of payload.edges) prereqIds[target].push(idAt(source))

  let missing = false
  for (const [id] of payload.nodes) {
    if (!positionMap.has(id)) {
      missing = true
      break
    }
  }
  if (missing) {
    const input: ConceptLayoutInput[] = payload.nodes.map(([id, , domain], i) => ({
      id,
      prerequisites: prereqIds[i],
      domain,
    }))
    const fallback = layoutDAG3D(input, { domainOrder: [...DOMAIN_ORDER] })
    for (const [id, p] of Object.entries(fallback.positions)) positionMap.set(id, p)
  }

  const onPathSet = onPathNodes ? new Set(onPathNodes) : null
  const nodes: RenderNode[] = payload.nodes.map(([id, label, domain, importance]) => ({
    id,
    name: label,
    domain,
    position: positionMap.get(id) ?? [0, 0, 0],
    status: statuses?.[id] ?? null,
    onPath: onPathSet?.has(id) ?? false,
    importance,
  }))

  const links: Link[] = payload.edges.map(([source, target]) => ({
    source: idAt(source),
    target: idAt(target),
  }))

  const nodeById = new Map<string, RenderNode>()
  for (const n of nodes) nodeById.set(n.id, n)

  const prereqsById = new Map<string, string[]>()
  const unlocksById = new Map<string, string[]>()
  for (const link of links) {
    const pre = prereqsById.get(link.target)
    if (pre) pre.push(link.source)
    else prereqsById.set(link.target, [link.source])
    const post = unlocksById.get(link.source)
    if (post) post.push(link.target)
    else unlocksById.set(link.source, [link.target])
  }

  const centroidSums = new Map<string, { x: number; y: number; z: number; n: number }>()
  for (const n of nodes) {
    const acc = centroidSums.get(n.domain) ?? { x: 0, y: 0, z: 0, n: 0 }
    acc.x += n.position[0]
    acc.y += n.position[1]
    acc.z += n.position[2]
    acc.n += 1
    centroidSums.set(n.domain, acc)
  }
  const domainCentroids: Record<string, Vec3> = {}
  for (const [domain, acc] of centroidSums) {
    domainCentroids[domain] = [acc.x / acc.n, acc.y / acc.n, acc.z / acc.n]
  }

  const byImportance = [...nodes].sort((a, b) => b.importance - a.importance || a.id.localeCompare(b.id))

  return { nodes, links, positionMap, nodeById, prereqsById, unlocksById, domainCentroids, byImportance }
}

// Landing default: an early, self-contained concept reads better than the
// alphabetically-first advanced one (abstract.group.def). Prefer the known
// arithmetic root, then the highest-unlock root, then anything.
function pickDefaultId(graph: BuiltGraph): string {
  if (graph.nodeById.has('arith.add.single')) return 'arith.add.single'
  const roots = graph.nodes.filter(n => (graph.prereqsById.get(n.id) ?? []).length === 0)
  const pool = roots.length > 0 ? roots : graph.nodes
  const best = [...pool].sort((a, b) => b.importance - a.importance || a.id.localeCompare(b.id))[0]
  return best?.id ?? ''
}

// ── Instanced nodes ─────────────────────────────────────

function NodeInstances({
  nodes,
  activeId,
  hoveredId,
  focused,
  focusIds,
  theme,
  isMobile,
  hoverEnabled,
  onHover,
  onSelect,
}: {
  nodes: RenderNode[]
  activeId: string
  hoveredId: string | null
  focused: boolean
  focusIds: Set<string> | null
  theme: 'dark' | 'light'
  isMobile: boolean
  hoverEnabled: boolean
  onHover: (id: string | null) => void
  onSelect: (id: string) => void
}) {
  const meshRef = useRef<THREE.InstancedMesh>(null)
  const hoveredRef = useRef<string | null>(null)
  const count = nodes.length
  const segments = isMobile ? 8 : 12
  const maxImportance = useMemo(
    () => Math.max(1, ...nodes.map(n => n.importance)),
    [nodes]
  )

  const { geometry, material } = useMemo(() => {
    const geometry = new THREE.SphereGeometry(NODE_RADIUS, segments, segments)
    const material = new THREE.MeshStandardMaterial({
      roughness: 0.32,
      metalness: 0.08,
      emissive: new THREE.Color('#ffffff'),
      emissiveIntensity: 0.18,
    })
    material.onBeforeCompile = shader => {
      shader.fragmentShader = shader.fragmentShader.replace(
        '#include <emissivemap_fragment>',
        '#include <emissivemap_fragment>\n#ifdef USE_INSTANCING_COLOR\n\ttotalEmissiveRadiance *= vColor;\n#endif'
      )
    }
    return { geometry, material }
  }, [segments])

  useEffect(() => {
    return () => {
      geometry.dispose()
      material.dispose()
    }
  }, [geometry, material])

  // Positions + importance-scaled sizes are stable per layout: upload once.
  useLayoutEffect(() => {
    const mesh = meshRef.current
    if (!mesh) return
    const dummy = new THREE.Object3D()
    const logMax = Math.log1p(maxImportance)
    for (let i = 0; i < count; i++) {
      const n = nodes[i]
      const norm = logMax > 0 ? Math.log1p(n.importance) / logMax : 0
      dummy.position.set(n.position[0], n.position[1], n.position[2])
      dummy.scale.setScalar(0.8 + 0.5 * norm)
      dummy.updateMatrix()
      mesh.setMatrixAt(i, dummy.matrix)
    }
    mesh.instanceMatrix.needsUpdate = true
  }, [nodes, count, maxImportance])

  // Colours depend on theme + selection/hover: refresh the buffer on change.
  useLayoutEffect(() => {
    const mesh = meshRef.current
    if (!mesh) return
    const col = new THREE.Color()
    const dimTarget = theme === 'dark' ? DIM_DARK : DIM_LIGHT
    for (let i = 0; i < count; i++) {
      const n = nodes[i]
      col.set(nodeDisplayColor(n, theme))
      if (n.onPath) col.lerp(GOLD, 0.35)
      const isActive = n.id === activeId
      const isHovered = n.id === hoveredId
      const inFocus = !focused || focusIds === null || focusIds.has(n.id)
      if (focused && !inFocus) {
        col.lerp(dimTarget, 0.78)
      } else if (isActive) {
        col.lerp(HIGHLIGHT, 0.5)
      } else if (isHovered) {
        col.lerp(HIGHLIGHT, 0.3)
      }
      mesh.setColorAt(i, col)
    }
    if (mesh.instanceColor) mesh.instanceColor.needsUpdate = true
  }, [nodes, theme, activeId, hoveredId, focused, focusIds, count])

  return (
    <instancedMesh
      ref={meshRef}
      args={[geometry, material, count]}
      frustumCulled={false}
      onPointerMove={
        hoverEnabled
          ? (e: ThreeEvent<PointerEvent>) => {
              e.stopPropagation()
              const id = e.instanceId != null ? nodes[e.instanceId]?.id ?? null : null
              if (id !== hoveredRef.current) {
                hoveredRef.current = id
                onHover(id)
              }
            }
          : undefined
      }
      onPointerOut={
        hoverEnabled
          ? () => {
              if (hoveredRef.current !== null) {
                hoveredRef.current = null
                onHover(null)
              }
            }
          : undefined
      }
      onClick={(e: ThreeEvent<MouseEvent>) => {
        e.stopPropagation()
        // Ignore drags (pan/rotate) so gestures never select by accident.
        if (e.delta > 6) return
        if (e.instanceId != null && nodes[e.instanceId]) onSelect(nodes[e.instanceId].id)
      }}
    />
  )
}

// ── Edges (one draw call, per-vertex alpha) ─────────────

function EdgeBatch({
  links,
  positionMap,
  activeId,
  focused,
  focusIds,
  theme,
  lod,
}: {
  links: Link[]
  positionMap: Map<string, Vec3>
  activeId: string
  focused: boolean
  focusIds: Set<string> | null
  theme: 'dark' | 'light'
  lod: LodTier
}) {
  const validLinks = useMemo(
    () => links.filter(l => positionMap.has(l.source) && positionMap.has(l.target)),
    [links, positionMap]
  )

  const geometry = useMemo(() => {
    const positions = new Float32Array(validLinks.length * 6)
    const colors = new Float32Array(validLinks.length * 6)
    const alphas = new Float32Array(validLinks.length * 2)
    validLinks.forEach((link, i) => {
      const sp = positionMap.get(link.source)!
      const tp = positionMap.get(link.target)!
      const idx = i * 6
      positions[idx] = sp[0]
      positions[idx + 1] = sp[1]
      positions[idx + 2] = sp[2]
      positions[idx + 3] = tp[0]
      positions[idx + 4] = tp[1]
      positions[idx + 5] = tp[2]
    })
    const geo = new THREE.BufferGeometry()
    geo.setAttribute('position', new THREE.BufferAttribute(positions, 3))
    geo.setAttribute('aColor', new THREE.BufferAttribute(colors, 3))
    geo.setAttribute('aAlpha', new THREE.BufferAttribute(alphas, 1))
    return geo
  }, [validLinks, positionMap])

  const material = useMemo(
    () =>
      new THREE.ShaderMaterial({
        transparent: true,
        depthWrite: false,
        vertexShader: `
          attribute float aAlpha;
          attribute vec3 aColor;
          varying float vAlpha;
          varying vec3 vColor;
          void main() {
            vAlpha = aAlpha;
            vColor = aColor;
            gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
          }
        `,
        fragmentShader: `
          varying float vAlpha;
          varying vec3 vColor;
          void main() {
            if (vAlpha <= 0.002) discard;
            gl_FragColor = vec4(vColor, vAlpha);
          }
        `,
      }),
    []
  )

  useEffect(() => {
    return () => {
      geometry.dispose()
      material.dispose()
    }
  }, [geometry, material])

  useLayoutEffect(() => {
    const colorAttr = geometry.getAttribute('aColor') as THREE.BufferAttribute
    const alphaAttr = geometry.getAttribute('aAlpha') as THREE.BufferAttribute
    const baseColor = new THREE.Color(theme === 'dark' ? LINK_COLOR : LINK_COLOR_LIGHT)
    const activeColor = new THREE.Color(theme === 'dark' ? ACTIVE_LINK_COLOR : ACTIVE_LINK_COLOR_LIGHT)
    const baseAlpha = theme === 'dark' ? EDGE_ALPHA[lod] : EDGE_ALPHA_LIGHT[lod]
    const dim = theme === 'dark' ? DIM_DARK : DIM_LIGHT

    validLinks.forEach((link, i) => {
      const touchesActive = link.source === activeId || link.target === activeId
      const inFocus = !focused || focusIds === null || (focusIds.has(link.source) && focusIds.has(link.target))
      let color = baseColor
      let alpha = baseAlpha
      if (focused && touchesActive) {
        color = activeColor
        alpha = FOCUS_EDGE_ALPHA
      } else if (focused && !inFocus) {
        color = baseColor.clone().lerp(dim, 0.7)
        alpha = baseAlpha * UNFOCUSED_EDGE_FACTOR
      }
      colorAttr.setXYZ(i * 2, color.r, color.g, color.b)
      colorAttr.setXYZ(i * 2 + 1, color.r, color.g, color.b)
      alphaAttr.setX(i * 2, alpha)
      alphaAttr.setX(i * 2 + 1, alpha)
    })
    colorAttr.needsUpdate = true
    alphaAttr.needsUpdate = true
  }, [geometry, validLinks, activeId, focused, focusIds, theme, lod])

  if (validLinks.length === 0) return null

  return <lineSegments geometry={geometry} material={material} frustumCulled={false} />
}

// ── Controls ────────────────────────────────────────────

function MapControls({
  controlsRef,
  autoRotate,
  onInteract,
  invalidate,
}: {
  controlsRef: React.MutableRefObject<OrbitControls | null>
  autoRotate: boolean
  onInteract: () => void
  invalidate: () => void
}) {
  const { camera, gl } = useThree()

  useEffect(() => {
    const controls = new OrbitControls(camera, gl.domElement)
    controls.enableDamping = true
    controls.dampingFactor = 0.08
    controls.enablePan = true
    controls.screenSpacePanning = true
    controls.minDistance = 8
    controls.maxDistance = 60
    controls.zoomSpeed = 0.9
    controls.rotateSpeed = 0.6
    controls.panSpeed = 0.8
    // Touch: one finger pans the map, two fingers pinch + pan. Tap selects.
    controls.touches = { ONE: THREE.TOUCH.PAN, TWO: THREE.TOUCH.DOLLY_PAN }
    controls.mouseButtons = { LEFT: THREE.MOUSE.PAN, MIDDLE: THREE.MOUSE.DOLLY, RIGHT: THREE.MOUSE.ROTATE }
    const onChange = () => invalidate()
    controls.addEventListener('change', onChange)
    controlsRef.current = controls

    const el = gl.domElement
    const mark = () => onInteract()
    el.addEventListener('pointerdown', mark, { once: true })
    el.addEventListener('wheel', mark, { once: true, passive: true })

    return () => {
      el.removeEventListener('pointerdown', mark)
      el.removeEventListener('wheel', mark)
      controls.removeEventListener('change', onChange)
      controls.dispose()
      controlsRef.current = null
    }
  }, [camera, gl, invalidate, controlsRef, onInteract])

  useFrame(() => {
    const controls = controlsRef.current
    if (!controls) return
    controls.update()
    // Keep the cloud reachable: clamp the pan target so a long drag cannot
    // lose the graph off-screen (camera moves with the target).
    const target = controls.target
    const maxPan = 30
    const len = Math.hypot(target.x, target.y, target.z)
    if (len > maxPan) {
      const scale = maxPan / len
      const nx = target.x * scale
      const ny = target.y * scale
      const nz = target.z * scale
      camera.position.x += nx - target.x
      camera.position.y += ny - target.y
      camera.position.z += nz - target.z
      target.set(nx, ny, nz)
      controls.update()
    }
  })

  useEffect(() => {
    const controls = controlsRef.current
    if (!controls) return
    controls.autoRotate = autoRotate
    controls.autoRotateSpeed = 0.55
    invalidate()
  }, [autoRotate, controlsRef, invalidate])

  return null
}

// ── Labels (single DOM overlay, imperative updates) ─────

interface LabelCandidate {
  key: string
  text: string
  position: Vec3
  kind: 'node' | 'domain'
  active: boolean
  dimmed: boolean
  priority: number
  color?: string
}

function LabelProjector({
  overlayRef,
  candidates,
}: {
  overlayRef: React.RefObject<HTMLDivElement>
  candidates: LabelCandidate[]
}) {
  const { camera, size } = useThree()
  const poolRef = useRef<HTMLSpanElement[]>([])
  const gridRef = useRef<Uint8Array>(new Uint8Array(0))
  const gridSizeRef = useRef({ cols: 0, rows: 0 })
  const projected = useRef(new THREE.Vector3())
  const viewScratch = useRef(new THREE.Vector3())

  useEffect(() => {
    const overlay = overlayRef.current
    if (!overlay) return
    const spans: HTMLSpanElement[] = []
    for (let i = 0; i < LABEL_POOL; i++) {
      const span = document.createElement('span')
      span.className = 'graph-label'
      span.style.display = 'none'
      overlay.appendChild(span)
      spans.push(span)
    }
    poolRef.current = spans
    return () => {
      for (const span of spans) span.remove()
      poolRef.current = []
    }
  }, [overlayRef])

  useFrame(() => {
    const overlay = overlayRef.current
    if (!overlay) return
    const pool = poolRef.current
    const width = size.width
    const height = size.height
    const cellW = 112
    const cellH = 18
    const cols = Math.max(1, Math.ceil(width / cellW))
    const rows = Math.max(1, Math.ceil(height / cellH))
    if (gridSizeRef.current.cols !== cols || gridSizeRef.current.rows !== rows) {
      gridSizeRef.current = { cols, rows }
      gridRef.current = new Uint8Array(cols * rows)
    } else {
      gridRef.current.fill(0)
    }
    const grid = gridRef.current
    const v = projected.current
    const cam = camera as THREE.PerspectiveCamera

    let used = 0
    for (const candidate of candidates) {
      if (used >= pool.length) break
      v.set(candidate.position[0], candidate.position[1], candidate.position[2])
      // Skip points behind the camera before projecting.
      const viewZ = viewScratch.current.copy(v).applyMatrix4(cam.matrixWorldInverse).z
      if (viewZ > -0.1) continue
      v.project(cam)
      if (v.z < -1 || v.z > 1) continue
      const x = (v.x * 0.5 + 0.5) * width
      const y = (-v.y * 0.5 + 0.5) * height
      if (x < -40 || x > width + 40 || y < -16 || y > height + 16) continue

      const textW = candidate.text.length * (candidate.kind === 'domain' ? 7.2 : 6) + 10
      const halfW = textW / 2
      const halfH = cellH / 2
      if (textW > width - 4) continue
      // Keep labels inside the overlay; anchor nodes above their dot and
      // domain labels on their centroid.
      const anchorX = Math.min(Math.max(x, halfW + 2), width - halfW - 2)
      const anchorY = candidate.kind === 'domain' ? y : y - 4
      const c0 = Math.max(0, Math.floor((anchorX - halfW) / cellW))
      const c1 = Math.min(cols - 1, Math.floor((anchorX + halfW) / cellW))
      const top = candidate.kind === 'domain' ? anchorY - halfH : anchorY - cellH
      const bottom = candidate.kind === 'domain' ? anchorY + halfH : anchorY
      const r0 = Math.max(0, Math.floor(top / cellH))
      const r1 = Math.min(rows - 1, Math.floor(bottom / cellH))
      let blocked = false
      for (let r = r0; r <= r1 && !blocked; r++) {
        for (let c = c0; c <= c1; c++) {
          if (grid[r * cols + c]) {
            blocked = true
            break
          }
        }
      }
      if (blocked) continue
      for (let r = r0; r <= r1; r++) {
        for (let c = c0; c <= c1; c++) grid[r * cols + c] = 1
      }

      const span = pool[used++]
      // Avoid redundant DOM mutations: text/class only change when the
      // candidate assigned to this pooled span changes.
      if (span.textContent !== candidate.text) span.textContent = candidate.text
      const className = `graph-label graph-label--${candidate.kind}${candidate.active ? ' is-active' : ''}${candidate.dimmed ? ' is-dim' : ''}`
      if (span.className !== className) span.className = className
      if (candidate.kind === 'domain') {
        span.style.setProperty('--label-color', candidate.color ?? 'var(--text-muted)')
      }
      span.style.display = 'block'
      span.style.transform =
        candidate.kind === 'domain'
          ? `translate3d(${anchorX.toFixed(1)}px, ${anchorY.toFixed(1)}px, 0) translate(-50%, -50%)`
          : `translate3d(${anchorX.toFixed(1)}px, ${anchorY.toFixed(1)}px, 0) translate(-50%, -100%)`
    }
    for (let i = used; i < pool.length; i++) pool[i].style.display = 'none'
  })

  return null
}

// ── Scene ───────────────────────────────────────────────

function GraphScene({
  graph,
  activeId,
  focused,
  focusIds,
  onSelect,
  theme,
  isMobile,
  isTouch,
  autoRotateBase,
  controlsRef,
  overlayRef,
  onInteract,
}: {
  graph: BuiltGraph
  activeId: string
  focused: boolean
  focusIds: Set<string> | null
  onSelect: (id: string) => void
  theme: 'dark' | 'light'
  isMobile: boolean
  isTouch: boolean
  autoRotateBase: boolean
  controlsRef: React.MutableRefObject<OrbitControls | null>
  overlayRef: React.RefObject<HTMLDivElement>
  onInteract: () => void
}) {
  const [hovered, setHovered] = useState<string | null>(null)
  const [lod, setLod] = useState<LodTier>('mid')
  const frameCount = useRef(0)
  const { gl, invalidate } = useThree()

  // LOD by apparent graph size (world units per pixel), not distance to
  // origin, so panning never changes detail level.
  useFrame(({ camera }) => {
    frameCount.current += 1
    if (frameCount.current % 8 !== 0) return
    const target = controlsRef.current?.target
    const distance = target ? camera.position.distanceTo(target) : camera.position.length()
    const fov = (camera as THREE.PerspectiveCamera).fov ?? 60
    const visibleHeight = 2 * Math.tan(((fov * Math.PI) / 180) / 2) * distance
    const tier: LodTier = visibleHeight > 46 ? 'far' : visibleHeight > 27 ? 'mid' : 'close'
    setLod(prev => (prev === tier ? prev : tier))
  })

  useEffect(() => {
    gl.domElement.style.cursor = hovered ? 'pointer' : 'default'
  }, [hovered, gl])

  const activeNode = graph.nodeById.get(activeId) ?? null

  const candidates = useMemo(() => {
    const out: LabelCandidate[] = []
    const seen = new Set<string>()
    const focusDomains = new Set<string>()
    if (focused && focusIds) {
      for (const id of focusIds) {
        const n = graph.nodeById.get(id)
        if (n) focusDomains.add(n.domain)
      }
    }
    const add = (
      key: string,
      text: string,
      position: Vec3,
      kind: LabelCandidate['kind'],
      priority: number,
      active = false,
      dimmed = false,
      color?: string
    ) => {
      if (seen.has(key)) return
      seen.add(key)
      out.push({ key, text, position, kind, priority, active, dimmed, color })
    }

    if (activeNode) add(`n:${activeNode.id}`, activeNode.name, activeNode.position, 'node', 0, true)
    if (hovered) {
      const n = graph.nodeById.get(hovered)
      if (n) add(`n:${n.id}`, n.name, n.position, 'node', 1)
    }
    if (focused && focusIds) {
      for (const id of focusIds) {
        if (id === activeId) continue
        const n = graph.nodeById.get(id)
        if (n) add(`n:${n.id}`, n.name, n.position, 'node', 2)
      }
    }
    if (lod !== 'close') {
      for (const [domain, centroid] of Object.entries(graph.domainCentroids)) {
        const dimmed = focused && focusIds !== null && !focusDomains.has(domain)
        add(`d:${domain}`, domainLabel(domain), centroid, 'domain', lod === 'far' ? 3 : 4, false, dimmed, domainColor(domain, theme))
      }
    }
    if (lod !== 'far') {
      const cap = lod === 'mid' ? 26 : 80
      for (const n of graph.byImportance.slice(0, cap)) {
        add(`n:${n.id}`, n.name, n.position, 'node', lod === 'mid' ? 5 : 6)
      }
    }
    return out.sort((a, b) => a.priority - b.priority)
  }, [graph, lod, activeNode, hovered, focused, focusIds, activeId, theme])

  const autoRotate = autoRotateBase && hovered === null && !focused

  return (
    <>
      <ambientLight intensity={0.65} />
      <directionalLight position={[8, 12, 10]} intensity={0.75} />

      <EdgeBatch
        links={graph.links}
        positionMap={graph.positionMap}
        activeId={activeId}
        focused={focused}
        focusIds={focusIds}
        theme={theme}
        lod={lod}
      />

      <NodeInstances
        nodes={graph.nodes}
        activeId={activeId}
        hoveredId={hovered}
        focused={focused}
        focusIds={focusIds}
        theme={theme}
        isMobile={isMobile}
        hoverEnabled={!isTouch}
        onHover={setHovered}
        onSelect={onSelect}
      />

      {activeNode && (
        <mesh position={activeNode.position}>
          <sphereGeometry args={[NODE_RADIUS * 1.7, 12, 12]} />
          <meshBasicMaterial color="#ffdd88" transparent opacity={focused ? 0.2 : 0.1} depthWrite={false} />
        </mesh>
      )}

      <MapControls
        controlsRef={controlsRef}
        autoRotate={autoRotate}
        onInteract={onInteract}
        invalidate={invalidate}
      />

      <LabelProjector overlayRef={overlayRef} candidates={candidates} />
    </>
  )
}

// ── Info panel ──────────────────────────────────────────

function InfoPanel({
  activeId,
  graph,
  theme,
  isMobile,
  expanded,
}: {
  activeId: string
  graph: BuiltGraph
  theme: 'dark' | 'light'
  isMobile: boolean
  expanded: boolean
}) {
  const concept = graph.nodeById.get(activeId)
  if (!concept) return null

  const color = domainColor(concept.domain, theme)
  const prereqConcepts = (graph.prereqsById.get(concept.id) ?? [])
    .map(id => graph.nodeById.get(id))
    .filter((c): c is RenderNode => Boolean(c))
  const unlockedBy = (graph.unlocksById.get(concept.id) ?? [])
    .map(id => graph.nodeById.get(id))
    .filter((c): c is RenderNode => Boolean(c))

  const body = (
    <InfoPanelBody
      concept={concept}
      color={color}
      prereqConcepts={prereqConcepts}
      unlockedBy={unlockedBy}
      isMobile={isMobile}
    />
  )

  const shellStyle = {
    background: 'transparent',
    borderTop: '0.5px solid var(--border)',
    padding: isMobile ? '0.5rem' : '0.75rem 1rem',
    borderRadius: 0,
  } as const

  const heading = (
    <>
      <div style={{ color: 'var(--text-muted)', fontSize: '12px', textTransform: 'uppercase', fontFamily: monoFont }}>selected concept</div>
      <div style={{ color: 'var(--text-primary)', fontSize: isMobile ? '16px' : '18px', marginTop: '4px', fontFamily: serifFont, fontWeight: 400 }}>{concept.name}</div>
    </>
  )

  // Mobile: collapsible via native <details>; collapsed until the user
  // selects a node (key resets it per selection). Desktop keeps the
  // always-open panel.
  if (isMobile) {
    return (
      <details open={expanded} key={activeId} style={shellStyle}>
        <summary style={{ cursor: 'pointer', listStyle: 'none' }}>
          {heading}
        </summary>
        {body}
      </details>
    )
  }

  return (
    <div style={shellStyle}>
      {heading}
      {body}
    </div>
  )
}

// Chip list with a native "+N more" expander so hub nodes don't push
// the mobile page several screens long. Everything stays reachable.
function ChipList({ title, items, isMobile }: { title: string; items: RenderNode[]; isMobile: boolean }) {
  if (items.length === 0) return null
  const cap = isMobile && items.length > 6 ? 6 : items.length
  const shown = items.slice(0, cap)
  const hiddenCount = items.length - shown.length
  return (
    <div style={{ marginTop: '8px' }}>
      <div style={{ color: 'var(--text-muted)', fontSize: '12px', textTransform: 'uppercase', marginBottom: '4px', fontFamily: monoFont }}>{title}</div>
      <div style={{ display: 'flex', gap: '6px', flexWrap: 'wrap' as const }}>
        {shown.map(c => (
          <span key={c.id} style={{ color: 'var(--text-secondary)', fontSize: '12px', fontFamily: monoFont }}>
            {c.name}
          </span>
        ))}
      </div>
      {hiddenCount > 0 && (
        <details style={{ marginTop: '4px' }}>
          <summary style={{ color: 'var(--accent-blue)', fontSize: '12px', fontFamily: monoFont, cursor: 'pointer' }}>
            +{hiddenCount} more
          </summary>
          <div style={{ display: 'flex', gap: '6px', flexWrap: 'wrap' as const, marginTop: '4px' }}>
            {items.slice(cap).map(c => (
              <span key={c.id} style={{ color: 'var(--text-secondary)', fontSize: '12px', fontFamily: monoFont }}>
                {c.name}
              </span>
            ))}
          </div>
        </details>
      )}
    </div>
  )
}

function InfoPanelBody({ concept, color, prereqConcepts, unlockedBy, isMobile }: {
  concept: RenderNode
  color: string
  prereqConcepts: RenderNode[]
  unlockedBy: RenderNode[]
  isMobile: boolean
}) {
  return (
    <>
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginTop: '8px', flexWrap: 'wrap' as const }}>
        <span style={{ color: 'var(--text-secondary)', fontSize: '12px', fontFamily: monoFont }}>
          <span style={{ color }}>●</span> {concept.domain.replace(/_/g, ' ')}
        </span>
        {concept.status && (
          <span style={{
            color: concept.status === 'mastered' ? 'var(--accent-teal)' : concept.status === 'learning' ? '#e8a849' : 'var(--text-muted)',
            fontSize: '12px',
            textTransform: 'uppercase',
            fontFamily: monoFont,
          }}>
            {concept.status}
          </span>
        )}
        {concept.onPath && (
          <span style={{ color: 'var(--accent-blue)', fontSize: '12px', fontFamily: monoFont }}>on path</span>
        )}
      </div>
      <ChipList title="prerequisites" items={prereqConcepts} isMobile={isMobile} />
      <ChipList title="unlocks" items={unlockedBy} isMobile={isMobile} />
      <div style={{ marginTop: '10px', display: 'flex', gap: '8px' }}>
        <a
          href={`/concept?id=${encodeURIComponent(concept.id)}`}
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
    </>
  )
}

// ── Root ────────────────────────────────────────────────

export default function MathConceptGraph3D({
  theme = 'dark',
  conceptStatuses,
  onPathNodes,
  onNodeSelect,
}: MathConceptGraph3DProps) {
  const isMobile = useIsMobile()
  const reducedMotion = usePrefersReducedMotion()
  const isTouch = useCoarsePointer()
  const [inView, setInView] = useState(true)
  const [domFocused, setDomFocused] = useState(false)
  const [interacted, setInteracted] = useState(false)
  const [payload, setPayload] = useState<GraphPayload | null>(null)
  const [positions, setPositions] = useState<Array<[string, Vec3]> | null>(null)
  const [failed, setFailed] = useState(false)
  const wrapRef = useRef<HTMLDivElement>(null)
  const overlayRef = useRef<HTMLDivElement>(null)
  const controlsRef = useRef<OrbitControls | null>(null)

  // Fetch graph payload + positions once; both are cached single-flight and
  // shared with the static poster. The corpus never enters the page bundle.
  useEffect(() => {
    let live = true
    Promise.all([loadGraphPayload(), loadPositionEntries()]).then(([graphPayload, entries]) => {
      if (!live) return
      if (!graphPayload) {
        setFailed(true)
        return
      }
      setPayload(graphPayload)
      setPositions(entries)
    })
    return () => {
      live = false
    }
  }, [])

  const graph = useMemo(
    () => (payload && positions ? buildGraph(payload, positions, conceptStatuses, onPathNodes) : null),
    [payload, positions, conceptStatuses, onPathNodes]
  )

  const [activeId, setActiveId] = useState('')
  const [focused, setFocused] = useState(false)
  // Mobile panel stays collapsed until the user actually selects a node.
  const [userSelected, setUserSelected] = useState(false)

  useEffect(() => {
    if (graph && graph.nodes.length > 0 && !graph.nodeById.has(activeId)) {
      setActiveId(pickDefaultId(graph))
      setFocused(false)
      setUserSelected(false)
    }
  }, [graph, activeId])

  const onNodeSelectRef = useRef(onNodeSelect)
  useEffect(() => {
    onNodeSelectRef.current = onNodeSelect
  }, [onNodeSelect])

  const handleSelect = useCallback((id: string) => {
    setActiveId(id)
    setFocused(true)
    setUserSelected(true)
    onNodeSelectRef.current?.(id)
  }, [])

  const pointerDown = useRef<{ x: number; y: number } | null>(null)

  // Clear focus on background clicks, but not when the pointer was panning
  // or rotating (same drag threshold as node selection).
  const handleMiss = useCallback((event?: MouseEvent) => {
    const start = pointerDown.current
    if (event && start && Math.hypot(event.clientX - start.x, event.clientY - start.y) > 6) return
    setFocused(false)
    setUserSelected(false)
  }, [])

  const handleInteract = useCallback(() => {
    setInteracted(true)
  }, [])

  // Pause rendering entirely when scrolled offscreen (frameloop 'never').
  useEffect(() => {
    const el = wrapRef.current
    if (!el || typeof IntersectionObserver === 'undefined') return
    const obs = new IntersectionObserver(entries => setInView(entries[0]?.isIntersecting ?? true), { threshold: 0 })
    obs.observe(el)
    return () => obs.disconnect()
  }, [])

  // Focus mode: selected node + its direct prerequisites/unlocks stay lit,
  // everything else dims (never removed). null = nothing selected yet.
  const focusIds = useMemo(() => {
    if (!graph || !focused || !activeId) return null
    const set = new Set<string>([activeId])
    for (const id of graph.prereqsById.get(activeId) ?? []) set.add(id)
    for (const id of graph.unlocksById.get(activeId) ?? []) set.add(id)
    return set
  }, [graph, focused, activeId])

  // Stable traversal order for keyboard browsing: domain, then label.
  const orderedIds = useMemo(
    () =>
      graph
        ? [...graph.nodes]
            .sort((a, b) => a.domain.localeCompare(b.domain) || a.name.localeCompare(b.name))
            .map(c => c.id)
        : [],
    [graph]
  )

  const handleGraphKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (orderedIds.length === 0) return
      const idx = Math.max(0, orderedIds.indexOf(activeId))
      let next: number | null = null
      if (e.key === 'ArrowRight' || e.key === 'ArrowDown') {
        next = Math.min(idx + 1, orderedIds.length - 1)
      } else if (e.key === 'ArrowLeft' || e.key === 'ArrowUp') {
        next = Math.max(idx - 1, 0)
      } else if (e.key === 'Home') {
        next = 0
      } else if (e.key === 'End') {
        next = orderedIds.length - 1
      } else if (e.key === 'Escape') {
        setFocused(false)
        setUserSelected(false)
      } else if (e.key === '+' || e.key === '=' || e.key === '-' || e.key === '_') {
        e.preventDefault()
        handleInteract()
        const controls = controlsRef.current
        const camera = controls?.object as THREE.PerspectiveCamera | undefined
        if (!camera) return
        const target = controls?.target ?? new THREE.Vector3()
        const distance = camera.position.distanceTo(target)
        const zoomInKey = e.key === '+' || e.key === '='
        const clamped = THREE.MathUtils.clamp(
          zoomInKey ? distance * 0.85 : distance * 1.18,
          controls?.minDistance ?? 8,
          controls?.maxDistance ?? 60
        )
        const direction = camera.position.clone().sub(target).normalize()
        camera.position.copy(target).addScaledVector(direction, clamped)
        controls?.update()
      }
      if (next !== null && orderedIds[next] !== activeId) {
        e.preventDefault()
        handleSelect(orderedIds[next])
      }
    },
    [orderedIds, activeId, handleSelect, handleInteract]
  )

  const autoRotateBase = !isMobile && !reducedMotion && !interacted && !domFocused && inView
  const frameloop = !inView ? 'never' : autoRotateBase ? 'always' : 'demand'

  return (
    <div ref={wrapRef}>
      <div
        className="concept-graph-3d"
        tabIndex={0}
        role="application"
        aria-label="Concept map. Use the left and right arrow keys to browse concepts, plus and minus to zoom, and Home or End to jump to the first or last concept."
        onKeyDown={handleGraphKeyDown}
        onFocus={() => setDomFocused(true)}
        onBlur={() => setDomFocused(false)}
        onPointerDown={e => {
          pointerDown.current = { x: e.clientX, y: e.clientY }
        }}
        style={{
          position: 'relative',
          height: isMobile ? '320px' : '520px',
          width: '100%',
          borderRadius: 0,
          overflow: 'hidden',
          border: '0.5px solid var(--border)',
          background: 'var(--graph-surface)',
        }}
      >
        <span className="sr-only">
          Interactive three-dimensional concept map. Arrow keys move between
          concepts, plus and minus zoom, and details of the selected concept
          appear below the map.
        </span>
        {graph ? (
          <>
            <Canvas
              camera={{ position: [0, 0, 30], fov: 60 }}
              gl={{ alpha: true, premultipliedAlpha: true, powerPreference: 'high-performance', antialias: !isMobile }}
              dpr={isMobile ? 1 : [1, 1.5]}
              frameloop={frameloop}
              onPointerMissed={handleMiss}
              onCreated={({ gl }) => {
                const canvas = gl.domElement as HTMLCanvasElement
                const onLost = (e: Event) => {
                  e.preventDefault()
                  console.warn('WebGL context lost — will attempt restore')
                }
                const onRestored = () => console.warn('WebGL context restored')
                canvas.addEventListener('webglcontextlost', onLost, false)
                canvas.addEventListener('webglcontextrestored', onRestored, false)
              }}
            >
              <GraphScene
                graph={graph}
                activeId={activeId}
                focused={focused}
                focusIds={focusIds}
                onSelect={handleSelect}
                theme={theme}
                isMobile={isMobile}
                isTouch={isTouch}
                autoRotateBase={autoRotateBase}
                controlsRef={controlsRef}
                overlayRef={overlayRef}
                onInteract={handleInteract}
              />
            </Canvas>
            <div ref={overlayRef} className="graph-label-layer" aria-hidden="true" />
          </>
        ) : (
          <div
            style={{
              height: '100%',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: 'var(--text-muted)',
              fontFamily: monoFont,
              fontSize: '13px',
            }}
          >
            {failed ? 'GRAPH UNAVAILABLE' : 'LOADING GRAPH'}
          </div>
        )}
      </div>
      {graph && <InfoPanel activeId={activeId} graph={graph} theme={theme} isMobile={isMobile} expanded={userSelected} />}
    </div>
  )
}
