import { Update } from "./models";
import { createServer } from "../../../backend/server";

// @ts-ignore
const server = createServer(window.api_urls.admin);

const getSettings = async () => {
    const { data } = await server.get<Settings>("/settings");
    return data;
};

const patchSettings = async (settings: Partial<Settings>) => {
    const { data } = await server.patch("/settings", settings);
    return data;
};

const getUpdate = async () => {
    const { data } = await server.get<Update>("/update");
    return data;
};

const getDatabases = async () => {
    const { data } = await server.get<string>("/db/dbms");
    return data;
};

const migrateDatabase = async (dbms: string) => {
    const body = { dbms };
    const { data } = await server.post("/db/dbms", body);
    return data;
};

export const API = {
    getSettings,
    patchSettings,
    getUpdate,
    getDatabases,
    migrateDatabase,
};
