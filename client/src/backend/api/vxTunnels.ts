import { createServer } from "../server";

// @ts-ignore
const server = createServer(window.api_urls.tunnels);

const installProvider = (provider: string) => {
    return server.post(`/provider/${provider}/install`);
};

const providerRoutes = (provider: string) => ({
    install: () => installProvider(provider),
});

export const vxTunnelsRoutes = {
    provider: providerRoutes,
};
