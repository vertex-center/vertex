export default class SSE {
    public events?: EventSource;

    public constructor(url: string) {
        console.log("SSE opened");
        this.events = new EventSource(url);
        this.events.onerror = console.error;
    }

    public on(key: string, handler: (e: any) => void) {
        this.events.addEventListener(key, (e) => {
            console.log("%c SSE ", "background-color:orange;color:black;", {
                event: key,
                data: e.data,
            });
            handler(e);
        });
    }

    public close() {
        console.log("SSE closed");
        this.events?.close();
    }
}
