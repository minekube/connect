---
description: The connect-player attribute contract for login, auth, and premium-detection plugins that run beside Minekube Connect.
---

# Login Plugin Integration

_This page is the contract for authors of login, auth, or premium-detection plugins - LibreLogin, AuthMe, nLogin, and
anything similar - that run on a proxy or server which also runs Minekube Connect. It is written for plugin authors;
the operator-facing behaviour (which plugin shapes conflict, and the `login-reassert` floor) is on
[Login and Auth Plugins](/guide/login-plugins)._

Connect authenticates players at its edge and hands your proxy an already-verified, offline-mode connection. Connect
publishes a stable marker on that connection - the `connect-player` Netty channel attribute - so your plugin can detect
a Connect-tunneled player and skip its own login flow. This is the same pattern many plugins already implement for
Floodgate's `floodgate-player` attribute.

## Why the marker exists

Connect terminates the real client connection at the Minekube edge, performs the Mojang online-mode handshake there, and
then relays the player's traffic into a local channel inside your proxy. On that inner connection, Connect marks the
player **offline mode** and supplies the real Mojang UUID, username, and skin properties itself.

There is no second Mojang session for the proxy to verify. So any plugin that later forces **online mode** on a
Connect-tunneled connection makes the proxy send an `EncryptionRequest` that can never be answered - the client's login
was already finished at the Connect edge, so the player's login hangs and never completes. The same applies to rewriting
the game profile after Connect has set it: the player ends up with the wrong UUID and no skin.

