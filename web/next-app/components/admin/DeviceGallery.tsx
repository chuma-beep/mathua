'use client'

import { useEffect, useRef } from 'react'
import * as THREE from 'three'
import { useTheme } from '@/hooks/useTheme'
import type { AdminLocationDevice } from '../../lib/api'

// A small turntable of device models, built from plain three like the globe, so
// the bundle gains no renderer wrapper. The screen fills and the floor arc both
// encode the same number: this class's share of device visits this week. The
// models are a gallery, not a chart, so the counts and percentages are printed
// underneath in text where a screen reader can read them.
//
// Device classes overlap: a learner who signs in on a phone and a laptop counts
// once in each class. The percentages are therefore shares of device visits,
// not of people, and the caption says so.

const ORDER = ['mobile', 'tablet', 'desktop'] as const
type Kind = (typeof ORDER)[number]

const LABEL: Record<Kind, string> = {
  mobile: 'Mobile',
  tablet: 'Tablet',
  desktop: 'Desktop',
}

interface Palette {
  body: string
  edge: string
  screen: string
  fill: string
  detail: string
  ring: string
  floor: string
  backdrop: string
}

const PALETTE: Record<'dark' | 'light', Palette> = {
  dark: {
    body: '#22344f',
    edge: '#6b93bf',
    screen: '#0d2036',
    fill: '#38bdf8',
    detail: '#7dd3fc',
    ring: '#60a5fa',
    floor: '#1e3a8a',
    backdrop: '#0b1526',
  },
  light: {
    body: '#b9d0e6',
    edge: '#6f90b3',
    screen: '#ffffff',
    fill: '#2563eb',
    detail: '#93c5fd',
    ring: '#3b82f6',
    floor: '#bfdbfe',
    backdrop: '#e9f1fa',
  },
}

function roundedShape(w: number, h: number, r: number): THREE.Shape {
  const x = -w / 2
  const y = -h / 2
  const shape = new THREE.Shape()
  shape.moveTo(x + r, y)
  shape.lineTo(x + w - r, y)
  shape.quadraticCurveTo(x + w, y, x + w, y + r)
  shape.lineTo(x + w, y + h - r)
  shape.quadraticCurveTo(x + w, y + h, x + w - r, y + h)
  shape.lineTo(x + r, y + h)
  shape.quadraticCurveTo(x, y + h, x, y + h - r)
  shape.lineTo(x, y + r)
  shape.quadraticCurveTo(x, y, x + r, y)
  return shape
}

// A rounded slab: width and height in the plane, thickness along z. The body is
// bevelled rather than perfectly sharp so the edge catches the key light.
function slab(w: number, h: number, d: number, radius: number): THREE.ExtrudeGeometry {
  const geometry = new THREE.ExtrudeGeometry(roundedShape(w, h, Math.min(radius, w / 2, h / 2)), {
    depth: d,
    bevelEnabled: true,
    bevelThickness: 0.012,
    bevelSize: 0.012,
    bevelSegments: 2,
    curveSegments: 10,
  })
  geometry.center()
  return geometry
}

interface Built {
  group: THREE.Group
  ringRadius: number
}

