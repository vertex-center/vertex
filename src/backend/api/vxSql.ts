import { createServer } from "../server";

// @ts-ignore
const server = createServer(window.api_urls.sql);

const installDbms = (dbms: string) => {
    return server.post(`/dbms/${dbms}/install`);
};

const getContainer = async (uuid: string) => {
    const { data } = await server.get(`/container/${uuid}`);
    return data;
};

const containerRoutes = (uuid: string) => ({
    get: () => getContainer(uuid),
});

const dbmsRoutes = (dbms: string) => ({
    install: () => installDbms(dbms),
});

export const vxSqlRoutes = {
    container: containerRoutes,
    dbms: dbmsRoutes,
};