That is also why a plugin which starts its **own** login handshake from the login packet cannot work on a tunneled
connection: there is no live client side left to complete it. See the packet-level row in
[the general rule for any login plugin](/guide/login-plugins#the-general-rule-for-any-login-plugin).

The `connect-player` attribute lets you detect that situation before you act.

## The contract

### `connect-player` - channel attribute

| | |
|---|---|
| Attribute name | `connect-player` (exact, permanent) |
| Netty key | `io.netty.util.AttributeKey.valueOf("connect-player")` |
| Constant | `com.minekube.connect.api.ConnectAttributes.CONNECT_PLAYER` |
| Value type | `com.minekube.connect.api.player.ConnectPlayer` (never `null` when the attribute is set) |
| Set on | The proxy-facing channel, at channel creation, before any platform login event fires |
| Platforms | Velocity, BungeeCord, Spigot - one set-site covers all three |
| Set when | The session is authenticated by Connect, i.e. `Auth#isPassthrough()` is `false` |
| Absent when | The connection is not Connect-tunneled, **or** it is a passthrough session |

Passthrough sessions are deliberately left unmarked: Connect did not authenticate them, so they should still go through
your plugin's normal login flow. An offline-mode endpoint, and any endpoint with `allow-offline-mode-players` enabled,
serves passthrough sessions - on those endpoints the marker is **not** available as an integration point, which is the
case that most often surprises plugin authors. See [Offline Mode](/guide/offline-mode).

Because the marker is set when the channel is created, it is already present in Velocity's `PreLoginEvent` and
`GameProfileRequestEvent`, in BungeeCord's `PreLoginEvent`, and in Spigot's login handling - including handlers
registered at the very earliest priority.

### `ConnectApi` - UUID lookups for online players

Once a player has finished logging in, `com.minekube.connect.api.ConnectApi` answers the same question by UUID:

```java
ConnectApi api = ConnectApi.getInstance();
boolean tunneled = api.isConnectPlayer(uuid);   // is this online player tunneled by Connect?
ConnectPlayer player = api.getPlayer(uuid);     // null if not
```

These are the right tool for post-login checks (session tracking, limbo routing, command gating). They are **not**
usable during login, because the player is not registered yet - that is what the channel attribute is for.

## How to integrate

If your plugin already exempts Floodgate players, the Connect exemption is the identical shape. Add it at every point
where you would otherwise override the connection's authentication decision or identity.

**Presence check, with no dependency on Connect** - Netty interns attribute keys by name, so this compiles and works
whether or not Connect is installed:

```java
private static final AttributeKey<Object> CONNECT_PLAYER =
        AttributeKey.valueOf("connect-player");

private static boolean isConnectPlayer(Channel channel) {
    return channel.hasAttr(CONNECT_PLAYER) && channel.attr(CONNECT_PLAYER).get() != null;
}
```

**Velocity - pre-login.** Reach the channel the same way you already do for Floodgate
(`LoginInboundConnection#delegate` → `InitialInboundConnection` → `MinecraftConnection` → `Channel`), then bail out:

```java
@Subscribe(order = PostOrder.LAST)
public void onPreLogin(PreLoginEvent event) {
    Channel channel = channelOf(event.getConnection());
    if (channel != null && isConnectPlayer(channel)) {
        return; // externally authenticated by Connect - never force online mode
    }
    ...
}
```

**Velocity - game profile.** Do not rebuild the profile from the original one; Connect has already put the player's real
UUID and skin properties in it.

**BungeeCord - pre-login.** Same check before any `setOnlineMode(true)`.

**Post-login / session tracking.** Use `ConnectApi.isConnectPlayer(player.getUniqueId())` to skip your own
authorization tracking, limbo routing, and command blocking, the same way you skip them for Floodgate players.

**Optional: read the player.** If you add Connect's `api` artifact as a dependency you can use
`ConnectAttributes.CONNECT_PLAYER` directly and read the `ConnectPlayer` value for the player's UUID, username, game
profile, and `Auth`. That is why the attribute value is a `ConnectPlayer` rather than a bare UUID - presence alone costs
you nothing, and the full identity is there when you want it.

### Login packet handlers

The marker covers plugins that decide **what to do with** the connection. It does not help a plugin that instead
intercepts the login **packet** and runs its own handshake (a premium/online-mode autologin that cancels `LoginStart`,
sends its own `EncryptionRequest` and checks the Mojang sessionserver): Connect's client login is already finished at the
edge, so that handshake has no peer to answer it. Such a plugin needs its premium/online-mode mode disabled for Connect
traffic, or a login path that does not start a second handshake. See
[the general rule for any login plugin](/guide/login-plugins#the-general-rule-for-any-login-plugin).

## Connect defends its own decision (operators)

Not every login plugin exempts Connect, so Connect does not rely on one. On proxies it registers a second pre-login
handler that restores its own decision if something changed it, which is what keeps a "force online mode for everyone"
plugin from hanging Connect players. That handler is a floor, not a veto, and it only ever acts on connections Connect
itself authenticated.

The config keys, the ordering guarantees, and what the floor cannot do are documented once, for operators, on
[Login and Auth Plugins](/guide/login-plugins#re-assert-configuration). That detail is deliberately not repeated here.
Note the scope: `login-reassert` is implemented on Velocity and BungeeCord only. On a plain Paper/Spigot server there is
no re-assert floor at all, so the `connect-player` marker above is the integration path and a Spigot login plugin is only
safe if it ignores tunneled players or checks the marker itself.

## Stability commitment

**`connect-player` is a permanent public contract.**

Minekube commits to this explicitly: once external plugins depend on the name `connect-player`, it will **never be
removed or renamed**. The attribute name, the fact that it is set before any login event, and the rule that it is absent
for passthrough sessions are all part of Connect's published API surface - treat them exactly like a published method
signature. Changes may only be additive.

This is a deliberate commitment, not an accident of the implementation. A future contributor must not treat
`connect-player` as an internal detail that is free to be refactored, renamed, or moved to a later point in the
connection lifecycle. Doing so silently breaks every third-party login plugin that integrates with Connect, with no
compile error anywhere to warn about it.

The commitment is enforced in Connect's own test suite, which pins both the exact name and the set-site, so a rename or a
moved marker fails its build.

## Checking a login plugin for compatibility

- **Not compatible** if it starts its own login handshake from the login packet (its premium/online-mode autologin), or
  if it can force online mode at pre-login (`forceOnlineMode()` / `setOnlineMode(true)`), or if it rewrites the game
  profile after Connect has set it.
- **Compatible by design** if it only ever acts on offline-mode connections, never forces online mode, and does not
  intercept the login packet itself.
- **If it has a Floodgate exemption**, ask its author to extend that exemption to Connect's `connect-player` attribute -
  the code is the same three lines. That only applies where Connect authenticated the session; for passthrough
  sessions the marker is absent by design.

## Related source

Connect's plugin code is public in [minekube/connect-java](https://github.com/minekube/connect-java):

- `api/src/main/java/com/minekube/connect/api/ConnectAttributes.java` - the attribute key and its Javadoc contract
- `api/src/main/java/com/minekube/connect/api/ConnectApi.java` - `isConnectPlayer(UUID)` / `getPlayer(UUID)`
- `api/src/main/java/com/minekube/connect/api/player/Auth.java` - `isPassthrough()`
- `core/src/main/java/com/minekube/connect/network/netty/LocalServerChannelWrapper.java` - the single set-site
- `docs/login-plugin-integration.md` - the connector-repository copy of this page; keep the two in step

_This page is mirrored by `docs/login-plugin-integration.md` in
[minekube/connect-java](https://github.com/minekube/connect-java). Change one, change the other._
