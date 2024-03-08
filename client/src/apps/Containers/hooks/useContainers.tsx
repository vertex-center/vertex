import { useMutation, useQuery } from "@tanstack/react-query";
import { API } from "../backend/api";
import { ContainerFilters } from "../backend/models";

export function useContainersTags() {
    const queryTags = useQuery({
        queryKey: ["containers", "tags"],
        queryFn: API.getAllTags,
    });
    const { data: tags } = queryTags;
    return { tags, ...queryTags };
}

export function useContainers(query: ContainerFilters) {
    const queryContainers = useQuery({
        queryKey: ["containers", query],
        queryFn: () => API.getContainers(query),
    });
    const { data: containers } = queryContainers;
    return { containers, ...queryContainers };
}

export function useReloadContainer() {
    const mutationReload = useMutation({
        mutationKey: ["containers"],
        mutationFn: (id: string) => API.reloadContainer(id),
    });
    return {
        ...mutationReload,
        reloadContainer: mutationReload.mutate,
    };
}
