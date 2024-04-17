# Minekube Browser - API Authentication

## Overview
Authentication in the Minekube Browser API is handled through API keys. This section explains how to obtain and use your API keys to access the API securely.

## Generating API Keys
1. **Login to Dashboard:** Access the [Minekube Dashboard](https://dashboard.minekube.com) with your user credentials.
2. **Navigate to API Section:** Locate the API management section.
3. **Create New Key:** Follow the prompts to generate a new API key.
4. **Set Permissions:** Assign the necessary permissions based on your application needs.

## Using API Keys
- **Include in Headers:** Always include your API key in the request headers:
  ```http
  Authorization: Bearer <your_api_key>

### Endpoint Protection
- **Unauthorized Access:** All protected endpoints will return a `401 Unauthorized` response if a valid API key is not included in the request headers.

## Rate Limiting
- **Purpose:** To prevent abuse and ensure fair usage, API requests are rate-limited. View specific rate limit details in the [Rate Limits](ratelimits.md) section.
- **Limits:** Rate limits are defined by the number of requests per minute, varying by API key.

## Security Best Practices:
::: tip

**Ensure secure storage of your API keys to prevent unwanted abuse. View our [Security Guide](../../guide/protections.html) for more info**
:::
### 1. Key Rotation
- **Security Recommendation:** Periodic rotation of API keys is advised to maintain security.
- **Regeneration:** Keys can be regenerated through the dashboard as needed.

### 2. IP Whitelisting
- **Added Security:** API keys can be associated with specific IP addresses or ranges, restricting access to whitelisted IPs only.

### 3. Fine Tuned Access
- **Access Control:** API keys are assigned different permission levels to control access to various endpoints or features within your [Minekube Organization](../../guide/quick-start.md).

## Next Steps
- **API Versioning:** Detailed documentation on versioning for proper use of API keys is available in the [Versioning](./versions.md) section.

- **API Endpoints:** View available endpoints and example responses in [Endpoints](./endpoints.md) section.