function buildDevice(kind: Kind, palette: Palette, share: number): Built {
  const group = new THREE.Group()
  const body = new THREE.MeshStandardMaterial({
    color: palette.body,
    metalness: 0.5,
    roughness: 0.36,
  })
  const screen = new THREE.MeshBasicMaterial({ color: palette.screen })
  const fill = new THREE.MeshBasicMaterial({ color: palette.fill, transparent: true, opacity: 0.92 })
  const detail = new THREE.MeshBasicMaterial({ color: palette.detail, transparent: true, opacity: 0.45 })
  const line = new THREE.LineBasicMaterial({ color: palette.edge, transparent: true, opacity: 0.55 })

  const outline = (geometry: THREE.BufferGeometry, mesh: THREE.Mesh) => {
    const edges = new THREE.LineSegments(new THREE.EdgesGeometry(geometry, 24), line)
    mesh.add(edges)
  }

  // The glass sits proud of the bevelled shell: the extrusion is deeper than
  // its nominal thickness, so a plane at depth/2 would be swallowed by the body
  // and the screen would read as an unlit slab.
  const panel = (w: number, h: number, z: number, parent: THREE.Object3D, y: number) => {
    const glass = new THREE.Mesh(new THREE.PlaneGeometry(w, h), screen)
    glass.position.set(0, y, z)
    parent.add(glass)
    // The share bar grows from the bottom edge of the glass, so a class with a
    // small share still shows a thin lit sliver rather than an empty screen.
    const litHeight = Math.max(0.02, (h - 0.06) * share)
    const lit = new THREE.Mesh(new THREE.PlaneGeometry(w - 0.04, litHeight), fill)
    lit.position.set(0, y - h / 2 + 0.03 + litHeight / 2, z + 0.001)
    parent.add(lit)
    // Two short bars above the lit area suggest an app view without drawing
    // text into a canvas that has no font to draw with.
    const bar = (fraction: number, rise: number) => {
      const mesh = new THREE.Mesh(new THREE.PlaneGeometry(w * fraction, 0.018), detail)
      mesh.position.set(0, y + h / 4 + rise, z + 0.002)
      parent.add(mesh)
    }
    bar(0.62, 0.06)
    bar(0.42, 0)
    return glass
  }

  if (kind === 'desktop') {
    const base = new THREE.Mesh(slab(1.22, 0.8, 0.07, 0.07), body)
    base.rotation.x = -Math.PI / 2
    base.position.set(0, -0.585, 0.1)
    group.add(base)
    outline(slab(1.22, 0.8, 0.07, 0.07), base)

    // The lid grows up from the hinge so its bottom edge stays on the base
    // however far it is tipped back.
    const lid = new THREE.Group()
    lid.position.set(0, -0.55, -0.3)
    lid.rotation.x = -0.2
    group.add(lid)
    const lidShell = new THREE.Mesh(slab(1.22, 0.78, 0.05, 0.05), body)
    lidShell.position.set(0, 0.39, 0)
    lid.add(lidShell)
    panel(1.1, 0.64, 0.045, lid, 0.39)
    return { group, ringRadius: 0.78 }
  }

  if (kind === 'tablet') {
    const shell = new THREE.Mesh(slab(0.78, 1.08, 0.065, 0.07), body)
    shell.position.set(0, -0.08, 0)
    group.add(shell)
    panel(0.68, 0.92, 0.052, group, -0.08)
    return { group, ringRadius: 0.55 }
  }

  const shell = new THREE.Mesh(slab(0.5, 1.16, 0.065, 0.09), body)
  shell.position.set(0, -0.04, 0)
  group.add(shell)
  panel(0.42, 1.0, 0.052, group, -0.04)
  return { group, ringRadius: 0.38 }
}

function ring(radius: number, palette: Palette, arc: number): THREE.Group {
  const spin = new THREE.Group()
  spin.rotation.y = Math.PI / 2
  const track = new THREE.Mesh(
    new THREE.RingGeometry(radius - 0.02, radius, 64),
    new THREE.MeshBasicMaterial({ color: palette.floor, transparent: true, opacity: 0.28, side: THREE.DoubleSide }),
  )
  track.rotation.x = -Math.PI / 2
  track.position.y = -0.612
  spin.add(track)
  if (arc > 0.001) {
    const lit = new THREE.Mesh(
      new THREE.TorusGeometry(radius - 0.01, 0.012, 8, 64, Math.PI * 2 * arc),
      new THREE.MeshBasicMaterial({ color: palette.ring, transparent: true, opacity: 0.95 }),
    )
    lit.rotation.x = -Math.PI / 2
    lit.position.y = -0.61
    spin.add(lit)
  }
  return spin
}

