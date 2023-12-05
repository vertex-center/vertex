import { AuthCredentials, Credentials, Email, User } from "./models";

import { createServer } from "../../../backend/server";

const server = createServer("7502");

const login = async (credentials: AuthCredentials) => {
    const Authorization = `Basic ${btoa(
        credentials.username + ":" + credentials.password
    )}`;
    const { data } = await server.post(
        "/login",
        {},
        { headers: { Authorization } }
    );
    return data;
};

const register = async (credentials: AuthCredentials) => {
    const Authorization = `Basic ${btoa(
        credentials.username + ":" + credentials.password
    )}`;
    const { data } = await server.post(
        "/register",
        {},
        { headers: { Authorization } }
    );
    return data;
};

const logout = async () => {
    const { data } = await server.post("/logout");
    return data;
};

const getCurrentUser = async () => {
    const { data } = await server.get<User>("/user");
    return data;
};

const patchCurrentUser = async (user: Partial<User>) => {
    const { data } = await server.patch("/user", user);
    return data;
};

const getCredentialsCurrentUser = async () => {
    const { data } = await server.get<Credentials[]>("/user/credentials");
    return data;
};

const getEmailsCurrentUser = async () => {
    const { data } = await server.get<Email[]>("/user/emails");
    return data;
};

const postEmailCurrentUser = async (email: Partial<Email>) => {
    const { data } = await server.post("/user/email", email);
    return data;
};

const deleteEmailCurrentUser = async (email: Partial<Email>) => {
    const { data } = await server.delete("/user/email", {
        data: email,
    });
    return data;
};

export const API = {
    login,
    register,
    logout,
    getCurrentUser,
    patchCurrentUser,
    getCredentialsCurrentUser,
    getEmailsCurrentUser,
    postEmailCurrentUser,
    deleteEmailCurrentUser,
};
