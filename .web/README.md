# Documentation Website

This website is built using [Vitepress](https://vitepress.vuejs.org/),
a modern static website generator for documentation.

## Setup

> You must have a recent version of Node.js (14+) installed.
> You may use [Volta](https://github.com/volta-cli/volta), a Node version manager,
> to install the latest version of Node and `yarn`.

```sh
$ curl https://get.volta.sh | bash
$ volta install node yarn
```

### Installation

Finally, you will need to install the Node.js dependencies for this project
using yarn or another package manager:

```sh
$ yarn install
```

### Local Development

```sh
$ yarn run dev/connect
```

This command starts a local development server and opens up a browser window.
Most changes are reflected live without having to restart the server.

### Build

```sh
$ yarn run build/connect
```

This command generates static content into the `dist` directory and can be served
using any static contents hosting service.

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
