# AuthSession API <VPBadge type='danger'>Coming soon</VPBadge>

> This feature is currently unavailable.
> You can view a draft of the custom Mojang AuthSession API here: https://github.com/minekube/mojang-multiauth.
> For support requests, join our Discord community. https://minekube.com/discord

This page explains the planned Connect AuthSession API and what it does not support today.

::: info Requirements
-> Requires a public IP address

-> Only useful if you want to use Connect with online mode players.
:::

## Current Status

The AuthSession API is not available in production yet. The example commands on this page are placeholders and should
not be used as working setup instructions.

Today, Connect supports online-mode players through connectors that integrate with Connect's current auth/session
handling, such as standard Gate with Connect enabled or the Connect Java Plugin on supported platforms.

Gate Lite behind Connect is different: Lite mode is a thin Java protocol reverse proxy. It forwards ping and player
authentication to the selected backend route instead of acting as a full proxy that can translate Connect's
edge-authenticated session into a backend online-mode login.

Because AuthSession and Connect passthrough support are not available yet, the following topology is not supported by
Gate Lite configuration alone:

`Connect -> Gate Lite -> Online Mode Backend`

Use standard Gate with Connect enabled or the Connect Java Plugin when you need online-mode players through Connect
today. Use Gate Lite behind Connect only when the backend auth model is compatible with a thin reverse-proxy hop.

## How it Works

The planned Connect AuthSession API is an adaptation of the Mojang AuthSession API. It is intended to securely
authenticate online-mode players from the Connect Network to your server/proxy.

The API distributes Mojang Sessionserver's "hasJoined" requests from online mode servers/proxies possibly behind
Minekube Connect,
Minehut Network, and regular Mojang.

## Using Connect AuthSession API

The following examples are not active setup instructions yet. They show the type of configuration that may be needed
once AuthSession support is released.

### Velocity

::: code-group

```shell [Terminal]
java -Dmojang.sessionserver=[MINEKUBE AUTHSESSION API COMING SOON] -jar velocity.jar
```

:::

#### Gate

TBD. Standard Gate has first-class Connect support today. Gate Lite backend-route AuthSession/passthrough support is not
available yet.

#### Waterfall

::: code-group

```shell [Terminal]
java -Dwaterfall.auth.url="[MINEKUBE AUTHSESSION API COMING SOON]?username=%s&serverId=%s%s" -jar waterfall.jar
```

:::

#### Paper

> Note: These instructions only apply if you are running Paper standalone and **NOT** under a proxy.
> If you are using Gate, Velocity, Waterfall, or Lilypad, you should have already configured this, and you can safely skip
> this section.

Add the following CLI flags to your start script:

```shell
-Dminecraft.api.auth.host=[MINEKUBE AUTHSESSION API COMING SOON]
-Dminecraft.api.account.host=[MINEKUBE AUTHSESSION API COMING SOON]
-Dminecraft.api.services.host=[MINEKUBE AUTHSESSION API COMING SOON]
-Dminecraft.api.session.host=[MINEKUBE AUTHSESSION API COMING SOON]
```

When you are done, your script may look something like the following:

::: code-group

```shell [Terminal]
java -Dminecraft.api.auth.host=[MINEKUBE AUTHSESSION API COMING SOON] \
     -Dminecraft.api.account.host=[MINEKUBE AUTHSESSION API COMING SOON] \
     -Dminecraft.api.services.host=[MINEKUBE AUTHSESSION API COMING SOON] \
     -Dminecraft.api.session.host=[MINEKUBE AUTHSESSION API COMING SOON] \
     -jar paper.jar
```

:::

In addition, ensure that you have set `enforce-secure-profile` to `false` in
your [server.properties](http://server.properties) file.

#### Lilypad

Set the following environment variable:

`LILYPAD_MOJANG_SESSIONSERVER_URL` to `[MINEKUBE AUTHSESSION API COMING SOON]`

#### BungeeCord

This proxy type is currently not supported due
to [an issue in BungeeCord](https://github.com/SpigotMC/BungeeCord/pull/3201); please use Gate, Velocity or Waterfall
instead.

## Compatibility

Only the listed software allows changing the Mojang AuthSession API url.
Bungeecord is not supported.
