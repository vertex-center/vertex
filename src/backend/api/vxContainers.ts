import { Container, ContainerQuery, Containers } from "../../models/container";
import { Env, Service } from "../../models/service";
import { DockerContainerInfo } from "../../models/docker";

import { createServer } from "../server";

// @ts-ignore
const server = createServer(window.api_urls.containers);

const getAllContainers = async () => {
    const { data } = await server.get<Containers>(`/containers`);
    return data;
};

const getAllTags = async () => {
    const { data } = await server.get<string[]>(`/containers/tags`);
    return data;
};

const searchContainers = async (query: ContainerQuery) => {
    const { data } = await server.get<Containers>(`/containers/search`, {
        params: query,
    });
    return data;
};

const installService = async (serviceId: string) => {
    const { data } = await server.post(`/service/${serviceId}/install`);
    return data;
};

const getAllServices = async () => {
    const { data } = await server.get<Service[]>(`/services`);
    return data;
};

const getContainer = async (id: string) => {
    const { data } = await server.get<Container>(`/container/${id}`);
    return data;
};

const deleteContainer = (id: string) => {
    return server.delete(`/container/${id}`);
};

const startContainer = (id: string) => {
    return server.post(`/container/${id}/start`);
};

const stopContainer = (id: string) => {
    return server.post(`/container/${id}/stop`);
};

const patchContainer = (id: string, params: any) => {
    return server.patch(`/container/${id}`, params);
};

const getLogs = async (id: string) => {
    const { data } = await server.get(`/container/${id}/logs`);
    return data;
};

const saveEnv = (id: string, env: Env) => {
    return server.patch(`/container/${id}/environment`, env);
};

const getDocker = async (id: string) => {
    const { data } = await server.get<DockerContainerInfo>(
        `/container/${id}/docker`
    );
    return data;
};

const recreateDocker = (id: string) => {
    return server.post(`/container/${id}/docker/recreate`);
};

const updateService = (id: string) => {
    return server.post(`/container/${id}/update/service`);
};

const getVersions = async (id: string, cache?: boolean) => {
    const { data } = await server.get<string[]>(
        `/container/${id}/versions?reload=${!cache}`
    );
    return data;
};

const containersRoutes = {
    all: getAllContainers,
    tags: getAllTags,
    search: searchContainers,
};

const containerRoutes = (id: string) => {
    return {
        get: () => getContainer(id),
        delete: () => deleteContainer(id),
        start: () => startContainer(id),
        stop: () => stopContainer(id),
        patch: (params: any) => patchContainer(id, params),
        logs: {
            get: () => getLogs(id),
        },
        env: {
            save: (env: Env) => saveEnv(id, env),
        },
        docker: {
            get: () => getDocker(id),
            recreate: () => recreateDocker(id),
        },
        update: {
            service: () => updateService(id),
        },
        versions: {
            get: (cache?: boolean) => getVersions(id, cache),
        },
    };
};

const servicesRoutes = {
    all: getAllServices,
};

const serviceRoutes = (service_id: string) => ({
    install: () => installService(service_id),
});

export const vxContainersRoutes = {
    containers: containersRoutes,
    container: containerRoutes,
    services: servicesRoutes,
    service: serviceRoutes,
};
