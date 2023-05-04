import axios from "axios";

export type LogLine = {
    id: number;
    kind: string;
    message: string;
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
    dependencies?: { [name: string]: boolean };
};

export type Instance = Service & {
    uuid: string;
    status: string;
    env: { [key: string]: string };
    use_docker?: boolean;
    use_releases?: boolean;
    launch_on_startup?: boolean;
};

export type DockerContainerInfo = {
    id: string;
    name: string;
    image: string;
    platform: string;
};

export type Package = {
    name: string;
    description?: string;
    homepage?: string;
    license?: string;
    check?: string;
    install?: { [pm: string]: string };
};

export type Dependency = Package & {
    installed: boolean;
};

export type About = {
    version: string;
    commit: string;
    date: string;
};

export type Updates = {
    last_checked: string;
    items: Update[];
};

export type Update = {
    id: string;
    name: string;
    current_version: string;
    latest_version: string;
    needs_restart?: boolean;
};

export type Dependencies = { [id: string]: Dependency };

export type Instances = { [uuid: string]: Instance };

export function route(path: string) {
    return `http://localhost:6130/api${path}`;
}

export async function getInstances(): Promise<Instances> {
    return new Promise((resolve, reject) => {
        axios
            .get(route("/instances"))
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function getAvailableServices(): Promise<Service[]> {
    return new Promise((resolve, reject) => {
        axios
            .get(route("/services/available"))
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function downloadService(
    repository: string,
    use_docker?: boolean,
    use_releases?: boolean
) {
    return new Promise((resolve, reject) => {
        axios
            .post(route("/services/download"), {
                repository,
                use_docker,
                use_releases,
            })
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function getInstance(uuid: string): Promise<Instance> {
    return new Promise((resolve, reject) => {
        axios
            .get(route(`/instance/${uuid}`))
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function deleteInstance(uuid: string) {
    return new Promise((resolve, reject) => {
        axios
            .delete(route(`/instance/${uuid}`))
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function patchInstance(uuid: string, params: any) {
    return new Promise((resolve, reject) => {
        axios
            .patch(route(`/instance/${uuid}`), params)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function startInstance(uuid: string) {
    return new Promise((resolve, reject) => {
        axios
            .post(route(`/instance/${uuid}/start`))
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function stopInstance(uuid: string) {
    return new Promise((resolve, reject) => {
        axios
            .post(route(`/instance/${uuid}/stop`))
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function saveInstanceEnv(uuid: string, env: Env) {
    return new Promise((resolve, reject) => {
        axios
            .patch(route(`/instance/${uuid}/environment`), env)
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function getInstanceDependencies(
    uuid: string
): Promise<Dependencies> {
    return new Promise((resolve, reject) => {
        axios
            .get(route(`/instance/${uuid}/dependencies`))
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function getInstanceDockerContainerInfo(
    uuid: string
): Promise<DockerContainerInfo> {
    return new Promise((resolve, reject) => {
        axios
            .get(route(`/instance/${uuid}/docker`))
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function installPackages(packages) {
    return new Promise((resolve, reject) => {
        axios
            .post(route(`/packages/install`), { packages })
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function getAbout(): Promise<About> {
    return new Promise((resolve, reject) => {
        axios
            .get(route("/about"))
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function getUpdates(reload?: boolean): Promise<Updates> {
    return new Promise((resolve, reject) => {
        axios
            .get(route(`/updates${reload ? "?reload=true" : ""}`))
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}

export async function executeUpdates(updates: { name: string }[]) {
    return new Promise((resolve, reject) => {
        axios
            .post(route("/updates"), { updates })
            .then((res) => resolve(res.data))
            .catch((err) => reject(err));
    });
}
