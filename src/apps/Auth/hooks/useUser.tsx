import { useQuery } from "@tanstack/react-query";
import { api } from "../../../backend/api/backend";

export default function useUser(username?: string) {
    const query = useQuery({
        queryKey: ["user", username],
        queryFn: api.auth.user().get,
        retry: (failureCount, error) => {
            // Don't retry too much if the error was caused by an authentication issue
            // @ts-ignore
            if (error?.response?.status === 401) return failureCount < 2;
            return failureCount < 3;
        },
    });
    const { data: user, isLoading: isLoadingUser, error: errorUser } = query;
    return { user, isLoadingUser, errorUser };
}
