import { api } from "../backend/backend";
import { useQuery } from "@tanstack/react-query";

export default function useInstance(uuid?: string) {
    const queryInstance = useQuery({
        queryKey: ["instances", uuid],
        queryFn: api.vxInstances.instance(uuid).get,
    });
    return { instance: queryInstance.data, ...queryInstance };
}
