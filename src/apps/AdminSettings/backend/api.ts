import { server } from "../../../backend/api/backend";
import { CPU, Host, SSHKey, Update } from "./models";

const getHost = async () => {
    const { data } = await server.get<Host>("/app/admin/hardware/host");
    return data;
};

const getCPUs = async () => {
    const { data } = await server.get<CPU[]>("/app/admin/hardware/cpus");
    return data;
};

const getSSHKeys = async () => {
    const { data } = await server.get<SSHKey[]>("/app/admin/ssh");
    return data;
};

export type AddSSHKeyBody = {
    authorized_key: string;
    username: string;
};

const addSSHKey = async (body: AddSSHKeyBody) => {
    const { data } = await server.post("/app/admin/ssh", body);
    return data;
};

export type DeleteSSHKeyBody = {
    fingerprint: string;
    username: string;
};

const deleteSSHKey = async (body: DeleteSSHKeyBody) => {
    const { data } = await server.delete("/app/admin/ssh", { data: body });
    return data;
};

const getSSHUsers = async () => {
    const { data } = await server.get<string[]>("/app/admin/ssh/users");
    return data;
};

const getSettings = async () => {
    const { data } = await server.get<Settings>("/app/admin/settings");
    return data;
};

const patchSettings = async (settings: Partial<Settings>) => {
    const { data } = await server.patch("/app/admin/settings", settings);
    return data;
};

const getUpdate = async () => {
    const { data } = await server.get<Update>("/app/admin/update");
    return data;
};

const installUpdate = async () => {
    const { data } = await server.post("/app/admin/update");
    return data;
};

const getDatabases = async () => {
    const { data } = await server.get<string>("/app/admin/db/dbms");
    return data;
};

const migrateDatabase = async (dbms: string) => {
    const body = { dbms };
    const { data } = await server.post("/app/admin/db/dbms", body);
    return data;
};

const reboot = async () => {
    const { data } = await server.post("/app/admin/hardware/reboot");
    return data;
};

export const API = {
    getHost,
    getCPUs,
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
    reboot,
};
