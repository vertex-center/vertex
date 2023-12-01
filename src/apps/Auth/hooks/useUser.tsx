import { useQuery } from "@tanstack/react-query";
import { api } from "../../../backend/api/backend";

export default function useUser(username?: string) {
    const query = useQuery({
        queryKey: ["user", username],
        queryFn: api.auth.user().get,
    });
    const { data: user, isLoading: isLoadingUser, error: errorUser } = query;
    return { user, isLoadingUser, errorUser };
}
