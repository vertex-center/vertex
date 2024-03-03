import { useQuery } from "@tanstack/react-query";
import { API } from "../backend/api";

export const useSettings = () => {
    const query = useQuery({
        queryKey: ["settings"],
        queryFn: API.getSettings,
    });
    const {
        data: settings,
        isLoading: isLoadingSettings,
        error: errorSettings,
    } = query;
    return { settings, isLoadingSettings, errorSettings };
};
