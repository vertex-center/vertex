import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { api, setAuthToken } from "../../../backend/api/backend";
import { Credentials } from "../../../models/auth";

export const useRegister = (
    options: UseMutationOptions<unknown, unknown, Credentials>
) => {
    const { onSuccess, ...others } = options;
    const mutation = useMutation({
        mutationKey: ["auth_register"],
        mutationFn: api.auth.register,
        onSuccess: (...args) => {
            const data: any = args[0];
            setAuthToken(data?.token);
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
