import {defineCollection} from 'astro:content'
import {z} from 'astro/zod'
import {docsLoader} from '@astrojs/starlight/loaders'
import {docsSchema} from '@astrojs/starlight/schema'

const author = z.union([
  z.string(),
  z.object({name: z.string(), role: z.string().optional(), href: z.string().optional(), imageUrl: z.string().optional()}),
])

export const collections = {
  docs: defineCollection({
    loader: docsLoader(),
    schema: docsSchema({
      extend: z.object({
        // shown next to the page title
        titleBadge: z.object({text: z.string(), type: z.string().default('tip')}).optional(),
        // blog posts
        date: z.coerce.date().optional(),
        category: z.string().optional(),
        author: author.optional(),
        imageUrl: z.string().optional(),
        imageAlt: z.string().optional(),
        imageWidth: z.number().optional(),
        imageHeight: z.number().optional(),
        // changelog
        latestBatch: z.string().optional(),
        subtext: z.string().optional(),
      }),
    }),
  }),
}
