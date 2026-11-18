import {defineRouteMiddleware} from '@astrojs/starlight/route-data'

// Blog posts are articles, not docs: no docs sidebar, keep the "On this page" outline.
export const onRequest = defineRouteMiddleware((context) => {
  const route = context.locals.starlightRoute
  if (route.id.startsWith('blog/') && route.id !== 'blog/index') route.hasSidebar = false
})
