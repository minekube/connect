# Brower API Endpoints

The Minekube Browser API provides a range of endpoints to interact with the platform:

## Server Listing

### `GET /servers`
- **Description:** Retrieve a list of all public servers.

### `GET /servers/{id}`
- **Description:** Get details of a specific server by its ID.

## User Authentication

### `POST /auth/login`
- **Description:** Authenticate a user and retrieve an access token.

### `GET /auth/me`
- **Description:** Get the currently authenticated user's details.

### `PUT /auth/me`
- **Description:** Update the currently authenticated user's details.

## Server Management

### `GET /my-servers`
- **Description:** Retrieve a list of servers owned by the authenticated user.

### `GET /my-servers/{id}`
- **Description:** Get details of a specific server owned by the authenticated user.

### `PUT /my-servers/{id}`
- **Description:** Update the details of a specific server owned by the authenticated user.

## Player Statistics

### `GET /stats/players`
- **Description:** Retrieve overall player statistics, including total players, online players, and unique players.

### `GET /stats/players/top`
- **Description:** Get a list of the top servers by player count.

## Error Responses

- **400 Bad Request:** Invalid request format or missing parameters.
- **401 Unauthorized:** Authentication failed or invalid API key.
- **403 Forbidden:** Access forbidden due to insufficient permissions.
- **404 Not Found:** Requested resource not found.
- **500 Internal Server Error:** An unexpected server error occurred.
