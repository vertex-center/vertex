import { useQuery } from "@tanstack/react-query";
import { API } from "../backend/api";

export const useHost = () => {
    const {
        data: host,
        error: errorHost,
        isLoading: isLoadingHost,
    } = useQuery({
        queryKey: ["hardware_host"],
        queryFn: API.getHost,
    });
    return { host, errorHost, isLoadingHost };
};
