'use client'

import { Component, useEffect, useMemo, useState, type ReactNode } from 'react'
import dynamic from 'next/dynamic'
import FigureFrame from '../FigureFrame'
import { Slider, Toggle } from '../Controls'
import { solid3dSchema } from '../schemas'
import { lateralArea, radius, volume } from '../math/solid'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

// The renderers are separate chunks on purpose: the Three.js scene is roughly
// 120 KB gzipped and a lesson without a 3D step must not pay for it. The
// dispatcher below decides which one to load, and never loads both.
const SolidWebGL = dynamic(() => import('./solid/SolidWebGL'), { ssr: false })
const SolidWireframe = dynamic(() => import('./solid/SolidWireframe'), { ssr: false })

/** Exported so the choice can be tested without a browser. */
export function pickRenderer(webgl: boolean, forceFallback = false): 'webgl' | 'wireframe' {
  return webgl && !forceFallback ? 'webgl' : 'wireframe'
}

class RendererBoundary extends Component<{ fallback: ReactNode; children: ReactNode }, { failed: boolean }> {
  state = { failed: false }

  static getDerivedStateFromError() {
    return { failed: true }
  }

  render() {
    return this.state.failed ? this.props.fallback : this.props.children
  }
}

export default function Solid3D({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => solid3dSchema.parse(params), [params])
  const [theta, setTheta] = useState(0.6)
  const [phi, setPhi] = useState(0.35)
  const [cut, setCut] = useState(defaults.cut)
  const [showCut, setShowCut] = useState(defaults.showCut)
  const [webgl, setWebgl] = useState<boolean | null>(null)

  // Capability is checked after mount: reading it during render would make the
  // server and client disagree about which renderer to draw.
  useEffect(() => {
    let supported = false
    try {
      const canvas = document.createElement('canvas')
      supported = Boolean(canvas.getContext('webgl2') ?? canvas.getContext('webgl'))
    } catch {
      supported = false
    }
    setWebgl(supported)
  }, [])

  const reset = () => {
    setTheta(0.6)
    setPhi(0.35)
    setCut(defaults.cut)
    setShowCut(defaults.showCut)
  }

  const orbit = (dTheta: number, dPhi: number) => {
    setTheta((value) => value + dTheta)
    setPhi((value) => Math.max(-1.2, Math.min(1.2, value + dPhi)))
  }

  const view = { theta, phi, cut, showCut, onOrbit: orbit }
  const wireframe = <SolidWireframe {...view} />
  const choice = webgl === null ? null : pickRenderer(webgl)

  return (
    <FigureFrame
      title={title ?? 'A solid of revolution, cut anywhere'}
      hint="Drag to turn the solid, or use the angle sliders. Move the cut to read the disk it makes."
      caption={
        caption ??
        'The body is the surface of revolution of r(x) = 1 + 0.35 sin(1.6x) about the x axis. The teal disk is the cross-section at the cut, and the volume is the integral of its area.'
      }
      onReset={reset}
      readout={[
        { label: 'r(cut)', value: fmt(radius(cut), 3) },
        { label: 'cross-section', value: fmt(Math.PI * radius(cut) ** 2, 3) },
        { label: 'volume', value: fmt(volume(), 3) },
        { label: 'lateral area', value: fmt(lateralArea(), 3) },
        { label: 'cut', value: fmt(cut, 2) },
        { label: 'renderer', value: choice ?? 'checking' },
      ]}
    >
      <div className="rounded-none border border-mathua-border bg-mathua-bg" style={{ minHeight: 320 }}>
        {choice === null ? (
          <div className="h-[320px] w-full bg-mathua-surface" aria-hidden="true" />
        ) : choice === 'webgl' ? (
          <RendererBoundary fallback={wireframe}>
            <SolidWebGL {...view} />
          </RendererBoundary>
        ) : (
          wireframe
        )}
      </div>

      <div className="mt-3 space-y-2">
        <Slider label="turn" value={Math.round((theta * 180) / Math.PI)} min={-180} max={180} step={1} onChange={(v) => setTheta((v * Math.PI) / 180)} suffix="°" />
        <Slider label="tilt" value={Math.round((phi * 180) / Math.PI)} min={-60} max={60} step={1} onChange={(v) => setPhi((v * Math.PI) / 180)} suffix="°" />
        <Slider label="cut" value={cut} min={0} max={4} step={0.05} onChange={setCut} />
        <Toggle label="Show the cross-section" checked={showCut} onChange={setShowCut} />
      </div>
    </FigureFrame>
  )
}
