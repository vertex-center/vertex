import { createServer } from "../server";

// @ts-ignore
const server = createServer(window.api_urls.reverse_proxy);

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
