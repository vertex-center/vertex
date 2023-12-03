import { useQuery } from "@tanstack/react-query";
import { API } from "../backend/api";

export const useCPUs = () => {
    const {
        data: cpus,
        error: errorCPUs,
        isLoading: isLoadingCPUs,
    } = useQuery({
        queryKey: ["hardware_cpus"],
        queryFn: API.getCPUs,
    });
    return { cpus, errorCPUs, isLoadingCPUs };
};
