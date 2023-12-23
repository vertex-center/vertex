import { SSHKey, Update } from "./models";
import { createServer } from "../../../backend/server";

// @ts-ignore
const server = createServer(window.api_urls.admin);

const getSSHKeys = async () => {
    const { data } = await server.get<SSHKey[]>("/ssh");
    return data;
};

export type AddSSHKeyBody = {
    authorized_key: string;
    username: string;
};

const addSSHKey = async (body: AddSSHKeyBody) => {
    const { data } = await server.post("/ssh", body);
    return data;
};

export type DeleteSSHKeyBody = {
    fingerprint: string;
    username: string;
};

const deleteSSHKey = async (body: DeleteSSHKeyBody) => {
    const { data } = await server.delete("/ssh", { data: body });
    return data;
};

const getSSHUsers = async () => {
    const { data } = await server.get<string[]>("/ssh/users");
    return data;
};

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

const installUpdate = async () => {
    const { data } = await server.post("/update");
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
    getSSHKeys,
    addSSHKey,
    deleteSSHKey,
    getSSHUsers,
    getSettings,
    patchSettings,
    getUpdate,
    installUpdate,
    getDatabases,
    migrateDatabase,
};
