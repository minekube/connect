# Server Protections

Server protections are a set of features that protect your server from malicious players and attacks.

## DDoS Protection

When hiding your server behind Connect, the network tries to protect your server from DDoS attacks by
using a distributed network of proxies to route player traffic through.
This means that if one of the proxies in the network is under attack,
other proxies will take over real player traffic and protect your server from being overloaded.

## Bot Protection

Connect uses a set of heuristics to detect bots and malicious players.
If a player is detected as a bot, they will be disconnected from the network.

Methods used to detect bots:
- Player movement
- Human verification captcha
- Heuristic to check for sending too many packets per second

One great method to make it harder for spambots is to enable online-mode for your server to 
only allow players to join your server if they have a valid Minecraft account.

::: tip

Protecting your server from bots is a never ending battle.
If you have a good idea on how to improve the bot protection, please let us know on our [Discord](https://minekube.com/discord).

:::

::: warning This feature is not yet available.
:::
