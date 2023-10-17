import { api } from "../../../backend/api/backend";
import { useQuery } from "@tanstack/react-query";

export default function useContainer(uuid?: string) {
    const queryContainer = useQuery({
        queryKey: ["containers", uuid],
        queryFn: api.vxContainers.container(uuid).get,
    });
    return { container: queryContainer.data, ...queryContainer };
}

export function useContainerLogs(uuid?: string) {
    const queryLogs = useQuery({
        queryKey: ["container_logs", uuid],
        queryFn: api.vxContainers.container(uuid).logs.get,
    });
    return { logs: queryLogs.data, ...queryLogs };
}

export function useDockerInfo(uuid?: string) {
    const queryDockerInfo = useQuery({
        queryKey: ["container_docker", uuid],
        queryFn: api.vxContainers.container(uuid).docker.get,
    });
    return { dockerInfo: queryDockerInfo.data, ...queryDockerInfo };
}
