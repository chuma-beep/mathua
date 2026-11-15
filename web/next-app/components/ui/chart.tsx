'use client'

import * as React from 'react'
import * as RechartsPrimitive from 'recharts'
import type { TooltipContentProps } from 'recharts'

import { cn } from '@/lib/utils'

// shadcn's chart wrapper around Recharts. The config maps each data key to a
// label and colour so tooltips and legends never repeat those strings, and the
// generated style block exposes `--color-<key>` to the chart primitives.
export type ChartConfig = {
  [key: string]: {
    label?: React.ReactNode
    color?: string
  }
}

type ChartContextProps = {
  config: ChartConfig
}

const ChartContext = React.createContext<ChartContextProps | null>(null)

function useChart() {
  const context = React.useContext(ChartContext)
  if (!context) {
    throw new Error('useChart must be used within a <ChartContainer />')
  }
  return context
}

function ChartContainer({
  id,
  className,
  children,
  config,
  ...props
}: React.ComponentProps<'div'> & {
  config: ChartConfig
  children: React.ComponentProps<typeof RechartsPrimitive.ResponsiveContainer>['children']
}) {
  const uniqueId = React.useId()
  const chartId = `chart-${id || uniqueId.replace(/:/g, '')}`
  const chartStyle = Object.fromEntries(
    Object.entries(config)
      .filter(([, item]) => item.color)
      .map(([key, item]) => [`--color-${key}`, item.color]),
  ) as React.CSSProperties

  return (
    <ChartContext.Provider value={{ config }}>
      <div
        data-chart={chartId}
        style={chartStyle}
        className={cn(
          'flex aspect-auto justify-center text-[10px] [&_.recharts-cartesian-axis-tick_text]:fill-mathua-muted [&_.recharts-cartesian-grid_line]:stroke-mathua-border [&_.recharts-curve.recharts-tooltip-cursor]:stroke-mathua-border [&_.recharts-layer]:outline-none [&_.recharts-surface]:outline-none',
          className,
        )}
        {...props}
      >
        <RechartsPrimitive.ResponsiveContainer>
          {children}
        </RechartsPrimitive.ResponsiveContainer>
      </div>
    </ChartContext.Provider>
  )
}

const ChartTooltip = RechartsPrimitive.Tooltip

function ChartTooltipContent({
  active,
  payload,
  className,
  indicator = 'dot',
  hideLabel = false,
  hideIndicator = false,
  label,
  labelFormatter,
  labelClassName,
  color,
  nameKey,
  labelKey,
}: Partial<TooltipContentProps<number, string>> &
  React.ComponentProps<'div'> & {
    hideLabel?: boolean
    hideIndicator?: boolean
    indicator?: 'line' | 'dot' | 'dashed'
    nameKey?: string
    labelKey?: string
  }) {
  const { config } = useChart()

  const tooltipLabel = React.useMemo(() => {
    if (hideLabel || !payload?.length) return null
    const [item] = payload
    const key = `${labelKey || item?.dataKey || item?.name || 'value'}`
    const itemConfig = getPayloadConfigFromPayload(config, item, key)
    const value =
      !labelKey && typeof label === 'string'
        ? config[label as keyof typeof config]?.label || label
        : itemConfig?.label

    if (labelFormatter) {
      return <div className={cn('font-mono text-[10px] text-mathua-secondary', labelClassName)}>{labelFormatter(value, payload)}</div>
    }

    if (!value) return null
    return <div className={cn('font-mono text-[10px] text-mathua-secondary', labelClassName)}>{value}</div>
  }, [label, labelFormatter, payload, hideLabel, labelClassName, config, labelKey])

  if (!active || !payload?.length) return null

  const nestLabel = payload.length === 1 && indicator !== 'dot'

  return (
    <div
      className={cn(
        'grid min-w-[9rem] items-start gap-1.5 border border-mathua-border-strong bg-mathua-surface px-2.5 py-2 font-mono text-[10px] shadow-none',
        className,
      )}
    >
      {!nestLabel ? tooltipLabel : null}
      <div className="grid gap-1">
        {payload.map((item, index) => {
          const key = `${nameKey || item.name || item.dataKey || 'value'}`
          const itemConfig = getPayloadConfigFromPayload(config, item, key)
          const indicatorColor = color || item.payload?.fill || item.color

          return (
            <div key={`${item.dataKey}-${index}`} className="flex w-full items-center gap-2">
              {!hideIndicator && (
                <span
                  className={cn(
                    'inline-block shrink-0',
                    indicator === 'dot' && 'h-1.5 w-1.5',
                    indicator === 'line' && 'h-0.5 w-3',
                    indicator === 'dashed' && 'h-0 w-3 border-t border-dashed',
                  )}
                  style={{ background: indicator === 'dashed' ? 'transparent' : indicatorColor, borderColor: indicatorColor }}
                />
              )}
              <span className="flex-1 text-mathua-muted">{itemConfig?.label ?? item.name}</span>
              <span className="tabular-nums text-mathua-primary">
                {typeof item.value === 'number' ? item.value.toLocaleString() : item.value}
              </span>
            </div>
          )
        })}
      </div>
    </div>
  )
}

const ChartLegend = RechartsPrimitive.Legend

function ChartLegendContent({
  className,
  hideIcon = false,
  payload,
  verticalAlign = 'bottom',
  nameKey,
}: React.ComponentProps<'div'> & {
  hideIcon?: boolean
  nameKey?: string
  verticalAlign?: 'top' | 'bottom'
  payload?: Array<{ value?: string; dataKey?: string; color?: string }>
}) {
  const { config } = useChart()
  if (!payload?.length) return null

  return (
    <div
      className={cn(
        'flex flex-wrap items-center justify-center gap-x-4 gap-y-1 pt-2',
        verticalAlign === 'top' ? 'pb-2 pt-0' : 'pt-2',
        className,
      )}
    >
      {payload.map(item => {
        const key = `${nameKey || item.dataKey || 'value'}`
        const itemConfig = config[key as keyof typeof config]
        return (
          <div key={String(item.value)} className="flex items-center gap-1.5 font-mono text-[9px] text-mathua-muted">
            {!hideIcon && <span className="h-2 w-2 shrink-0" style={{ background: item.color ?? itemConfig?.color }} />}
            {itemConfig?.label ?? item.value}
          </div>
        )
      })}
    </div>
  )
}

function getPayloadConfigFromPayload(config: ChartConfig, payload: unknown, key: string) {
  if (typeof payload !== 'object' || payload === null) return undefined
  const payloadObject = payload as Record<string, unknown>
  const nested = payloadObject.payload as Record<string, unknown> | undefined
  let configLabelKey: string = key

  if (key in payloadObject && typeof payloadObject[key] === 'string') {
    configLabelKey = payloadObject[key] as string
  } else if (nested && key in nested && typeof nested[key] === 'string') {
    configLabelKey = nested[key] as string
  }

  return configLabelKey in config ? config[configLabelKey] : config[key]
}

export { ChartContainer, ChartTooltip, ChartTooltipContent, ChartLegend, ChartLegendContent, useChart }
