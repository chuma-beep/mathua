// What both renderers receive. The parent owns the camera and the cut so the
// WebGL scene and the SVG fallback expose the same controls, the same readout
// and the same keyboard path; a renderer only draws and reports drag deltas.
export interface SolidViewProps {
  theta: number
  phi: number
  cut: number
  showCut: boolean
  onOrbit: (dTheta: number, dPhi: number) => void
}
