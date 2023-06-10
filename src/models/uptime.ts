export type Uptime = {
    name: string;
    ping_url?: string;
    current: string;
    interval_seconds: number;
    remaining_seconds: number;
    history: {
        status: string;
    }[];
};
