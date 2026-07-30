# Archived developers hostname

`developers.minekube.com` is an archived hostname retained only for compatibility with external links that cannot be updated.
The `redirect-minekube-developers` Worker preserves its existing catch-all and gives known immovable download links real HTTP redirects.

The Connect Spigot download points at the mutable `latest` GitHub release asset maintained by the Connect Java release workflow.
The hostname is only a compatibility hop and is not the artifact authority.

Run the local behavior test:

```sh
mise exec node@24 -- node --test .web/workers/developers-redirect/worker.test.mjs
```

After merging a reviewed change, deploy the existing Worker:

```sh
mise exec node@24 -- npx --yes wrangler@4.115.0 deploy --config .web/workers/developers-redirect/wrangler.jsonc
```

Verify the live download:

```sh
mise exec node@24 -- node .web/scripts/check-plugin-download.mjs
```
