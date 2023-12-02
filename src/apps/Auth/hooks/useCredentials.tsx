import { useQuery } from "@tanstack/react-query";
import { api } from "../../../backend/api/backend";

export const useCredentials = () => {
    const query = useQuery({
        queryKey: ["credentials"],
        queryFn: api.auth.user().credentials.get,
    });
    const {
        data: credentials,
        isLoading: isLoadingCredentials,
        error: errorCredentials,
    } = query;
    return { credentials, isLoadingCredentials, errorCredentials };
};
