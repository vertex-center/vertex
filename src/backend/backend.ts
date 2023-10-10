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
    about: () => server.get<About>("/about"),
    hardware: () => server.get<Hardware>("/hardware"),

    vxInstances: {
        instances: {
            all: () => server.get<Instances>("/app/vx-instances/instances"),
            search: (query: InstanceQuery) =>
                server.get<Instances>("/app/vx-instances/instances/search", {
                    params: query,
                }),
            checkForUpdates: () =>
                server.get<Instances>(
                    "/app/vx-instances/instances/checkupdates"
                ),
        },

        service: (service_id: string) => ({
            install: () =>
                server.post(`/app/vx-instances/service/${service_id}/install`),
        }),

        services: {
            all: () => server.get<Service[]>("/app/vx-instances/services"),
        },

        instance: (id: string) => ({
            get: () => server.get<Instance>(`/app/vx-instances/instance/${id}`),
            delete: () => server.delete(`/app/vx-instances/instance/${id}`),
            start: () => server.post(`/app/vx-instances/instance/${id}/start`),
            stop: () => server.post(`/app/vx-instances/instance/${id}/stop`),
            patch: (params: any) =>
                server.patch(`/app/vx-instances/instance/${id}`, params),

            logs: {
                get: () => server.get(`/app/vx-instances/instance/${id}/logs`),
            },

            env: {
                save: (env: Env) =>
                    server.patch(
                        `/app/vx-instances/instance/${id}/environment`,
                        env
                    ),
            },

            docker: {
                get: () =>
                    server.get<DockerContainerInfo>(
                        `/app/vx-instances/instance/${id}/docker`
                    ),
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
                get: (cache?: boolean) =>
                    server.get<string[]>(
                        `/app/vx-instances/instance/${id}/versions?reload=${!cache}`
                    ),
            },
        }),
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

    vxTunnels: {
        provider: (provider: string) => ({
            install: () =>
                server.post(`/app/vx-tunnels/provider/${provider}/install`),
        }),
    },

    vxMonitoring: {
        metrics: () => server.get<Metric[]>(`/vx-monitoring`),
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
            get: () => server.get(`/app/vx-sql/instance/${uuid}`),
        }),
        dbms: (dbms: string) => ({
            install: () => server.post(`/app/vx-sql/dbms/${dbms}/install`),
        }),
    },

    dependencies: {
        get: (reload?: boolean) =>
            server.get<DependenciesUpdate>(
                `/dependencies${reload ? "?reload=true" : ""}`
            ),
        install: (
            updates: {
                name: string;
            }[]
        ) => server.post(`/dependencies/update`, { updates }),
    },

    settings: {
        get: () => server.get<Settings>("/settings"),
        patch: (settings: Partial<Settings>) =>
            server.patch("/settings", settings),
    },

    vxReverseProxy: {
        redirects: {
            get: () =>
                server.get<ProxyRedirects>("/app/vx-reverse-proxy/redirects"),
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
