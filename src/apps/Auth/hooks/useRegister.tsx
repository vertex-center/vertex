import {
    useMutation,
    UseMutationOptions,
    useQueryClient,
} from "@tanstack/react-query";
import { api, setAuthToken } from "../../../backend/api/backend";
import { AuthCredentials } from "../../../models/auth";

export const useRegister = (
    options: UseMutationOptions<unknown, unknown, AuthCredentials>
) => {
    const { onSuccess, ...others } = options;
    const queryClient = useQueryClient();
    const mutation = useMutation({
        mutationKey: ["auth_register"],
        mutationFn: api.auth.register,
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
