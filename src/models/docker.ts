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
        virtual_size?: number;
        tags?: string[];
    };
};
