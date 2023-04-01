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

<img width="1515" alt="image" src="https://user-images.githubusercontent.com/12123721/229261331-b2c8de51-f88a-458d-9eff-0ea73c83a0ad.png">

Vertex is designed to be installed on low-powered computers like Raspberry Pi, so avoiding systems like Docker and prioritizing compiled languages like Go, Rust or C++.

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

<img width="822" alt="image" src="https://user-images.githubusercontent.com/12123721/229261568-399c5970-600b-4023-96b1-bddd7feba81e.png">

There are 3 methods to install new services:

1. **From Marketplace.** These are services officially [released](https://github.com/vertex-center/vertex-services) by the Vertex team. Simply select the service you want from the UI, and install.
2. **From Git.** For unofficial services, you can easily provide a link to the repository. Vertex will clone and start the project automatically.
3. **From Local storage.** This method is useful while contributing to Vertex. You can select the path of a service on your computer, and Vertex will handle it automatically.

With all these methods, Vertex aims to be a flexible way to manage all your services.

## License

[Vertex](https://github.com/vertex-center/vertex) is released under the [MIT License](./LICENSE.md).
