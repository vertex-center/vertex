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

> [!NOTE]
> Everything about the installation process is available in the [Vertex Documentation](https://docs.vertex.arra.red/).

Vertex can be installed easily using Docker. We provide three different infrastructure setups:

- **Bundle**: Includes all the services in a minimal number of containers. This is the recommended setup for small deployments.
- **Microservices**: A more complex setup that separates the services into different containers. This is the recommended setup for large deployments, if you want to scale the services independently, or if you want more reliability.
- **Development**: A setup for development purposes. This is the recommended setup if you want to contribute to the project.

### Method A - Bundle

1. Be sure to have [Docker](https://docs.docker.com/get-docker/) installed and running.

2. Download the [docker-compose.yml](https://github.com/vertex-center/vertex/blob/main/docker/docker-compose.yml) and the [setup_postgres.sh](https://github.com/vertex-center/vertex/blob/main/docker/multidb/setup_postgres.sh) files.

   You should have the following directory structure:

    ```plaintext
    docker-compose.yml
    multidb/
        setup_postgres.sh
    ```

3. In a terminal, in the same directory as the `docker-compose.yml` file, run the following command:

    ```bash
    docker-compose up
    ```

4. Open [http://localhost:6133](http://localhost:6133) in your browser and start using Vertex!

### Method B - Microservices

1. Be sure to have [Docker](https://docs.docker.com/get-docker/) installed and running.

2. Download the [micro.docker-compose.yml](https://github.com/vertex-center/vertex/blob/main/docker/micro.docker-compose.yml) file.

3. In a terminal, in the same directory as the `micro.docker-compose.yml` file, run the following command:

    ```bash
    docker-compose -f micro.docker-compose.yml up
    ```

4. Open [http://localhost:7518](http://localhost:7518) in your browser and start using Vertex!

### Method C - Install for development

1. Be sure to have [Docker](https://docs.docker.com/get-docker/) installed and running.

2. Clone the repository:

    ```bash
    git clone https://github.com/vertex-center/vertex
    cd vertex
    ```
   
3. Run the following command:

    ```bash
    make run-dev
    ```

4. Open [http://localhost:5173](http://localhost:5173) in your browser and start using Vertex!

## License

[Vertex](https://github.com/vertex-center/vertex) is released under the [MIT License](./LICENSE.md).
