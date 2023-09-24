export type SSHKey = {
    type: string;
    fingerprint_sha_256: string;
};

export type SSHKeys = SSHKey[];
