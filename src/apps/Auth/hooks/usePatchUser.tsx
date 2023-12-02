import { api } from "../../../backend/api/backend";
import { useMutation, UseMutationOptions } from "@tanstack/react-query";

export const usePatchUser = (
    options: UseMutationOptions<unknown, unknown, Partial<User>>
) => {
    const mutation = useMutation({
        mutationKey: ["user"],
        mutationFn: api.auth.user().patch,
        ...options,
    });
    const {
        mutate: patchUser,
        isLoading: isPatchingUser,
        error: errorPatchUser,
    } = mutation;
    return { patchUser, isPatchingUser, errorPatchUser };
};
