export type DockerContainerInfo = {
    container?: {
        id?: string;
        name?: string;
        platform?: string;
    };
    image?: {
        id?: string;
        architecture?: string;
        os?: string;
        size?: number;
        tags?: string[];
    };
};
