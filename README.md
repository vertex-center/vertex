<p align="center">
    <img height="96" src="https://github.com/vertex-center/vertex-design/raw/main/logos/transparent/vertex_logo_transparent.png" alt="Vertex logo" />
</p>
<h1 align="center">Vertex</h1>

<p align="center">
<img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/vertex-center/vertex?color=DE3C4B&labelColor=1E212B&style=for-the-badge">
<img alt="GitHub license" src="https://img.shields.io/github/license/vertex-center/vertex?color=DE3C4B&labelColor=1E212B&style=for-the-badge">
<img alt="GitHub contributors" src="https://img.shields.io/github/contributors/vertex-center/vertex?color=DE3C4B&labelColor=1E212B&style=for-the-badge">
</p>

---

> [!WARNING]
> Vertex is currently in development. While it is usable, and I am using it in production, there may be breaking changes before the v1.0 release. You can keep following the project on [Discord](https://discord.gg/tGZV6X6ZJh), or by starring the repository.

---

## About

### Vertex

Vertex is a program that allows you to create your self-hosted server easily. Install, configure and start new services in less than a minute.

Vertex is designed to be installed on low-powered computers like Raspberry Pi, so prioritizing compiled languages like Go, Rust or C++.

<img src="https://github.com/vertex-center/docs/assets/12123721/abbce3bc-01ef-4d86-b79e-0eefe57e08ce" alt="Vertex screenshot" />
<img src="https://github.com/vertex-center/docs/assets/12123721/f0cfe161-e015-4eee-86fc-ffffb9235d4e" alt="Vertex screenshot" />

## Features

- Install containers in one click from templates!
- Manage your containers easily (env, ports, etc)
- Receive alerts on Discord when a container is down
- Easy setup on Kubernetes with [Helm Charts](https://github.com/vertex-center/charts)
- _And more to come! (Database, Monitoring, etc)_

## Installation

Vertex can be installed in multiple ways:
- **From binaries**: Pre-compiled binaries are available for Linux, macOS and Windows
- **With Docker**: Docker images are available on GitHub Container Registry

Each of these methods can be used in two modes:
- **Bundled**: All apps are included in one unique binary (or container). This is the simplest and lightest way to install Vertex.
- **Microservices**: Each app is a separate binary (or container). This is the recommended way if you want to scale your installation horizontally.

Everything about the installation process is available in the [Vertex Documentation](https://docs.vertex.arra.red/).

## License

[Vertex](https://github.com/vertex-center/vertex) is released under the [MIT License](./LICENSE.md).
