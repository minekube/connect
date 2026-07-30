# Login and Auth Plugins

_This page explains how Connect changes the login flow on your proxy, which login plugins conflict with it, and how to
tell in advance whether a given plugin will work._

## How Connect Changes Login

For players using Connect's normal authenticated flow, Connect authenticates them at the
[Connect edge](/guide/#the-connect-network) before they ever reach your proxy or server. By the time the connection
arrives at your [connector](/guide/connectors/), the Mojang session has already been verified and consumed. This page's
UUID and skin guidance applies to those players; unauthenticated players allowed through with
`allow-offline-mode-players` do not have a Mojang identity for Connect to supply.

For an authenticated Connect player, the Connect plugin therefore does two things during login:

- at pre-login, it marks the connection **offline mode**, so your proxy does not send an encryption request
- it supplies the player's real Mojang UUID and skin properties itself, out of band

This is not the same as running an [offline-mode](/guide/offline-mode) server. The player is authenticated - just
somewhere else. There is no second Mojang session left for your proxy to verify, so **a Connect-tunneled connection will
hang if an online-mode decision reaches the proxy**. The player gets no kick message and never finishes joining. Connect
v0.13.1 and newer protect the documented LibreLogin path by default; arbitrary-plugin guarantees depend on proxy ordering
support.

## LibreLogin with Premium Autologin

LibreLogin's premium autologin handler can overwrite Connect's offline-mode decision and force the connection online.
On Connect versions before v0.13.1, the proxy then sends an encryption request into a tunnel whose real client is
attached to the Connect edge, and the login never completes.

**LibreLogin works out of the box with Connect v0.13.1 or newer.** Connect now re-asserts its own offline-mode decision
after other login plugins have run, so premium autologin stays enabled and the login completes. On Velocity, Connect also
restores the player's skin properties by default while leaving LibreLogin's database UUID unchanged.

On BungeeCord, the login and premium autologin fix still applies. During pre-login, BungeeCord exposes no
profile-properties API on the pending connection, so `restore-full-profile` can re-assert the Mojang UUID and username
but cannot restore skin properties through this setting.

### Re-assert Configuration

The shipped defaults are:

```yaml
login-reassert:
  enabled: true
  restore-full-profile: false
```

`login-reassert.enabled` is the optional off-switch. Set it to `false` only if you deliberately want another plugin to
change Connect's login decision after Connect. This restores the previous last-writer-wins behavior; a plugin that
forces online mode can make Connect players hang during login again.

By default, Connect preserves the UUID chosen by the login plugin. This keeps the plugin's database lookups consistent,
but code that calls Connect's `isConnectPlayer(uuid)` or `getPlayer(uuid)` API may not recognise the player under that
different UUID.

To restore Connect's full profile, including the Mojang UUID and username, set
`login-reassert.restore-full-profile` to `true`. **Before enabling it with LibreLogin, set
`new-uuid-creator: MOJANG` in LibreLogin's config.** With LibreLogin's default `CRACKED` UUID creator, full-profile
restoration makes its database lookups miss and its join handlers fail.

## The General Rule for Any Login Plugin

The conflict is not specific to one plugin. Check any login or auth plugin against this rule:

Connect v0.13.1+ guarantees the LibreLogin path described above. For arbitrary login plugins, strict after-all protection on
Velocity requires the numeric-priority API; on legacy Velocity, Connect uses `PostOrder.LAST`, so another `LAST` handler
can still run after it depending on plugin load order.

| Plugin behavior | Result with Connect |
| --- | --- |
| Only acts on offline-mode connections, never forces online mode | Compatible by design |
| Can force online mode at pre-login | Connect v0.13.1+ re-asserts its offline-mode decision by default; on legacy Velocity, arbitrary plugins may still depend on plugin load order; older versions, or a disabled re-assert, can hang during login |
| Rewrites the game profile after Connect has set it | Connect v0.13.1+ restores skin properties on Velocity by default; the plugin's UUID remains unless full-profile restoration is enabled |

If a plugin exposes a Floodgate-style "skip externally authenticated players" exemption, its author can cooperate with
Connect by including Connect's `connect-player` connection marker in that exemption, alongside the Floodgate one. That
requires no dependency on Connect classes.

See the [Compatibility Matrix](/guide/compatibility#proxy-and-login-plugins) for the wider set of proxy and login plugin
combinations that need extra care.
