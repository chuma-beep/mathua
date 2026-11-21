'use client'

import { useEffect, useRef } from 'react'
import * as THREE from 'three'
import ThreeGlobe from 'three-globe'
import { COUNTRY_CENTROIDS } from '../../lib/countryCentroids'
import { useTheme } from '../../hooks/useTheme'

export type GlobeCountry = {
  code: string
  active_7d: number
}

type CountryFeature = {
  type: string
  properties: Record<string, unknown>
  geometry: unknown
}

// The polygon data is 488 KB and never changes, so the module keeps one
// promise for every mount instead of refetching on a theme toggle.
let countriesPromise: Promise<CountryFeature[]> | null = null
function loadCountries(): Promise<CountryFeature[]> {
  if (!countriesPromise) {
    countriesPromise = fetch('/globe/countries.geojson')
      .then((res) => {
        if (!res.ok) throw new Error(`countries ${res.status}`)
        return res.json()
      })
      .then((data: { features?: CountryFeature[] }) =>
        Array.isArray(data.features) ? data.features : [],
      )
      // A failed polygon fetch leaves the plain sphere and markers, which is
      // still a usable summary; the table below is the authoritative data.
      .catch(() => [])
  }
  return countriesPromise
}

const PALETTE = {
  dark: {
    ocean: '#0b1526',
    oceanEmissive: '#091120',
    land: '#5a9bdc',
    atmosphere: '#3b82f6',
    marker: '#a5f3fc',
    ring: '165, 243, 252',
  },
  light: {
    ocean: '#e9f1fa',
    oceanEmissive: '#dde9f6',
    land: '#8fb8dd',
    atmosphere: '#93c5fd',
    marker: '#1d4ed8',
    ring: '37, 99, 235',
  },
}

