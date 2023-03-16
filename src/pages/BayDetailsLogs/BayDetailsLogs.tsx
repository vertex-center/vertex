import Logs from "../../components/Logs/Logs";
import { useEffect, useState } from "react";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../backend/sse";
import { useParams } from "react-router-dom";

export default function BayDetailsLogs() {
    const { uuid } = useParams();

    const [logs, setLogs] = useState<any[]>([]);

    useEffect(() => {
        const sse = registerSSE(`http://localhost:6130/service/${uuid}/events`);

        const onStdout = (e) => {
            const logLine = JSON.parse(e.data);
            setLogs((logs) => [...logs, logLine]);
        };

        const onStderr = (e) => {
            const logLine = JSON.parse(e.data);
            setLogs((logs) => [...logs, logLine]);
        };

        registerSSEEvent(sse, "stdout", onStdout);
        registerSSEEvent(sse, "stderr", onStderr);

        return () => {
            unregisterSSEEvent(sse, "stdout", onStdout);
            unregisterSSEEvent(sse, "stderr", onStderr);

            unregisterSSE(sse);
        };
    }, [uuid]);

    return <Logs lines={logs} />;
}
