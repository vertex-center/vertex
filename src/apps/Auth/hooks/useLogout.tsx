import {
    useMutation,
    UseMutationOptions,
    useQueryClient,
} from "@tanstack/react-query";
import { setAuthToken } from "../../../backend/api/backend";
import { API } from "../backend/api";

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
