import {
    useMutation,
    UseMutationOptions,
    useQueryClient,
} from "@tanstack/react-query";
import { API } from "../backend/api";
import { setAuthToken } from "../../../backend/server";

export const useLogout = (options: UseMutationOptions) => {
    const { onSuccess, ...others } = options;
    const queryClient = useQueryClient();
    const mutation = useMutation({
        mutationKey: ["auth_logout"],
        mutationFn: API.logout,
        onSuccess: (...args) => {
            setAuthToken(undefined);
            queryClient.invalidateQueries(["user"]);
            options.onSuccess?.(...args);
        },
        ...others,
    });
    const {
        mutate: logout,
        isLoading: isLoggingOut,
        error: errorLogout,
    } = mutation;
    return { logout, isLoggingOut, errorLogout };
};
