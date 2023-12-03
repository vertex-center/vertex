import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { API } from "../backend/api";

export const useReboot = (options: UseMutationOptions) => {
    const {
        mutate: reboot,
        isLoading: isRebooting,
        error: errorReboot,
    } = useMutation({
        mutationKey: ["admin_reboot"],
        mutationFn: API.reboot,
        ...options,
    });
    return { reboot, isRebooting, errorReboot };
};
