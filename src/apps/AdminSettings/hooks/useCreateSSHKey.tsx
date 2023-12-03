import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { AddSSHKeyBody, API } from "../backend/api";

export const useCreateSSHKey = (
    options: UseMutationOptions<unknown, unknown, AddSSHKeyBody>
) => {
    const {
        mutate: createKey,
        isLoading: isCreatingKey,
        error: errorCreateKey,
        reset: resetCreateKey,
    } = useMutation({
        mutationKey: ["admin_ssh_keys"],
        mutationFn: API.addSSHKey,
        ...options,
    });
    return {
        createKey,
        isCreatingKey,
        errorCreateKey,
        resetCreateKey,
    };
};
