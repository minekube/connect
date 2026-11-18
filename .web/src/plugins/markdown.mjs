// Markdown features carried over from the VitePress site.
import fs from 'node:fs'
import path from 'node:path'

const includeRe = /^<!--\s*@include:\s*(.+?)\s*-->$/

/** Reads an included markdown file without its frontmatter, expanding nested includes. */
export function readInclude(fromFile, rel) {
  const file = path.resolve(path.dirname(fromFile), rel)
  const src = fs.readFileSync(file, 'utf8').replace(/\r\n/g, '\n').replace(/^---\n[\s\S]*?\n---\n/, '')
  return {file, src}
}

/** Text-level include expansion, used by the plain-markdown (.md / llms) routes. */
export function expandIncludes(src, fromFile) {
  return src.replace(/^<!--\s*@include:\s*(.+?)\s*-->$/gm, (_, rel) => {
    const inc = readInclude(fromFile, rel)
    return expandIncludes(inc.src, inc.file)
  })
}

/**
 * `<!--@include: ./file.md-->` splices another markdown file in place, and
 * `:::code-group` / `:::details[Title]` containers render as tabs / <details>.
 */
export function remarkCompat() {
  const processor = this
  return (tree, vfile) => {
    expand(tree, vfile.path)
    vitepressHeadingIds(tree)

    function expand(node, fromFile) {
      if (!node.children) return
      for (let i = 0; i < node.children.length; i++) {
        const child = node.children[i]
        const m = child.type === 'html' && child.value.trim().match(includeRe)
        if (m) {
          const inc = readInclude(fromFile, m[1])
          const sub = processor.parse(inc.src)
          expand(sub, inc.file)
          node.children.splice(i, 1, ...sub.children)
          i += sub.children.length - 1
          continue
        }
        if (child.type === 'html') child.value = badgeHtml(child.value)
        if (child.type === 'containerDirective' && child.name === 'code-group') codeGroup(child)
        else if (child.type === 'containerDirective' && child.name === 'details') details(child)
        expand(child, fromFile)
      }
    }
  }
}

let groupId = 0

function codeGroup(node) {
  const id = `cg-${++groupId}`
  const codes = node.children.filter(c => c.type === 'code')
  const labels = codes.map((c, i) => {
    const t = c.meta?.match(/title="([^"]*)"/)
    if (t) c.meta = c.meta.replace(t[0], '').trim()
    return t ? t[1] : c.lang || `Tab ${i + 1}`
  })
  const esc = s => s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/"/g, '&quot;')
  node.data = {hName: 'div', hProperties: {className: ['code-group'], 'data-code-group': id}}
  node.children = [
    {
      type: 'html',
      value: `<div class="code-group-tabs" role="tablist">${labels.map((l, i) =>
        `<button type="button" role="tab" aria-selected="${i === 0}" data-tab="${i}">${esc(l)}</button>`).join('')}</div>`,
    },
    ...codes.map((c, i) => ({
      type: 'codeGroupPanel',
      data: {hName: 'div', hProperties: {className: ['code-group-panel'], role: 'tabpanel', hidden: i !== 0, 'data-panel': i}},
      children: [c],
    })),
  ]
}

function details(node) {
  const label = node.children[0]?.data?.directiveLabel ? node.children.shift() : null
  node.data = {hName: 'details', hProperties: {className: ['details-block']}}
  node.children.unshift({
    type: 'detailsSummary',
    data: {hName: 'summary'},
    children: label ? label.children : [{type: 'text', value: 'Details'}],
  })
}

/** `<Badge type="tip">Text</Badge>` raw HTML (open and close tags arrive as separate nodes) -> styled span. */
export function badgeHtml(html) {
  return html
    .replace(/<Badge(?:\s+type="(\w+)")?\s*>/g, (_, t) => `<span class="mk-badge mk-badge-${t || 'tip'}">`)
    .replace(/<\/Badge>/g, '</span>')
}

// VitePress heading slugs (@mdit-vue/shared slugify), so every existing
// /page#anchor link keeps working. Astro keeps ids that are already set.
const rSpecial = /[\s~`!@#$%^&*()\-_+=[\]{}|\\;:"'“”‘’<>,.?/]+/g
export const vitepressSlug = (s) => s.normalize('NFKD').replace(/[̀-ͯ]/g, '').replace(/[\u0000-\u001f]/g, '')
  .replace(rSpecial, '-').replace(/-{2,}/g, '-').replace(/^-+|-+$/g, '').replace(/^(\d)/, '_$1').toLowerCase()

/** Sets heading ids in mdast, ahead of the rehype slugger Starlight runs. */
function vitepressHeadingIds(tree) {
  const seen = new Map()
  // VitePress slugged text and inline code, never inline HTML (badges).
  const text = (n) => n.type === 'text' || n.type === 'inlineCode' ? n.value : n.type === 'html' ? '' : (n.children ?? []).map(text).join('')
  ;(function visit(n) {
    if (n.type === 'heading') {
      const base = vitepressSlug(text(n).trim())
      const count = seen.get(base) ?? 0
      seen.set(base, count + 1)
      n.data = {...n.data, hProperties: {...n.data?.hProperties, id: count ? `${base}-${count}` : base}}
    }
    n.children?.forEach(visit)
  })(tree)
}
