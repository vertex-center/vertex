import axios from "axios";

export type LogLine = {
    id: number;
    kind: string;
    message: string;
};

export type Logs = {
    lines: LogLine[];
};

export type Env = { [key: string]: string };

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

export type Instance = Service & {
    status: string;
    logs: Logs;
    env: { [key: string]: string };
};

export type Instances = { [uuid: string]: Instance };

export async function getInstances(): Promise<Instances> {
    return new Promise((resolve, reject) => {
        axios
            .get("http://localhost:6130/instances")
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

export async function downloadService(repository: string) {
    return new Promise((resolve, reject) => {
        axios
            .post("http://localhost:6130/services/download", { repository })
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function getInstance(uuid: string) {
    return new Promise((resolve, reject) => {
        axios
            .get(`http://localhost:6130/instance/${uuid}`)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function deleteInstance(uuid: string) {
    return new Promise((resolve, reject) => {
        axios
            .delete(`http://localhost:6130/instance/${uuid}`)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function startInstance(uuid: string) {
    return new Promise((resolve, reject) => {
        axios
            .post(`http://localhost:6130/instance/${uuid}/start`)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function stopInstance(uuid: string) {
    return new Promise((resolve, reject) => {
        axios
            .post(`http://localhost:6130/instance/${uuid}/stop`)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function saveInstanceEnv(uuid: string, env: Env) {
    return new Promise((resolve, reject) => {
        axios
            .patch(`http://localhost:6130/instance/${uuid}/environment`, env)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}
