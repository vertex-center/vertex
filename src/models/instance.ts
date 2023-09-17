import { Service } from "./service";

export type InstallMethod = "script" | "release" | "docker";

export type InstanceQuery = {
    features?: string[];
};

export type InstanceUpdate = {
    current_version: string;
    latest_version: string;
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
    types: string[];
    names: string[];
};

export type Instance = Service & {
    uuid: string;
    status: string;
    features: Features;
    env: { [key: string]: string };
    databases: DatabaseEnvironment[];
    install_method?: InstallMethod;
    launch_on_startup?: boolean;
    display_name?: string;
    update?: InstanceUpdate;
};

export type Instances = { [uuid: string]: Instance };
