import axios from "axios";

export type Service = {
    id: string;
    name: string;
    repository: string;
    description: string;
};

export type InstalledService = Service & {
    status: string;
};

export type InstalledServices = { [id: string]: InstalledService };

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

export async function startService(service: InstalledService) {
    return new Promise((resolve, reject) => {
        axios
            .post(`http://localhost:6130/service/${service.id}/start`)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function stopService(service: InstalledService) {
    return new Promise((resolve, reject) => {
        axios
            .post(`http://localhost:6130/service/${service.id}/stop`)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}
