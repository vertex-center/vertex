import { v4 as uuidv4 } from "uuid";

type SSE = {
    eventSource: EventSource;
    url: string;
    watchers: number;
};

const allSSE: { [uuid: string]: SSE } = {};

export function registerSSE(url: string): string {
    let uuid = Object.keys(allSSE).find((uuid) => allSSE[uuid].url === url);

    if (uuid !== undefined) {
        allSSE[uuid].watchers++;
        console.log(allSSE[uuid].watchers, "registered on", uuid);
        return uuid;
    }

    uuid = uuidv4();
    // @ts-ignore
    const eventSource = new EventSource(`${window.apiURL}/api${url}`);
    allSSE[uuid] = { url, eventSource, watchers: 1 };
    console.log("SSE", uuid, "opened.");
    console.log(allSSE[uuid].watchers, "registered on", uuid);

    return uuid;
}

export function unregisterSSE(uuid: string) {
    allSSE[uuid].watchers--;
    console.log(allSSE[uuid].watchers, "registered on", uuid);
    if (allSSE[uuid].watchers === 0) {
        console.log("SSE", uuid, "closed.");
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
