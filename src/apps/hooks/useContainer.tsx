import { api } from "../../backend/api/backend";
import { useQuery } from "@tanstack/react-query";

export default function useContainer(uuid?: string) {
    const queryContainer = useQuery({
        queryKey: ["containers", uuid],
        queryFn: api.vxContainers.container(uuid).get,
    });
    return { container: queryContainer.data, ...queryContainer };
}
