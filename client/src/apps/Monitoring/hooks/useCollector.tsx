import { api } from "../../../backend/api/backend";
import { useQuery } from "@tanstack/react-query";

export const useCollector = (name: string) => {
    const {
        data: collector,
        isLoading: isLoadingCollector,
        error: errorCollector,
    } = useQuery({
        queryKey: ["monitoring_collector", name],
        queryFn: api.vxMonitoring.collector(name).get,
    });

    return {
        collector,
        isLoadingCollector,
        errorCollector,
    };
};
