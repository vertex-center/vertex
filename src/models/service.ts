export type Env = { [key: string]: string };

export type EnvVariable = {
    type: string;
    name: string;
    display_name: string;
    secret: boolean;
    default: string;
    description: string;
};

export type URL = {
    name: string;
    port: string;
    home?: string;
    ping?: string;
    kind: string;
};

export type ServiceMethodScript = {
    file: string;
    dependencies?: { [name: string]: boolean };
};

export type ServiceMethodRelease = {
    dependencies?: { [name: string]: boolean };
};

export type ServiceMethodDocker = {
    image?: string;
    dockerfile?: string;
    ports?: string[];
    volumes?: { [key: string]: string };
};

export type ServiceMethods = {
    script?: ServiceMethodScript;
    release?: ServiceMethodRelease;
    docker?: ServiceMethodDocker;
};

export type Service = {
    id: string;
    name: string;
    repository: string;
    description: string;
    environment: EnvVariable[];
    urls?: URL[];
    methods?: ServiceMethods;
};
