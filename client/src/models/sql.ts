type SqlDBMS = {
    username?: string;
    password?: string;
    databases?: SqlDatabase[];
};

type SqlDatabase = {
    name: string;
};
