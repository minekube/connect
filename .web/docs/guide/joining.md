---
description: Addresses, ports, and authentication behavior for Java, Bedrock, and offline Java players using Connect.
---

# How to join your Connect Server

<!--@include: ./includes/joining.md-->

## Who Can Join

This matrix applies when players use a Connect address and the server owner runs a normal Connect plugin or standard
Gate connector:

| Player identity | Accepted by default | Identity used by the endpoint | Server-owner action |
| --- | --- | --- | --- |
| Java client authenticated with a Java Edition account | Yes | Verified Java UUID and profile | None beyond the normal Connect setup |
| Bedrock client authenticated with Microsoft/Xbox and authoritatively linked to Java | Yes | Linked Java UUID and profile | None |
| Bedrock client authenticated with Microsoft/Xbox but not linked to Java | Yes | Stable native Bedrock/XUID-derived profile | None |
| Offline/cracked Java client | No | Self-asserted offline Java name and UUID | Explicitly [allow offline-mode players](/guide/offline-mode#allowing-offline-java-on-an-endpoint) |
| Bedrock client without valid Microsoft/Xbox authentication | No | None; Connect rejects the session | Unsupported on Connect-managed Bedrock |

Account linking is optional for Connect-managed Bedrock. When Connect can verify an authoritative Java link, it
preserves that Java identity; otherwise, the player's verified Xbox identity is enough. A backend plugin that is
specifically written to require a premium Java UUID can still reject a native Bedrock profile, so diagnose that plugin
rather than enabling offline mode or installing another Geyser instance.

## What the Server Owner Configures

For authenticated Java and Microsoft/Xbox-authenticated Bedrock players, configure only the normal Connect connector.
Connect handles Bedrock translation at the edge. Do not install backend Geyser or Floodgate, open a backend UDP port,
or set Gate `bedrock: true` for Connect-routed players.

Offline/cracked Java is a separate, explicit opt-in. The [Offline Mode guide](/guide/offline-mode) shows the setting for
both the Connect Java Plugin and Gate.
