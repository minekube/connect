import type {APIRoute} from 'astro'
import {llmsPages, mdPath, site, toMarkdown} from '../plugins/llms'

export const GET: APIRoute = async () => {
  const pages = await llmsPages()
  const parts = pages.map((e) => `---\nurl: ${site}${mdPath(e.id)}\n---\n\n${toMarkdown(e)}`)
  return new Response(parts.join('\n'), {headers: {'content-type': 'text/plain; charset=utf-8'}})
}
