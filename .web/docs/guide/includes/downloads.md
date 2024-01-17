## Downloading the Connect Plugin

The installation steps are the same as for every other Minecraft plugin
_(Download jar, put in plugins folder, start server)_. Download the Connect plugin right here!

| <VPBadge>Spigot/PaperMC Plugin</VPBadge>                                                                   | <VPBadge>Velocity Plugin</VPBadge>                                                                             | <VPBadge>BungeeCord Plugin</VPBadge>                                                                       |
|------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------|
| [connect-spigot.jar](https://github.com/minekube/connect-java/releases/download/latest/connect-spigot.jar) | [connect-velocity.jar](https://github.com/minekube/connect-java/releases/download/latest/connect-velocity.jar) | [connect-bungee.jar](https://github.com/minekube/connect-java/releases/download/latest/connect-bungee.jar) |

Ready to experience Minekube Connect? Download the latest stable release for your platform!
:point_up_2:

![Console showing public domain](/images/terminal-log.png)

[_What does the Connect Plugin do?_](/guide/#the-connect-plugin)

## Disabling "enforce secure player profiles" <VPBadge type='danger'>Required</VPBadge>

Since Minecraft 1.19 the `enforce-secure-profile` property was introduced.
Players joining through the [Connect Network](/guide/#the-connect-network) to your server won't be able to join if this
setting
is enabled. It is safe to disable this setting as it only affects chat messages.

::: code-group

```properties [server.properties <VPBadge>Spigot/PaperMC</VPBadge>]
enforce-secure-profile=false

# If you disable online-mode, then enforce-secure-profile has no effect
online-mode=true
```

```toml [velocity.toml <VPBadge>Velocity</VPBadge>]
force-key-authentication = false

# If you disable online-mode, then force-key-authentication has no effect
online-mode = true
```

```yaml [config.yml <VPBadge>BungeeCord</VPBadge>]
enforce_secure_profile=false

  # If you disable online-mode, then enforce-secure-profile has no effect
online_mode=true
```

:::
