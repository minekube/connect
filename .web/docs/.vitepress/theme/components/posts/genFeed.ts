import path from 'path'
import { writeFileSync } from 'fs'
import { Feed } from 'feed'
import { createContentLoader, type SiteConfig } from 'vitepress'
import {Post} from "./posts.data";
import {ogUrl} from "../../../config";

const baseUrl = ogUrl

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
      'Copyright (c) 2021-present, Yuxi (Evan) You and blog contributors'
  })

  const posts = (await createContentLoader('blog/*.md', {
      excerpt: true,
      render: true,
      transform(raw): Post[] {
        return raw.filter(({ url }) => !url.endsWith('/')) // Exclude 'index.md'
      }
  }).load())
  console.log(posts)

  posts.sort(
    (a, b) =>
      +new Date(b.frontmatter.date as string) -
      +new Date(a.frontmatter.date as string)
  )

  for (const { url, excerpt, frontmatter, html } of posts) {
    console.log(html)
    feed.addItem({
      title: frontmatter.title,
      id: `${baseUrl}${url}`,
      link: `${baseUrl}${url}`,
      description: excerpt,
      content: html?.replaceAll('&ZeroWidthSpace;', ''),
      author: [
        {
          name: frontmatter.author,
          link: frontmatter.twitter
            ? `https://twitter.com/${frontmatter.twitter}`
            : undefined
        }
      ],
      date: frontmatter.date
    })
  }

  writeFileSync(path.join(config.outDir, 'feed.rss'), feed.rss2())
}
