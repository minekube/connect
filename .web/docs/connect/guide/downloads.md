## Downloading the Connect Plugin

The installation steps are the same as for every other Minecraft plugin
_(Download jar, put in plugins folder, start server)_. The only difference is that
you can download the Connect plugin right here, instead of SpigotMC.


| Downloads                                                                                              |                                                                                                    |                                                                                                    |
|--------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------|
| [Spigot/PaperMC](https://github.com/minekube/connect-java/releases/download/latest/connect-spigot.jar) | [Velocity](https://github.com/minekube/connect-java/releases/download/latest/connect-velocity.jar) | [BungeeCord](https://github.com/minekube/connect-java/releases/download/latest/connect-bungee.jar) |

Ready to experience Minekube Connect? Download the latest stable release for your platform!
:point_up_2:

## Disabling `enforce secure player profiles`

::: details Why is this necessary?

Since Minecraft 1.19 the `enforce-secure-profile` property was introduced.
Players joining through the Connect network to your server won't be able to join if this setting
is enabled. It is safe to disable this setting as it only affects chat messages.

:::

::: warning Spigot/PaperMC's `server.properties`

You need to disable `enforce-secure-profile` property

```text 
...
enable-rcon=false
online-mode=true // [!code focus]
enforce-secure-profile=false // [!code focus]
enable-status=true
...
```

or simply disable `online-mode` which also disables `enforce-secure-profile`

```text 
...
enable-rcon=false
online-mode=false // [!code focus]
enforce-secure-profile=true
...
```

:::


- If using Velocity set `force-key-authentication` to `false` in `velocity.toml`.
- If using BungeeCord set `enforce_secure_profile` to `false` in `config.yml`.
