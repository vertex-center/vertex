import { Port } from "../backend/models";
import { API } from "../backend/api";
import { useMutation, UseMutationOptions } from "@tanstack/react-query";

export function usePatchPort(
    options?: UseMutationOptions<unknown, unknown, Port>
) {
    const {
        mutate: patchPort,
        mutateAsync: patchPortAsync,
        ...rest
    } = useMutation({
        ...options,
        mutationFn: (port: Port) => API.patchPort(port),
    });
    return { patchPort, patchPortAsync, ...rest };
}
