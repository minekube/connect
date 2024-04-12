import { createContentLoader } from 'vitepress'

export interface Post {
  title: string
  url: string
  imageUrl: string
  date: {
    time: number
    string: string
  }
  excerpt: string | undefined,
  category: string,
  author: {
    name: string,
    role: string,
    href: string
    imageUrl: string
  }
}

declare const data: Post[]
export { data }

export default createContentLoader('blog/*.md', {
  excerpt: true,
  transform(raw): Post[] {
    return raw
      .filter(({ url }) => !url.endsWith('/')) // Exclude 'index.md'
      .map(({ url, frontmatter, excerpt }) => ({
        ...frontmatter,
        url,
        excerpt,
        date: formatDate(frontmatter.date),
      }))
      .sort((a, b) => b.date.time - a.date.time)
  }
})

function formatDate(raw: string): Post['date'] {
  const date = new Date(raw)
  date.setUTCHours(12)
  return {
    time: +date,
    string: date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  }
}