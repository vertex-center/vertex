import { Instance, InstanceQuery, Instances } from "../../models/instance";
import { Env, Service } from "../../models/service";
import { DockerContainerInfo } from "../../models/docker";
import { server } from "./backend";

const BASE_URL = "/app/vx-instances";

const getAllInstances = async () => {
    const { data } = await server.get<Instances>(`${BASE_URL}/instances`);
    return data;
};

const searchInstances = async (query: InstanceQuery) => {
    const { data } = await server.get<Instances>(
        `${BASE_URL}/instances/search`,
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

const getInstance = async (id: string) => {
    const { data } = await server.get<Instance>(`${BASE_URL}/instance/${id}`);
    return data;
};

const deleteInstance = (id: string) => {
    return server.delete(`${BASE_URL}/instance/${id}`);
};

const startInstance = (id: string) => {
    return server.post(`${BASE_URL}/instance/${id}/start`);
};

const stopInstance = (id: string) => {
    return server.post(`${BASE_URL}/instance/${id}/stop`);
};

const patchInstance = (id: string, params: any) => {
    return server.patch(`${BASE_URL}/instance/${id}`, params);
};

const getLogs = async (id: string) => {
    const { data } = await server.get(`${BASE_URL}/instance/${id}/logs`);
    return data;
};

const saveEnv = (id: string, env: Env) => {
    return server.patch(`${BASE_URL}/instance/${id}/environment`, env);
};

const getDocker = async (id: string) => {
    const { data } = await server.get<DockerContainerInfo>(
        `${BASE_URL}/instance/${id}/docker`
    );
    return data;
};

const recreateDocker = (id: string) => {
    return server.post(`${BASE_URL}/instance/${id}/docker/recreate`);
};

const updateService = (id: string) => {
    return server.post(`${BASE_URL}/instance/${id}/update/service`);
};

const getVersions = async (id: string, cache?: boolean) => {
    const { data } = await server.get<string[]>(
        `${BASE_URL}/instance/${id}/versions?reload=${!cache}`
    );
    return data;
};

const instancesRoutes = {
    all: getAllInstances,
    search: searchInstances,
};

const instanceRoutes = (id: string) => {
    return {
        get: () => getInstance(id),
        delete: () => deleteInstance(id),
        start: () => startInstance(id),
        stop: () => stopInstance(id),
        patch: (params: any) => patchInstance(id, params),
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

export const vxInstancesRoutes = {
    instances: instancesRoutes,
    instance: instanceRoutes,
    services: servicesRoutes,
    service: serviceRoutes,
};
