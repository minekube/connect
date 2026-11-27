import rss from '@astrojs/rss'
import type {APIRoute} from 'astro'
import {getCollection} from 'astro:content'
import {site} from '../plugins/llms'

// Feed readers drop unknown markup, which would silently lose the product and
// status markers, so badges are flattened to <strong> in the feed only.
const flatten = (html = '') => html
  .replace(/<!--[\s\S]*?-->/g, '')
  .replace(/<span class="mk-badge[^"]*">([\s\S]*?)<\/span>/g, '<strong>$1</strong>')

export const GET: APIRoute = async () => {
  const batches = (await getCollection('docs', (e) => /^changelog\/20/.test(e.id)))
    .sort((a, b) => +b.data.date! - +a.data.date!)
  return rss({
    title: 'Minekube Changelog',
    // Keep in sync with the scope wording in changelog/index.md.
    description: 'User-visible changes across the Minekube platform. Complete coverage for the Connect plugin, GeyserLite and Craftless since June 4, 2026; selective for Gate and the hosted Connect service.',
    site: `${site}/changelog/`,
    customData: `<language>en</language>`,
    items: batches.map((b) => ({
      title: b.data.title,
      link: `${site}/${b.id}`,
      pubDate: b.data.date,
      description: b.data.description,
      content: flatten(b.rendered?.html),
      author: 'Minekube',
    })),
  })
}
