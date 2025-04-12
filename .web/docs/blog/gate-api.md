---
layout: Post
title: 'Control Your Minecraft Proxy with the New Gate API'
category: Engineering
date: 2024-12-12
imageUrl: '/public/blog/gate-api/preview.png'
author:
  name: Gate Team
  role: Engineering
  href: 'https://github.com/minekube'
  imageUrl: 'https://github.com/minekube.png'
---

Building tools around your Minecraft proxy shouldn't require a PhD in plugin development. Today we're launching Gate API - a simple way to control your Gate proxy from any programming language. Want to move players between servers? Check server status? Build an admin panel? Now you can do it with a few lines of code.

## Where we started

Gate has always been about making Minecraft server management easier. But when it came to building tools around it, you had two options: write a plugin or parse log files. Neither was fun.

We wanted something better. Something that would let you:

- Build admin panels without wrestling with plugins
- Create Discord bots that actually work
- Automate server management without parsing logs
- Do it all in the language you love

## Enter Gate API

Gate API is our answer to these problems. It's a modern HTTP/gRPC API that speaks your language (literally - we support TypeScript, Python, Rust, and more).

Getting started is simple:

```yaml
api:
  enabled: true
  bind: localhost:8080
```

That's it. Your Gate proxy is now API-enabled.

## What can you build?

Let's look at some real examples. Here's what you can do with just a few lines of code:

### Managing Players

```typescript
// List all online players
const { players } = await client.listPlayers({});
console.log(`${players.length} players online`);

// Move a player to another server
await client.connectPlayer({
    player: "Notch",
    server: "minigames"
});
```

### Server Management

```typescript
// List servers and their status
const { servers } = await client.listServers({});
servers.forEach(server => {
    console.log(`${server.name}: ${server.players} players`);
});

// Add a new server
await client.registerServer({
    name: "minigames",
    address: "10.0.0.2:25565"
});
```

### Real Examples

Here's a simple example that shows server status:

```typescript
async function getServerStatus() {
    const { servers } = await client.listServers({});
    return servers.map(s => `${s.name}: ${s.players} players`).join('\n');
}
```

Want to auto-balance players across servers?

```typescript
async function balanceServers(from: string, to: string) {
    const { players } = await client.listPlayers({
        servers: [from]
    });

    // Move half the players
    const half = Math.floor(players.length / 2);
    for (let i = 0; i < half; i++) {
        await client.connectPlayer({
            player: players[i].username,
            server: to
        });
    }
}
```

### Web Apps

Yes, you can even use it in the browser (with proper security, of course):

```typescript
const client = createClient(GateService, createConnectTransport({
    baseUrl: 'https://your-secured-api.example.com'
}));

function ServerList() {
    const [servers, setServers] = useState([]);

    useEffect(() => {
        async function update() {
            const { servers } = await client.listServers({});
            setServers(servers);
        }
        update();
        return setInterval(update, 5000);
    }, []);

    return servers.map(server =>
        <div key={server.name}>
            {server.name}: {server.players} players
    </div>
);
}
```

⚠️ **Important**: Always secure your API before exposing it to the web!

## Running in the browser

Yes, you can even call Gate API directly from web apps! The API works seamlessly in browsers:

```typescript
import { createClient, createConnectTransport } from '@connectrpc/connect-web';
import { GateService } from '@buf/minekube_gate.connectrpc_es/gate/v1/gate_connect';

const client = createClient(
    GateService,
    createConnectTransport({
        baseUrl: 'http://localhost:8080',
    })
);
```

But remember: Gate API is powerful. It can do anything your proxy can do. Before exposing it to the web, put an auth proxy in front of it (unless you want random internet users controlling your servers).

## What's next?

This is just the beginning. We're working on:

- WebSocket support for real-time events
- More management endpoints
- Better monitoring capabilities

Want to try it out? Check out our [API documentation](https://gate.minekube.com/developers/api/). We can't wait to see what you build!

Got questions? Join us on [Discord](https://minekube.com/discord) - we're always around to help.
