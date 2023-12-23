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
