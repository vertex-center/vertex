import { useQuery } from "@tanstack/react-query";
import { api } from "../../../backend/api/backend";

export const useSettings = () => {
    const query = useQuery({
        queryKey: ["settings"],
        queryFn: api.settings.get,
    });
    const {
        data: settings,
        isLoading: isLoadingSettings,
        error: errorSettings,
    } = query;
    return { settings, isLoadingSettings, errorSettings };
};
