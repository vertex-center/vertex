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

export type Features = {
    databases: DatabaseFeature[];
};

export type DatabaseFeature = {
    type?: string;
    port?: string;
    username?: string;
    password?: string;
};

export type DatabaseEnvironment = {
    display_name: string;
    types?: string[];
    names?: { [name: string]: string };
};

export type Service = {
    id: string;
    name: string;
    repository: string;
    description: string;
    color?: string;
    icon?: string;
    features: Features;
    environment: EnvVariable[];
    databases: { [name: string]: DatabaseEnvironment };
    urls?: URL[];
    methods?: ServiceMethods;
};
