import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { api } from "../../../backend/api/backend";

export const useDBMSMutation = (
    options: UseMutationOptions<unknown, unknown, string>
) => {
    const mutation = useMutation({
        mutationKey: ["admin_db_dbms"],
        mutationFn: api.admin.data.dbms.migrate,
        ...options,
    });
    const {
        mutate: migrate,
        isLoading: isMigrating,
        error: errorMigrate,
    } = mutation;
    return { migrate, isMigrating, errorMigrate };
};
