import axios from "axios";
import { Console } from "../logging/logging";

export const createServer = (baseURL: string) => {
    const s = axios.create({ baseURL });

    s.interceptors.request.use(async (config) => {
        if (!config.headers.Authorization) {
            config.headers.Authorization = `Bearer ${getAuthToken()}`;
        }
        return config;
    });

    s.interceptors.request.use((req) => {
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

    return s;
};

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
