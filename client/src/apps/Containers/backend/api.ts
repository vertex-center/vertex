import { createServer } from "../../../backend/server";
import {
    Container,
    ContainerFilters,
    Containers,
    EnvVariables,
    Port,
    Tags,
} from "./models";
import { DockerContainerInfo } from "../../../models/docker";
import { Template } from "./template";

// @ts-ignore
const server = createServer(window.api_urls.containers);

const getContainers = async (query?: ContainerFilters) => {
    const { data } = await server.get<Containers>(`/containers`, {
        params: query,
    });
    return data;
};

const getAllTags = async () => {
    const { data } = await server.get<Tags>(`/tags`);
    return data;
};

const getAllTemplates = async () => {
    const { data } = await server.get<Template[]>(`/templates`);
    return data;
};

const getContainer = async (id: string) => {
    const { data } = await server.get<Container>(`/containers/${id}`);
    return data;
};

export type CreateContainerOptions = {
    template_id?: string;
    image?: string;
    image_tag?: string;
};

const createContainer = async (options: CreateContainerOptions) => {
    const { data } = await server.post(`/containers`, options);
    return data;
};

const deleteContainer = (id: string) => {
    return server.delete(`/containers/${id}`);
};

const startContainer = (id: string) => {
    return server.post(`/containers/${id}/start`);
};

const stopContainer = (id: string) => {
    return server.post(`/containers/${id}/stop`);
};

const patchContainer = (id: string, params: any) => {
    return server.patch(`/containers/${id}`, params);
};

const getLogs = async (id: string) => {
    const { data } = await server.get(`/containers/${id}/logs`);
    return data;
};

const getContainerEnvironment = async (id: string) => {
    const { data } = await server.get<EnvVariables>(
        `/containers/${id}/environment`
    );
    return data;
};

const saveEnv = (id: string, env: EnvVariables) => {
    return server.patch(`/containers/${id}/environment`, { env });
};

const getContainerPorts = async (id: string) => {
    const { data } = await server.get(`/containers/${id}/ports`);
    return data;
};

const saveContainerPorts = (id: string, ports: Port[]) => {
    return server.patch(`/containers/${id}/ports`, { ports });
};

const getDocker = async (id: string) => {
    const { data } = await server.get<DockerContainerInfo>(
        `/containers/${id}/docker`
    );
    return data;
};

const recreateDocker = (id: string) => {
    return server.post(`/containers/${id}/docker/recreate`);
};

const updateService = (id: string) => {
    return server.post(`/containers/${id}/update/service`);
};

const getVersions = async (id: string, cache?: boolean) => {
    const { data } = await server.get<string[]>(
        `/containers/${id}/versions?cache=${cache}`
    );
    return data;
};

export const API = {
    getContainer,
    getContainers,
    getAllTags,
    createContainer,
    deleteContainer,
    startContainer,
    stopContainer,
    patchContainer,
    getLogs,
    getContainerEnvironment,
    saveEnv,
    getContainerPorts,
    saveContainerPorts,
    getDockerInfo: getDocker,
    recreateDocker,
    updateService,
    getVersions,
    getAllTemplates,
};
