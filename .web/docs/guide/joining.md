## Joining your Server

Every server has a unique configurable endpoint name that directly reflects
the domain players can join the server with.
If you leave this field empty, Connect will use a temporary random endpoint name
for your server provided by our [random name service](https://randomname.minekube.net/).

You can always update that endpoint name in `plugins/connect/config.yml`

```yaml 
endpoint: your-server-name
```

> The environment variable `CONNECT_ENDPOINT` takes precedence over the endpoint value
> in the configuration file.

### Joining with free provided Public Domain

After installing Connect plugin and starting your server
your will see the free public domain for your server that looks like
`<your server name>.play.minekube.net`.

![Console showing public domain](/images/terminal-log.png)

Ping requests are also mirrored to the endpoint server.


### Joining from Browser Hub

Players can also discover your server from the in-game
[Browser Hub](advertising#browser-hub) at `minekube.net`
and can join with the in-game UIs or with the `/browser join <your-server-name>` command.

### Joining from your custom domain

Go to [Server Domains](domains) to learn how to use your own domain for your server.
