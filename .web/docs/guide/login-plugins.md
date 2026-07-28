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
somewhere else. There is no second Mojang session left for your proxy to verify, so **any plugin that forces online mode
on a Connect-tunneled connection will hang that player's login**. The player gets no kick message and never finishes
joining.

## Known Incompatible: LibreLogin with Premium Autologin

LibreLogin registers its pre-login handler to run after Connect's and overwrites Connect's decision. For any player it
considers premium, it forces online mode on the connection. Your proxy then sends an encryption request into a tunnel
whose real client is attached to the Connect edge, and the login state machine never completes. This is the reported
"player falls into limbo and never joins" symptom.

### Workaround

Disable premium autologin in LibreLogin, so that no user has a premium UUID. LibreLogin then only ever forces offline
mode, which agrees with what Connect already decided, and logins complete.

### What the Workaround Does Not Fix

Even with premium autologin disabled, LibreLogin still replaces the game profile Connect provides with its own database
UUID and drops Connect's skin properties. Two consequences you should expect:

- **default skins** for players who joined through Connect
- **plugins that call Connect's `isConnectPlayer(uuid)` API stop recognising tunneled players**, because they look the
  player up by a UUID that no longer matches

The workaround makes logins complete. It does not restore Connect's identity handling. If skins or Connect-API-based
integrations matter to your setup, this combination is not fully working yet.

A fix along these lines has been proposed upstream to LibreLogin, following the exemption mechanism LibreLogin already
implements for other externally authenticated players. This page will link the upstream tracking item once it is
available.

<!-- Maintainers: when the upstream LibreLogin change lands, replace the paragraph above with the tracking link and
update the workaround section to state the fixed versions. -->

## The General Rule for Any Login Plugin

The conflict is not specific to one plugin. Check any login or auth plugin against this rule:

| Plugin behavior | Result with Connect |
| --- | --- |
| Only acts on offline-mode connections, never forces online mode | Compatible by design |
| Can force online mode at pre-login | Hangs the login of every Connect-tunneled player it applies to |
| Rewrites the game profile after Connect has set it | Login completes, but Connect's UUID and skin are lost |

If a plugin exposes a Floodgate-style "skip externally authenticated players" exemption, its author can cooperate with
Connect by including Connect's `connect-player` connection marker in that exemption, alongside the Floodgate one. That
requires no dependency on Connect classes.

See the [Compatibility Matrix](/guide/compatibility#proxy-and-login-plugins) for the wider set of proxy and login plugin
combinations that need extra care.
