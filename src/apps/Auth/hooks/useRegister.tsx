import {
    useMutation,
    UseMutationOptions,
    useQueryClient,
} from "@tanstack/react-query";
import { API } from "../backend/api";
import { AuthCredentials } from "../backend/models";
import { setAuthToken } from "../../../backend/server";

export const useRegister = (
    options: UseMutationOptions<unknown, unknown, AuthCredentials>
) => {
    const { onSuccess, ...others } = options;
    const queryClient = useQueryClient();
    const mutation = useMutation({
        mutationKey: ["auth_register"],
        mutationFn: API.register,
        onSuccess: (...args) => {
            const data: any = args[0];
            setAuthToken(data?.token);
            queryClient.invalidateQueries(["user"]);
            options.onSuccess?.(...args);
        },
        ...others,
    });
    const {
        mutate: register,
        isLoading: isRegistering,
        error: errorRegister,
    } = mutation;
    return { register, isRegistering, errorRegister };
};
