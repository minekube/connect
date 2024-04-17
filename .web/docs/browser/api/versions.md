# Minekube Browser - API Versioning

## Overview
To ensure backward compatibility and smooth updates, the Minekube Browser API follows a clear versioning system.

## Version Numbering
API versions are indicated by a semantic version number, e.g., `v1`, `v2`. This number is included in the base URL of each API endpoint.

## Version Stability
- **Major Updates:** Transitioning from `v1` to `v2` may introduce breaking changes.
- **Minor Updates:** Upgrades such as `v1.0` to `v1.1` will remain backward-compatible.

## Deprecation Policy
Deprecated endpoints or features will be clearly marked in the documentation, including the expected removal date.

## Version Migration
When a new major version is released, a migration path will be provided to help users transition smoothly to the new version.

## Version Header
API requests can include an optional `Accept-Version` header to specify the desired API version. If not provided, the latest stable version will be used.

## Versioning in Responses
API responses may include an `API-Version` header, indicating the version used to process the request.

## Version-Specific Documentation
The API documentation will provide separate sections for each major version, outlining any changes or additions.

## Release Notes
Release notes will be provided for each API version, detailing any changes, improvements, or fixes.

## Sunset Policy
Deprecated versions will continue to be supported for a defined period, after which they may be sunset. Advance notice will be provided before sunsetting a version.

## Client Libraries
Official client libraries will be updated to support new API versions, and older versions will be maintained for a reasonable period.

## Version Negotiation
The API may support version negotiation in the future, allowing clients to specify their preferred version and receive a response in that version.
