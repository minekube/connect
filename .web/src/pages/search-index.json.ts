// Lightweight search index used by the search dialog when Pagefind is not available
// (the dev server, or if the Pagefind bundle fails to load). Pagefind stays the main engine.
import {getCollection} from 'astro:content'
import path from 'node:path'
import {expandIncludes, vitepressSlug} from '../plugins/markdown.mjs'

const skip = [/^guide\/includes\//, /^changelog\/20/]

export async function GET() {
  const docs = await getCollection('docs', (e) => !skip.some((r) => r.test(e.id)))
  const pages = docs.map((e) => {
    const body = expandIncludes(e.body ?? '', path.resolve(e.filePath ?? ''))
      .replace(/<Badge[^>]*>[\s\S]*?<\/Badge>/g, '')
      .replace(/^import .*$/gm, '')
    const headings = [...body.matchAll(/^#{2,3}\s+(.+)$/gm)].map((m) => {
      const t = m[1].replace(/<[^>]+>/g, '').replace(/`/g, '').replace(/\[([^\]]+)\]\([^)]*\)/g, '$1').trim()
      return {t, id: vitepressSlug(t)}
    })
    const text = body
      .replace(/```[\s\S]*?```/g, ' ')
      .replace(/<!--[\s\S]*?-->/g, ' ')
      .replace(/<[^>]+>/g, ' ')
      .replace(/\[([^\]]+)\]\([^)]*\)/g, '$1')
      .replace(/[#*_>`|]+/g, ' ')
      .replace(/:::\w*(\[[^\]]*\])?/g, ' ')
      .replace(/\s+/g, ' ')
      .trim()
      .slice(0, 6000)
    const isIndex = /[\\/]index\.mdx?$/.test(e.filePath ?? '')
    const url = e.id === '' || e.id === 'index' ? '/' : isIndex ? `/${e.id}/` : `/${e.id}`
    return {url, title: e.data.title, desc: e.data.description ?? '', headings, text}
  })
  return new Response(JSON.stringify(pages), {headers: {'Content-Type': 'application/json; charset=utf-8'}})
}
