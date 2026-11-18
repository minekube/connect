import type {APIRoute, GetStaticPaths} from 'astro'
import {llmsPages, toMarkdown} from '../plugins/llms'

export const getStaticPaths: GetStaticPaths = async () =>
  (await llmsPages()).map((entry) => ({params: {slug: entry.id}, props: {entry}}))

export const GET: APIRoute = ({props}) =>
  new Response(toMarkdown(props.entry), {headers: {'content-type': 'text/markdown; charset=utf-8'}})
