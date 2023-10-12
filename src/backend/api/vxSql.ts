import { server } from "./backend";

const BASE_URL = "/app/vx-sql";

const installDbms = (dbms: string) => {
    return server.post(`${BASE_URL}/dbms/${dbms}/install`);
};

const getInstance = async (uuid: string) => {
    const { data } = await server.get(`${BASE_URL}/instance/${uuid}`);
    return data;
};

const instanceRoutes = (uuid: string) => ({
    get: () => getInstance(uuid),
});

const dbmsRoutes = (dbms: string) => ({
    install: () => installDbms(dbms),
});

export const vxSqlRoutes = {
    instance: instanceRoutes,
    dbms: dbmsRoutes,
};
