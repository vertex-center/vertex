export type AuthCredentials = {
    username: string;
    password: string;
};

export type Credentials = {
    name: string;
    description: string;
};

export type User = {
    id: number;
    username: string;
};

export type Email = {
    id: number;
    user_id: number;
    email: string;
    created_at: number;
    updated_at: number;
    deleted_at: number;
};
