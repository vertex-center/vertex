export type About = {
    version: string;
    commit: string;
    date: string;

    os?: string;
    arch?: string;
};

export type Update = {
    baseline: Baseline;
};

export type Baseline = {
    date: string;
    version: string;
    description: string;
    vertex: string;
    vertex_client: string;
    vertex_server: string;
};
