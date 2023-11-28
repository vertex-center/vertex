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
import { Credentials } from "../../models/auth";

export const server = axios.create({
    // @ts-ignore
    baseURL: `${window.apiURL}/api`,
    headers: {
        Authorization: `Bearer ${getAuthToken()}`,
    },
});

export function setAuthToken(token: string) {
    // Set cookie
    document.cookie = `vertex_auth_token=${token};path=/`;
    server.defaults.headers.common["Authorization"] = `Bearer ${token}`;
}

function getAuthToken() {
    return document?.cookie
        ?.split(";")
        ?.find((c) => c.trim().startsWith("vertex_auth_token="))
        ?.slice("vertex_auth_token=".length);
}

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
    const { data } = await server.get<Update>("/admin/update");
    return data;
};

const installUpdate = async () => {
    const { data } = await server.post("/admin/update");
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

    admin: {
        data: {
            dbms: {
                get: async () => {
                    const { data } = await server.get<string>("/admin/db/dbms");
                    return data;
                },
                migrate: async (dbms: string) => {
                    const { data } = await server.post("/admin/db/dbms", {
                        dbms,
                    });
                    return data;
                },
            },
        },
    },

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
            const { data } = await server.get<Settings>("/admin/settings");
            return data;
        },
        patch: (settings: Partial<Settings>) =>
            server.patch("/admin/settings", settings),
    },

    auth: {
        login: async (credentials: Credentials) => {
            const Authorization = `Basic ${btoa(
                credentials.username + ":" + credentials.password
            )}`;
            const { data } = await server.post(
                "/auth/login",
                {},
                { headers: { Authorization } }
            );
            return data;
        },
        register: async (credentials: Credentials) => {
            const Authorization = `Basic ${btoa(
                credentials.username + ":" + credentials.password
            )}`;
            const { data } = await server.post(
                "/auth/register",
                {},
                { headers: { Authorization } }
            );
            return data;
        },
        logout: async () => {
            const { data } = await server.post("/auth/logout");
            return data;
        },
    },
};
