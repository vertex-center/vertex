export type ContainerFilters = {
    features?: string[];
    tags?: string[];
};

export type ContainerUpdate = {
    current_version: string;
    latest_version: string;
};

export type ServiceUpdate = {
    available?: boolean;
};

export type Tags = Tag[];
export type Tag = {
    id: string;
    container_id: string;
    name: string;
};

export type EnvVariable = {
    id: string;
    container_id: string;
    type: string;
    name: string;
    display_name?: string;
    value: string;
    default?: string;
    description?: string;
    secret?: boolean;
};

export type Port = {
    id?: string;
    container_id: string;
    in: string;
    out: string;
};

export type Containers = Container[];
export type Container = {
    id: string;
    template_id: string;
    user_id: string;
    image: string;
    image_tag: string;
    status: string;
    launch_on_startup: boolean;
    name: string;
    description?: string;
    color?: string;
    icon?: string;
    command?: string;
    environment: EnvVariable[];
    capabilities: {
        container_id: string;
        name: string;
    };
    ports: {
        container_id: string;
        in: string;
        out: string;
    }[];
    volumes: {
        container_id: string;
        in: string;
        out: string;
    }[];
    sysctls: {
        container_id: string;
        name: string;
        value: string;
    }[];
    tags: Tags;

    update?: ContainerUpdate;
    service_update?: ServiceUpdate;
    databases?: { [key: string]: string };
};
