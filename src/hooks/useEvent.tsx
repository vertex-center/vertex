import { useEffect } from "react";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../backend/sse";

type ServerEvent = (e: MessageEvent) => void;

export function useServerEvent(
    route: string,
    events: {
        [name: string]: ServerEvent;
    }
) {
    useEffect(() => {
        console.log("useServerEvent", route, events);

        const sse = registerSSE(route);

        Object.entries(events).forEach(([name, event]) => {
            registerSSEEvent(sse, name, event);
        });

        return () => {
            Object.entries(events).forEach(([name, event]) => {
                unregisterSSEEvent(sse, name, event);
            });
            unregisterSSE(sse);
        };
    }, [route]);
}
