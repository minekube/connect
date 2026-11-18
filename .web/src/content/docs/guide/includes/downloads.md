---
title: "Downloads (shared section)"
pagefind: false
---
## Downloading the Connect Plugin

The installation steps are the same as for every other Minecraft plugin
_(Download jar, put in plugins folder, start server)_. Download the Connect plugin right here!

| <Badge type="tip">Spigot/PaperMC Plugin</Badge>                                                                   | <Badge type="tip">Velocity Plugin</Badge>                                                                             | <Badge type="tip">BungeeCord Plugin</Badge>                                                                       |
|------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------|
| [connect-spigot.jar](https://github.com/minekube/connect-java/releases/download/latest/connect-spigot.jar) | [connect-velocity.jar](https://github.com/minekube/connect-java/releases/download/latest/connect-velocity.jar) | [connect-bungee.jar](https://github.com/minekube/connect-java/releases/download/latest/connect-bungee.jar) |

Ready to experience Minekube Connect? Download the latest stable release for your platform!
:point_up_2:

![Console showing public domain](/images/terminal-log.png)

## Disabling "enforce secure player profiles" <Badge type="danger">Required</Badge>

Since Minecraft 1.19 the `enforce-secure-profile` property was introduced.
Players joining through the [Connect Network](/guide/#the-connect-network) to your server won't be able to join if this
setting
is enabled. It is safe to disable this setting as it only affects chat messages.

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
enforce_secure_profile=false

  # If you disable online-mode, then enforce-secure-profile has no effect
online_mode=true
```

:::
