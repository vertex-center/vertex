import { createServer } from "../server";

const server = createServer("7514");

const installProvider = (provider: string) => {
    return server.post(`/provider/${provider}/install`);
};

const providerRoutes = (provider: string) => ({
    install: () => installProvider(provider),
});

export const vxTunnelsRoutes = {
    provider: providerRoutes,
};
