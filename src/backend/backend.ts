import axios from "axios";

export type EnvVariable = {
    type: string;
    name: string;
    display_name: string;
    secret: boolean;
    default: string;
    description: string;
};

export type Service = {
    id: string;
    name: string;
    repository: string;
    description: string;
    environment: EnvVariable[];
};

export type InstalledService = Service & {
    status: string;
};

export type InstalledServices = { [uuid: string]: InstalledService };

export async function getInstalledServices(): Promise<InstalledServices> {
    return new Promise((resolve, reject) => {
        axios
            .get("http://localhost:6130/services")
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function getAvailableServices(): Promise<Service[]> {
    return new Promise((resolve, reject) => {
        axios
            .get("http://localhost:6130/services/available")
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function postDownloadService(service: Service) {
    return new Promise((resolve, reject) => {
        axios
            .post("http://localhost:6130/services/download", { service })
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function getService(uuid: string) {
    return new Promise((resolve, reject) => {
        axios
            .get(`http://localhost:6130/service/${uuid}`)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function startService(uuid: string) {
    return new Promise((resolve, reject) => {
        axios
            .post(`http://localhost:6130/service/${uuid}/start`)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function stopService(uuid: string) {
    return new Promise((resolve, reject) => {
        axios
            .post(`http://localhost:6130/service/${uuid}/stop`)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}
