import axios from "axios";
import {
    InstallMethod,
    Instance,
    InstanceQuery,
    Instances,
} from "../models/instance";
import { Env, Service } from "../models/service";
import { Dependencies } from "../models/dependency";
import { DockerContainerInfo } from "../models/docker";
import { About } from "../models/about";
import { Updates } from "../models/update";

type InstallServiceParams = {
    method: InstallMethod;
    service_id: string;
};

const server = axios.create({
    // @ts-ignore
    baseURL: `${window.apiURL}/api`,
});

export const api = {
    about: {
        get: () => server.get<About>("/about"),
    },

    instances: {
        get: () => server.get<Instances>("/instances"),
        search: (query: InstanceQuery) =>
            server.get<Instances>("/instances/search", { params: query }),
        checkForUpdates: () => server.get<Instances>("/instances/checkupdates"),
    },

    services: {
        install: (params: InstallServiceParams) =>
            server.post("/services/install", params),

        available: {
            get: () => server.get<Service[]>("/services/available"),
        },
    },

    instance: {
        get: (id: string) => server.get<Instance>(`/instance/${id}`),
        delete: (id: string) => server.delete(`/instance/${id}`),
        start: (id: string) => server.post(`/instance/${id}/start`),
        stop: (id: string) => server.post(`/instance/${id}/stop`),
        patch: (id: string, params: any) =>
            server.patch(`/instance/${id}`, params),

        logs: {
            get: (id: string) => server.get(`/instance/${id}/logs`),
        },

        env: {
            save: (id: string, env: Env) =>
                server.patch(`/instance/${id}/environment`, env),
        },

        dependencies: {
            get: (id: string) =>
                server.get<Dependencies>(`/instance/${id}/dependencies`),
        },

        docker: {
            get: (id: string) =>
                server.get<DockerContainerInfo>(`/instance/${id}/docker`),
            recreate: (id: string) =>
                server.post(`/instance/${id}/docker/recreate`),
        },
    },

    packages: {
        install: (packages: any) =>
            server.post(`/packages/install`, { packages }),
    },

    updates: {
        get: (reload?: boolean) =>
            server.get<Updates>(`/updates${reload ? "?reload=true" : ""}`),
        install: (updates: { name: string }[]) =>
            server.post(`/updates`, { updates }),
    },

    settings: {
        get: () => server.get<Settings>("/settings"),
        patch: (settings: Partial<Settings>) =>
            server.patch("/settings", settings),
    },

    proxy: {
        redirects: {
            get: () => server.get<ProxyRedirects>("/proxy/redirects"),
            add: (source: string, target: string) =>
                server.post("/proxy/redirect", { source, target }),
            delete: (id: string) => server.delete(`/proxy/redirect/${id}`),
        },
    },
};
