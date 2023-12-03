import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { API } from "../backend/api";
import { User } from "../backend/models";

export const usePatchUser = (
    options: UseMutationOptions<unknown, unknown, Partial<User>>
) => {
    const mutation = useMutation({
        mutationKey: ["user"],
        mutationFn: API.patchCurrentUser,
        ...options,
    });
    const {
        mutate: patchUser,
        isLoading: isPatchingUser,
        error: errorPatchUser,
    } = mutation;
    return { patchUser, isPatchingUser, errorPatchUser };
};
