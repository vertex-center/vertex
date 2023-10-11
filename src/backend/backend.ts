import axios from "axios";
import { Instance, InstanceQuery, Instances } from "../models/instance";
import { Env, Service } from "../models/service";
import { Dependencies as DependenciesUpdate } from "../models/update";
import { DockerContainerInfo } from "../models/docker";
import { About } from "../models/about";
import { Hardware } from "../models/hardware";
import { SSHKeys } from "../models/security";
import { Metric } from "../models/metrics";

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
    about: async () => {
        const { data } = await server.get<About>("/about");
        return data;
    },
    hardware: async () => {
        const { data } = await server.get<Hardware>("/hardware");
        return data;
    },

    vxInstances: {
        instances: {
            all: async () => {
                const { data } = await server.get<Instances>(
                    "/app/vx-instances/instances"
                );
                return data;
            },
            search: async (query: InstanceQuery) => {
                const { data } = await server.get<Instances>(
                    "/app/vx-instances/instances/search",
                    { params: query }
                );
                return data;
            },
        },

        service: (service_id: string) => ({
            install: async () => {
                const { data } = await server.post(
                    `/app/vx-instances/service/${service_id}/install`
                );
                return data;
            },
        }),

        services: {
            all: async () => {
                const { data } = await server.get<Service[]>(
                    "/app/vx-instances/services"
                );
                return data;
            },
        },

        instance: (id: string) => ({
            get: async () => {
                const { data } = await server.get<Instance>(
                    `/app/vx-instances/instance/${id}`
                );
                return data;
            },
            delete: () => server.delete(`/app/vx-instances/instance/${id}`),
            start: () => server.post(`/app/vx-instances/instance/${id}/start`),
            stop: () => server.post(`/app/vx-instances/instance/${id}/stop`),
            patch: (params: any) =>
                server.patch(`/app/vx-instances/instance/${id}`, params),

            logs: {
                get: async () => {
                    const { data } = await server.get(
                        `/app/vx-instances/instance/${id}/logs`
                    );
                    return data;
                },
            },

            env: {
                save: (env: Env) =>
                    server.patch(
                        `/app/vx-instances/instance/${id}/environment`,
                        env
                    ),
            },

            docker: {
                get: async () => {
                    const { data } = await server.get<DockerContainerInfo>(
                        `/app/vx-instances/instance/${id}/docker`
                    );
                    return data;
                },
                recreate: () =>
                    server.post(
                        `/app/vx-instances/instance/${id}/docker/recreate`
                    ),
            },

            update: {
                service: () =>
                    server.post(
                        `/app/vx-instances/instance/${id}/update/service`
                    ),
            },

            versions: {
                get: async (cache?: boolean) => {
                    const { data } = await server.get<string[]>(
                        `/app/vx-instances/instance/${id}/versions?reload=${!cache}`
                    );
                    return data;
                },
            },
        }),
    },

    security: {
        ssh: {
            get: async () => {
                const { data } = await server.get<SSHKeys>("/security/ssh");
                return data;
            },
            add: (authorized_key: string) =>
                server.post("/security/ssh", { authorized_key }),
            delete: (fingerprint: string) =>
                server.delete(`/security/ssh/${fingerprint}`),
        },
    },

    vxTunnels: {
        provider: (provider: string) => ({
            install: () =>
                server.post(`/app/vx-tunnels/provider/${provider}/install`),
        }),
    },

    vxMonitoring: {
        metrics: async () => {
            const { data } = await server.get<Metric[]>(
                `/app/vx-monitoring/metrics`
            );
            return data;
        },
        collector: (collector: string) => ({
            install: () =>
                server.post(
                    `/app/vx-monitoring/collector/${collector}/install`
                ),
        }),
        visualizer: (visualizer: string) => ({
            install: () =>
                server.post(
                    `/app/vx-monitoring/visualizer/${visualizer}/install`
                ),
        }),
    },

    vxSql: {
        instance: (uuid: string) => ({
            get: async () => {
                const { data } = await server.get(
                    `/app/vx-sql/instance/${uuid}`
                );
                return data;
            },
        }),
        dbms: (dbms: string) => ({
            install: () => server.post(`/app/vx-sql/dbms/${dbms}/install`),
        }),
    },

    dependencies: {
        get: async (reload?: boolean) => {
            const { data } = await server.get<DependenciesUpdate>(
                `/dependencies${reload ? "?reload=true" : ""}`
            );
            return data;
        },
        install: (
            updates: {
                name: string;
            }[]
        ) => server.post(`/dependencies/update`, { updates }),
    },

    settings: {
        get: async () => {
            const { data } = await server.get<Settings>("/settings");
            return data;
        },
        patch: (settings: Partial<Settings>) =>
            server.patch("/settings", settings),
    },

    vxReverseProxy: {
        redirects: {
            get: async () => {
                const { data } = await server.get<ProxyRedirects>(
                    "/app/vx-reverse-proxy/redirects"
                );
                return data;
            },
            add: (source: string, target: string) =>
                server.post("/app/vx-reverse-proxy/redirect", {
                    source,
                    target,
                }),
            delete: (id: string) =>
                server.delete(`/app/vx-reverse-proxy/redirect/${id}`),
        },
    },
};
