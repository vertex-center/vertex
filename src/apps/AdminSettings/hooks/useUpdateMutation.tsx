import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { API } from "../backend/api";

export const useUpdateMutation = (options: UseMutationOptions) => {
    const mutation = useMutation({
        mutationKey: ["updates"],
        mutationFn: API.installUpdate,
        ...options,
    });
    const {
        mutate: installUpdate,
        isLoading: isInstallingUpdate,
        error: errorInstallUpdate,
    } = mutation;
    return { installUpdate, isInstallingUpdate, errorInstallUpdate };
};
