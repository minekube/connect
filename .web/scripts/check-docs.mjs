import {existsSync, readFileSync} from 'node:fs'
import {resolve} from 'node:path'

const root = new URL('..', import.meta.url)

function readDoc(path) {
  return readFileSync(resolve(root.pathname, path), 'utf8')
}

function assertIncludes(content, expected, file) {
  if (!content.includes(expected)) {
    throw new Error(`${file} is missing required docs coverage: ${expected}`)
  }
}

function assertAll(file, required) {
  const content = readDoc(file)

  for (const expected of required) {
    assertIncludes(content, expected, file)
  }
}

function assertNotIncludes(content, unexpected, file) {
  if (content.includes(unexpected)) {
    throw new Error(`${file} contains retired guidance: ${unexpected}`)
  }
}

function assertMissing(path) {
  if (existsSync(resolve(root.pathname, path))) {
    throw new Error(`${path} must be removed because the public Developer API is unavailable`)
  }
}

assertAll('docs/guide/bedrock.md', [
  'Connect-managed Bedrock',
  'self-hosted Gate direct Bedrock',
  'backend Floodgate API',
  'online-mode',
  'custom domain',
  'Gate Lite',
  'Paper/Velocity/Bungee connector',
  'Do not enable backend Gate `bedrock: true` for Connect-managed Bedrock',
  'use the normal Bedrock port `19132`',
  'It does not mean you need to open UDP `19132`',
  'official Microsoft/Xbox Bedrock authentication',
  'Bedrock Identity Enforcement',
  'metadata-url',
  'endpoint and organization',
  'Discord support response draft',
  'without owning or linking Java Edition',
  'stable profile derived from the verified Bedrock XUID',
  'Without valid Microsoft/Xbox authentication',
])

assertAll('docs/guide/offline-mode.md', [
  'Connect-managed Bedrock identity',
  'official Microsoft/Xbox Bedrock auth',
  'Do not enable either offline-mode option just because a Bedrock player is joining through Connect',
  'cracked.minekube.net',
  'cracked.minekube.com',
  'allowOfflineModePlayers: true',
  'Connect offline mode applies to Java Edition',
])

assertNotIncludes(
  readDoc('docs/guide/offline-mode.md'),
  '`offline.minekube.net`',
  'docs/guide/offline-mode.md',
)

assertAll('docs/guide/joining.md', [
  '## Who Can Join',
  'Stable native Bedrock/XUID-derived profile',
  'Bedrock client without valid Microsoft/Xbox authentication',
  'What the Server Owner Configures',
])

assertAll('docs/guide/includes/joining.md', [
  '`25565` (normally omitted)',
  '`19132`',
  'cracked.minekube.net',
  'cracked.minekube.com',
])

assertAll('docs/guide/quick-start.md', [
  'Microsoft/Xbox-authenticated Bedrock players can join with or without a linked Java account',
  'Offline/cracked Java players require the endpoint owner to opt in',
  'Bedrock players without valid Microsoft/Xbox authentication cannot join',
])

assertAll('docs/guide/connectors/plugin.md', [
  'Bedrock Identity',
  'metadata-url',
  'enforcement: warn',
])

assertAll('docs/guide/compatibility.md', [
  'Velocity snapshots',
  'NLogin',
  'AuthMe',
  'FastLogin',
  'Gate Lite with an online-mode backend',
  'AuthSession/passthrough for Lite backend routes is not available today',
  'Arclight/Ketting/Forge hybrids',
  'FabricProxy-Lite',
  'CrossStitch',
  'Polymer',
  'NeoForge 1.21.x / Proxy-Compatible-Forge through Connect',
  'Gate v0.69.1 fixed the open-bundle transition',
  'connect-java/issues/141',
  'Connect does not currently tunnel the separate UDP transport',
  'connect/issues/159',
])

assertAll('docs/guide/login-plugins.md', [
  'Connect v0.13.1 or newer',
  'premium autologin stays enabled',
  'login-reassert.enabled',
  'login-reassert.restore-full-profile',
  'new-uuid-creator: MOJANG',
  'On Velocity, Connect also',
  "restores the player's skin properties by default",
  'login-reassert:\n  enabled: true\n  restore-full-profile: false',
  'During pre-login, BungeeCord exposes no\nprofile-properties API on the pending connection, so `restore-full-profile` can re-assert the Mojang UUID and username\nbut cannot restore skin properties through this setting.',
  'protect the documented LibreLogin path by default; arbitrary-plugin guarantees depend on proxy ordering\nsupport',
  'Connect v0.13.1+ guarantees the LibreLogin path described above.',
  'strict after-all protection on\nVelocity requires the numeric-priority API;',
  'on legacy Velocity, Connect uses `PostOrder.LAST`',
  'another `LAST` handler\ncan still run after it depending on plugin load order.',
])

assertNotIncludes(
  readDoc('docs/guide/login-plugins.md'),
  'Disable premium autologin in LibreLogin',
  'docs/guide/login-plugins.md',
)

assertAll('docs/guide/compatibility.md', [
  'Use Connect v0.13.1 or newer and keep premium autologin enabled',
])

assertNotIncludes(
  readDoc('docs/guide/compatibility.md'),
  'Known incompatible',
  'docs/guide/compatibility.md',
)

assertAll('docs/guide/connectors/gate.md', [
  'Current behavior',
  'Not supported today',
  'Connect -> Gate Lite -> Online Mode Backend',
  'Connect passthrough/AuthSession support',
])

assertAll('docs/guide/auth-api.md', [
  'The AuthSession API is not available in production yet',
  'Gate Lite behind Connect is different',
  'Use standard Gate with Connect enabled or the Connect Java Plugin',
])