export default function DeviceGallery({ devices }: { devices: AdminLocationDevice[] }) {
  const host = useRef<HTMLDivElement>(null)
  const { theme, mounted } = useTheme()

  const counts = ORDER.map((kind) => devices.find((d) => d.device === kind)?.active_7d ?? 0)
  const total = counts.reduce((sum, n) => sum + n, 0)
  const signature = counts.join(':')

  useEffect(() => {
    const node = host.current
    if (!node || !mounted) return

    const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    const palette = PALETTE[theme === 'dark' ? 'dark' : 'light']

    const scene = new THREE.Scene()
    const camera = new THREE.PerspectiveCamera(32, 1, 0.1, 30)
    camera.position.set(0, 0.04, 4.4)
    camera.lookAt(0, -0.04, 0)

    const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true })
    renderer.setPixelRatio(Math.min(2, window.devicePixelRatio || 1))
    renderer.outputColorSpace = THREE.SRGBColorSpace
    renderer.domElement.setAttribute('aria-hidden', 'true')
    renderer.domElement.className = 'h-full w-full touch-none cursor-grab active:cursor-grabbing'
    node.appendChild(renderer.domElement)

    scene.add(new THREE.AmbientLight(0xffffff, theme === 'dark' ? 1.5 : 0.85))
    const key = new THREE.DirectionalLight(0xffffff, theme === 'dark' ? 2.4 : 1.0)
    key.position.set(2.4, 3.4, 3.6)
    scene.add(key)
    const rim = new THREE.DirectionalLight(0x93c5fd, theme === 'dark' ? 1.6 : 0.35)
    rim.position.set(-3, 0.6, -2.4)
    scene.add(rim)

    const floor = new THREE.Mesh(
      new THREE.CircleGeometry(2.6, 64),
      new THREE.MeshBasicMaterial({ color: palette.backdrop, transparent: true, opacity: theme === 'dark' ? 0.5 : 0.35 }),
    )
    floor.rotation.x = -Math.PI / 2
    floor.position.y = -0.618
    scene.add(floor)

    const turntable = new THREE.Group()
    scene.add(turntable)
    const bases: number[] = []
    const spreads = [-1.55, 0, 1.55]
    counts.forEach((count, index) => {
      const share = total > 0 ? count / total : 0
      const built = buildDevice(ORDER[index], palette, share)
      built.group.position.set(spreads[index], 0, 0)
      bases.push(built.group.position.y)
      turntable.add(built.group)
      turntable.add(
        (() => {
          const r = ring(built.ringRadius, palette, share)
          r.position.x = spreads[index]
          return r
        })(),
      )
    })

    let width = 0
    let height = 0
    const fit = () => {
      width = node.clientWidth || 1
      height = node.clientHeight || 1
      renderer.setSize(width, height, false)
      camera.aspect = width / height
      camera.updateProjectionMatrix()
      // On a narrow card the three models would run past the edges, so the
      // turntable shrinks rather than cropping the outer devices.
      const scale = Math.max(0.62, Math.min(1.32, camera.aspect / 2.4))
      turntable.scale.setScalar(scale)
    }
    fit()
    const observer = new ResizeObserver(fit)
    observer.observe(node)

    let yaw = 0
    let targetYaw = 0
    let dragX = 0
    let dragging = false
    const onDown = (event: PointerEvent) => {
      dragging = true
      dragX = event.clientX
      renderer.domElement.setPointerCapture(event.pointerId)
    }
    const onMove = (event: PointerEvent) => {
      if (!dragging) return
      targetYaw = Math.max(-0.5, Math.min(0.5, targetYaw + (event.clientX - dragX) * 0.007))
      dragX = event.clientX
    }
    const onUp = () => {
      dragging = false
    }
    renderer.domElement.addEventListener('pointerdown', onDown)
    renderer.domElement.addEventListener('pointermove', onMove)
    renderer.domElement.addEventListener('pointerup', onUp)
    renderer.domElement.addEventListener('pointercancel', onUp)

    const clock = new THREE.Clock()
    let frame = 0
    const animate = () => {
      frame = requestAnimationFrame(animate)
      const time = clock.getElapsedTime()
      yaw += (targetYaw - yaw) * 0.08
      if (!reduced) {
        // A slow sway keeps the metal edges alive; the drag offset is added on
        // top so the reader's angle is never fought by the idle motion.
        turntable.rotation.y = yaw + Math.sin(time * 0.32) * 0.05
        turntable.children.forEach((child, index) => {
          if (index % 2 !== 0) return
          child.position.y = bases[Math.floor(index / 2)] + Math.sin(time * 0.9 + index) * 0.014
        })
      } else {
        turntable.rotation.y = yaw
      }
      renderer.render(scene, camera)
    }
    animate()

    return () => {
      cancelAnimationFrame(frame)
      observer.disconnect()
      renderer.domElement.removeEventListener('pointerdown', onDown)
      renderer.domElement.removeEventListener('pointermove', onMove)
      renderer.domElement.removeEventListener('pointerup', onUp)
      renderer.domElement.removeEventListener('pointercancel', onUp)
      scene.traverse((object) => {
        const mesh = object as THREE.Mesh
        mesh.geometry?.dispose()
        const material = mesh.material
        if (Array.isArray(material)) material.forEach((m) => m.dispose())
        else material?.dispose()
      })
      renderer.dispose()
      renderer.forceContextLoss()
      if (renderer.domElement.parentNode === node) node.removeChild(renderer.domElement)
    }
  }, [mounted, theme, signature, counts, total])

  return (
    <div>
      <div className="relative mx-auto w-full max-w-[620px]">
        <div aria-hidden="true" className="pointer-events-none absolute inset-x-[14%] bottom-[8%] top-[18%] bg-mathua-blue-faint blur-2xl" />
        <div ref={host} data-testid="devices-3d" className="relative h-[218px] w-full sm:h-[236px]" />
      </div>
      <div className="mt-2 grid grid-cols-3 gap-2">
        {ORDER.map((kind, index) => {
          const count = counts[index]
          const share = total > 0 ? Math.round((count / total) * 100) : 0
          return (
            <div key={kind} className="text-center">
              <p className="font-mono text-[10px] uppercase tracking-[0.12em] text-mathua-muted">{LABEL[kind]}</p>
              <p className="mt-0.5 font-mono text-sm text-mathua-primary">{count.toLocaleString()}</p>
              <p className="font-mono text-[10px] text-mathua-muted">{share}% of device visits</p>
            </div>
          )
        })}
      </div>
      <p className="mt-3 font-mono text-[10px] text-mathua-muted">
        Screen fill and floor arc show the share of each class. A learner who signs in on two
        devices counts in both, so these shares are device visits, not people.
      </p>
      <p className="sr-only">
        Device gallery: {ORDER.map((kind, index) => `${LABEL[kind]} ${counts[index]}`).join(', ')}. Drag the
        models to turn the gallery.
      </p>
    </div>
  )
}
