export type SSHKey = {
    type: string;
    fingerprint_sha_256: string;
    username: string;
};

export type SSHKeys = SSHKey[];
