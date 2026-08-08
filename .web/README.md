# Documentation Website

This website is built using [VitePress](https://vitepress.dev/),
a modern static website generator for documentation.

## Setup

Use Node.js 24 and the repository's pinned Yarn release. With [mise](https://mise.jdx.dev/):

```sh
$ mise exec node@24 -- corepack yarn install
```

### Installation

Finally, you will need to install the Node.js dependencies for this project
using the pinned package manager:

```sh
$ mise exec node@24 -- corepack yarn install
```

### Local Development

```sh
$ mise exec node@24 -- corepack yarn dev
```

This command starts a local development server and opens up a browser window.
Most changes are reflected live without having to restart the server.

### Build

```sh
$ mise exec node@24 -- corepack yarn check:docs
$ mise exec node@24 -- corepack yarn build
$ mise exec node@24 -- corepack yarn check:changelog-layout
```

The build generates the VitePress site, `/llms.txt`, `/llms-full.txt`, and a Markdown `.md` route for every included
documentation page. It then validates that the LLM index only links to generated files. The changelog layout check uses
an installed Chrome browser to verify rail continuity, rail/node alignment, sticky-date boundaries, product filtering,
and the mobile layout against the built site.

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
