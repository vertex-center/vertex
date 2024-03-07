import { Port, PortFilters } from "../backend/models";
import { API } from "../backend/api";
import {
    useMutation,
    UseMutationOptions,
    useQuery,
} from "@tanstack/react-query";

export function usePorts(filters?: PortFilters) {
    const queryPorts = useQuery({
        queryKey: ["ports", filters],
        queryFn: () => API.getPorts(filters),
    });
    return {
        ...queryPorts,
        ports: queryPorts.data,
        isLoadingPorts: queryPorts.isLoading,
        errorPorts: queryPorts.error,
    };
}

export function usePatchPort(
    options?: UseMutationOptions<unknown, unknown, Port>
) {
    const {
        mutate: patchPort,
        mutateAsync: patchPortAsync,
        ...rest
    } = useMutation({
        ...options,
        mutationFn: (port: Port) => API.patchPort(port),
    });
    return { patchPort, patchPortAsync, ...rest };
}

export function useDeletePort(
    options?: UseMutationOptions<unknown, unknown, string>
) {
    const {
        mutate: deletePort,
        mutateAsync: deletePortAsync,
        ...rest
    } = useMutation({
        ...options,
        mutationFn: (id: string) => API.deletePort(id),
    });
    return { deletePort, deletePortAsync, ...rest };
}

export function useCreatePort(
    options?: UseMutationOptions<unknown, unknown, Port>
) {
    const {
        mutate: createPort,
        mutateAsync: createPortAsync,
        ...rest
    } = useMutation({
        ...options,
        mutationFn: (port: Port) => API.createPort(port),
    });
    return { createPort, createPortAsync, ...rest };
}
