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
    <div className={`w-full max-w-full min-w-0 overflow-hidden ${className}`}>
      <h3 className="font-serif text-[1.1rem] font-normal text-mathua-primary mb-3 tracking-tight">
        Levels of mastery
      </h3>
      <div className="w-full max-w-full overflow-x-auto overscroll-x-contain -mx-4 px-4 sm:mx-0 sm:px-0">
        <table className="w-full min-w-[280px] border-collapse border-none bg-transparent">
          <tbody>
            {levels.map((level) => (
              <tr
                key={level.num}
                className="border-b-[0.5px] border-mathua-border bg-transparent"
                style={level.elite ? { borderTop: '1px solid var(--accent-blue)' } : undefined}
              >
                <td className="font-serif text-[13px] text-mathua-blue py-2 pr-3 whitespace-nowrap align-baseline">
                  {parseInt(level.num)}.
                </td>
                <td className="font-serif text-[14px] sm:text-base text-mathua-primary py-2 align-baseline min-w-0">
                  <span className="block truncate sm:whitespace-nowrap">{level.name}</span>
                </td>
                <td className="font-serif text-[12px] text-mathua-muted text-right py-2 pl-3 whitespace-nowrap align-baseline">
                  {level.range} concepts
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
