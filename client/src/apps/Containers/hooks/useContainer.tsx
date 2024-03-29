import {
    useMutation,
    UseMutationOptions,
    useQuery,
} from "@tanstack/react-query";
import { API } from "../backend/api";

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

export function useRecreateContainer(
    options?: UseMutationOptions<unknown, unknown, string>
) {
    const queryRecreate = useMutation({
        ...options,
        mutationFn: (id: string) => API.recreateDocker(id),
    });
    return {
        ...queryRecreate,
        recreateContainer: queryRecreate.mutate,
        isPendingRecreate: queryRecreate.isPending,
        errorRecreate: queryRecreate.error,
    };
}
