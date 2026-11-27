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
  home page (`index.mdx`), `plans.mdx`, `team.mdx` and `credits.md`. Paths map
  1:1 to URLs.
- `src/components/` - Starlight overrides (`Header`, `Head`, `PageTitle`,
  `MarkdownContent`, `Footer`, `Search`, `ThemeSwitch`), the home page
  (`Home.astro`), changelog shell, blog list, team and the Vue components for
  plans, positions and the globe (`vue/`).
- `src/plugins/markdown.mjs` - the markdown features kept from the VitePress
  site: `<!--@include: ./file.md-->`, `:::code-group`, `:::details[Title]`,
  `<Badge type="tip|note|caution|danger">` and VitePress-compatible heading ids.
- `src/pages/` - `feed.rss`, `changelog.rss`, `llms.txt`, `llms-full.txt`, the
  per-page `.md` routes and `search-index.json` (search fallback, see below).
- `public/` - static files served as-is (`_headers`, `_redirects`, images).

Search is a custom dialog (`src/components/Search.astro`, Ctrl/⌘ K or `/`) on
top of Pagefind, which `astro build` indexes. When Pagefind is not available
(the dev server), it searches `/search-index.json` instead, so search works in
`yarn dev` too.

Asides use Starlight syntax: `:::note`, `:::tip`, `:::caution`, `:::danger`,
with an optional title in brackets (`:::tip[Title]`).

### Docs structure

The sidebar (`astro.config.mjs`) follows the order a reader needs things, the same split
Cloudflare Docs and VitePress use:

- **Get started** - overview, quick start, downloads, joining. Enough to go from zero to a
  running server.
- **Connectors** - one page per Connector (Gate, Java plugin) plus an overview to choose.
- **Concepts** - how Connect works: tunnels, edge locations, protection. Explains, no steps.
- **Guides** - one task per page, titled by what the reader wants to do ("Use a custom
  domain"). Steps first, background after.
- **Reference** - lookup material: compatibility matrix, integration and API details.
- **About Connect** - background and plans (why, use cases, roadmap).

Page titles are short and in sentence case. URLs never change when a page moves group; if
a page must be renamed, keep the old URL working.

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
