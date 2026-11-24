---
title: "Bedrock players"
description: Let Microsoft/Xbox-authenticated Bedrock players join Connect endpoints with or without a linked Java account.
---

Bedrock players can join your endpoint through the same public addresses Java players already use:

- `<endpoint>.play.minekube.net`
- custom domains attached to the endpoint
- endpoint routing through the Connect network

The Connect edge translates Bedrock traffic before it reaches your connector, so the usual Connect setup stays the same
for Paper/Spigot, Velocity, BungeeCord, and Gate connectors.

<div class="mk-arch-wrap not-content">
<div class="mk-arch" role="img" aria-label="Java players on port 25565 and Bedrock players on port 19132 join the same Connect address; the Connect edge translates Bedrock and tunnels everyone to your Java server or proxy.">
<div class="a-col"><div class="a-node sm"><svg viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><rect width="20" height="14" x="2" y="3" rx="2"/><path d="M8 21h8"/><path d="M12 17v4"/></svg><strong>Java</strong><small>port 25565</small></div><div class="a-node sm"><svg viewBox="0 0 24 24" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><rect x="5" y="2" width="14" height="20" rx="2"/><path d="M12 18h.01"/></svg><strong>Bedrock</strong><small>port 19132</small></div></div>
<div class="a-link"><span>join</span></div>
<div class="a-node a-edge"><img src="/minekube-logo.png" alt="" width="28" height="28"><strong>Connect edge</strong><small>name.play.minekube.net</small><em class="a-chip">Bedrock translation</em></div>
<div class="a-link a-tunnel"><span>tunnel</span></div>
<div class="a-zone"><p class="a-zone-label">Your server · no Geyser, no UDP port</p><div class="a-stack"><div class="a-node"><img src="/images/logos/paper.png" alt="" width="28" height="28"><strong>Your server</strong><small>Java only</small><em class="a-chip">Connect plugin or Gate</em></div></div></div>
</div>
</div>

:::tip[Managed default]
Microsoft/Xbox-authenticated Bedrock players can join Connect endpoints without owning or linking Java Edition. If an
authoritative Java link exists, Connect preserves the linked Java profile; otherwise it uses the player's verified
native Bedrock/XUID identity. No Bedrock-specific server configuration is required.
:::

## At a glance

| Your situation | Path | What you do |
| --- | --- | --- |
| Players join through a Connect address or a custom domain attached to your endpoint | Connect-managed Bedrock | Nothing Bedrock-specific. Players use port `19132`. |
| Players join a Gate address you run yourself, outside Connect | self-hosted Gate direct Bedrock | Add `bedrock: true` to that standard Gate instance. |
| You run Gate Lite | Gate Lite behind Connect only | Gate Lite cannot be a direct Bedrock listener. Use standard Gate for that. |

- **Works:** Java and Bedrock on the same endpoint subdomain or custom domain, players with or without a linked Java
  account, and the compatibility data Connect forwards to backend Floodgate API users.
- **Does not work:** Bedrock players without valid Microsoft/Xbox authentication, plugins that require a local
  Floodgate/Geyser runtime, and Gate Lite as a direct Bedrock listener.

## Connect-routed Bedrock

Use this path when players join through Connect addresses such as `<endpoint>.play.minekube.net`, `minekube.net`, or a
custom domain attached to your endpoint.

For Bedrock clients, add the same Connect hostname as the server address and use the normal Bedrock port `19132`:

- **Server Address:** `<endpoint>.play.minekube.net`
- **Port:** `19132`

A custom domain attached to the endpoint also works with port `19132`.

This is a player-facing Connect edge port. It does not mean you need to open UDP `19132` on your home network, VPS,
backend server, or connector host.

You do not need to:

- install Geyser on your backend
- install Floodgate on your backend
- open a local UDP Bedrock port
- set `bedrock: true` in your backend Gate just for Connect players

This applies to the Connect Java plugin and to Gate used as a Connect connector.

Do not enable backend Gate `bedrock: true` for Connect-managed Bedrock. That setting is for self-hosted Gate direct
Bedrock listeners and can create a second Bedrock path that support then has to diagnose separately.

## Direct self-hosted Gate Bedrock

Use this path when Bedrock players connect directly to a Gate address you operate yourself, outside the Connect network.
Standard Gate can then manage Bedrock locally:

