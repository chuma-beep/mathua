import fs from 'node:fs'
import { fromXml } from 'xast-util-from-xml'
let src = fs.readFileSync('scripts/ingest/ingest_pretext.mjs','utf8')
src = src.replace(/^await main\(\)\s*$/m, 'export { transform }')
fs.writeFileSync('scripts/ingest/_t.mjs', src)
const { transform } = await import('./_t.mjs')

const resp = await fetch('https://raw.githubusercontent.com/PCCMathSAC/orcca/edition/src/section-equations-and-inequalities-with-fractions.ptx')
const tree = fromXml(await resp.text())
transform(tree)
const section = tree.children.find(c => c.name === 'section')
function find(node, needle, path='section') {
  const s = JSON.stringify(node.type==='text' ? node.value : node.name ?? '')
  if (node.type === 'text' && node.value?.includes(needle)) { console.log('PATH:', path); return true }
  let hit = false
  for (const c of node.children ?? []) {
    if (find(c, needle, path + '>' + (c.name ?? '#text'))) { hit = true; break }
  }
  return hit
}
find(section, '\\amp')
fs.unlinkSync('scripts/ingest/_t.mjs')
