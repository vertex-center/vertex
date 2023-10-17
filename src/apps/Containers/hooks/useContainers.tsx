import { useQuery } from "@tanstack/react-query";
import { api } from "../../../backend/api/backend";
import { ContainerQuery } from "../../../models/container";

export function useContainersTags() {
    const queryTags = useQuery({
        queryKey: ["containers", "tags"],
        queryFn: api.vxContainers.containers.tags,
    });
    const { data: tags } = queryTags;
    return { tags, ...queryTags };
}

export function useContainers(query: ContainerQuery) {
    const queryContainers = useQuery({
        queryKey: ["containers", query],
        queryFn: () => api.vxContainers.containers.search(query),
    });
    const { data: containers } = queryContainers;
    return { containers, ...queryContainers };
}
