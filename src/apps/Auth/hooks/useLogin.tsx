import {
    useMutation,
    UseMutationOptions,
    useQueryClient,
} from "@tanstack/react-query";
import { AuthCredentials } from "../backend/models";
import { API } from "../backend/api";
import { setAuthToken } from "../../../backend/server";

export const useLogin = (
    options: UseMutationOptions<unknown, unknown, AuthCredentials>
) => {
    const { onSuccess, ...others } = options;
    const queryClient = useQueryClient();
    const mutation = useMutation({
        mutationKey: ["auth_login"],
        mutationFn: API.login,
        onSuccess: (...args) => {
            const data: any = args[0];
            setAuthToken(data?.token);
            queryClient.invalidateQueries(["user"]);
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
