// MathLive renders with the KaTeX font family and fetches the .woff2 files from
// `<origin>/fonts` at runtime — it does not bundle them. With `output: 'export'`
// there is no server to serve them, so they are copied into `public/fonts` at
// build time and end up in the exported `out/` like any other static asset.
//
// Runs from `prebuild` and `predev`, because a dev server that 404s the fonts
// renders every fraction with a fallback face and the page looks broken in a way
// that is very hard to diagnose from the console.
//
// Idempotent and cheap: it compares against a stamp file and exits if the
// destination is already current.

import { cp, mkdir, readFile, writeFile, access } from 'node:fs/promises'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const here = dirname(fileURLToPath(import.meta.url))
const root = join(here, '..')
const src = join(root, 'node_modules', 'mathlive', 'fonts')
const dest = join(root, 'public', 'fonts')
const stamp = join(dest, '.mathlive-version')

async function main() {
  const version = JSON.parse(await readFile(join(root, 'node_modules', 'mathlive', 'package.json'), 'utf8')).version

  try {
    await access(dest)
    if ((await readFile(stamp, 'utf8')) === version) {
      console.log(`mathlive fonts: ${version} already in public/fonts`)
      return
    }
  } catch {
    // No stamp yet — fall through and copy.
  }

  await mkdir(dest, { recursive: true })
  await cp(src, dest, { recursive: true })
  await writeFile(stamp, version)
  console.log(`mathlive fonts: copied ${version} to public/fonts`)
}

main().catch(err => {
  console.error('mathlive fonts: copy failed —', err.message)
  // A missing font is a rendering regression, not a build failure, so warn loudly
  // and let the build continue rather than blocking a deploy on an asset copy.
  process.exitCode = 0
})