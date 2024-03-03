import { vxTunnelsRoutes } from "./vxTunnels";
import { vxMonitoringRoutes } from "./vxMonitoring";
import { vxSqlRoutes } from "./vxSql";
import { vxReverseProxyRoutes } from "./vxReverseProxy";
import { VertexApp } from "../../models/app";
import { createServer } from "../server";

// @ts-ignore
export const server = createServer(window.api_urls.vertex);

export const api = {
    vxTunnels: vxTunnelsRoutes,
    vxMonitoring: vxMonitoringRoutes,
    vxSql: vxSqlRoutes,
    vxReverseProxy: vxReverseProxyRoutes,

    apps: {
        all: async () => {
            const { data } = await server.get<VertexApp[]>("/apps");
            return data;
        },
    },
};
