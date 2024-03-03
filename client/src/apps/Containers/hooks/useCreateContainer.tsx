import {
    useMutation,
    UseMutationOptions,
    useQueryClient,
} from "@tanstack/react-query";
import { API, CreateContainerOptions } from "../backend/api";

export const useCreateContainer = (
    options: UseMutationOptions<unknown, unknown, CreateContainerOptions>
) => {
    const { onSuccess, ...others } = options;
    const queryClient = useQueryClient();
    const mutation = useMutation({
        mutationKey: ["containers_create"],
        mutationFn: API.createContainer,
        onSuccess: (...args) => {
            queryClient.invalidateQueries({ queryKey: ["containers"] });
            options.onSuccess?.(...args);
        },
        ...others,
    });
    const {
        mutate: createContainer,
        isPending: isCreatingContainer,
        error: errorCreatingContainer,
    } = mutation;
    return { createContainer, isCreatingContainer, errorCreatingContainer };
};
