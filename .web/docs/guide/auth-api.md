# AuthSession API <VPBadge type='danger'>Coming soon</VPBadge>

> This feature is currently unavailable.
> You can view a draft of the custom Mojang AuthSession API here: https://github.com/minekube/mojang-multiauth.
> For support requests, join our Discord community. https://minekube.com/discord

This guide will help you configure the Connect auth session server across different platforms and
understand the managed reverse proxy by Minekube Connect.

::: info Requirements
-> Requires a public IP address

-> Only useful if you want to use Connect with online mode players.
:::

## How it Works

The Connect's AuthSession API is an innovative adaptation of the Mojang AuthSession API. It's designed to securely
authenticate online mode players from the Connect Network to your server/proxy.

The API distributes Mojang Sessionserver's "hasJoined" requests from online mode servers/proxies possibly behind
Minekube Connect,
Minehut Network, and regular Mojang.

## Using Connect AuthSession API

Follow the instructions for your software below.

### Velocity

::: code-group

```shell [Terminal]
java -Dmojang.sessionserver=[MINEKUBE AUTHSESSION API COMMING SOON] -jar velocity.jar
```

:::

#### Gate

TBD

#### Waterfall

::: code-group

```shell [Terminal]
java -Dwaterfall.auth.url="[MINEKUBE AUTHSESSION API COMMING SOON]?username=%s&serverId=%s%s" -jar waterfall.jar
```

:::

#### Paper

> Note: These instructions only apply if you are running Paper standalone and **NOT** under a proxy.
> If you are using Gate, Velocity, Waterfall, or Lilypad, you should have already configured this, and you can safely skip
> this section.

Add the following CLI flags to your start script:

```shell
-Dminecraft.api.auth.host=[MINEKUBE AUTHSESSION API COMMING SOON]
-Dminecraft.api.account.host=[MINEKUBE AUTHSESSION API COMMING SOON]
-Dminecraft.api.services.host=[MINEKUBE AUTHSESSION API COMMING SOON]
-Dminecraft.api.session.host=[MINEKUBE AUTHSESSION API COMMING SOON]
```

When you are done, your script may look something like the following:

::: code-group

```shell [Terminal]
java -Dminecraft.api.auth.host=[MINEKUBE AUTHSESSION API COMMING SOON] \
     -Dminecraft.api.account.host=[MINEKUBE AUTHSESSION API COMMING SOON] \
     -Dminecraft.api.services.host=[MINEKUBE AUTHSESSION API COMMING SOON] \
     -Dminecraft.api.session.host=[MINEKUBE AUTHSESSION API COMMING SOON] \
     -jar paper.jar
```

:::

In addition, ensure that you have set `enforce-secure-profile` to `false` in
your [server.properties](http://server.properties) file.

#### Lilypad

Set the following environment variable:

`LILYPAD_MOJANG_SESSIONSERVER_URL` to `[MINEKUBE AUTHSESSION API COMMING SOON]`

#### BungeeCord

This proxy type is currently not supported due
to [an issue in BungeeCord](https://github.com/SpigotMC/BungeeCord/pull/3201); please use Gate, Velocity or Waterfall
instead.

## Compatibility

Only the listed software allows changing the Mojang AuthSession API url.
Bungeecord is not supported.