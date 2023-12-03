import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { API } from "../backend/api";

export const useMigrateDatabase = (
    options: UseMutationOptions<unknown, unknown, string>
) => {
    const mutation = useMutation({
        mutationKey: ["admin_db_dbms"],
        mutationFn: API.migrateDatabase,
        ...options,
    });
    const {
        mutate: migrate,
        isLoading: isMigrating,
        error: errorMigrate,
    } = mutation;
    return { migrate, isMigrating, errorMigrate };
};