assertAll('docs/changelog/index.md', [
  'layout: page',
  'title: Changelog',
  'latestBatch: August 9 – August 24, 2026',
  '<ChangelogLanding>',
  '<details class="changelog-policy">',
  // The raw anchor is load-bearing: a markdown link to /changelog.rss fails the
  // VitePress dead-link check, and a plain href is rewritten by the client router.
  '<a href="/changelog.rss" target="_blank" rel="noreferrer">RSS</a>',
  // The entry policy is owned by the comms home and reproduced verbatim. It is
  // not a link-required rule: the next clause states the internal-repository
  // exception, and the page ships unlinked hosted entries under it. Do not
  // reword either half without that owner.
  'Entries link to the public release that contains the change. Parts of the platform, including the hosted Connect service, are developed in internal repositories; those entries carry a date and a description and are marked *internal repository, no public link*.',
  // The scope statement is testable on purpose: it names where coverage is
  // complete and where it is selective, instead of claiming completeness the
  // page cannot back. Do not restore a blanket "every change" claim.
  'For the Connect plugin, GeyserLite and Craftless it is complete from June 4, 2026',
  'For Gate and the hosted Connect service it is selective',
  '<!--@include: ./2026-08-24.md-->',
  '<!--@include: ./2026-08-08.md-->',
  '<!--@include: ./2026-07-27.md-->',
])

assertAll('docs/.vitepress/theme/index.ts', [
  'import ChangelogLanding from "./components/changelog/ChangelogLanding.vue"',
  "app.component('ChangelogLanding', ChangelogLanding)",
])

assertAll('docs/.vitepress/theme/components/changelog/ChangelogLanding.vue', [
  'Filter changelog by product',
  'Subscribe to RSS',
  'View LLM docs',
  "const products = ['All products', 'Connect', 'Connect plugin', 'Gate', 'GeyserLite', 'Craftless']",
  "group.className = 'changelog-date-group'",
  "entry.classList.add('changelog-entry')",
  "visibleEntries.at(-1)?.setAttribute('data-rail-end', '')",
  'changelog-entry-node',
  'changelog-backdrop-dots',
  'position: sticky',
])

assertAll('docs/index.md', [
  'title: Product Changelog',
  'link: /changelog/',
])

assertAll('docs/changelog/2026-08-08.md', [
  'date: 2026-08-08',
  'Every Connect endpoint now accepts Microsoft/Xbox-authenticated Bedrock players without requiring a linked Java account by default',
  'connect-java/releases/tag/0.15.3',
  'gate/releases/tag/v0.71.1',
  'geyserlite/releases/tag/v0.5.2',
])

assertAll('docs/changelog/2026-08-24.md', [
  'date: 2026-08-24',
  // Every product in the batch links a release that has downloadable assets.
  'connect-java/releases/tag/0.15.9',
  'gate/releases/tag/v0.71.3',
  'geyserlite/releases/tag/v0.5.21',
  'craftless/releases/tag/v0.3.7',
])

assertAll('docs/changelog/2026-07-27.md', [
  'date: 2026-07-27',
  // Entries with no public citation say so in the slot a release link occupies.
  '(Internal repository, no public link.)',
])

assertAll('docs/.vitepress/theme/components/posts/genFeed.ts', [
  "createContentLoader('changelog/*.md'",
  "'changelog.rss'",
  // Feed readers strip unknown elements, which would drop the badge marker.
  'function stripComponents',
  // A subscriber never sees the page's scope paragraph - the channel
  // description sets their expectation once and permanently, so it has to
  // carry the same scope statement. Page and feed ship together or neither.
  "description: 'User-visible changes across the Minekube platform. Complete coverage for the Connect plugin, GeyserLite and Craftless since June 4, 2026; selective for Gate and the hosted Connect service.'",
])

assertAll('docs/public/_redirects', [
  '/guide/changelog /changelog/ 301',
])

const vitepressConfig = readDoc('docs/.vitepress/config.ts')
assertNotIncludes(vitepressConfig, "text: 'Developers API'", 'docs/.vitepress/config.ts')
assertNotIncludes(vitepressConfig, "link: '/guide/api/", 'docs/.vitepress/config.ts')
assertAll('docs/.vitepress/config.ts', [
  "import llmstxt from 'vitepress-plugin-llms'",
  'domain: ogUrl',
  "'guide/includes/*'",
  "'changelog/20*.md'",
])

assertAll('docs/public/_headers', [
  'rel="llms-txt"',
  'rel="llms-full-txt"',
  'Content-Type: text/markdown; charset=utf-8',
])

assertAll('docs/public/_redirects', [
  '/.well-known/llms.txt /llms.txt 301',
  '/.well-known/llms-full.txt /llms-full.txt 301',
])

assertAll('package.json', [
  'vitepress-plugin-llms',
  'node scripts/check-llms.mjs',
])

for (const path of [
  'docs/guide/api/index.md',
  'docs/guide/api/clients.md',
  'docs/guide/api/authentication.md',
  'docs/guide/api/super-endpoints.md',
  'docs/guide/api/examples.md',
  'docs/guide/api/javaexample/SamplePlugin.java',
  'docs/guide/api/goexample/example_test.go',
  'docs/guide/api/goexample/go.mod',
  'docs/guide/api/goexample/go.sum',
]) {
  assertMissing(path)
}

for (const path of ['docs/guide/index.md', 'docs/guide/adoption-plan.md']) {
  assertNotIncludes(readDoc(path), '/guide/api/', path)
}

console.log('Docs content assertions passed.')
