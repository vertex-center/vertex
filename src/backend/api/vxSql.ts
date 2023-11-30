import { server } from "./backend";

const BASE_URL = "/app/sql";

const installDbms = (dbms: string) => {
    return server.post(`${BASE_URL}/dbms/${dbms}/install`);
};

const getContainer = async (uuid: string) => {
    const { data } = await server.get(`${BASE_URL}/container/${uuid}`);
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
