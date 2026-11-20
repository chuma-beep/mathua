'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Slider } from '../Controls'
import { probabilityTreeSchema } from '../schemas'
import { probabilityTree } from '../math/probability'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 300

export default function ProbabilityTree({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => probabilityTreeSchema.parse(params), [params])
  const [prior, setPrior] = useState(defaults.prior)
  const [sensitivity, setSensitivity] = useState(defaults.sensitivity)
  const [specificity, setSpecificity] = useState(defaults.specificity)

  const tree = probabilityTree(prior, sensitivity, specificity)
  const reset = () => {
    setPrior(defaults.prior)
    setSensitivity(defaults.sensitivity)
    setSpecificity(defaults.specificity)
  }

  if (!tree) {
    return (
      <FigureFrame
        title={title ?? 'A test result updates a prior'}
        readout={[{ label: 'status', value: 'invalid probabilities' }]}
      >
        <p className="font-mono text-[10px] text-mathua-muted">Probabilities must be finite numbers.</p>
      </FigureFrame>
    )
  }

  const branch = (x1: number, y1: number, x2: number, y2: number, key: string) => (
    <line key={key} x1={x1} y1={y1} x2={x2} y2={y2} stroke="var(--text-muted)" strokeWidth="1" />
  )
  const node = (x: number, y: number, key: string, fill = 'var(--bg)') => (
    <circle key={key} cx={x} cy={y} r="2.5" fill={fill} stroke="var(--text-muted)" strokeWidth="1" />
  )
  const label = (x: number, y: number, text: string, key: string, anchor: 'start' | 'middle' | 'end' = 'start') => (
    <text key={key} x={x} y={y} textAnchor={anchor} fontSize="8" fill="var(--text-primary)">
      {text}
    </text>
  )

  return (
    <FigureFrame
      title={title ?? 'A test result updates a prior'}
      hint="Change the prior, the sensitivity or the specificity. The four leaves are the joint probabilities, and the posterior is read from the positive leaves."
      caption={
        caption ??
        'The tree is the sample space split first by disease status and then by test result. Each leaf is a product of the branches above it. P(D | +) is the disease-positive leaf divided by all positive leaves, so a rare disease can stay rare after a positive test.'
      }
      onReset={reset}
      readout={[
        { label: 'prior P(D)', value: fmt(tree.prior, 4) },
        { label: 'sensitivity P(+|D)', value: fmt(tree.sensitivity, 4) },
        { label: 'specificity P(−|¬D)', value: fmt(tree.specificity, 4) },
        { label: 'P(+)', value: fmt(tree.positive, 4) },
        { label: 'P(D∩+)', value: fmt(tree.diseasePositive, 4) },
        { label: 'P(¬D∩+)', value: fmt(tree.healthyPositive, 4) },
        { label: 'P(D∩−)', value: fmt(tree.diseaseNegative, 4) },
        { label: 'P(¬D∩−)', value: fmt(tree.healthyNegative, 4) },
        { label: 'P(D|+)', value: fmt(tree.diseaseGivenPositive, 4) },
        { label: 'P(¬D|+)', value: fmt(tree.healthyGivenPositive, 4) },
        { label: 'P(D|−)', value: fmt(tree.diseaseGivenNegative, 4) },
        { label: 'leaves sum', value: fmt(tree.positive + tree.negative, 4) },
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`Probability tree with prior ${fmt(tree.prior, 3)}, sensitivity ${fmt(tree.sensitivity, 3)} and specificity ${fmt(tree.specificity, 3)}`}
      >
        {branch(34, 150, 132, 78, 'root-disease')}
        {branch(34, 150, 132, 222, 'root-healthy')}
        {branch(132, 78, 238, 42, 'disease-positive')}
        {branch(132, 78, 238, 114, 'disease-negative')}
        {branch(132, 222, 238, 186, 'healthy-positive')}
        {branch(132, 222, 238, 258, 'healthy-negative')}

        {node(34, 150, 'root')}
        {node(132, 78, 'disease')}
        {node(132, 222, 'healthy')}
        {node(238, 42, 'disease-positive-node', 'var(--accent-blue)')}
        {node(238, 114, 'disease-negative-node')}
        {node(238, 186, 'healthy-positive-node', 'var(--accent-blue)')}
        {node(238, 258, 'healthy-negative-node')}

        {label(8, 146, 'start', 'start-label')}
        {label(70, 106, `D ${fmt(tree.disease, 3)}`, 'disease-label')}
        {label(70, 212, `¬D ${fmt(tree.healthy, 3)}`, 'healthy-label')}
        {label(164, 52, `+ ${fmt(tree.positiveGivenDisease, 3)}`, 'disease-positive-label')}
        {label(164, 112, `− ${fmt(tree.negativeGivenDisease, 3)}`, 'disease-negative-label')}
        {label(164, 196, `+ ${fmt(tree.positiveGivenHealthy, 3)}`, 'healthy-positive-label')}
        {label(164, 258, `− ${fmt(tree.negativeGivenHealthy, 3)}`, 'healthy-negative-label')}

        {label(244, 38, `D∩+ ${fmt(tree.diseasePositive, 4)}`, 'disease-positive-leaf')}
        {label(244, 110, `D∩− ${fmt(tree.diseaseNegative, 4)}`, 'disease-negative-leaf')}
        {label(244, 182, `¬D∩+ ${fmt(tree.healthyPositive, 4)}`, 'healthy-positive-leaf')}
        {label(244, 254, `¬D∩− ${fmt(tree.healthyNegative, 4)}`, 'healthy-negative-leaf')}

        <text x={244} y={286} fontSize="8" fill="var(--accent-blue)">
          P(+) = {fmt(tree.positive, 4)}
        </text>
      </svg>

      <div className="mt-3 grid grid-cols-1 gap-1.5">
        <Slider label="Prior" value={prior} min={0} max={1} step={0.01} onChange={setPrior} />
        <Slider
          label="Sensitivity"
          value={sensitivity}
          min={0}
          max={1}
          step={0.01}
          onChange={setSensitivity}
        />
        <Slider
          label="Specificity"
          value={specificity}
          min={0}
          max={1}
          step={0.01}
          onChange={setSpecificity}
        />
      </div>
    </FigureFrame>
  )
}
