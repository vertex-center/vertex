import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { API, DeleteSSHKeyBody } from "../backend/api";

export const useDeleteSSHKey = (
    options: UseMutationOptions<unknown, unknown, DeleteSSHKeyBody>
) => {
    const {
        mutate: deleteKey,
        isLoading: isDeletingKey,
        error: errorDeleteKey,
    } = useMutation({
        mutationKey: ["admin_ssh_keys"],
        mutationFn: API.deleteSSHKey,
        ...options,
    });
    return { deleteKey, isDeletingKey, errorDeleteKey };
};
