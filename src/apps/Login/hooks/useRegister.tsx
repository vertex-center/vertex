import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { api } from "../../../backend/api/backend";
import { Credentials } from "../../../models/auth";

export const useRegister = (
    options: UseMutationOptions<unknown, unknown, Credentials>
) => {
    const mutation = useMutation({
        mutationKey: ["auth_register"],
        mutationFn: api.auth.register,
        ...options,
    });
    const {
        mutate: register,
        isLoading: isRegistering,
        error: errorRegister,
    } = mutation;
    return { register, isRegistering, errorRegister };
};
