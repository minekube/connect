# About Connect Tunnels
<VPBadge>Networking knowledge recommended</VPBadge>

_Don't believe in magic? This page explains how [Connect Tunnels](/guide/#connect-tunnels) enable player connections beyond private networks._

[[toc]]

## Introduction

Publicly routable IP addresses are required to allow inbound connections from the Internet (outside your home network).
Your ISP (Internet Service Provider) assigns you a public IP address, but sometimes it's dynamic and changes over time.

Further, your router has a NAT (Network Address Translation) that allows multiple devices to share the same
public IP address in your home network. This is why you can have multiple devices connected to the Internet
at the same time, their private IP addresses are translated to your public IP address when communicating with the Internet.

**However, your home firewall blocks inbound connections by default, including Minecraft connections on port `25565`.**

## Connect to the Rescue

Connect bypasses this whole problem by establishing outbound only tunnel connections from your server
to publicly routable edge proxies in the [Connect Network](/guide/#the-connect-network).
This way, you don't have to configure port forwarding on your router.

Connect edge proxies accept inbound connections from players and routes them through encrypted tunnels to your server.
The [Connectors](/guide/connectors/) establish a tunnel connection to the edge proxy for every player that joins your
server.

Attentive readers are now wondering how the plugin knows when a player wants to join and from which edge proxy.
This is where the [Connect Session Service](/guide/#connect-sessions) comes in.
The plugin watches for player session proposals from this service and establishes a tunnel connection to the correct edge proxy.

## The Join Flow

This is how the join flow looks for every player:

1. Connect plugin watches for player session proposals
2. A player on an edge proxy wants to join your server and sends a session proposal
3. Your server accepts the proposal and establishes a tunnel to the edge proxy to take over the player session

This flow is quite fast and happens within a few milliseconds.
The player doesn't notice any difference and can play on the performant tunnel to your server as usual.

## Connectors

Check out the [Connectors Overview](/guide/connectors/) to learn more about Connectors.

---

**If you find these concepts intimidating, don't worry! The [Quick Start](/guide/quick-start) only requires basic
knowledge, and you should be able to follow along without being an expert in any of these.**

