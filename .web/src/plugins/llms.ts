import {getCollection, type CollectionEntry} from 'astro:content'
import path from 'node:path'
import {expandIncludes} from './markdown.mjs'

export const site = 'https://connect.minekube.com'

// Same exclusions as the old vitepress-plugin-llms config: include fragments,
// dated changelog partials, and placeholder or duplicate pages. MDX pages that
// are only components (home, plans, team, blog index) carry no prose.
const ignore = [/^guide\/includes\//, /^changelog\/20/, /^guide\/downloads$/, /^guide\/games$/, /^(index|plans|team|blog)$/]

export async function llmsPages() {
  const docs = await getCollection('docs')
  return docs.filter((e) => e.id !== '' && !ignore.some((r) => r.test(e.id)))
    .sort((a, b) => a.id.localeCompare(b.id))
}

export const mdPath = (id: string) => `/${id}.md`

/** Plain markdown for LLM routes: includes expanded, components flattened. */
export function toMarkdown(e: CollectionEntry<'docs'>) {
  let body = expandIncludes(e.body ?? '', path.resolve(e.filePath ?? ''))
  body = body
    .replace(/^import .*$/gm, '')
    .replace(/^<[A-Z]\w*[^>]*\/>$/gm, '')
    .replace(/<Badge(?:\s+type="\w+")?\s*>([\s\S]*?)<\/Badge>/g, '**$1**')
    .replace(/<!--[\s\S]*?-->/g, '')
    .replace(/\n{3,}/g, '\n\n')
  const desc = e.data.description ? `\n> ${e.data.description}\n` : ''
  return `# ${e.data.title}\n${desc}\n${body.trim()}\n`
}
