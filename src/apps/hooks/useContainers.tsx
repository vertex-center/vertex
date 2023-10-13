import { useQuery } from "@tanstack/react-query";
import { api } from "../../backend/api/backend";

export const useContainersTags = () => {
    const queryTags = useQuery({
        queryKey: ["containers", "tags"],
        queryFn: api.vxContainers.containers.tags,
    });
    const { data: tags } = queryTags;
    return { tags, ...queryTags };
};
