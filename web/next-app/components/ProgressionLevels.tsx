interface Level {
  num: string
  name: string
  range: string
  elite?: boolean
}

interface ProgressionLevelsProps {
  levels: Level[]
  className?: string
}

export default function ProgressionLevels({ levels, className = '' }: ProgressionLevelsProps) {
  return (
    <ul className={`list-none ${className}`}>
      {levels.map((level) => (
        <li
          key={level.num}
          className={`flex items-center gap-4 py-2.5 border-b border-mathua-border text-sm text-mathua-primary ${
            level.elite
              ? 'bg-mathua-blue/10 rounded-md px-3 -mx-1'
              : ''
          }`}
        >
          <span className="font-mono text-[11px] text-mathua-muted min-w-[24px]">
            {level.num}
          </span>
          <span
            className={`flex-1 font-medium ${
              level.elite ? 'text-mathua-blue' : ''
            }`}
          >
            {level.name}
          </span>
          <span className="font-mono text-[11px] text-mathua-muted text-right">
            {level.range}
          </span>
        </li>
      ))}
    </ul>
  )
}
