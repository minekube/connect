import rss from '@astrojs/rss'
import type {APIRoute} from 'astro'
import {getCollection} from 'astro:content'
import {site} from '../plugins/llms'

const authorName = (a: unknown) => (typeof a === 'string' ? a : (a as {name?: string})?.name ?? 'Minekube')

export const GET: APIRoute = async () => {
  const posts = (await getCollection('docs', (e) => e.id.startsWith('blog/') && e.data.date))
    .sort((a, b) => +b.data.date! - +a.data.date!)
  return rss({
    title: 'The Minekube Blog',
    description: 'The official blog for the Minekube platform',
    site,
    customData: `<language>en</language><copyright>Copyright (c) ${new Date().getFullYear()}, Minekube and contributors</copyright><image><url>${site}/minekube-logo.png</url><title>The Minekube Blog</title><link>${site}</link></image>`,
    items: posts.map((p) => ({
      title: p.data.title,
      link: `${site}/${p.id}`,
      pubDate: p.data.date,
      description: p.data.description,
      content: p.rendered?.html,
      author: authorName(p.data.author),
      enclosure: p.data.imageUrl ? {url: new URL(p.data.imageUrl, site).href, length: 0, type: p.data.imageUrl.endsWith('.png') ? 'image/png' : 'image/jpeg'} : undefined,
    })),
  })
}
