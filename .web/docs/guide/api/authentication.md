# Authentication

_Clients and [Endpoints](/guide/#connect-endpoints) authenticate to Connect API
with an API token and the Endpoint name so that Connect knows who is making requests
and what permissions you have._

---

If you are using the Connect API through the [Java API provided by Connect Plugin](/guide/api/clients#provided-by-connect-plugin)
then you do not write authentication code.

## Required Headers

All requests to the Connect API require the following headers:
- `Connect-Endpoint` - The name of the endpoint the token belongs to.
- `Authorization: Bearer <ENDPOINT_TOKEN>` - The token of the endpoint you are connecting to.

## Endpoint Names

Connect has the concept of globally unique **endpoint names** to identify your server even after restarts.
The [Connectors](/guide/connectors/) use a token file to authenticate that you
own an **endpoint name** in the [Connect Network](/guide/#the-connect-network).

::: code-group
```json [plugins/connect/token.json]
{"token":"T-ozinikukmabrpyzogjjl"}
```
:::

The token and endpoint name have a direct relationship.
If you lose your `token.json` file and haven't imported your endpoint to the
[Connect Dashboard](https://app.minekube.com), your endpoint
name is lost, and you will have to use a new name. If you have imported your endpoint
you can always reset your token in the dashboard.

## Super Endpoints

Checkout [Super Endpoints guide](/guide/api/super-endpoints) for authorizing other endpoints to act on your Endpoint's behalf.
