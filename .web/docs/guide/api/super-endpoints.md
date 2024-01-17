# Super Endpoints

Super Endpoints are endpoints that are authorized to control another endpoint.
They are configured in the [Connect Plugin](/guide/connectors/plugin) config `plugins/connect/config.yml`.

## Authorization

To allow other [Endpoints](/guide/#connect-endpoints) to control your Endpoint you can add them to your
super endpoints list in your [Connect Plugin](/guide/connectors/plugin) config.

Super endpoints are authorized to control this endpoint via Connect API and can e.g. disconnect players from this endpoint,
send messages to players, etc. You can add as many super endpoint names as you want.

## Configuration Examples

### Authorize all

Let's say you have a `hub` server, a `survival` server, and a `pvp` server.
You own the config of all servers and want every server to be able to send players to each other.

::: code-group
```yaml [PvP Endpoint]
# As the pvp endpoint you want to allow the hub 
# and survival server to send players to you.
endpoint: my-pvp
super-endpoints:
 - my-hub
 - my-survival
```

```yaml [Hub Endpoint]
# As the hub endpoint you want to allow the pvp 
# and survival server to send players to you.
endpoint: my-hub
super-endpoints:
 - my-pvp
 - my-survival
```
```yaml [Survival Endpoint]
# As the survival endpoint you want to allow the pvp
# and hub server to send players to you.
endpoint: my-survival
super-endpoints:
 - my-pvp
 - my-hub
```
:::

![](https://mermaid.ink/svg/pako:eNqFkV1PwyAUhv8KObttk7kP64hZQst655Vemd7gYCsZFEKhWrf9d3F2ThPNLkj4eJ6Tw3n3sDZcAIatY7ZG-VPVINT6Xgmk-7QNrpMdU2gjlcKj2xWdFmUS353ZCTy6KQkl-fmcvkruazyxb79q1OFl0MuCTmf0opdZPh-Pr-i2s4NOFllWzC_63WRGFqu_9c81yPdpujyQ4Gvj5LtoDz__dQWLrQ_E9yD-xSABLZxmksdh7k89gK-FFhXguOXM7SqommPkWPDmsW_WgL0LIoFgOfOCShYz0IA3TLXxVnDpjXv4SucUUgKWNc_GnJnjB1LOlmk)
