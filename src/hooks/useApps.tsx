import { useQuery } from "@tanstack/react-query";
import { api } from "../backend/api/backend";

export const useApps = () => {
    const queryApps = useQuery({
        queryKey: ["apps"],
        queryFn: api.apps.all,
    });
    const { data: apps } = queryApps;
    return { apps };
};
