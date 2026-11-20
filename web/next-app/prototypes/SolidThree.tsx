'use client'

import { useEffect, useRef, useState } from 'react'
import {
  AmbientLight,
  MeshNormalMaterial,
  BufferAttribute,
  BufferGeometry,
  CircleGeometry,
  DirectionalLight,
  LineBasicMaterial,
  LineLoop,
  Mesh,
  MeshLambertMaterial,
  PerspectiveCamera,
  Scene,
  WebGLRenderer,
} from 'three'
import { radius, surface, X_MAX, X_MIN } from './solid'

/**
 * Option A: Three.js.
 *
 * Rendered on demand rather than in a loop. The figure is a still solid that the
 * learner turns, so a requestAnimationFrame loop would spend the phone's battery
 * drawing an unchanged picture sixty times a second. Every state change schedules
 * exactly one frame.
 */
export default function SolidThree({ onReady }: { onReady?: (ms: number) => void }) {
  const host = useRef<HTMLDivElement | null>(null)
  const [theta, setTheta] = useState(0.7)
  const [phi, setPhi] = useState(1.15)
  const [cut, setCut] = useState(true)
  const [cutX, setCutX] = useState(2)

  useEffect(() => {
    const element = host.current
    if (!element) return

    const started = performance.now()
    const scene = new Scene()
    const camera = new PerspectiveCamera(38, 1, 0.1, 100)
    const renderer = new WebGLRenderer({ antialias: true, alpha: true })
    renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2))
    const width = element.clientWidth || 320
    renderer.setSize(width, width, false)
    element.appendChild(renderer.domElement)

    const mesh = surface()
    const geometry = new BufferGeometry()
    geometry.setAttribute('position', new BufferAttribute(mesh.positions, 3))
    geometry.setAttribute('normal', new BufferAttribute(mesh.normals, 3))
    geometry.setIndex(new BufferAttribute(mesh.indices, 1))
    // DoubleSide because the learner may turn the solid to look into an open end;
    // with culling on, that view is a hole rather than the inside of the surface.
    const material = new MeshLambertMaterial({ color: '#4a86c8', side: 2 })
    const body = new Mesh(geometry, material)
    scene.add(body)
    const center = (X_MIN + X_MAX) / 2

    // Caps. The lathe surface is open at both ends, so without these the camera
    // looks straight into the tube and the picture reads as a flat silhouette.
    const capMaterial = new MeshLambertMaterial({ color: '#4a86c8', side: 2 })
    const capStart = new Mesh(new CircleGeometry(radius(X_MIN), 64), capMaterial)
    capStart.rotation.y = -Math.PI / 2
    capStart.position.x = -center
    const capEnd = new Mesh(new CircleGeometry(radius(X_MAX), 64), capMaterial)
    capEnd.rotation.y = Math.PI / 2
    capEnd.position.x = X_MAX - center
    scene.add(capStart, capEnd)

    // Ambient plus key, kept below 1.6 together: at 3.6 the blue channel clipped
    // and the whole solid rendered as one flat tone.
    scene.add(new AmbientLight(0xffffff, 0.55))
    const key = new DirectionalLight(0xffffff, 1.05)
    // Off to one side, not along the camera axis: a key light behind the camera
    // makes N·L almost constant over the visible surface and the solid reads flat.
    key.position.set(-5, 3, -6)
    scene.add(key)

    // The cross-section: a circle of radius r(cutX) in the plane x = cutX.
    const r = radius(cutX)
    const circle: number[] = []
    for (let i = 0; i < 64; i++) {
      const a = (2 * Math.PI * i) / 64
      circle.push(0, r * Math.cos(a), r * Math.sin(a))
    }
    const sectionGeometry = new BufferGeometry()
    sectionGeometry.setAttribute('position', new BufferAttribute(new Float32Array(circle), 3))
    const sectionMaterial = new LineBasicMaterial({ color: '#2f8f6f' })
    const section = new LineLoop(sectionGeometry, sectionMaterial)
    const diskGeometry = new CircleGeometry(r, 64)
    const diskMaterial = new MeshLambertMaterial({
      color: '#2f8f6f',
      transparent: true,
      opacity: 0.3,
      side: 2,
    })
    const disk = new Mesh(diskGeometry, diskMaterial)
    disk.rotation.y = Math.PI / 2
    scene.add(section, disk)

    // The geometry is already centred on x = 0, so the camera looks at the
    // origin; looking at the model-space centre as well would shift the solid
    // twice and push it out of frame.
    const place = () => {
      const distance = 9
      camera.position.set(
        distance * Math.sin(phi) * Math.cos(theta),
        distance * Math.cos(phi),
        distance * Math.sin(phi) * Math.sin(theta),
      )
      camera.lookAt(0, 0, 0)
    }

    const draw = () => {
      section.position.x = cutX - center
      disk.position.x = cutX - center
      section.visible = cut
      disk.visible = cut
      place()
      renderer.render(scene, camera)
    }

    draw()
    onReady?.(performance.now() - started)

    const observer = new ResizeObserver(() => {
      const next = element.clientWidth || 320
      renderer.setSize(next, next, false)
      camera.aspect = 1
      camera.updateProjectionMatrix()
      draw()
    })
    observer.observe(element)

    // Expose the parts the React effect below needs without recreating them.
    const api = { draw, dispose }
    ;(element as HTMLDivElement & { __solid?: typeof api }).__solid = api

    function dispose() {
      observer.disconnect()
      geometry.dispose()
      material.dispose()
      capMaterial.dispose()
      sectionGeometry.dispose()
      sectionMaterial.dispose()
      diskGeometry.dispose()
      diskMaterial.dispose()
      renderer.dispose()
      renderer.domElement.remove()
    }

    return dispose
    // Mount once; interaction redraws through the second effect.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    const element = host.current as (HTMLDivElement & { __solid?: { draw: () => void } }) | null
    element?.__solid?.draw()
  }, [theta, phi, cut, cutX])

  const turn = (dTheta: number, dPhi: number) => {
    setTheta((t) => t + dTheta)
    setPhi((p) => Math.min(Math.PI - 0.15, Math.max(0.15, p + dPhi)))
  }

  return (
    <div>
      <div
        ref={host}
        data-proto="three"
        role="img"
        aria-label="A solid of revolution with a circular cross-section"
        tabIndex={0}
        onKeyDown={(event) => {
          if (event.key === 'ArrowLeft') turn(-0.12, 0)
          else if (event.key === 'ArrowRight') turn(0.12, 0)
          else if (event.key === 'ArrowUp') turn(0, -0.1)
          else if (event.key === 'ArrowDown') turn(0, 0.1)
          else return
          event.preventDefault()
        }}
        onPointerDown={(event) => {
          const start = { x: event.clientX, y: event.clientY }
          const move = (moveEvent: PointerEvent) => {
            turn((moveEvent.clientX - start.x) * 0.01, -(moveEvent.clientY - start.y) * 0.01)
            start.x = moveEvent.clientX
            start.y = moveEvent.clientY
          }
          const up = () => {
            window.removeEventListener('pointermove', move)
            window.removeEventListener('pointerup', up)
          }
          window.addEventListener('pointermove', move)
          window.addEventListener('pointerup', up)
        }}
        className="mx-auto w-full max-w-[420px] cursor-grab touch-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-mathua-blue active:cursor-grabbing"
      />
      <div className="mt-2 flex flex-wrap items-center gap-2 font-mono text-[10px] text-mathua-secondary">
        <label className="flex items-center gap-1">
          <input type="checkbox" checked={cut} onChange={(e) => setCut(e.target.checked)} />
          cross-section
        </label>
        <input
          type="range"
          aria-label="cross-section position"
          min={X_MIN}
          max={X_MAX}
          step={0.02}
          value={cutX}
          onChange={(e) => setCutX(Number(e.target.value))}
          className="min-w-[110px] flex-1"
        />
        <output className="tabular-nums">{cutX.toFixed(2)}</output>
      </div>
    </div>
  )
}
