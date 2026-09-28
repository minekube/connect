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
  // The backend identity shape is the support surface for "wrong UUID" and
  // "wrong username prefix" reports: both values are intentional, so the page
  // has to show the concrete forms and the no-action-required stance. The
  // Floodgate contrast is what explains why old UUID-keyed data does not carry
  // over on a first Connect join.
  'What a Bedrock player looks like at your backend',
  '`_<gamertag>`',
  'RFC 4122 version-5',
  'appears as `_icedRyan`',
  'No action is required for either value',
  'Do not rewrite it',
  "Floodgate's `00000000-0000-0000-XUID` key is a different identifier",
  'does not follow a player from a plain Geyser + Floodgate setup',
  // The dot-prefix rationale: Java's account alphabet is [A-Za-z0-9_], but a
  // proxy-supplied login name is validated by the backend server, and Paper's
  // own check accepts `.` while vanilla's is the narrower alphabet. The page
  // has to keep that distinction or it reintroduces the false "Paper rejects
  // the dot prefix" claim it was corrected for.
  'Connect normalizes the result into `[A-Za-z0-9_]`',
  "backend server's own validation",
  "Paper's check also accepts `.`",
  "vanilla's check is only the `[A-Za-z0-9_]` alphabet",
])

// Neither identifier is wrong or deprecated, and Connect does not offer a
// configurable prefix or a different UUID form. Restoring that framing is a
// contract change, not a docs edit.
assertNotIncludes(
  readDoc('docs/guide/bedrock.md'),
  'deprecated',
  'docs/guide/bedrock.md',
)

// The retired rationale claimed a Java-wide rule that Paper itself contradicts.
// Restoring it is a contract change, not a docs edit: the page must explain the
// backend-validated alphabet (Paper accepts `.`, vanilla is narrower).
assertNotIncludes(
  readDoc('docs/guide/bedrock.md'),
  'Java profile names only accept ASCII letters, digits, and underscores',
  'docs/guide/bedrock.md',
)

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
  // The matrix row is the entry point for the Bedrock identity shape, so the
  // link target has to keep pointing at the documented section.
  '[Stable native Bedrock/XUID-derived profile](/guide/bedrock#what-a-bedrock-player-looks-like-at-your-backend)',
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

// The "Endpoint Token" tip is the support surface for Watch auth failures: it must name the token
// file (not the config that holds the endpoint name) and the org-owned case, and it must not send
// users to fix their organization selection or keep the old "If your get" typo.
assertAll('docs/guide/quick-start.md', [
  'Resetting the token in the',
  'invalidates the previous token immediately',
  'never in the Connector config that holds the',
  '`CONNECT_TOKEN` environment variable',
  'a trailing space or newline is a different token',
  'the endpoint name belongs to an organization',
])

assertNotIncludes(
  readDoc('docs/guide/quick-start.md'),
  'If your get',
  'docs/guide/quick-start.md',
)

assertAll('docs/guide/connectors/plugin.md', [
  'Bedrock Identity',
  'metadata-url',
  'enforcement: warn',
  '## Endpoint Token',
  '`plugins/connect/token.json`',
])

assertNotIncludes(
  readDoc('docs/guide/connectors/plugin.md'),
  'Switch to the owning Minekube organization',
  'docs/guide/connectors/plugin.md',
)

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
  'connect-java/issues/140',
  'connect-java/pull/142',
  'PacketEvents-based Velocity plugins (including Sonar and some nLogin builds) | Fixed in Connect Java 0.15.5',
  'Connect does not currently tunnel the separate UDP transport',
  'connect/issues/159',
])

assertNotIncludes(
  readDoc('docs/guide/compatibility.md'),
  'Product incompatibility',
  'docs/guide/compatibility.md',
)

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
  'injects into the proxy\'s Netty pipeline after login',
  'such as PacketEvents-based plugins (Sonar, some nLogin builds)',
  'on an older install the action is an update to the current',
])

