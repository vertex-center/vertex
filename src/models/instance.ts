import { Service } from "./service";

export type Instance = Service & {
    uuid: string;
    status: string;
    env: { [key: string]: string };
    use_docker?: boolean;
    use_releases?: boolean;
    launch_on_startup?: boolean;
};

export type Instances = { [uuid: string]: Instance };
