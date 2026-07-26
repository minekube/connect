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

  await genChangelogFeed(config)
}

async function genChangelogFeed(config: SiteConfig) {
  const feed = new Feed({
    title: 'Minekube Changelog',
    description: 'Every user-visible change we ship across the Minekube platform',
    id: `${baseUrl}/changelog/`,
    link: `${baseUrl}/changelog/`,
    language: 'en',
    image: `${baseUrl}/minekube-logo.png`,
    favicon: `${baseUrl}/favicon.png`,
    copyright:
      `Copyright (c) ${new Date().getFullYear()}, Minekube and contributors`
  })

  const entries = await createContentLoader('changelog/*.md', {
    excerpt: true,
    render: true,
    transform(raw) {
      return raw.filter(({ url }) => !url.endsWith('/')) // Exclude 'index.md'
    }
  }).load()

  entries.sort(
    (a, b) =>
      +new Date(b.frontmatter.date as string) -
      +new Date(a.frontmatter.date as string)
  )

  for (const { url, excerpt, frontmatter, html } of entries) {
    feed.addItem({
      title: frontmatter.title,
      id: `${baseUrl}${url}`,
      link: `${baseUrl}${url}`,
      description: frontmatter.description ?? excerpt,
      content: stripComponents(html)?.replaceAll('&ZeroWidthSpace;', ''),
      author: [{ name: 'Minekube' }],
      date: new Date(frontmatter.date as string)
    })
  }

  writeFileSync(path.join(config.outDir, 'changelog.rss'), feed.rss2())
}

// createContentLoader renders markdown without the Vue runtime, so theme
// components such as <VPBadge> survive as literal tags. Feed readers strip
// unknown elements, which would silently drop the Self-hosted / Hosted marker,
// so flatten them to <strong> for the feed only.
function stripComponents(html?: string) {
  return html
    ?.replace(/<!--[\s\S]*?-->/g, '')
    .replace(
      /<VPBadge[^>]*>([\s\S]*?)<\/VPBadge>/g,
      (_match, text) => `<strong>${text}</strong>`
    )
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
