import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { api } from "../../../backend/api/backend";

export const useSettingsChannelMutation = (
    options: UseMutationOptions<unknown, unknown, boolean>
) => {
    const mutation = useMutation<unknown, unknown, boolean, unknown>({
        mutationKey: ["settings"],
        mutationFn: (beta?: boolean) =>
            api.settings.patch({
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
