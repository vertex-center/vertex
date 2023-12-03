import { useQuery } from "@tanstack/react-query";
import { API } from "../backend/api";

export const useSSHUsers = () => {
    const {
        data: sshUsers,
        error: sshUsersError,
        isLoading: isSSHUsersLoading,
    } = useQuery({
        queryKey: ["admin_ssh_users"],
        queryFn: API.getSSHUsers,
    });
    return { sshUsers, sshUsersError, isSSHUsersLoading };
};
