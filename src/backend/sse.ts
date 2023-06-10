import { v4 as uuidv4 } from "uuid";

type UUID = string;

type SSE = {
    eventSource: EventSource;
    url: string;
    watchers: number;
};

const allSSE: { [uuid: UUID]: SSE } = {};

export function registerSSE(url: string): UUID {
    let uuid = Object.keys(allSSE).find((uuid) => allSSE[uuid].url === url);

    if (uuid !== undefined) {
        allSSE[uuid].watchers++;
        console.log(allSSE[uuid].watchers, "registered on", uuid);
        return uuid;
    }

    uuid = uuidv4();
    const eventSource = new EventSource("http://localhost:6130/api" + url);
    allSSE[uuid] = { url, eventSource, watchers: 1 };
    console.log("SSE", uuid, "opened.");
    console.log(allSSE[uuid].watchers, "registered on", uuid);

    return uuid;
}

export function unregisterSSE(uuid: UUID) {
    allSSE[uuid].watchers--;
    console.log(allSSE[uuid].watchers, "registered on", uuid);
    if (allSSE[uuid].watchers === 0) {
        console.log("SSE", uuid, "closed.");
        allSSE[uuid].eventSource.close();
        delete allSSE[uuid];
    }
}

export function registerSSEEvent(
    uuid: UUID,
    key: string,
    handler: (e: any) => void
) {
    allSSE[uuid].eventSource.addEventListener(key, handler);
}

export function unregisterSSEEvent(
    uuid: UUID,
    key: string,
    handler: (e: any) => void
) {
    allSSE[uuid].eventSource.removeEventListener(key, handler);
}
