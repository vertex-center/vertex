import { createServer } from "../server";

const server = createServer("7508");

const getRedirects = async () => {
    const { data } = await server.get<ProxyRedirects>(`/redirects`);
    return data;
};

const addRedirect = (source: string, target: string) => {
    const data = { source, target };
    return server.post(`/redirect`, data);
};

const deleteRedirect = (id: string) => {
    return server.delete(`/redirect/${id}`);
};

const redirectsRoutes = {
    all: getRedirects,
    add: addRedirect,
    delete: deleteRedirect,
};

export const vxReverseProxyRoutes = {
    redirects: redirectsRoutes,
};
