export type Host = {
    os?: string;
    arch?: string;
    platform?: string;
    version?: string;
    name?: string;
};

export type Hardware = {
    host?: Host;
};
