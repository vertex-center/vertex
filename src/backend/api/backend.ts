import { About } from "../../models/about";
import { vxContainersRoutes } from "./vxContainers";
import { vxTunnelsRoutes } from "./vxTunnels";
import { vxMonitoringRoutes } from "./vxMonitoring";
import { vxSqlRoutes } from "./vxSql";
import { vxReverseProxyRoutes } from "./vxReverseProxy";
import { VertexApp } from "../../models/app";
import { vxServiceEditorRoutes } from "./vxServiceEditor";
import { createServer } from "../server";

// @ts-ignore
export const server = createServer(window?.apiPortVertex ?? "6130");

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
};
