import { v4 as uuidv4 } from "uuid";
import { EventSourcePolyfill } from "event-source-polyfill";

import { getAuthToken } from "./server";

type SSE = {
    eventSource: EventSourcePolyfill;
    url: string;
    watchers: number;
};

const allSSE: { [uuid: string]: SSE } = {};

export function registerSSE(port: string, url: string): string {
    let uuid = Object.keys(allSSE).find((uuid) => allSSE[uuid].url === url);

    if (uuid !== undefined) {
        allSSE[uuid].watchers++;
        return uuid;
    }

    uuid = uuidv4();
    const eventSource = new EventSourcePolyfill(
        // @ts-ignore
        `${window.apiURL}:${port}/api${url}`,
        {
            headers: {
                Authorization: `Bearer ${getAuthToken()}`,
            },
        }
    );
    allSSE[uuid] = { url, eventSource, watchers: 1 };

    return uuid;
}

export function unregisterSSE(uuid: string) {
    allSSE[uuid].watchers--;
    if (allSSE[uuid].watchers === 0) {
        allSSE[uuid].eventSource.close();
        delete allSSE[uuid];
    }
}

export function registerSSEEvent(
    uuid: string,
    key: string,
    handler: (e: any) => void
) {
    allSSE[uuid].eventSource.addEventListener(key, handler);
}

export function unregisterSSEEvent(
    uuid: string,
    key: string,
    handler: (e: any) => void
) {
    allSSE[uuid].eventSource.removeEventListener(key, handler);
}
