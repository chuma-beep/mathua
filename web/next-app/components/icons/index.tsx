'use client'

/**
 * The app's icon surface. Everything imports from here; nothing imports the icon libraries
 * directly.
 *
 * Two reasons this is one module rather than 20 imports at the call sites.
 *
 * **The box is guaranteed.** The animated icons render a `<div>` sized by a `size` prop, not an
 * `<svg>` sized by a `size-4` class. The collapsed 44px sidebar rail is only even because every
 * icon in it occupies the same box, and a Tailwind class that silently does nothing to a div is
 * how that rail would go ragged. Size is set here so it cannot be forgotten.
 *
 * **The colour is guaranteed.** The animated icons take a literal `color`, and the default is
 * not `currentColor`. Left alone they would render in whatever the library's default is, which
 * does not follow light and dark mode and is not a theme token — the exact failure ADR-036 was
 * written about. Every wrapper passes `currentColor`, so `text-mathua-*` on the parent does what
 * it says.
 *
 * Motion is on by default (`isAnimated` defaults to true in the library), so icons animate when
 * they mount and again on hover. That was a deliberate choice for this rollout.
 */

import { forwardRef } from 'react'
import {
  BookOpenIcon,
  ChartNoAxesColumnIcon,
  ChevronsLeftIcon,
  ChevronsRightIcon,
  CompassIcon as CompassAnimatedIcon,
  Dice5Icon,
  FileTextIcon,
  HandshakeIcon,
  HouseIcon,
  LayersIcon,
  LogOutIcon,
  NetworkIcon,
  PauseIcon,
  PencilIcon,
  PlayIcon,
  SettingsIcon,
  ShieldIcon,
  SunMoonIcon,
  TrophyIcon,
  UserIcon,
  ChevronDownIcon,
  ChevronUpIcon,
} from '@animateicons/react/lucide'

/**
 * What a caller may pass. Deliberately no `color`: the wrapper owns it, because a caller who can
 * set a literal colour is one edit away from an icon that ignores the theme.
 */
type IconProps = {
  size?: number
  duration?: number
  className?: string
  title?: string
  /**
   * Defaults to true (motion on mount, and again on hover). Set false for an icon whose motion
   * is driven by the caller: the compass spins forward on menu open and reverse on close, and a
   * mount animation would both add noise and fight that.
   */
  isAnimated?: boolean
}

export type IconHandle = { startAnimation: () => void; stopAnimation: () => void }

/** What every export of this module is. Nav lists hold these rather than re-declaring a shape. */
export type IconComponent = React.ForwardRefExoticComponent<
  IconProps & React.RefAttributes<IconHandle>
>

/**
 * One wrapper for every icon, so size and colour cannot be forgotten at a call site.
 *
 * The ref is forwarded because two icons are driven imperatively rather than on hover, and the
 * library exposes startAnimation/stopAnimation through it.
 */
/** What the library component itself accepts, which is a superset of what callers may pass. */
type LibraryProps = IconProps & {
  color?: string
  role?: string
  'aria-label'?: string
  'aria-hidden'?: boolean
}
type AnyIcon = React.ComponentType<LibraryProps & React.RefAttributes<IconHandle>>

function wrap(Icon: AnyIcon) {
  return forwardRef<IconHandle, IconProps>(function AnimatedIcon(
    { size = 16, duration = 1, className, title, isAnimated },
    ref,
  ) {
    return (
      <Icon
        ref={ref}
        size={size}
        duration={duration}
        color="currentColor"
        className={className}
        isAnimated={isAnimated}
        aria-hidden={title ? undefined : true}
        role={title ? 'img' : undefined}
        aria-label={title}
      />
    )
  })
}

// Navigation and chrome.
export const BookOpen = wrap(BookOpenIcon as never)
export const ChartNoAxesColumn = wrap(ChartNoAxesColumnIcon as never)
export const ChevronsLeft = wrap(ChevronsLeftIcon as never)
export const ChevronsRight = wrap(ChevronsRightIcon as never)
export const Compass = wrap(CompassAnimatedIcon as never)
export const FileText = wrap(FileTextIcon as never)
export const Handshake = wrap(HandshakeIcon as never)
export const House = wrap(HouseIcon as never)
export const Layers = wrap(LayersIcon as never)
export const LogOut = wrap(LogOutIcon as never)
export const Network = wrap(NetworkIcon as never)
export const Pause = wrap(PauseIcon as never)
export const Pencil = wrap(PencilIcon as never)
export const Play = wrap(PlayIcon as never)
export const Settings = wrap(SettingsIcon as never)
// The administrative entry in the profile rail. A shield rather than a key or a gear: the gear
// already means Settings, and a key would imply this is how you get in rather than what is
// behind it.
export const Shield = wrap(ShieldIcon as never)
export const SunMoon = wrap(SunMoonIcon as never)
export const Trophy = wrap(TrophyIcon as never)
export const User = wrap(UserIcon as never)
export const ChevronDown = wrap(ChevronDownIcon as never)
export const ChevronUp = wrap(ChevronUpIcon as never)

// Settings' avatar randomiser. The app used lucide's `Dices` (two dice) and there is no animated
// `Dices` in the set — `Dice5` is a single die. This is the one icon whose shape changes, and it
// is called out here rather than left to be noticed.
export const Dices = wrap(Dice5Icon as never)

/** Every icon this module exports, so a guard can check the surface without parsing JSX. */
export const ICON_NAMES = [
  'BookOpen',
  'ChartNoAxesColumn',
  'ChevronLeft',
  'ChevronsLeft',
  'ChevronsRight',
  'Compass',
  'Dices',
  'FileText',
  'Handshake',
  'House',
  'Layers',
  'LogOut',
  'Network',
  'Pause',
  'Pencil',
  'Play',
  'Settings',
  'Shield',
  'SunMoon',
  'Trophy',
  'User',
  'ChevronDown',
  'ChevronUp',
] as const
