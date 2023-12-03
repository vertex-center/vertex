import { AuthCredentials, Credentials, User } from "./models";
import { server } from "../../../backend/api/backend";

const login = async (credentials: AuthCredentials) => {
    const Authorization = `Basic ${btoa(
        credentials.username + ":" + credentials.password
    )}`;
    const { data } = await server.post(
        "/app/auth/login",
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
        "/app/auth/register",
        {},
        { headers: { Authorization } }
    );
    return data;
};

const logout = async () => {
    const { data } = await server.post("/app/auth/logout");
    return data;
};

const getCurrentUser = async () => {
    const { data } = await server.get<User>("/app/auth/user");
    return data;
};

const patchCurrentUser = async (user: Partial<User>) => {
    const { data } = await server.patch("/app/auth/user", user);
    return data;
};

const getCredentialsCurrentUser = async () => {
    const { data } = await server.get<Credentials[]>(
        "/app/auth/user/credentials"
    );
    return data;
};

export const API = {
    login,
    register,
    logout,
    getCurrentUser,
    patchCurrentUser,
    getCredentialsCurrentUser,
};
