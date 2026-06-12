# Bedrock Support

Connect lets Bedrock players join connected endpoints through the same public addresses Java players already use:

- `<endpoint>.play.minekube.net`
- custom domains attached to the endpoint
- endpoint routing through the Connect network

For Connect-routed players, Bedrock translation is handled by the Connect edge before traffic reaches your connector.
That means the usual Connect setup stays the same for Paper/Spigot, Velocity, BungeeCord, and Gate connectors.

## Which Setup Do I Need?

| Setup | What to configure | Who handles Bedrock translation |
| --- | --- | --- |
| Connect Plugin on Paper/Spigot, Velocity, or BungeeCord | Install and run the Connect plugin normally | Connect edge |
| Standard Gate with `connect` enabled | Enable Connect normally | Connect edge for Connect-routed players |
| Standard Gate with direct Bedrock clients | Add `bedrock: true` to Gate | Your Gate instance |
| Standard Gate with both Connect and direct Bedrock | Enable Connect and `bedrock: true` | Connect edge for Connect addresses, Gate for direct Bedrock address |
| Gate Lite behind Connect | Configure Gate Lite as a Java connector | Connect edge for Connect-routed Bedrock players |
| Gate Lite as a direct Bedrock listener | Use standard Gate instead | Not supported in Gate Lite |

## Connect-Routed Bedrock

Use this path when players join through Connect addresses such as `<endpoint>.play.minekube.net`, `minekube.net`, or a
custom domain attached to your endpoint.

You do not need to:

- install Geyser on your backend
- install Floodgate on your backend
- open a local UDP Bedrock port
- set `bedrock: true` in your backend Gate just for Connect players

This applies to the Connect Java plugin and to Gate used as a Connect connector.

## Direct Self-Hosted Gate Bedrock

Use this path when Bedrock players connect directly to a Gate address you operate yourself, outside the Connect network.
In that case, standard Gate can manage Bedrock locally:

```yaml
bedrock: true
```

With that setting, Gate starts and manages the Bedrock translation runtime for direct Bedrock clients. This local setting
does not change how Connect-routed Bedrock players are handled.

## Domains

Connect domains are shared by Java and Bedrock players. If a custom domain already routes Java players to your endpoint,
Bedrock players use the same domain after it is attached to the endpoint in the dashboard.

Use a separate direct Bedrock address only when you intentionally run local Bedrock support on your own standard Gate
instance.

## Account Linking and Online Mode

Connect-routed Bedrock players arrive through the managed Bedrock path and are represented through the compatibility
layer before they reach your connector. Servers that require a linked Java account may still reject Bedrock players who
have not linked their account. In that case, the player needs to link their Bedrock account with the Java account expected
by the target server.

If you are troubleshooting online-mode behavior, first identify whether the player joined through a Connect address or a
direct self-hosted Gate Bedrock address. The required configuration is different for those two paths.
