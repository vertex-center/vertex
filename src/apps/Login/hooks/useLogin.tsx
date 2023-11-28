import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { api, setAuthToken } from "../../../backend/api/backend";
import { Credentials } from "../../../models/auth";

export const useLogin = (
    options: UseMutationOptions<unknown, unknown, Credentials>
) => {
    const { onSuccess, ...others } = options;
    const mutation = useMutation({
        mutationKey: ["auth_login"],
        mutationFn: api.auth.login,
        onSuccess: (...args) => {
            const data: any = args[0];
            setAuthToken(data?.token);
            options.onSuccess?.(...args);
        },
        ...others,
    });
    const {
        mutate: login,
        isLoading: isLoggingIn,
        error: errorLogin,
    } = mutation;
    return { login, isLoggingIn, errorLogin };
};
