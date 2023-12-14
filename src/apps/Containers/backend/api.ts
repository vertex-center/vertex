import { createServer } from "../../../backend/server";
import {
    Container,
    ContainerQuery,
    Containers,
    EnvVariables,
    Tags,
} from "./models";
import { DockerContainerInfo } from "../../../models/docker";
import { Service } from "./service";

// @ts-ignore
const server = createServer(window.api_urls.containers);

const getAllContainers = async () => {
    const { data } = await server.get<Containers>(`/containers`);
    return data;
};

const getAllTags = async () => {
    const { data } = await server.get<Tags>(`/tags`);
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

const saveEnv = (id: string, env: EnvVariables) => {
    return server.patch(`/container/${id}/environment`, { env });
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

export const API = {
    getContainer,
    getAllContainers,
    getAllTags,
    searchContainers,
    deleteContainer,
    startContainer,
    stopContainer,
    patchContainer,
    getLogs,
    saveEnv,
    getDockerInfo: getDocker,
    recreateDocker,
    updateService,
    getVersions,
    installService,
    getAllServices,
};
