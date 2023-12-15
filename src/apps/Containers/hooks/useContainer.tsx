import { useQuery } from "@tanstack/react-query";
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
