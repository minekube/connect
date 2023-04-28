# Authentication

_How do client, servers and players authenticate to Connect?_

## Required Headers

All requests to the Connect API require the following headers:
- `Connect-Endpoint` - The name of the endpoint the token belongs to.
- `Connect-Token` - The token of the endpoint you are connecting to.

## Unique Server Endpoint Name

Connect has the concept of unique endpoint names to identify your server even after restarts.
The Connect Plugin uses a token file to authenticate that you own an endpoint name in the Connect Network.

```json plugins/connect/token.json
{"token":"T-ozinikukmabrpyzogjjl"}
```

The token and endpoint name have a direct relationship.
If you lose your `token.json` file your endpoint name is lost, and you will have to use a new name.

[//]: # (## Super Endpoints)

[//]: # ()
[//]: # (TODO)
