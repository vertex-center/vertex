import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { API } from "../backend/api";

export const usePatchSettings = (
    options: UseMutationOptions<unknown, unknown, Partial<Settings>>
) => {
    const mutation = useMutation({
        mutationKey: ["settings"],
        mutationFn: API.patchSettings,
        ...options,
    });
    const {
        mutate: patchSettings,
        isPending: isPatchingSettings,
        error: errorPatchingSettings,
    } = mutation;
    return { patchSettings, isPatchingSettings, errorPatchingSettings };
};
