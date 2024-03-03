import { useQuery } from "@tanstack/react-query";
import { API } from "../backend/api";

export const useUpdate = () => {
    const query = useQuery({
        queryKey: ["updates"],
        queryFn: API.getUpdate,
    });
    const {
        data: update,
        isLoading: isLoadingUpdate,
        error: errorUpdate,
    } = query;
    return { update, isLoadingUpdate, errorUpdate };
};
