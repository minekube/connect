# Gate Proxy - Connector Guide

Using Gate as a Connector for your Minecraft server or proxy.
It is the best supported, most frequently updated Connector with the most capabilities.

-> If you haven't already, checkout [Gate's Quick Start](https://gate.minekube.com/guide/quick-start) first,
it's just a tiny binary.

**Table of Contents**
[[TOC]]

## Enable Connect in Gate configuration

First of all, in the Gate configuration, you:

1. enable Connect mode
2. choose your endpoint name

-> See [Enabling Connect](https://gate.minekube.com/guide/connect) page for how to do it.

## Gate (standard)

For an introduction to Gate, see [Gate Docs](https://gate.minekube.com/guide).

As one of the most capable open source Minecraft proxy in the world, Gate has first class support for Connect with
online-mode.

-> To customize Gate checkout the [starter plugin template](https://github.com/minekube/gate-plugin-template)

## Gate Lite Mode

For an introduction to Gate Lite mode, see [Gate Lite Docs](https://gate.minekube.com/guide/lite).

::: tip Coming Soon -> Allowing online mode in Gate Lite backend routes

`Connect -> Gate Lite -> Online Mode Backend`

```yaml
connect:
  enforcePassthrough: true
```

This feature will allow to offload Connect player authentication to your Lite backend routes,
by enabling pass-through connections.

:::

If your backend server behind Gate Lite is online mode and let Connect authenticate players upfront,
you need to use Connect's [AuthSession API](../auth-api)
to allow online mode players from Connect Network to join your backend server.

## Example Setups

Minekube Connect advances the way players connect and developers architect secure Minecraft servers and networks.
Let's take a look at some common example setups.

### #1 Example: Online mode Velocity

`Connect -> [ Gate Lite -> [ Velocity -> Papers ] ]`

> We could only use Gate -> Paper without Velocity, but we wanted to show how to use Connect with an existing proxy. (
> e.g. you need Velocity plugins)

- We have Velocity in online mode running on `localhost:25577` and want to use Connect through Gate.
- Because Velocity is a proxy itself, we need to use Gate Lite mode for simple reverse proxying.
- We enable `lite` in Gate config and add a wildcard route `*` to forward all traffic to Velocity at `localhost:25577`.
- We enable `connect` mode in Gate config and choose a name for our endpoint.
- To allow online mode players from Connect Network to join our online mode Velocity server, we run Velocity with
  `-Dmojang.sessionserver=` to use Connect's [AuthSession API](../auth-api).
- Done! We can now join to our Velocity server at `<endpoint>.play.minekube.net` and online mode players from Connect
  Network
  can join as well. All without using port forwarding nor the [~~Java Plugi~~n](plugin.md).

> Velocity here could have been replaced by another proxy like BungeeCord/Waterfall or even Gate standard
> (reduce Gate layers: try to set up Example #2 instead).

### #2 Example: Paper online or offline

`Connect -> [ Gate -> Papers ]`

- We have Paper in online or offline mode running on `localhost:25565` and want to use Connect through Gate.
- This time we chose to not enable `lite` mode in Gate config, because we want to use Gate's full proxying capabilities
  like switching servers with `/server`
- We enable `connect` mode in Gate config and choose a name for our endpoint.
- Done! We can now join to our Paper server at `<endpoint>.play.minekube.net` and online mode players from Connect
  Network can join as well thanks to Gate's first class Connect online-mode support. All without using port forwarding
  nor the [~~Java Plugi~~n](plugin.md) nor the [~~AuthSession API~~](../auth-api).

## Getting Support

If you have any questions or need help, simply post a support request in the Minekube Community Discord
forum!. https://minekube.com/discord