import axios from "axios";
import { InstallMethod, Instance, Instances } from "../models/instance";
import { Env, Service } from "../models/service";
import { Dependencies } from "../models/dependency";
import { DockerContainerInfo } from "../models/docker";
import { Uptime } from "../models/uptime";
import { About } from "../models/about";
import { Updates } from "../models/update";

type InstallServiceParams = {
    method: InstallMethod;
    service_id: string;
};

const api = axios.create({
    baseURL: "http://localhost:6130/api",
});

export const getAbout = async () => api.get<About>("/about");
export const getInstances = async () => api.get<Instances>("/instances");
export const getAvailableServices = async () =>
    api.get<Service[]>("/services/available");
export const installService = async (params: InstallServiceParams) =>
    api.post("/services/install", params);
export const getInstance = async (uuid: string) =>
    api.get<Instance>(`/instance/${uuid}`);
export const deleteInstance = async (uuid: string) =>
    api.delete(`/instance/${uuid}`);
export const patchInstance = async (uuid: string, params: any) =>
    api.patch(`/instance/${uuid}`, params);
export const startInstance = async (uuid: string) =>
    api.post(`/instance/${uuid}/start`);
export const stopInstance = async (uuid: string) =>
    api.post(`/instance/${uuid}/stop`);
export const getLatestLogs = async (uuid: string) =>
    api.get(`/instance/${uuid}/logs`);
export const saveInstanceEnv = async (uuid: string, env: Env) =>
    api.patch(`/instance/${uuid}/environment`, env);
export const getInstanceDependencies = async (uuid: string) =>
    api.get<Dependencies>(`/instance/${uuid}/dependencies`);
export const getInstanceDockerContainerInfo = async (uuid: string) =>
    api.get<DockerContainerInfo>(`/instance/${uuid}/docker`);
export const getInstanceStatus = async (uuid: string) =>
    api.get<Uptime[]>(`/instance/${uuid}/status`);
export const installPackages = async (packages) =>
    api.post(`/packages/install`, { packages });
export const getUpdates = async (reload?: boolean) =>
    api.get<Updates>(`/updates${reload ? "?reload=true" : ""}`);
export const executeUpdates = async (updates: { name: string }[]) =>
    api.post("/updates", { updates });
