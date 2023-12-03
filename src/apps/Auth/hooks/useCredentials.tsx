import { useQuery } from "@tanstack/react-query";
import { API } from "../backend/api";

export const useCredentials = () => {
    const query = useQuery({
        queryKey: ["credentials"],
        queryFn: API.getCredentialsCurrentUser,
    });
    const {
        data: credentials,
        isLoading: isLoadingCredentials,
        error: errorCredentials,
    } = query;
    return { credentials, isLoadingCredentials, errorCredentials };
};
