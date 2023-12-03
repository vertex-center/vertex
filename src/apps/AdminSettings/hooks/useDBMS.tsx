import { useQuery } from "@tanstack/react-query";
import { API } from "../backend/api";

export const useDBMS = () => {
    const query = useQuery({
        queryKey: ["admin_db_dbms"],
        queryFn: API.getDatabases,
    });
    const { data: dbms, isLoading: isLoadingDbms, error: errorDbms } = query;
    return { dbms, isLoadingDbms, errorDbms };
};
