import {defineRouteMiddleware} from '@astrojs/starlight/route-data'

// Blog posts are articles, not docs: one centred reading column, no docs sidebar and no outline.
export const onRequest = defineRouteMiddleware((context) => {
  const route = context.locals.starlightRoute
  if (route.id.startsWith('blog/') && route.id !== 'blog/index') {
    route.hasSidebar = false
    route.toc = undefined
  }
})
