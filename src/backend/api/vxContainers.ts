import { Container, ContainerQuery, Containers } from "../../models/container";
import { Env, Service } from "../../models/service";
import { DockerContainerInfo } from "../../models/docker";
import { server } from "./backend";

const BASE_URL = "/app/vx-containers";

const getAllContainers = async () => {
    const { data } = await server.get<Containers>(`${BASE_URL}/containers`);
    return data;
};

const getAllTags = async () => {
    const { data } = await server.get<string[]>(`${BASE_URL}/containers/tags`);
    return data;
};

const searchContainers = async (query: ContainerQuery) => {
    const { data } = await server.get<Containers>(
        `${BASE_URL}/containers/search`,
        { params: query }
    );
    return data;
};

const installService = async (serviceId: string) => {
    const { data } = await server.post(
        `${BASE_URL}/service/${serviceId}/install`
    );
    return data;
};

const getAllServices = async () => {
    const { data } = await server.get<Service[]>(`${BASE_URL}/services`);
    return data;
};

const getContainer = async (id: string) => {
    const { data } = await server.get<Container>(`${BASE_URL}/container/${id}`);
    return data;
};

const deleteContainer = (id: string) => {
    return server.delete(`${BASE_URL}/container/${id}`);
};

const startContainer = (id: string) => {
    return server.post(`${BASE_URL}/container/${id}/start`);
};

const stopContainer = (id: string) => {
    return server.post(`${BASE_URL}/container/${id}/stop`);
};

const patchContainer = (id: string, params: any) => {
    return server.patch(`${BASE_URL}/container/${id}`, params);
};

const getLogs = async (id: string) => {
    const { data } = await server.get(`${BASE_URL}/container/${id}/logs`);
    return data;
};

const saveEnv = (id: string, env: Env) => {
    return server.patch(`${BASE_URL}/container/${id}/environment`, env);
};

const getDocker = async (id: string) => {
    const { data } = await server.get<DockerContainerInfo>(
        `${BASE_URL}/container/${id}/docker`
    );
    return data;
};

const recreateDocker = (id: string) => {
    return server.post(`${BASE_URL}/container/${id}/docker/recreate`);
};

const updateService = (id: string) => {
    return server.post(`${BASE_URL}/container/${id}/update/service`);
};

const getVersions = async (id: string, cache?: boolean) => {
    const { data } = await server.get<string[]>(
        `${BASE_URL}/container/${id}/versions?reload=${!cache}`
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
            get: () => getVersions(id),
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
