import axios from "axios";
import { About } from "../../models/about";
import { SSHKeys } from "../../models/security";
import { vxContainersRoutes } from "./vxContainers";
import { vxTunnelsRoutes } from "./vxTunnels";
import { vxMonitoringRoutes } from "./vxMonitoring";
import { vxSqlRoutes } from "./vxSql";
import { vxReverseProxyRoutes } from "./vxReverseProxy";
import { VertexApp } from "../../models/app";
import { Console } from "../../logging/logging";
import { Update } from "../../models/update";
import { vxServiceEditorRoutes } from "./vxServiceEditor";
import { CPU, Host } from "../../models/hardware";

export const server = axios.create({
    // @ts-ignore
    baseURL: `${window.apiURL}/api`,
});

// server.interceptors.response.use(async (response) => {
//     if (process.env.NODE_ENV === "development")
//         await new Promise((resolve) => setTimeout(resolve, 1000));
//
//     return response;
// });

server.interceptors.request.use((req) => {
    if (!req) return;

    const info = {
        url: req.url,
        method: req.method,
    };

    if (req.data) info["data"] = req.data;
    if (req.params) info["params"] = req.params;

    Console.request("Sending request\n%O", info);

    return req;
});

const getAbout = async () => {
    const { data } = await server.get<About>("/about");
    return data;
};

const getHost = async () => {
    const { data } = await server.get<Host>("/hardware/host");
    return data;
};

const getCPUs = async () => {
    const { data } = await server.get<CPU[]>("/hardware/cpus");
    return data;
};

const getUpdate = async () => {
    const { data } = await server.get<Update>("/update");
    return data;
};

const installUpdate = async () => {
    const { data } = await server.post("/update");
    return data;
};

export const api = {
    about: getAbout,
    hardware: {
        host: getHost,
        cpus: getCPUs,
    },

    vxContainers: vxContainersRoutes,
    vxTunnels: vxTunnelsRoutes,
    vxMonitoring: vxMonitoringRoutes,
    vxSql: vxSqlRoutes,
    vxReverseProxy: vxReverseProxyRoutes,
    vxServiceEditor: vxServiceEditorRoutes,

    apps: {
        all: async () => {
            const { data } = await server.get<VertexApp[]>("/apps");
            return data;
        },
    },

    security: {
        ssh: {
            get: async () => {
                const { data } = await server.get<SSHKeys>("/security/ssh");
                return data;
            },
            add: (authorized_key: string, username: string) =>
                server.post("/security/ssh", { authorized_key, username }),
            delete: (fingerprint: string, username: string) =>
                server.delete("/security/ssh", {
                    data: { fingerprint, username },
                }),
            users: async () => {
                const { data } = await server.get<string[]>(
                    "/security/ssh/users"
                );
                return data;
            },
        },
    },

    update: {
        get: getUpdate,
        install: installUpdate,
    },

    settings: {
        get: async () => {
            const { data } = await server.get<Settings>("/settings");
            return data;
        },
        patch: (settings: Partial<Settings>) =>
            server.patch("/settings", settings),
    },
};
