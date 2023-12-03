export type Host = {
    hostname?: string;
    uptime?: number;
    boot_time?: number;
    procs?: number;
    os?: string;
    platform?: string;
    platform_family?: string;
    platform_version?: string;
    kernel_version?: string;
    kernel_arch?: string;
    virtualization_system?: string;
    virtualization_role?: string;
    host_id?: string;
};

export type CPU = {
    count?: number;
    vendor_id?: string;
    family?: string;
    model?: string;
    stepping?: number;
    physical_id?: string;
    core_id?: string;
    cores_count?: number;
    model_name?: string;
    mhz?: number;
    cache_size?: number;
    flags?: string[];
    microcode?: string;
};

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
