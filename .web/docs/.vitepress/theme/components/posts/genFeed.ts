import path from 'path'
import { writeFileSync } from 'fs'
import { Feed } from 'feed'
import { createContentLoader, type SiteConfig } from 'vitepress'
import {Post} from "./posts.data";

const baseUrl = 'https://connect.minekube.com'

export async function genFeed(config: SiteConfig) {
  const feed = new Feed({
    title: 'The Minekube Blog',
    description: 'The official blog for the Minekube platform',
    id: baseUrl,
    link: baseUrl,
    language: 'en',
    image: `${baseUrl}/minekube-logo.png`,
    favicon: `${baseUrl}/favicon.png`,
    copyright:
      `Copyright (c) ${new Date().getFullYear()}, Minekube and contributors`
  })

  const posts = (await createContentLoader('blog/*.md', {
      excerpt: true,
      render: true,
      transform(raw): Post[] {
        return raw.filter(({ url }) => !url.endsWith('/')) // Exclude 'index.md'
      }
  }).load())

  posts.sort(
    (a, b) =>
      +new Date(b.frontmatter.date as string) -
      +new Date(a.frontmatter.date as string)
  )

  for (const { url, excerpt, frontmatter, html } of posts) {
    const author = normalizeAuthor(frontmatter.author)

    feed.addItem({
      title: frontmatter.title,
      id: `${baseUrl}${url}`,
      link: `${baseUrl}${url}`,
      description: frontmatter.description ?? excerpt,
      content: html?.replaceAll('&ZeroWidthSpace;', ''),
      image: frontmatter.imageUrl ? new URL(frontmatter.imageUrl, baseUrl).href : undefined,
      author: [
        {
          name: author.name,
          link: author.href
        }
      ],
      date: frontmatter.date
    })
  }

  writeFileSync(path.join(config.outDir, 'feed.rss'), feed.rss2())
}

function normalizeAuthor(author: unknown) {
  if (author && typeof author === 'object' && 'name' in author) {
    const data = author as { name?: unknown, href?: unknown }
    return {
      name: typeof data.name === 'string' ? data.name : 'Minekube',
      href: typeof data.href === 'string' ? data.href : undefined
    }
  }

  return {
    name: typeof author === 'string' ? author : 'Minekube',
    href: undefined
  }
}
