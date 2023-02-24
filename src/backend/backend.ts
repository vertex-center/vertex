import axios from "axios";

export type Service = {
    id: string;
    name: string;
    repository: string;
    description: string;
};

export async function getInstalledServices(): Promise<Service[]> {
    return new Promise((resolve, reject) => {
        axios
            .get("http://localhost:6130/installed")
            .then((res) => {
                resolve(res.data);
            })
            .catch((err) => {
                reject(err);
            });
    });
}

export async function getAvailableServices(): Promise<Service[]> {
    return new Promise((resolve, reject) => {
        axios
            .get("http://localhost:6130/available")
            .then((res) => {
                resolve(res.data);
            })
            .catch((err) => {
                reject(err);
            });
    });
}
