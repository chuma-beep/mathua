'use client'

import { useEffect, useRef } from 'react'
import * as THREE from 'three'
import { CENTER_X, X_MAX, X_MIN, radius, surface } from '../../math/solid'
import type { SolidViewProps } from './props'

/**
 * The WebGL renderer. It is imported only through the dispatcher's dynamic
 * import, so a lesson with no 3D step never downloads three.js.
 *
 * Nothing here animates: a camera change or a cut change renders one frame and
 * stops. There is no requestAnimationFrame loop, which is why a figure sitting
 * on screen costs nothing while the learner reads the paragraph above it.
 */
export default function SolidWebGL({ theta, phi, cut, showCut, onOrbit }: SolidViewProps) {
  const host = useRef<HTMLDivElement | null>(null)
  const scene = useRef<{
    renderer: THREE.WebGLRenderer
    scene: THREE.Scene
    camera: THREE.PerspectiveCamera
    group: THREE.Group
    section: THREE.Mesh
    ring: THREE.LineLoop
    render: () => void
  } | null>(null)
  const dragging = useRef(false)

  useEffect(() => {
    const element = host.current
    if (!element) return

    const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true })
    renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2))
    element.appendChild(renderer.domElement)

    const world = new THREE.Scene()
    const camera = new THREE.PerspectiveCamera(35, 1, 0.1, 100)
    camera.position.set(0, 0, 9)
    camera.lookAt(0, 0, 0)

    world.add(new THREE.AmbientLight(0xffffff, 0.55))
    const key = new THREE.DirectionalLight(0xffffff, 1.05)
    key.position.set(-5, 3, -6)
    world.add(key)

    const group = new THREE.Group()
    world.add(group)

    const mesh = surface(48, 64)
    const geometry = new THREE.BufferGeometry()
    geometry.setAttribute('position', new THREE.Float32BufferAttribute(mesh.positions, 3))
    geometry.setAttribute('normal', new THREE.Float32BufferAttribute(mesh.normals, 3))
    geometry.setIndex(mesh.indices)
    const material = new THREE.MeshLambertMaterial({ color: 0x4a86c8, side: THREE.DoubleSide })
    group.add(new THREE.Mesh(geometry, material))

    // Without caps the camera looks into an open tube rather than at a solid.
    const capGeometry = (radiusAt: number) => new THREE.CircleGeometry(radiusAt, 64)
    const capMaterial = new THREE.MeshLambertMaterial({ color: 0x3f78b4, side: THREE.DoubleSide })
    const leftCap = new THREE.Mesh(capGeometry(radius(X_MIN)), capMaterial)
    leftCap.rotation.y = -Math.PI / 2
    leftCap.position.x = X_MIN - CENTER_X
    group.add(leftCap)
    const rightCap = new THREE.Mesh(capGeometry(radius(X_MAX)), capMaterial)
    rightCap.rotation.y = Math.PI / 2
    rightCap.position.x = X_MAX - CENTER_X
    group.add(rightCap)

    const section = new THREE.Mesh(
      new THREE.CircleGeometry(1, 64),
      new THREE.MeshBasicMaterial({ color: 0x2a9d8f, transparent: true, opacity: 0.28, side: THREE.DoubleSide }),
    )
    section.rotation.y = Math.PI / 2
    group.add(section)

    const ringPoints = Array.from({ length: 64 }, (_, i) => {
      const angle = (2 * Math.PI * i) / 64
      return new THREE.Vector3(0, Math.cos(angle), Math.sin(angle))
    })
    const ring = new THREE.LineLoop(
      new THREE.BufferGeometry().setFromPoints(ringPoints),
      new THREE.LineBasicMaterial({ color: 0x2a9d8f }),
    )
    group.add(ring)

    const render = () => renderer.render(world, camera)

    const resize = () => {
      const size = element.clientWidth || 320
      renderer.setSize(size, size, false)
      renderer.domElement.style.width = '100%'
      renderer.domElement.style.height = 'auto'
      camera.aspect = 1
      camera.updateProjectionMatrix()
      render()
    }
    resize()
    const observer = typeof ResizeObserver === 'undefined' ? null : new ResizeObserver(resize)
    observer?.observe(element)

    scene.current = { renderer, scene: world, camera, group, section, ring, render }

    return () => {
      observer?.disconnect()
      geometry.dispose()
      material.dispose()
      capMaterial.dispose()
      leftCap.geometry.dispose()
      rightCap.geometry.dispose()
      section.geometry.dispose()
      ;(section.material as THREE.Material).dispose()
      ring.geometry.dispose()
      ;(ring.material as THREE.Material).dispose()
      renderer.dispose()
      renderer.domElement.remove()
      scene.current = null
    }
  }, [])

  useEffect(() => {
    const current = scene.current
    if (!current) return
    current.group.rotation.y = -theta
    current.group.rotation.x = phi
    current.section.visible = showCut
    current.ring.visible = showCut
    const r = radius(cut)
    current.section.scale.set(r, r, 1)
    current.section.position.x = cut - CENTER_X
    current.ring.position.x = cut - CENTER_X
    current.ring.scale.set(1, r, r)
    current.render()
  }, [theta, phi, cut, showCut])

  return (
    <div
      ref={host}
      className="mx-auto w-full max-w-[420px] cursor-grab touch-none active:cursor-grabbing [&>canvas]:block [&>canvas]:rounded-none"
      role="img"
      aria-label="A solid of revolution rendered in three dimensions"
      onPointerDown={(event) => {
        dragging.current = true
        event.currentTarget.setPointerCapture(event.pointerId)
      }}
      onPointerMove={(event) => {
        if (!dragging.current) return
        onOrbit(event.movementX * 0.01, event.movementY * 0.01)
      }}
      onPointerUp={() => {
        dragging.current = false
      }}
      onPointerCancel={() => {
        dragging.current = false
      }}
    />
  )
}
