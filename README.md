<a name="readme-top"></a>

<!--
*** Thanks for checking out the Gate Plugin Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]

<br />
<div align="center">
  <a href="https://github.com/minekube/gate-plugin-template">
    <img src="https://raw.githubusercontent.com/minekube/gate-plugin-template/main/assets/hero.png" alt="Logo" width="128" height="128">
  </a>

<h3 align="center">Gate Ping Protocol</h3>

  <p align="center">
    A Gate proxy plugin that customizes server ping responses with configurable protocol windows.
    <br />
    <br />
    <a href="https://gate.minekube.com/developers/"><strong>Gate docs »</strong></a>
    ·
    <a href="https://minekube.com/discord">Discord</a>
    ·
    <a href="https://github.com/minekube/gate/issues">File an issue</a>
  </p>
</div>

## Overview

Gate Ping Protocol hooks into the [Minekube Gate](https://github.com/minekube/gate) `PingEvent` and rewrites the version banner that Minecraft clients see in the server list. The list of supported protocol numbers and friendly version names is driven entirely by `protocol.yml`, making it easy to advertise exactly the versions your network supports without touching code.

### Highlights

- Keeps the reported protocol number in sync with the connecting client when possible, preventing red "Incompatible" warnings.
- Builds readable version ranges automatically (for example `1.19.3-1.20.1`).
- Falls back to Mojang defaults or `v<protocol>` labels when a custom name is omitted.
- Logs and skips any malformed entries so a single typo does not crash the proxy.

## Quick Start

1. **Clone** the repo: `git clone <repo-url>` and `cd gate_ping_protocol`.
2. **Copy the sample config**: `cp protocol.example.yml protocol.yml` (or start from scratch).
3. **Describe your supported versions** inside `protocol.yml` (details below).
4. **Run Gate** with the plugin enabled: `go run .` (add `-d` for verbose debugging).
5. **Ping from your Minecraft client** and verify that the server list now advertises the configured versions.

> The project ships with `config.yml` for a minimal proxy setup and `Makefile` helpers such as `make lint` and `make test` for local CI parity.

## Configuration (`protocol.yml`)

Example: (`protocol.example.yml`)

```yaml
protocols:
  - number: 763
    names: ["1.20.1"]
  - number: 762
    names: ["1.19.4"]
  - number: 761
    names: ["1.19.3"]
```

| Field   | Type      | Required | Description |
|---------|-----------|----------|-------------|
| `number` | integer  | ✔ | Minecraft protocol number. Leaving it out skips the entry and logs a warning. |
| `names`  | string[] | ✖ | Friendly labels shown to players. When empty, the plugin uses Gate's built-in mapping or falls back to `v<number>`. |

**Ordering matters:** the plugin iterates from top to bottom to build contiguous ranges. List protocols from oldest to newest (ascending numbers) so the range formatter can detect gaps correctly.

### How the plugin responds

1. When a client pings, Gate emits a `PingEvent` that the plugin intercepts.
2. The response banner (`ping.Version`) is rewritten with:
   - `Name`: the formatted range string derived from `protocols` (e.g., `1.19.3-1.20.1`).
   - `Protocol`: the client's own protocol number if it exists in the configured list; otherwise, the first supported protocol to keep the server selectable.
3. If the client requests a protocol outside the configured set, it still sees the advertised range but is assigned the first supported protocol so the proxy remains joinable.

## Running & Development

- `go run .` — start Gate with the plugin.
- `go run . -d` — enable debug logging (handy for watching ping negotiations).
- `make lint`, `make test` — optional helpers for CI parity.
- `Dockerfile` — build a container image (`docker build -t gate-ping-protocol .`).

### Docker Image

Use the published container if you prefer running Gate without installing Go locally:

```bash
docker pull ghcr.io/watermelon1024/gate-ping-protocol:latest
```

Example run command (mount your `protocol.yml` and expose Gate's default port):

```bash
docker run --rm \
  -p 25565:25565 \
  -v "$(pwd)/protocol.yml:/app/protocol.yml:ro" \
  ghcr.io/watermelon1024/gate-ping-protocol:latest
```

Customize the port mapping or config volume path as needed for your deployment.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

[contributors-shield]: https://img.shields.io/github/contributors/minekube/gate.svg?style=for-the-badge

[contributors-url]: https://github.com/minekube/gate/graphs/contributors

[forks-shield]: https://img.shields.io/github/forks/minekube/gate-plugin-template.svg?style=for-the-badge

[forks-url]: https://github.com/minekube/gate-plugin-template/network/members

[stars-shield]: https://img.shields.io/github/stars/minekube/gate.svg?style=for-the-badge

[stars-url]: https://github.com/minekube/gate-plugin-template/stargazers

[issues-shield]: https://img.shields.io/github/issues/minekube/gate.svg?style=for-the-badge

[issues-url]: https://github.com/minekube/gate-plugin-template/issues

[license-shield]: https://img.shields.io/github/license/minekube/gate.svg?style=for-the-badge

[license-url]: https://github.com/minekube/gate/blob/master/LICENSE

[product-screenshot]: https://github.com/minekube/gate/raw/master/.web/docs/public/og-image.png