assertNotIncludes(
  readDoc('docs/guide/login-plugins.md'),
  'Disable premium autologin in LibreLogin',
  'docs/guide/login-plugins.md',
)

// The rule table only described the pre-login conflict. Two shapes support actually hits were
// missing: a plugin that takes the login *packet* over and runs its own handshake (the
// `login-reassert` floor cannot help - it restores a decision, it cannot answer a handshake), and
// the passthrough rule that makes the `connect-player` exemption unusable on an offline-mode
// endpoint. Pin both, plus the link to the published contract.
assertAll('docs/guide/login-plugins.md', [
  'Hooks the login **packet** and runs its own authentication handshake',
  'no live client side to complete it',
  'not answer a handshake another plugin is waiting on',
  'which is why the first row alone is not enough to clear an auth plugin',
  'That exemption only fires where Connect authenticated the session',
  'Connect deliberately leaves `connect-player` unset',
  'A support answer cannot',
  'a plain Paper/Spigot server',
  '`login-reassert` is implemented on Velocity and BungeeCord',
  '[Login Plugin Integration](/guide/login-plugin-integration)',
])

// The bare "never forces online mode -> Compatible by design" row invited the wrong conclusion
// for packet-level plugins, which never touch the proxy's online-mode setting. The row has to
// keep the packet clause; dropping it is a regression, not a wording preference.
assertNotIncludes(
  readDoc('docs/guide/login-plugins.md'),
  'never forces online mode | Compatible by design',
  'docs/guide/login-plugins.md',
)

// The stall is bounded and cause-blind, not silent. Measured live on Paper 26.3-49 with
// packetevents 2.14.0 + connect-spigot 0.15.15 (control/cancel pair differing only by a cancel
// flag): the client gets nothing for ~30 s and is then dropped by the backend's own login
// timeout, with a line that names no cause. "usually with no kick message and nothing in the
// server log" overstated it and must not come back; the condition that keeps the bound honest
// (it belongs to the server and exists only while the handshake reaches it, since Connect has
// no login deadline of its own) is pinned together with it.
assertAll('docs/guide/login-plugins.md', [
  "the player is eventually dropped by your server's own login timeout",
  '`Took too long to log in`',
  'a line that names no cause',
  'is easily mistaken for a client or network problem',
  'Connect has no login deadline of its own',
  'It therefore only exists while the',
  'a plugin that swallows the handshake as well as the login packet',
  'nothing times out',
])

assertNotIncludes(
  readDoc('docs/guide/login-plugins.md'),
  'nothing in the server log',
  'docs/guide/login-plugins.md',
)

// The stall is cause-blind only up to connector 0.15.15. connect-spigot 0.15.16 added
// `connect_login_stall_watchdog`, which sits upstream of the data handler that consumes
// LOGIN_START and emits ONE default-verbosity line ~10 s after the handshake when the login start
// never arrives. The public text has to say so, in the same place the cause-blind timeout is
// described, or operators and support keep telling users that nothing names the cause while the
// installed connector already named it. The sentence is version-gated on purpose (`0.15.16 and
// newer`) so it is not retrofit advice for older connectors, and it must keep the no-fix half:
// the line makes the stall visible, it does not end it and it does not complete the login.
// Dropping the gate, the line text, or the no-fix clause fails this check.
assertAll('docs/guide/login-plugins.md', [
  'On connect-spigot **0.15.16 and newer** the connector logs one named line for this stall',
  '`Connect tunneled login stalled: no LOGIN_START reached the connector within 10000 ms`, with the player,',
  'session and endpoint), so the stall can be attributed instead of guessed at - the line makes the stall visible, it does',
  'not end it and it does not complete the login.',
])

