import { Service } from "./service";

export type InstallMethod = "script" | "release" | "docker";

export type InstanceQuery = {
    features?: string[];
};

export type InstanceUpdate = {
    current_version: string;
    latest_version: string;
};

export type Instance = {
    service: Service;
    uuid: string;
    status: string;
    environment: { [key: string]: string };
    install_method?: InstallMethod;
    launch_on_startup?: boolean;
    display_name?: string;
    update?: InstanceUpdate;
};

export type Instances = { [uuid: string]: Instance };
