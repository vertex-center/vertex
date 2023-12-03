import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { API } from "../backend/api";

export const usePatchSettings = (
    options: UseMutationOptions<unknown, unknown, boolean>
) => {
    const mutation = useMutation<unknown, unknown, boolean, unknown>({
        mutationKey: ["settings"],
        mutationFn: (beta?: boolean) =>
            API.patchSettings({
                updates_channel: beta ? "beta" : "stable",
            }),
        ...options,
    });
    const {
        mutate: setChannel,
        isLoading: isSettingChannel,
        error: errorSetChannel,
    } = mutation;
    return { setChannel, isSettingChannel, errorSetChannel };
};
