---
title: "Joining (shared section)"
pagefind: false
---
## Joining your server

Every server has a unique, configurable [endpoint](/guide/#connect-endpoints) name, which is also part of
the domain players join the server with.
If you do not set one, Connect uses a temporary random endpoint name
from the [Random Name Service](https://randomname.minekube.net/).

You can change the endpoint name in the Connector config at any time:

:::code-group
```yaml title="plugins/connect/config.yml"
endpoint: your-server-name
```

```yaml title="Gate config.yml"
connect:
  name: your-server-name
```
:::

> The environment variable `CONNECT_ENDPOINT` takes precedence over the configuration file.

### Joining with the free public domain

After you install the Connector and start your server, the console shows
the [free public domain](/guide/domains) for your server, which looks like
`<endpoint>.play.minekube.net`.

:::code-group
```shell title="Server Console"
[connect] Enabling connect vX.Y.Z
[connect] Enpoint name: live-beru
[connect] Your public address: live-beru.play.minekube.net
```
:::

Use that same endpoint hostname in either edition:

| Client | Server address | Port |
| --- | --- | --- |
| Java Edition | `<endpoint>.play.minekube.net` | `25565` (normally omitted) |
| Bedrock Edition | `<endpoint>.play.minekube.net` | `19132` |

Verified custom domains work the same way: Java and Bedrock use the same hostname and their edition's normal port.

When the Connector starts, the console prints the address:

```text title="Server console"
[connect] Enabling connect v0.2.13
[connect] Endpoint name: live-beru
[connect] Your public address: live-beru.play.minekube.net
```

Ping requests are forwarded to your server as well.

### Joining from the Browser Hub

Players can also discover your server from the in-game
[Browser Hub](/guide/advertising#browser-hub) at `minekube.net`
and can join with the in-game UIs or with the `/browser join <your-server-name>` command. Offline/cracked Java players
use `cracked.minekube.net` or `cracked.minekube.com` instead and can only enter endpoints that explicitly allow them.

See [Who can join](/guide/joining#who-can-join) for the Java, Bedrock, account-linking, and offline-mode identity matrix.
