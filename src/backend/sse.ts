export default class SSE {
    public events?: EventSource;

    public constructor() {
        console.log("SSE opened");
        this.events = new EventSource("http://localhost:6130/events");
        this.events.onmessage = (e: any) => console.log(e);
        this.events.onerror = (e: any) => console.error(e);
    }

    public on(key: string, handler: (e: any) => void) {
        this.events.addEventListener(key, handler);
    }

    public close() {
        console.log("SSE closed");
        this.events?.close();
    }
}
