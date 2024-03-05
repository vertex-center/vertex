import {
    useMutation,
    UseMutationOptions,
    useQuery,
    useQueryClient,
} from "@tanstack/react-query";
import { API } from "../backend/api";
import { Port } from "../backend/models";

export default function useContainer(id?: string) {
    const queryContainer = useQuery({
        queryKey: ["containers", id],
        queryFn: () => API.getContainer(id),
    });
    return { container: queryContainer.data, ...queryContainer };
}

export function useContainerLogs(id?: string) {
    const queryLogs = useQuery({
        queryKey: ["container_logs", id],
        queryFn: () => API.getLogs(id),
    });
    return { logs: queryLogs.data, ...queryLogs };
}

export function useDockerInfo(id?: string) {
    const queryDockerInfo = useQuery({
        queryKey: ["container_docker", id],
        queryFn: () => API.getDockerInfo(id),
    });
    return { dockerInfo: queryDockerInfo.data, ...queryDockerInfo };
}

export function useContainerEnv(id?: string) {
    const queryEnv = useQuery({
        queryKey: ["container_env", id],
        queryFn: () => API.getContainerEnvironment(id),
    });
    return {
        env: queryEnv.data,
        isLoadingEnv: queryEnv.isLoading,
        errorEnv: queryEnv.error,
    };
}

export function useContainerPorts(id?: string) {
    const queryPorts = useQuery({
        queryKey: ["container_ports", id],
        queryFn: () => API.getContainerPorts(id),
    });
    return {
        ports: queryPorts.data,
        isLoadingPorts: queryPorts.isLoading,
        errorPorts: queryPorts.error,
    };
}

export function useSaveContainerPorts(
    id?: string,
    options?: UseMutationOptions<unknown, unknown, Port[]>
) {
    const queryClient = useQueryClient();
    const { mutate: savePorts, ...rest } = useMutation({
        ...options,
        mutationFn: (ports: Port[]) => API.saveContainerPorts(id, ports),
        onSettled: async (...args) => {
            await queryClient.invalidateQueries({
                queryKey: ["container_ports", id],
            });
            options?.onSettled?.(...args);
        },
    });
    return { savePorts, ...rest };
}
