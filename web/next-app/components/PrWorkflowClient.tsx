'use client'

import dynamic from 'next/dynamic'
import Loading from './Loading'

// Next 16 rejects `ssr: false` in a Server Component, so the dynamic import
// lives in this small client boundary. The workflow diagram measures its own
// container, which is why it stays client-only.
const PrWorkflow = dynamic(() => import('./PrWorkflow'), {
  ssr: false,
  loading: () => (
    <div style={{
      height: 400, border: '0.5px solid var(--border)',
      display: 'flex', alignItems: 'center', justifyContent: 'center',
      color: 'var(--text-muted)', fontFamily: "var(--font-jetbrains-mono), 'JetBrains Mono', monospace", fontSize: '13px',
    }}>
      <Loading label="LOADING WORKFLOW" />
    </div>
  ),
})

export default PrWorkflow
