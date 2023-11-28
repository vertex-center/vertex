import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { api } from "../../../backend/api/backend";

export const useLogout = (options: UseMutationOptions) => {
    const mutation = useMutation({
        mutationKey: ["auth_logout"],
        mutationFn: api.auth.logout,
        ...options,
    });
    const {
        mutate: logout,
        isLoading: isLoggingOut,
        error: errorLogout,
    } = mutation;
    return { logout, isLoggingOut, errorLogout };
};
