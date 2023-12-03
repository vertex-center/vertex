import axios from "axios";
import { About } from "../../models/about";
import { vxContainersRoutes } from "./vxContainers";
import { vxTunnelsRoutes } from "./vxTunnels";
import { vxMonitoringRoutes } from "./vxMonitoring";
import { vxSqlRoutes } from "./vxSql";
import { vxReverseProxyRoutes } from "./vxReverseProxy";
import { VertexApp } from "../../models/app";
import { Console } from "../../logging/logging";
import { vxServiceEditorRoutes } from "./vxServiceEditor";
import { AuthCredentials, Credentials } from "../../models/auth";

export const server = axios.create({
    // @ts-ignore
    baseURL: `${window.apiURL}/api`,
});

export function setAuthToken(token?: string) {
    if (token === undefined) {
        // delete cookie
        document.cookie = "vertex_auth_token=;Max-Age=-99999999;path=/";
        return;
    }
    const expires = new Date();
    expires.setTime(expires.getTime() + 60 * 60 * 24 * 365);
    document.cookie = `vertex_auth_token=${token};path=/;SameSite=Lax;expires=${expires.toUTCString()}`;
}

export function getAuthToken() {
    return document?.cookie
        ?.split(";")
        ?.find((c) => c.trim().startsWith("vertex_auth_token="))
        ?.replace("vertex_auth_token=", "");
}

server.interceptors.request.use(async (config) => {
    if (!config.headers.Authorization) {
        config.headers.Authorization = `Bearer ${getAuthToken()}`;
    }
    return config;
});

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

export const api = {
    about: getAbout,

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

    auth: {
        login: async (credentials: AuthCredentials) => {
            const Authorization = `Basic ${btoa(
                credentials.username + ":" + credentials.password
            )}`;
            const { data } = await server.post(
                "/app/auth/login",
                {},
                { headers: { Authorization } }
            );
            return data;
        },
        register: async (credentials: AuthCredentials) => {
            const Authorization = `Basic ${btoa(
                credentials.username + ":" + credentials.password
            )}`;
            const { data } = await server.post(
                "/app/auth/register",
                {},
                { headers: { Authorization } }
            );
            return data;
        },
        logout: async () => {
            const { data } = await server.post("/app/auth/logout");
            return data;
        },
        user: (username?: string) => {
            if (username !== undefined) {
                username = `/${username}`;
            } else {
                username = "";
            }

            return {
                credentials: {
                    get: async () => {
                        const { data } = await server.get<Credentials[]>(
                            `/app/auth/user${username}/credentials`
                        );
                        return data;
                    },
                },
                get: async () => {
                    const { data } = await server.get<User>(
                        `/app/auth/user${username}`
                    );
                    return data;
                },
                patch: async (user: Partial<User>) => {
                    const { data } = await server.patch(
                        `/app/auth/user${username}`,
                        user
                    );
                    return data;
                },
            };
        },
    },
};
