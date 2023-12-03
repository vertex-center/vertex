import { useQuery } from "@tanstack/react-query";
import { API } from "../backend/api";

export const useSSHKeys = () => {
    const {
        data: sshKeys,
        error: keysError,
        isLoading: isKeysLoading,
    } = useQuery({
        queryKey: ["admin_ssh_keys"],
        queryFn: API.getSSHKeys,
    });
    return { sshKeys, keysError, isKeysLoading };
};
