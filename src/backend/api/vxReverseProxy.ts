import { server } from "./backend";

const BASE_URL = `/app/reverse-proxy`;

const getRedirects = async () => {
    const { data } = await server.get<ProxyRedirects>(`${BASE_URL}/redirects`);
    return data;
};

const addRedirect = (source: string, target: string) => {
    const data = { source, target };
    return server.post(`${BASE_URL}/redirect`, data);
};

const deleteRedirect = (id: string) => {
    return server.delete(`${BASE_URL}/redirect/${id}`);
};

const redirectsRoutes = {
    all: getRedirects,
    add: addRedirect,
    delete: deleteRedirect,
};

export const vxReverseProxyRoutes = {
    redirects: redirectsRoutes,
};
