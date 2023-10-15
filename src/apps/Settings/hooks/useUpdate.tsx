import { useQuery } from "@tanstack/react-query";
import { api } from "../../../backend/api/backend";

export const useUpdate = () => {
    const query = useQuery({
        queryKey: ["updates"],
        queryFn: api.update.get,
    });
    const {
        data: update,
        isLoading: isLoadingUpdate,
        error: errorUpdate,
    } = query;
    return { update, isLoadingUpdate, errorUpdate };
};
