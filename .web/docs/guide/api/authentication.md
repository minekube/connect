# Authentication

_Clients and [Endpoints](/guide/#connect-endpoints) authenticate to Connect API
with an API token and the Endpoint name so that Connect knows who is making requests
and what permissions you have._

---

If you are using the Connect API through the [Java API provided by Connect Plugin](/guide/api/#java-and-kotlin-client),
you do not write authentication code.

## Required Headers

All requests to the Connect API require the following headers:
- `Connect-Endpoint` - The name of the endpoint the token belongs to.
- `Connect-Token` - The token of the endpoint you are connecting to.

## Endpoint Names

Connect has the concept of globally unique **endpoint names** to identify your server even after restarts.
The [Connect Plugin](/guide/#the-connect-plugin) uses a token file to authenticate that you
own an **endpoint name** in the [Connect Network](/guide/#the-connect-network).

::: code-group
```json [plugins/connect/token.json]
{"token":"T-ozinikukmabrpyzogjjl"}
```
:::

The token and endpoint name have a direct relationship.
If you lose your `token.json` file your endpoint name is lost, and you will have to use a new name.

[//]: # (## Super Endpoints)

[//]: # ()
[//]: # (TODO)