// The `connect-player` attribute contract lived only in the connector repository, so plugin
// authors reading the public guide could not find it (the route 404'd). It is a published,
// permanent contract: pin the exact attribute name, the Netty key, the constant, the single
// set-site, the passthrough-absent rule and the no-rename sentence, so a rename or a
// "passthrough is marked too" rewrite fails the build instead of shipping silently.
assertAll('docs/guide/login-plugin-integration.md', [
  '`connect-player` (exact, permanent)',
  'io.netty.util.AttributeKey.valueOf("connect-player")',
  'com.minekube.connect.api.ConnectAttributes.CONNECT_PLAYER',
  'com.minekube.connect.api.player.ConnectPlayer',
  'Velocity, BungeeCord, Spigot - one set-site covers all three',
  '`Auth#isPassthrough()` is `false`',
  'it is a passthrough session',
  '**`connect-player` is a permanent public contract.**',
  'it will **never be\nremoved or renamed**',
  'them exactly like a published method',
  'ConnectApi.getInstance()',
  'isConnectPlayer(uuid)',
  'no re-assert floor at all',
  'the code is the same three lines',
  // The ordering fact behind the packet-level row: the connector's own login path is triggered
  // by the login start packet, and a packet-level listener observes it first (live pipeline
  // dump: PacketEvents' decoder ahead of the vanilla decoder, both ahead of the connector's
  // data handler), so a plugin that consumes it leaves the login pending with no decision for a
  // re-assert to restore. Losing either half re-opens the "just re-assert it" answer.
  'the connector completes its own login only when',
  'receives the login start packet',
  'runs **upstream** of that handler',
  'leaves the login pending with nothing to complete it',
  'there is no login decision left to restore',
  "the only thing that ends the stall is the server's own login timeout",
])

// Plugin authors get the same fact on the contract page, next to the sentence it qualifies: from
// connect-spigot 0.15.16 the connector names this stall itself, one line, ~10 s after the
// handshake. It is version-gated and it still is not a fix - authors must not read it as Connect
// completing the login, and support must not read it as a reason to stop attributing the case.
assertAll('docs/guide/login-plugin-integration.md', [
  'On connect-spigot **0.15.16 and newer** the connector logs one named line for this stall',
  '`Connect tunneled login stalled: no LOGIN_START reached the connector within 10000 ms`, with the player,',
  'session and endpoint), so the stall can be attributed instead of guessed at - the line makes the stall visible, it does',
  'not end it and it does not complete the login.',
])

// The support-facing matrix needs the packet-level row too: it is where support answers start,
// and "High"/"do not upgrade Connect" are the two facts that keep a case from being misrouted.
assertAll('docs/guide/compatibility.md', [
  'Mojang sessionserver check - for example LoginTo 4.0.x with its premium mode enabled) | High |',
  'Packet-level premium autologin',
  'do not answer it with a Connect version upgrade',
  'the `connect-player` exemption only exists for Connect-authenticated sessions',
  '[integration contract](/guide/login-plugin-integration)',
])

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
  'latestBatch: September 15 – September 27, 2026',
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
  '<!--@include: ./2026-09-27.md-->',

  '<!--@include: ./2026-09-14.md-->',
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

assertAll('docs/changelog/2026-09-27.md', [
  'date: 2026-09-27',
  // Every product in the batch links a release that has downloadable assets.
  'connect-java/releases/tag/0.15.15',
  'geyserlite/releases/tag/v0.5.31',
  'gate/releases/tag/v0.74.0',
  'gate/releases/tag/v0.74.27',
])

assertAll('docs/changelog/2026-09-14.md', [
  'date: 2026-09-14',
  // Every product in the batch links a release that has downloadable assets.
  'connect-java/releases/tag/0.15.13',
  'gate/releases/tag/v0.73.13',
  'geyserlite/releases/tag/v0.5.27',
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
  // The published contract has to be reachable from the docs nav, next to the page it belongs
  // to; a page nobody can navigate to is still unpublished.
  "link: '/guide/login-plugin-integration'",
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
