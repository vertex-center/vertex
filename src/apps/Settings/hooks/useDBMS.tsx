import { useQuery } from "@tanstack/react-query";
import { api } from "../../../backend/api/backend";

export const useDBMS = () => {
    const query = useQuery({
        queryKey: ["admin_db_dbms"],
        queryFn: api.admin.data.dbms.get,
    });
    const { data: dbms, isLoading: isLoadingDbms, error: errorDbms } = query;
    return { dbms, isLoadingDbms, errorDbms };
};
