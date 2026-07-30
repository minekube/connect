# Archived developers hostname

`developers.minekube.com` is an archived hostname retained only for compatibility with external links that cannot be updated.
The `redirect-minekube-developers` Worker preserves its existing catch-all and gives known immovable download links real HTTP redirects.

The Connect Spigot download points at the mutable `latest` GitHub release asset maintained by the Connect Java release workflow.
The hostname is only a compatibility hop and is not the artifact authority.

Run the local behavior test:

```sh
mise exec node@24 -- node --test worker.test.mjs
```

After merging a reviewed change, deploy the existing Worker:

```sh
mise exec node@24 -- npx --yes wrangler@4.115.0 deploy --config wrangler.jsonc
```

Verify the live download from `.web`:

```sh
mise exec node@24 -- node scripts/check-plugin-download.mjs
```