// A polygon globe with the Natural Earth outlines, one point per country and
// a pulsing ring for the busiest few. The canvas is decorative: the country
// table below carries the same numbers as text.
//
// Rotation stops under prefers-reduced-motion. Dragging still works because
// that motion is the reader's, not the page's.
export default function CountryGlobe({ countries }: { countries: GlobeCountry[] }) {
  const canvas = useRef<HTMLCanvasElement>(null)
  const { theme, mounted } = useTheme()
  const signature = countries.map((c) => `${c.code}:${c.active_7d}`).join(',')

  useEffect(() => {
    const element = canvas.current
    if (!element || !mounted) return

    const dark = theme === 'dark'
    const palette = dark ? PALETTE.dark : PALETTE.light
    const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    const max = Math.max(1, ...countries.map((c) => c.active_7d))

    const markers = countries.flatMap((c) => {
      const [lat, lng] = COUNTRY_CENTROIDS[c.code] ?? []
      if (lat === undefined || lng === undefined) return []
      return [{
        lat,
        lng,
        radius: 0.45 + 1.25 * Math.sqrt(c.active_7d / max),
      }]
    })
    const rings = [...markers].slice(0, 5).map((m) => ({
      ...m,
      radius: 2.5 + 4.5 * Math.sqrt((m.radius - 0.45) / 1.25),
    }))

    const scene = new THREE.Scene()
    const camera = new THREE.PerspectiveCamera(40, 1, 1, 3000)
    camera.position.set(0, 0, 360)

    const renderer = new THREE.WebGLRenderer({ canvas: element, antialias: true, alpha: true })
    renderer.setPixelRatio(Math.min(2, window.devicePixelRatio || 1))
    renderer.setClearColor(0x000000, 0)

    const group = new THREE.Group()
    group.rotation.set(0.32, -0.55, 0)
    scene.add(group)

    const globe = new ThreeGlobe({ animateIn: false })
      .globeMaterial(new THREE.MeshPhongMaterial({
        color: palette.ocean,
        emissive: palette.oceanEmissive,
        emissiveIntensity: dark ? 0.45 : 0.35,
        shininess: 14,
      }))
      .showAtmosphere(true)
      .atmosphereColor(palette.atmosphere)
      .atmosphereAltitude(0.18)
      // Fine hexagons at a small margin tile into solid continents. At the
      // original 0.72 margin the land read as a faint texture on dark.
      .hexPolygonsData([])
      .hexPolygonResolution(3)
      .hexPolygonMargin(0.16)
      .hexPolygonUseDots(false)
      // Accessors treat a bare string as a property name, so the palette
      // colors need a function even though they are constant.
      .hexPolygonColor(() => palette.land)
      .pointsData(markers)
      .pointLat('lat')
      .pointLng('lng')
      .pointAltitude(0.012)
      .pointRadius('radius')
      .pointColor(() => palette.marker)
      .ringsData(rings)
      .ringLat('lat')
      .ringLng('lng')
      .ringAltitude(0.01)
      .ringMaxRadius('radius')
      .ringPropagationSpeed(1.6)
      .ringRepeatPeriod(1800)
      .ringColor(() => (t: number) => `rgba(${palette.ring}, ${Math.max(0, 1 - t)})`)
    group.add(globe)

    scene.add(new THREE.AmbientLight(0xffffff, dark ? 1.8 : 2.1))
    const key = new THREE.DirectionalLight(0xffffff, dark ? 1.4 : 1.0)
    key.position.set(-160, 90, 180)
    scene.add(key)
    const rim = new THREE.DirectionalLight(0x60a5fa, dark ? 1.2 : 0.4)
    rim.position.set(150, -70, -160)
    scene.add(rim)

    let live = true
    let needsRender = true
    let dragging = false
    let lastX = 0
    let lastY = 0

    void loadCountries().then((features) => {
      if (!live || features.length === 0) return
      globe.hexPolygonsData(features)
      needsRender = true
    })

    const resize = () => {
      const width = Math.max(1, element.clientWidth)
      const height = Math.max(1, element.clientHeight)
      renderer.setSize(width, height, false)
      camera.aspect = width / height
      camera.updateProjectionMatrix()
      needsRender = true
    }
    const observer = new ResizeObserver(resize)
    observer.observe(element)
    resize()

    const onDown = (event: PointerEvent) => {
      dragging = true
      lastX = event.clientX
      lastY = event.clientY
      element.setPointerCapture(event.pointerId)
    }
    const onMove = (event: PointerEvent) => {
      if (!dragging) return
      group.rotation.y += (event.clientX - lastX) * 0.005
      group.rotation.x = Math.max(-0.9, Math.min(0.9, group.rotation.x + (event.clientY - lastY) * 0.005))
      lastX = event.clientX
      lastY = event.clientY
      needsRender = true
    }
    const onUp = (event: PointerEvent) => {
      dragging = false
      if (element.hasPointerCapture(event.pointerId)) element.releasePointerCapture(event.pointerId)
    }
    element.addEventListener('pointerdown', onDown)
    element.addEventListener('pointermove', onMove)
    element.addEventListener('pointerup', onUp)
    element.addEventListener('pointercancel', onUp)

    let frame = 0
    const tick = () => {
      frame = requestAnimationFrame(tick)
      if (!reduced && !dragging) {
        group.rotation.y += 0.0016
        needsRender = true
      }
      if (needsRender) {
        renderer.render(scene, camera)
        needsRender = false
      }
    }
    tick()

    return () => {
      live = false
      cancelAnimationFrame(frame)
      observer.disconnect()
      element.removeEventListener('pointerdown', onDown)
      element.removeEventListener('pointermove', onMove)
      element.removeEventListener('pointerup', onUp)
      element.removeEventListener('pointercancel', onUp)
      scene.traverse((object) => {
        const mesh = object as THREE.Mesh
        mesh.geometry?.dispose()
        const material = mesh.material
        if (Array.isArray(material)) material.forEach((m) => m.dispose())
        else material?.dispose()
      })
      renderer.dispose()
    }
    // signature carries every value the markers depend on; countries is a new
    // array on each parent render and would rebuild the WebGL context.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [theme, mounted, signature])

  return (
    <div className="relative mx-auto aspect-square w-full max-w-[480px]">
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-[8%] rounded-full bg-mathua-blue-faint blur-2xl"
      />
      <canvas
        ref={canvas}
        aria-hidden="true"
        className="relative h-full w-full touch-none cursor-grab active:cursor-grabbing"
      />
      <p className="sr-only">
        A globe with the countries that have known accounts, a point sized by this week&apos;s
        active learners and a ring for the busiest countries. The table below lists the same
        countries with their account and activity counts.
      </p>
    </div>
  )
}