```yaml
bedrock: true
```

Gate then starts and manages the Bedrock translation runtime for direct Bedrock clients. This local setting does not
change how Connect-routed Bedrock players are handled.

## Domains

Java and Bedrock players share Connect domains. If a custom domain already routes Java players to your endpoint, Bedrock
players use the same domain after it is attached to the endpoint in the dashboard.

| Client | Normal port |
| --- | --- |
| Java | `25565` |
| Bedrock | `19132` |

Use a separate direct Bedrock address only when you intentionally run local Bedrock support on your own standard Gate
instance.

## Support behavior matrix

<div class="mk-setups not-content">
<div class="mk-setup"><p class="mk-setup-title">Connect-managed Bedrock with Paper/Velocity/Bungee connector</p><dl><dt>What to configure</dt><dd>Install the Connect plugin normally. Do not enable backend Gate <code>bedrock: true</code> for Connect-managed Bedrock.</dd><dt>Bedrock translation</dt><dd>Connect edge</dd><dt>Backend Floodgate API</dt><dd>Supported only for the compatibility data Connect forwards. If a plugin expects a local Floodgate/Geyser runtime, collect logs and plugin names.</dd><dt>Online-mode and accounts</dt><dd>Java online-mode stays on the Java path. Microsoft/Xbox-authenticated Bedrock players do not need a linked Java account by default.</dd><dt>Custom domains</dt><dd>Supported. Bedrock and Java can use the same endpoint subdomain or custom domain.</dd></dl></div>
<div class="mk-setup"><p class="mk-setup-title">Connect-managed Bedrock with standard Gate <code>connect</code> enabled</p><dl><dt>What to configure</dt><dd>Enable the Gate Connect connector normally. Do not add local Bedrock settings only for Connect players.</dd><dt>Bedrock translation</dt><dd>Connect edge</dd><dt>Backend Floodgate API</dt><dd>Same Connect-managed compatibility behavior as plugin connectors.</dd><dt>Online-mode and accounts</dt><dd>Microsoft/Xbox-authenticated Bedrock players do not need a linked Java account by default. Do not switch the backend to offline mode only for Bedrock.</dd><dt>Custom domains</dt><dd>Supported. Connect addresses and custom domains stay shared.</dd></dl></div>
<div class="mk-setup"><p class="mk-setup-title">self-hosted Gate direct Bedrock</p><dl><dt>What to configure</dt><dd>Add <code>bedrock: true</code> to the standard Gate instance that Bedrock clients join directly.</dd><dt>Bedrock translation</dt><dd>Your Gate instance</dd><dt>Backend Floodgate API</dt><dd>Local Gate/Floodgate behavior belongs to that Gate setup, not Connect-managed Bedrock.</dd><dt>Online-mode and accounts</dt><dd>Follow Gate's direct Bedrock and account-linking requirements.</dd><dt>Custom domains</dt><dd>Use a separate direct Bedrock hostname or port unless you intentionally route that hostname outside Connect.</dd></dl></div>
<div class="mk-setup"><p class="mk-setup-title">Standard Gate with both Connect and direct Bedrock</p><dl><dt>What to configure</dt><dd>Enable Connect for Connect addresses and <code>bedrock: true</code> only for the separate direct Bedrock listener.</dd><dt>Bedrock translation</dt><dd>Connect edge for Connect addresses, Gate for the direct Bedrock address</dd><dt>Backend Floodgate API</dt><dd>Diagnose by join address. Connect-managed joins use Connect compatibility data; direct joins use local Gate behavior.</dd><dt>Online-mode and accounts</dt><dd>Ask which address the player used before changing online-mode or linking settings.</dd><dt>Custom domains</dt><dd>Supported when DNS clearly separates Connect-routed names from direct Gate names.</dd></dl></div>
<div class="mk-setup"><p class="mk-setup-title">Gate Lite behind Connect</p><dl><dt>What to configure</dt><dd>Configure Gate Lite as a Java connector behind Connect.</dd><dt>Bedrock translation</dt><dd>Connect edge</dd><dt>Backend Floodgate API</dt><dd>Treat backend Floodgate API reports as Connect-managed compatibility reports.</dd><dt>Online-mode and accounts</dt><dd>Gate Lite is not the direct Bedrock authority; keep auth decisions in the backend/proxy path.</dd><dt>Custom domains</dt><dd>Supported for Connect-routed endpoint domains.</dd></dl></div>
<div class="mk-setup no"><p class="mk-setup-title">Gate Lite as a direct Bedrock listener</p><dl><dt>What to configure</dt><dd>Use standard Gate instead.</dd><dt>Bedrock translation</dt><dd>Not supported in Gate Lite</dd><dt>Backend Floodgate API</dt><dd>Not supported.</dd><dt>Online-mode and accounts</dt><dd>Not supported.</dd><dt>Custom domains</dt><dd>Not supported.</dd></dl></div>
</div>

