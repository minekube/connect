// Starlight renders every page through one [...slug] route, so with
// build.format 'preserve' guide/connectors/index.md lands in guide/connectors.html.
// Move those back to guide/connectors/index.html, the layout VitePress produced
// and the one Cloudflare serves at /guide/connectors/.
import fs from 'node:fs'
import path from 'node:path'
import {fileURLToPath} from 'node:url'

export function indexPages() {
  return {
    name: 'index-pages',
    hooks: {
      'astro:build:done': ({dir}) => {
        const out = fileURLToPath(dir)
        const docs = fileURLToPath(new URL('../content/docs/', import.meta.url))
        const walk = (d) => fs.readdirSync(d, {withFileTypes: true}).flatMap((e) =>
          e.isDirectory() ? [path.join(d, e.name), ...walk(path.join(d, e.name))] : [])
        for (const srcDir of walk(docs)) {
          if (!fs.readdirSync(srcDir).some((f) => /^index\.mdx?$/.test(f))) continue
          const rel = path.relative(docs, srcDir)
          const flat = path.join(out, `${rel}.html`)
          if (!fs.existsSync(flat)) throw new Error(`index-pages: expected ${flat}`)
          fs.mkdirSync(path.join(out, rel), {recursive: true})
          fs.renameSync(flat, path.join(out, rel, 'index.html'))
        }
      },
    },
  }
}
