---
title: "Downloads (shared section)"
pagefind: false
---
## Downloading the Connect plugin

Install it like any other Minecraft plugin: download the jar for your platform, put it in the `plugins` folder and
start the server.

<div class="mk-dl-cards not-content">
<a class="mk-dl-card" href="https://github.com/minekube/connect-java/releases/download/latest/connect-spigot.jar" download><img src="/images/logos/paper.png" alt="" width="32" height="32"><span class="t"><strong>Paper / Spigot</strong><small>connect-spigot.jar</small></span><span class="i"><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 4v11m0 0-4.5-4.5M12 15l4.5-4.5M5 19.5h14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg></span></a>
<a class="mk-dl-card" href="https://github.com/minekube/connect-java/releases/download/latest/connect-velocity.jar" download><img src="/images/logos/velocity.png" alt="" width="32" height="32"><span class="t"><strong>Velocity</strong><small>connect-velocity.jar</small></span><span class="i"><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 4v11m0 0-4.5-4.5M12 15l4.5-4.5M5 19.5h14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg></span></a>
<a class="mk-dl-card" href="https://github.com/minekube/connect-java/releases/download/latest/connect-bungee.jar" download><img src="/images/logos/spigot.png" alt="" width="32" height="32"><span class="t"><strong>BungeeCord</strong><small>connect-bungee.jar</small></span><span class="i"><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 4v11m0 0-4.5-4.5M12 15l4.5-4.5M5 19.5h14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg></span></a>
<a class="mk-dl-card rec" href="https://gate.minekube.com/guide/connect" target="_blank" rel="noreferrer"><img src="/minekube-logo.png" alt="" width="32" height="32"><span class="t"><strong>Gate</strong><small>Built in, no download</small></span><span class="i"><svg viewBox="0 0 24 24" aria-hidden="true"><path d="M7 17 17 7M9 7h8v8" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg></span></a>
</div>

Each card downloads the latest stable release. Once the server starts, the console prints your public address.

## Disabling "enforce secure player profiles" <Badge type="danger">Required</Badge>

Minecraft 1.19 introduced the `enforce-secure-profile` property.
Players joining your server through the [Connect Network](/guide/#the-connect-network) cannot join while this setting
is enabled. It is safe to disable, because it only affects chat messages.

:::code-group

```properties title="server.properties (Spigot/PaperMC)"
enforce-secure-profile=false

# If you disable online-mode, then enforce-secure-profile has no effect
online-mode=true
```

```toml title="velocity.toml (Velocity)"
force-key-authentication = false

# If you disable online-mode, then force-key-authentication has no effect
online-mode = true
```

```yaml title="config.yml (BungeeCord)"
enforce_secure_profile: false

# If you disable online-mode, then enforce-secure-profile has no effect
online_mode: true
```

:::
