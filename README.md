<p align="center">
    <img height="96" src="https://github.com/vertex-center/vertex-design/raw/main/logos/transparent/vertex_logo_transparent.png" alt="Vertex logo" />
</p>
<h1 align="center">Vertex</h1>

<p align="center">
<img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/vertex-center/vertex?color=DE3C4B&labelColor=1E212B&style=for-the-badge">
<img alt="GitHub license" src="https://img.shields.io/github/license/vertex-center/vertex?color=DE3C4B&labelColor=1E212B&style=for-the-badge">
<img alt="GitHub contributors" src="https://img.shields.io/github/contributors/vertex-center/vertex?color=DE3C4B&labelColor=1E212B&style=for-the-badge">
</p>

## About

### Vertex

Vertex is a program that allows you to create your self-hosted server easily. Install, configure and start new services in less than a minute.

Vertex is designed to be installed on low-powered computers like Raspberry Pi, so avoiding systems like Docker and prioritizing compiled languages like Go, Rust or C++.

<img width="1884" alt="image" src="https://github.com/vertex-center/vertex/assets/12123721/a2b59a98-3b0d-4323-9db3-2078683eaf90">

<img width="1884" alt="image" src="https://github.com/vertex-center/vertex/assets/12123721/1b801d93-3ea0-4547-a6e6-4ccb3caf2c4c">

## Installation

### Method 1: From binaries

Binaries are released regularly. The latest release is available [here](https://github.com/vertex-center/vertex/releases/).

Decompress and execute the binary. Then, go to http://localhost:6130/. Enjoy!

### Method 2: Manual

1. Clone this repository
   ```bash
   git clone https://github.com/vertex-center/vertex
   cd vertex
   ```
2. Run Vertex
   ```bash
   go run .
   ```
3. Access from http://localhost:6130/

## Available services

All services are fetched from [Vertex Services](https://github.com/vertex-center/vertex-services). Installation can be done in one click, using Docker or package manager. Everything is handled automatically.

<img width="1884" alt="image" src="https://github.com/vertex-center/vertex/assets/12123721/30fef18d-ec3a-46e5-b73c-9c26d8a19a9f">

## License

[Vertex](https://github.com/vertex-center/vertex) is released under the [MIT License](./LICENSE.md).
