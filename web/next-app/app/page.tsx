'use client'

import dynamic from 'next/dynamic'
import { useState, useEffect } from 'react'
import './globals.css'

const MathConceptGraph3D = dynamic(() => import('../components/MathConceptGraph3D'), {
  ssr: false,
  loading: () => (
    <div style={{ 
      height: '520px', 
      background: 'var(--bg-card, #1a1a1a)', 
      borderRadius: '8px', 
      border: '0.5px solid var(--border, #333)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      color: 'var(--text-muted, #888)',
      fontSize: '14px'
    }}>Loading graph...</div>
  ),
})

export default function HomePage() {
  const [theme, setTheme] = useState<'dark' | 'light'>('dark')
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    const saved = localStorage.getItem('mathua-theme')
    const initialTheme = (saved === 'light' || saved === 'dark') ? saved : 'dark'
    setTheme(initialTheme)
    setMounted(true)
  }, [])

  useEffect(() => {
    if (!mounted) return
    if (theme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [theme, mounted])

  const toggleTheme = () => {
    const next = theme === 'dark' ? 'light' : 'dark'
    setTheme(next)
    localStorage.setItem('mathua-theme', next)
  }

  if (!mounted) {
    return <div style={{ background: '#0c0c0d', minHeight: '100vh' }} />
  }

  return (
    <div className="container">
      <button onClick={toggleTheme} className="theme-toggle" aria-label="Toggle theme">
        {theme === 'dark' ? '☀' : '☾'}
      </button>

      <section className="hero">
        <div className="hero-nav">
          <a href="/how-it-works" className="btn btn-outline">How it works</a>
        </div>
        <h1>Master the foundation.<br />Earn the abstraction.</h1>
        <p className="subtitle">
          Mathua is an open-source adaptive math learning engine. It never lets you advance until you have truly mastered the prerequisite — both speed and accuracy must be proven.
        </p>
        <div className="hero-cta">
         <a href="#" className="btn btn-blue">Open the web app</a>
         <a href="https://github.com/chuma-beep/mathua" className="btn btn-outline">View on GitHub</a>  
        </div>
        <div className="stat-badges">
          <span className="stat-badge">60 topics</span>
          <span className="stat-badge">120+ connections</span>
          <span className="stat-badge">9 fields</span>
          <span className="stat-badge">web + desktop</span>
          <span className="stat-badge">open source</span>
        </div>
        <div className="">
          <div className="max-w-[900px]  mt-0">
            <MathConceptGraph3D theme={theme} />
          </div>
        </div>
      </section>

      <hr className="section-divider" />

      <section className="card-section">
        <span className="section-label">How it works</span>
        <h2 className="text-center">A different kind of math system</h2>
        <div className="grid-three">
          <div className="info-card">
            <h4>Mastery gating</h4>
            <p>You cannot advance until your streak and response time both meet the threshold. Knowing the answer is not enough — you must know it fast.</p>
          </div>
          <div className="info-card">
            <h4>Concept graph</h4>
            <p>Every concept is a node with explicit prerequisites. The scheduler reads the graph and your progress to decide what you see next. Nothing is shown before its foundation is solid.</p>
          </div>
          <div className="info-card">
            <h4>Spaced repetition</h4>
            <p>Concepts you master resurface automatically before they decay. Reviews are woven into your session — there is no separate review mode.</p>
          </div>
        </div>
      </section>

      <section className="card-section">
        <span className="section-label">Mastery pipeline</span>
        <h2>Five states. One direction.</h2>
        <div className="pipeline-flow">
          <span className="pipeline-state">UNSEEN</span>
          <span className="pipeline-arrow">→</span>
          <span className="pipeline-state">LEARNING</span>
          <span className="pipeline-arrow">→</span>
          <span className="pipeline-state">PRACTICING</span>
          <span className="pipeline-arrow">→</span>
          <span className="pipeline-state mastered">MASTERED</span>
          <span className="pipeline-arrow">→</span>
          <span className="pipeline-state">DECAYING</span>
        </div>
       <div className="formula-container">
        <pre className="formula-block">
{`priority = (0.7 × days_since_last_seen) + (0.3 × (1 − mastery))
         + 5.0 if DECAYING  +  2.0 if newly_unlocked`}
        </pre>
      </div> 
    <p className="" style={{ textAlign: 'center' }}>
          The scheduler enforces three hard rules: prerequisites must be mastered before a concept unlocks, the same concept never appears twice in a row, and roughly 70% of each session is new material.
        </p>
      </section>

      <hr className="section-divider" />

      <section className="card-section">
        <span className="section-label">Curriculum</span>
        <h2>What Mathua covers</h2>
        <table className="domain-table">
          <thead>
            <tr><th>Domain</th><th style={{textAlign: 'right'}}>Concepts</th></tr>
          </thead>
          <tbody>
            <tr><td>Counting</td><td className="concept-count">7</td></tr>
            <tr><td>Arithmetic</td><td className="concept-count">39</td></tr>
            <tr><td>Fractions</td><td className="concept-count">18</td></tr>
            <tr><td>Pre-Algebra</td><td className="concept-count">17</td></tr>
            <tr><td>Algebra — v1.1 <span className="coming-soon">coming soon</span></td><td className="concept-count">—</td></tr>
          </tbody>
        </table>
        <p className="table-note">Problems are generated on demand — never stored. There is nothing to memorise.</p>
      </section>

      <hr className="section-divider" />

      <section className="card-section">
        <span className="section-label">Progression</span>
        <h2>Rank by mastery. Level by depth.</h2>
        <div className="grid-two">
          <div>
            <p className="body-text" style={{ textAlign: 'center' }}>
              The leaderboard resets every Monday at 00:00 UTC. Your score is calculated from three components:
            </p>
             <div className="formula-container">
            <pre className="formula-block">
       {`score = (mastered_count × 100)
          + speed_bonus
           + (current_streak × 10)`}
            </pre>

             </div> 

          </div>
          <div>
            <ul className="level-list">
              <li className="level-item"><span className="level-num">01</span><span className="level-name">Novice</span><span className="level-range">0–5</span></li>
              <li className="level-item"><span className="level-num">02</span><span className="level-name">Apprentice</span><span className="level-range">6–15</span></li>
              <li className="level-item"><span className="level-num">03</span><span className="level-name">Student</span><span className="level-range">16–25</span></li>
              <li className="level-item"><span className="level-num">04</span><span className="level-name">Scholar</span><span className="level-range">26–35</span></li>
              <li className="level-item"><span className="level-num">05</span><span className="level-name">Adept</span><span className="level-range">36–45</span></li>
              <li className="level-item"><span className="level-num">06</span><span className="level-name">Expert</span><span className="level-range">46–55</span></li>
              <li className="level-item"><span className="level-num">07</span><span className="level-name">Master</span><span className="level-range">56–65</span></li>
              <li className="level-item"><span className="level-num">08</span><span className="level-name">Grandmaster</span><span className="level-range">66–75</span></li>
              <li className="level-item elite"><span className="level-num">09</span><span className="level-name">Math Architect</span><span className="level-range">76–81</span></li>
            </ul>
          </div>
        </div>
      </section>

      <hr className="section-divider" />

      <section className="card-section">
        <span className="section-label">Platforms</span>
        <h2>One engine. Two ways to run it.</h2>
        <div className="grid-two">
          <div className="client-card">
            <div className="client-type">Web</div>
            <h3>Browser</h3>
            <ul>
              <li>Runs in any modern browser</li>
              <li>React frontend with KaTeX math rendering</li>
              <li>Account required — progress syncs across devices</li>
              <li>Global weekly leaderboard</li>
              <li>Interactive concept graph view</li>
            </ul>
          </div>
          <div className="client-card">
            <div className="client-type">Desktop TUI</div>
            <h3>Terminal</h3>
            <ul>
              <li>Single binary download — no runtime dependencies</li>
              <li>Bubble Tea terminal interface</li>
              <li>SQLite storage — all data stays on your machine</li>
              <li>No account needed</li>
              <li>Fully offline after download</li>
            </ul>
          </div>
        </div>
     <div className="formula-container">
        <pre className="architecture-block">
{`┌──────────────────┐    ┌──────────────────┐
│  Web (React)    │    │  TUI (Bubble Tea)│
└────────┬─────────┘    └────────┬─────────┘
         ▼                       ▼
┌──────────────────────────────────────┐
│            API Layer               │
└──────────────────┬───────────────────┘
                   ▼
┌──────────────────────────────────────┐
│        Scheduling Engine           │
└──────────────────┬────────────────┘
                   ▼
┌──────────────────────────────────────┐
│   Grading & Problem Generation    │
└──────────────────┬───────────────────┘
                   ▼
┌──────────────────────────────────────┐
│  Data Layer (PostgreSQL / SQLite) │
└──────────────────────────────────────┘`}
        </pre>
       
     </div>
        </section>

      <hr className="section-divider" />

      <section className="card-section">
        <span className="section-label">Contributing</span>
        <h2>Built to be extended.</h2>
        <p className="body-text" style={{ textAlign: 'center' }}>
          Every concept is a JSON node. Every problem is a Go generator function. Every contribution goes through a graph validator that rejects cycles and orphaned nodes automatically.
        </p>
          <div className="formula-container">
        <pre className="json-block" style={{ textAlign: 'left' }}>
{`{
  "id":                "arith.add.multi",
  "label":             "Multi-digit addition",
  "domain":            "arithmetic",
  "prerequisites":     ["arith.add.single", "arith.add.carry"],
  "mastery_threshold": {
    "streak":       5,
    "max_seconds":  8
  }
}`}
        </pre>
        </div>
    <div className="contrib-button">
        <a href="https://github.com/chuma-beep/mathua/blob/main/CONTRIBUTING.md" className="btn btn-blue ">Read CONTRIBUTING.md →</a>         
         </div>
        </section>

      <footer className="footer">
        <div className="footer-grid">
          <div>
            <div className="footer-brand">Mathua</div>
            <div className="footer-tagline">Math Understanding Agent</div>
          </div>
          <div className="footer-links">
            <a href="#">Web App</a>
            <a href="https://github.com/chuma-beep/mathua">GitHub</a>
            <a href="https://github.com/chuma-beep/mathua/blob/main/docs/architecture.md">Docs</a>
            <a href="https://github.com/chuma-beep/mathua/blob/main/CONTRIBUTING.md">Contributing</a>
            <a href="#">Roadmap</a>
          </div>
        </div>
        <div className="footer-bottom">
          Inspired by the mastery-gating philosophy of Math Academy. No content or code from Math Academy is used.
        </div>
      </footer>
    </div>
  )
}
