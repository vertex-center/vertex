import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { api } from "../../../backend/api/backend";

export const useUpdateMutation = (options: UseMutationOptions) => {
    const mutation = useMutation({
        mutationKey: ["updates"],
        mutationFn: api.update.install,
        ...options,
    });
    const {
        mutate: installUpdate,
        isLoading: isInstallingUpdate,
        error: errorInstallUpdate,
    } = mutation;
    return { installUpdate, isInstallingUpdate, errorInstallUpdate };
};
