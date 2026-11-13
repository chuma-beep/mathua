import { existsSync, readFileSync, writeFileSync } from 'node:fs'

// The service worker's cache name is its version. A static version means every
// deploy ships byte-identical sw.js: the browser sees no update, never installs
// a new worker, and the activate handler that purges old caches never runs. The
// Dockerfile stamps public/sw.js for the Fly image, and the Cloudflare Worker
// deploy uploads the postbuild-stamped out/sw.js, so mathua.com always gets a
// version that changes with the build.
//
// This runs as `postbuild`, after `next build` has copied public/ into out/.
// Only the exported copy is rewritten; the source file stays stampable by the
// Dockerfile without a dirty working tree.
const file = 'out/sw.js'
if (!existsSync(file)) {
  console.log('stamp-sw: no out/sw.js (not a static export); nothing to stamp')
  process.exit(0)
}

const stamp = process.env.SW_VERSION || `v2-${new Date().toISOString().replace(/[-:]/g, '').replace(/\..*/, 'Z')}`
const source = readFileSync(file, 'utf8')
if (!/const VERSION = '[^']*'/.test(source)) {
  console.error('stamp-sw: VERSION declaration missing from out/sw.js')
  process.exit(1)
}
writeFileSync(file, source.replace(/const VERSION = '[^']*'/, `const VERSION = '${stamp}'`))
console.log(`stamp-sw: ${file} -> ${stamp}`)
