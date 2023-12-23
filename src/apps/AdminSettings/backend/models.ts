export type SSHKey = {
    type: string;
    fingerprint_sha_256: string;
    username: string;
};

export type Update = {
    baseline: Baseline;
    updating: boolean;
};

export type Baseline = {
    date: string;
    version: string;
    description: string;
    vertex: string;
    vertex_client: string;
    vertex_server: string;
};
