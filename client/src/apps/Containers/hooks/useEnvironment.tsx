import { API } from "../backend/api";
import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { EnvVariable } from "../backend/models";

export function usePatchEnv(
    options?: UseMutationOptions<unknown, unknown, EnvVariable>
) {
    const {
        mutate: patchEnv,
        mutateAsync: patchEnvAsync,
        ...rest
    } = useMutation({
        ...options,
        mutationFn: (port: EnvVariable) => API.patchEnv(port),
    });
    return { patchEnv, patchEnvAsync, ...rest };
}

export function useDeleteEnv(
    options?: UseMutationOptions<unknown, unknown, string>
) {
    const {
        mutate: deleteEnv,
        mutateAsync: deleteEnvAsync,
        ...rest
    } = useMutation({
        ...options,
        mutationFn: (id: string) => API.deleteEnv(id),
    });
    return { deleteEnv, deleteEnvAsync, ...rest };
}

export function useCreateEnv(
    options?: UseMutationOptions<unknown, unknown, EnvVariable>
) {
    const {
        mutate: createEnv,
        mutateAsync: createEnvAsync,
        ...rest
    } = useMutation({
        ...options,
        mutationFn: (port: EnvVariable) => API.createEnv(port),
    });
    return { createEnv, createEnvAsync, ...rest };
}
