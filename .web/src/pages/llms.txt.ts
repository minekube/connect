import type {APIRoute} from 'astro'
import {llmsPages, mdPath, site} from '../plugins/llms'

export const GET: APIRoute = async () => {
  const pages = await llmsPages()
  const lines = pages.map((e) => `- [${e.data.title}](${site}${mdPath(e.id)})${e.data.description ? `: ${e.data.description.replace(/\s+/g, ' ')}` : ''}`)
  return new Response(
    `# Minekube Connect\n\n> The Ingress Tunnel for Minecraft Servers\n\n## Table of Contents\n\n${lines.join('\n')}\n`,
    {headers: {'content-type': 'text/plain; charset=utf-8'}},
  )
}
