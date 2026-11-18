# Documentation Website

This website is built with [Astro Starlight](https://starlight.astro.build/),
a documentation framework on top of [Astro](https://astro.build/).

## Setup

Use Node.js 24 and the repository's pinned Yarn release. With [mise](https://mise.jdx.dev/):

```sh
$ mise exec node@24 -- corepack yarn install
```

### Local Development

```sh
$ mise exec node@24 -- corepack yarn dev
```

This starts a local development server at http://localhost:4321. Most changes
are reflected live without having to restart the server.

### Where things live

- `src/content/docs/` - every page: `guide/`, `blog/`, `changelog/`, plus the
  home page (`index.mdx`), `plans.mdx` and `team.mdx`. Paths map 1:1 to URLs.
- `src/components/` - Starlight overrides (`Header`, `Head`, `PageTitle`,
  `MarkdownContent`, `Footer`), the home page, changelog shell, blog list and
  the Vue components reused from the old site (`vue/`).
- `src/plugins/markdown.mjs` - the markdown features kept from the VitePress
  site: `<!--@include: ./file.md-->`, `:::code-group`, `:::details[Title]`,
  `<Badge type="tip|note|caution|danger">` and VitePress-compatible heading ids.
- `src/pages/` - `feed.rss`, `changelog.rss`, `llms.txt`, `llms-full.txt` and
  the per-page `.md` routes.
- `public/` - static files served as-is (`_headers`, `_redirects`, images).

Asides use Starlight syntax: `:::note`, `:::tip`, `:::caution`, `:::danger`,
with an optional title in brackets (`:::tip[Title]`).

### Build

```sh
$ mise exec node@24 -- corepack yarn check:docs
$ mise exec node@24 -- corepack yarn build
$ mise exec node@24 -- corepack yarn check:changelog-layout
```

The build generates the site into `dist/` with a Pagefind search index,
`/llms.txt`, `/llms-full.txt`, and a Markdown `.md` route for every included
documentation page. It then validates that the LLM index only links to
generated files. The changelog layout check uses an installed Chrome browser to
verify Cloudflare-style entry-by-entry rail continuity, exact rail/node/date
geometry, sticky entry-date boundaries, product filtering, and the rail-free
mobile layout against the built site.

### Deployment

Our docs are deployed as [Cloudflare Workers Static
Assets](https://developers.cloudflare.com/workers/static-assets/). Build the
site, deploy and validate the isolated `connect-docs-canary.minekube.com`
canary, then deploy the production
Worker only when the custom hostname is authorized to move:

```sh
$ yarn build
$ yarn deploy:worker:canary
# validate https://connect-docs-canary.minekube.com
$ yarn deploy:worker
```
