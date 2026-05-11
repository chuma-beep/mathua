interface FooterProps {
  className?: string
}

export default function Footer({ className = '' }: FooterProps) {
  return (
    <footer
      className={`bg-mathua-surface-elevated py-10 px-6 mt-10 border-t border-mathua-border ${className}`}
    >
      <div className="max-w-container mx-auto grid grid-cols-[1fr_2fr] gap-6 items-center max-md:grid-cols-1 max-md:text-center">
        <div>
          <div className="text-xl font-semibold text-mathua-blue">Mathua</div>
          <div className="font-mono text-[10px] text-mathua-muted uppercase tracking-[0.1em] mt-1">
            Math Understanding Agent
          </div>
        </div>
        <div className="flex gap-5 justify-end text-[13px] text-mathua-secondary max-md:justify-center max-md:flex-col max-md:items-center max-md:gap-3">
          <a href="#">Web App</a>
          <a href="https://github.com/chuma-beep/mathua">GitHub</a>
          <a href="https://github.com/chuma-beep/mathua/blob/main/docs/architecture.md">
            Docs
          </a>
          <a href="https://github.com/chuma-beep/mathua/blob/main/CONTRIBUTING.md">
            Contributing
          </a>
          <a href="#">Roadmap</a>
        </div>
      </div>
      <div className="max-w-container mx-auto mt-6 pt-5 border-t border-mathua-border text-center text-[11px] text-mathua-muted opacity-70">
        Inspired by the mastery-gating philosophy of Math Academy. No content or code from Math
        Academy is used.
      </div>
    </footer>
  )
}
