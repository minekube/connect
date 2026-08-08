---
description: Connect a Minecraft server and let authenticated Java and Bedrock players join through one public address.
---

# Quick Start

_This page explains the steps to use Connect with your Minecraft
server or network._

:::: info Prerequisites

- You have a Minecraft server running anywhere (locally or remote).
- You have a Java client, a Bedrock client signed in to Microsoft/Xbox, or an offline Java client for testing an
  endpoint that explicitly allows offline players.
::::

## Go to the Connect Dashboard! <VPBadge>Optional</VPBadge>

Proceed the `Quick Start` from there:

-> https://app.minekube.com

If you want to use the Connect Dashboard after the set-up, follow the steps below.
You can later import your Minecraft Endpoints to the dashboard to manage them.

## Step 1: Choose a Connector

Connectors are the link between your Minecraft Endpoints and the Connect Network.

-> [Go to Available Connectors](/guide/connectors/#available-connectors)

## Step 2: Launch your Connector

Follow the guide of the Connector you chose to install and launch it.
The first time the Connector is launched, the Endpoint will be protected with a token,
that you can reset in the Dashboard at any time.

-> If you face serious issues [get help](#getting-help) in the Discord.

## Done!

Your Minecraft server is now protected by the Connect Network and players can join it
at `<endpoint>.play.minekube.net`!

Java and Bedrock players use the same endpoint address. Bedrock support is handled by
the Connect edge for Connect-routed players, so Paper/Spigot, Velocity, and BungeeCord
plugin users do not need to install Geyser or change Gate Bedrock settings.

::: info What works with the normal setup

- Authenticated Java players can join.
- Microsoft/Xbox-authenticated Bedrock players can join with or without a linked Java account.
- Offline/cracked Java players require the endpoint owner to opt in.
- Bedrock players without valid Microsoft/Xbox authentication cannot join the managed Bedrock path.

See [Joining Servers](/guide/joining#who-can-join) for the complete address, port, identity, and configuration matrix.
:::

**Next Steps:** Consider to import your Endpoints to the Connect Dashboard.

## Getting Help

If you need help, join the Minekube Community Discord
and post a support request in the **#support** forum.

-> https://minekube.com/discord

::: tip Endpoint Token

If your get an authentication error from the Watch service, then try to reset the token
in the Dashboard and update it in your `connect.json` / `token.json` file of your Connector.

:::
