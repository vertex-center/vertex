import { useEffect } from "react";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../backend/sse";
import { Console } from "../logging/logging";

type ServerEvent = (e: MessageEvent) => void;

export function useServerEvent(
    route: string,
    events: {
        [name: string]: ServerEvent;
    },
    disabled?: boolean
) {
    useEffect(() => {
        if (disabled) {
            return;
        }

        const sse = registerSSE(route);

        Console.event("SSE connected\n%O", { id: sse });

        Object.entries(events).forEach(([name, event]) => {
            registerSSEEvent(sse, name, event);
        });

        return () => {
            Object.entries(events).forEach(([name, event]) => {
                unregisterSSEEvent(sse, name, event);
            });
            unregisterSSE(sse);

            Console.event("SSE disconnected\n%O", { id: sse });
        };
    }, [route, disabled]);
}