## Account linking and online mode

Connect-routed Bedrock players arrive through the managed Bedrock path and are represented through the compatibility
layer before they reach your connector. The default endpoint policy accepts a Microsoft/Xbox-authenticated Bedrock
identity without requiring a linked Java Edition account. This is separate from generic offline-mode Java access.

- With an authoritative account link, the endpoint receives the linked Java identity.
- Without a link, the endpoint receives a stable profile derived from the verified Bedrock XUID.
- Without valid Microsoft/Xbox authentication, the managed Bedrock session is rejected.

See the complete [player identity matrix](/guide/joining#who-can-join) for the corresponding Java and offline-Java paths.

### Rejection reasons

The exact rejection reason identifies the component that needs attention:

| Reason | Meaning | What to do |
| --- | --- | --- |
| `policy_linked_java_only` | The endpoint has an explicit linked-Java-only policy override. Connector settings such as `bedrock-identity.enforcement` cannot change this edge policy. | Include the endpoint name when asking Minekube support to review an unexpected override. |
| `scope_unavailable` | Connect cannot resolve a usable endpoint/organization identity scope. | Reinstalling Geyser or changing backend online mode will not repair endpoint ownership. |
| `capability_unavailable` | The selected connector did not prove the required Bedrock identity capability. | Update to a compatible connector build. Identity keys, URLs, and capability settings are automatic for the managed service; changing backend authentication does not repair a missing connector capability. |

### Troubleshooting online mode

First identify whether the player joined through a Connect address or a direct self-hosted Gate Bedrock address. The
required configuration is different for those two paths.

Keep online-mode decisions tied to the whole Java forwarding topology. A Connect-managed Bedrock report is not, by
itself, a reason to disable online-mode or allow offline-mode players on the backend.

Connect-managed Bedrock identity uses official Microsoft/Xbox Bedrock authentication at the Connect edge. It is not:

- a Minekube password system
- backend Geyser authentication
- the same thing as allowing generic offline-mode Java players

The connector can enforce the Bedrock identity that Connect already verified, while Java online-mode players continue to
use the normal Java session path.

## Bedrock Identity Enforcement

For the normal managed service, install or update your connector and keep its defaults. You do not need to configure
identity URLs, copy keys, or declare capabilities. Current Connect Java releases also supply the legacy identity defaults
when that section is absent from an older configuration file. Explicit identity overrides remain operator-controlled.

:::caution[Gate version compatibility]
Gate v0.74.28 and earlier do not include automatic verification of the managed v1 Bedrock handover. Upgrade to
[Gate v0.74.29 or later](https://github.com/minekube/gate/releases/tag/v0.74.29) for automatic Connect Bedrock identity support.
Keep the normal Connect configuration; no Bedrock or identity settings are needed.
Configuring v2 identity keys on an older Gate build does not make it compatible with the managed v1 handover.
:::

How it works:

1. The Connect edge verifies the player's Microsoft/Xbox Bedrock identity.
2. It signs a short-lived identity envelope and forwards the session to your connector.
3. Gate verifies this envelope before opening the player tunnel. It checks the signature, expiry,
   endpoint and organization scope, session binding, and replay protection before admitting the Bedrock identity. The
   Connect Java Plugin also provides a verifier for the managed handover.

Minekube publishes the current and previous verification keys for automatic key rotation. These are public verification
keys; private signing keys stay at the Connect edge.

Custom Watch services and explicit signed-principal v2 deployments use their own trust configuration. See the
[Connect Java identity reference](https://github.com/minekube/connect-java/blob/main/docs/bedrock-identity.md)
for those advanced deployments. They are separate from the normal managed setup.

## What a Bedrock player looks like at your backend

A Connect-managed Bedrock player has no Java Edition account, so the endpoint receives a profile that Connect derives
from the player's verified Xbox identity. Both parts have a fixed shape for every unlinked Bedrock player, and both are
intentional.

| Value | Form | Example |
| --- | --- | --- |
| Username | `_<gamertag>` | `icedRyan` appears as `_icedRyan` |
| UUID | RFC 4122 version-5, derived from the XUID | `xxxxxxxx-xxxx-5xxx-yxxx-xxxxxxxxxxxx` |

No action is required for either value. If a backend plugin rejects the player because the username or UUID is not a
Mojang one, that plugin is the thing to adjust. Rewriting the UUID or switching the backend to offline mode does not
make a Connect-managed Bedrock join produce a Mojang identity.

### Username

The Connect edge applies the configured Bedrock username prefix `.` to the gamertag, and
Connect normalizes the result into `[A-Za-z0-9_]`, capped at 16 characters, so that every backend implementation
accepts the login name. That is why the prefix reaches your backend as `_`, and why every space or non-ASCII character in the gamertag is encoded to `_` as
well. A gamertag like `icedRyan` appears as `_icedRyan`, and `.LLG icedRyan` appears as `_LLG_icedRyan`.

Why not keep the dot? The Java account alphabet is the narrower `[A-Za-z0-9_]`, but what constrains a proxy-supplied
login name is the backend server's own validation, not one Java-wide rule:

- Paper's check also accepts `.` and rejects spaces.
- vanilla's check is only the `[A-Za-z0-9_]` alphabet.
- Both cap the name at 16 characters.

Connect normalizes into the vanilla alphabet because it cannot know which implementation a backend runs, so the `.`
prefix is not a way to hand a dot to a backend that validates names the vanilla way.

### UUID

The UUID is a stable RFC 4122 version-5 UUID derived from the verified Bedrock XUID (the established XUID namespace,
SHA-1 hashed with the version-5 and variant bits applied). It deliberately does not look like a Mojang UUID, so a
version-5 value such as `xxxxxxxx-xxxx-5xxx-yxxx-xxxxxxxxxxxx` is the normal form and not a defect.

The same XUID always produces the same UUID, so bans, permissions, economy entries, and other UUID-keyed data stay
attached to the same player on every reconnect and on every endpoint.

When you are debugging:

- It is an identifier, not a secret and not a signature. Seeing it in logs is expected.
- Do not rewrite it. Connect checks the profile UUID of the Bedrock session against the XUID-derived value before it
  hands the signed Bedrock identity to your connector, so a component that substitutes a different UUID makes Connect
  reject the session instead of fixing anything.

### Floodgate's `00000000-0000-0000-XUID` key is a different identifier

Geyser and Floodgate identify a Bedrock connection with a second value: the XUID placed in the low 64 bits of an
otherwise zero UUID, conventionally written `00000000-0000-0000-XUID`. That value belongs to a local Geyser + Floodgate
runtime. A Connect-managed join does not deliver it to your backend.

This is why UUID-keyed data does not follow a player from a plain Geyser + Floodgate setup into Connect on the first
join. Both identifiers are stable and both come from the same verified XUID, but they are different values, and no
Connect setting renames one into the other.

## Whitelisting a Bedrock player

The whitelist is not part of Connect or Gate. Gate's built-in commands are `/server`, `/glist`, and `/send`, and
neither project ships a player whitelist, so the list that decides who may join lives on your backend server (or in a
whitelist plugin running on it).

### Add the entry by UUID

1. Let the player attempt one join. Vanilla and Paper log the identity your backend was handed, in the line
   `UUID of player _<gamertag> is <uuid>`, before the whitelist is consulted. An attempt that was refused for not
   being whitelisted still records the UUID.
2. Add that UUID to `whitelist.json` with the same name the log showed:

   ```json
   [
     { "uuid": "xxxxxxxx-xxxx-5xxx-yxxx-xxxxxxxxxxxx", "name": "_<gamertag>" }
   ]
   ```

3. Run `/whitelist reload`.

The whitelist matches on the UUID. The `name` field is what your server displays and what `/whitelist list` shows, so
copy it from the log rather than retyping the gamertag. Keep every other entry in the file: the whole file is re-read on
reload.

That UUID is the value your backend keys everything on (permissions, bans, economy, player data), so it is also the
value to use in a permissions plugin. It is stable across reconnects and across endpoints.

### Why `/whitelist add <name>` does not work

`/whitelist add <name>` cannot add a Connect-managed Bedrock player, because the command resolves the name through
Mojang's profile API. A Bedrock gamertag has no Java profile, and the name your backend receives for an unlinked
player is the normalized `_<gamertag>` form described above, which resolves to nothing either, so the add fails.

If the name you typed happens to be a real Java account, the command does succeed, but it whitelists that Java player,
whose UUID is a different one, so the Bedrock player still cannot join. No Connect setting changes either outcome, and
neither does reinstalling the connector.

### If no UUID line is logged

If the player is refused before that line is written, something earlier in the pipeline is deciding, such as a login
plugin or a whitelist enforced by a proxy. Look in that component's own log for the profile it was handed before
changing anything.

Turning the whitelist off is the last resort, and it is not a private action: while it is off, anyone who can reach that
backend may join. Keep the window short, do it on a backend that is not directly reachable from the internet, and turn
it back on as soon as you have the UUID.

### A Java link moves the whitelisted UUID

The entry has to match the identity your backend receives, not the gamertag:

| Player | Whitelist entry |
| --- | --- |
| Unlinked Bedrock player | The version-5 UUID derived from the XUID. |
| Linked player | The linked Java UUID, so an entry added before the link stops matching. |

A Bedrock entry never matches a Java player, and a Java entry never matches an unlinked Bedrock player.

If a player links a Java account later, move their entry to the linked Java UUID. Do not switch the backend to offline
mode to make a name-based whitelist accept a Bedrock player: that changes the identity of every Java player on the
backend, and it still does not create a Mojang profile for the gamertag.

## Backend Floodgate API compatibility

Some backend plugins query Floodgate-style player metadata to decide whether a player is from Bedrock, linked to Java,
or allowed past a login step. Connect-managed Bedrock forwards compatibility data for supported paths, but servers can
still fail when a plugin assumes it is talking to a local Floodgate/Geyser installation.

When Floodgate-dependent behavior fails, ask for:

- the exact join address the Bedrock player used
- connector type and version
- backend server type and version
- Floodgate, Geyser, auth, or login plugin names and versions
- the kick text and logs from Connect, the connector, proxy, and backend around the same timestamp
- whether the same player succeeds through a direct self-hosted Gate Bedrock address, if one exists
- whether `bedrock-identity.enforcement` is `disabled`, `warn`, or `require`
- Bedrock identity warnings such as missing envelope, invalid signature, expired identity, replayed nonce, policy
  mismatch, endpoint mismatch, or organization mismatch

Do not recommend installing backend Geyser/Floodgate or enabling Gate `bedrock: true` unless the user wants to operate a
separate self-hosted Gate direct Bedrock path.

Do not ask users to paste signed Bedrock identity envelopes, private keys, account tokens, or full player profile dumps
into public support channels. The useful support data is the join address, connector version, endpoint ownership, relevant
timestamps, sanitized logs, and the exact rejection reason.

## Discord support response draft

> Connect-managed Bedrock is separate from self-hosted Gate Bedrock. If your player joins through
> `<endpoint>.play.minekube.net`, `minekube.net`, or a custom domain attached to the endpoint, Bedrock clients should use
> that hostname with port `19132`. Do not install backend Geyser/Floodgate, open local UDP `19132`, or enable backend Gate
> `bedrock: true` just for that player. Connect handles Bedrock translation at the edge. If a Floodgate-dependent plugin
> rejects the player, send the join hostname, connector type/version, backend type/version, Floodgate/Geyser/auth plugin
> versions, the exact kick text, and logs from the Connect connector/proxy/backend at the same timestamp. Only use
> `bedrock: true` when Bedrock clients connect directly to your own standard Gate instance outside Connect. Unlinked
> Bedrock accounts are accepted by the default Connect endpoint policy; if the kick reason is
> `policy_linked_java_only`, include the endpoint name so Minekube support can review its explicit policy override.
