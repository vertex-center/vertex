import { Service } from "./service";

export type InstallMethod = "script" | "release" | "docker";

export type Instance = Service & {
    uuid: string;
    status: string;
    env: { [key: string]: string };
    install_method?: InstallMethod;
    use_releases?: boolean;
    launch_on_startup?: boolean;
};

export type Instances = { [uuid: string]: Instance };
