import axios from "axios";
import {
    InstallMethod,
    Instance,
    InstanceQuery,
    Instances,
} from "../models/instance";
import { Env, Service } from "../models/service";
import { Dependencies as DependenciesUpdate } from "../models/update";
import { DockerContainerInfo } from "../models/docker";
import { About } from "../models/about";
import { Hardware } from "../models/hardware";
import { SSHKeys } from "../models/security";
import { Metric } from "../models/metrics";

type InstallServiceParams = {
    method: InstallMethod;
    service_id: string;
};

const server = axios.create({
    // @ts-ignore
    baseURL: `${window.apiURL}/api`,
});

// server.interceptors.response.use(async (response) => {
//     if (process.env.NODE_ENV === "development")
//         await new Promise((resolve) => setTimeout(resolve, 1000));
//
//     return response;
// });

export const api = {
    about: {
        get: () => server.get<About>("/about"),
    },

    hardware: {
        get: () => server.get<Hardware>("/hardware"),
    },

    instances: {
        get: () => server.get<Instances>("/instances"),
        search: (query: InstanceQuery) =>
            server.get<Instances>("/instances/search", { params: query }),
        checkForUpdates: () => server.get<Instances>("/instances/checkupdates"),
    },

    security: {
        ssh: {
            get: () => server.get<SSHKeys>("/security/ssh"),
            add: (authorized_key: string) =>
                server.post("/security/ssh", { authorized_key }),
            delete: (fingerprint: string) =>
                server.delete(`/security/ssh/${fingerprint}`),
        },
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

        docker: {
            get: (id: string) =>
                server.get<DockerContainerInfo>(`/instance/${id}/docker`),
            recreate: (id: string) =>
                server.post(`/instance/${id}/docker/recreate`),
        },

        update: {
            service: (id: string) =>
                server.post(`/instance/${id}/update/service`),
        },

        versions: {
            get: (id: string, cache?: boolean) =>
                server.get<string[]>(
                    `/instance/${id}/versions?reload=${!cache}`
                ),
        },
    },

    metrics: {
        get: () => server.get<Metric[]>(`/metrics`),
        collector: (collector: string) => ({
            install: () =>
                server.post(`/metrics/collector/${collector}/install`),
        }),
        visualizer: (visualizer: string) => ({
            install: () =>
                server.post(`/metrics/visualizer/${visualizer}/install`),
        }),
    },

    dependencies: {
        get: (reload?: boolean) =>
            server.get<DependenciesUpdate>(
                `/dependencies${reload ? "?reload=true" : ""}`
            ),
        install: (updates: { name: string }[]) =>
            server.post(`/dependencies/update`, { updates }),
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
