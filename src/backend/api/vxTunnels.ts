import { server } from "./backend";

const BASE_URL = "/app/vx-tunnels";

const installProvider = (provider: string) => {
    return server.post(`${BASE_URL}/provider/${provider}/install`);
};

const providerRoutes = (provider: string) => ({
    install: () => installProvider(provider),
});

export const vxTunnelsRoutes = {
    provider: providerRoutes,
};
