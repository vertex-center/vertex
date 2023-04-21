import Logs from "../../components/Logs/Logs";
import { Fragment, useEffect, useState } from "react";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../backend/sse";
import { useParams } from "react-router-dom";
import { route } from "../../backend/backend";
import { Title } from "../../components/Text/Text";
import useInstance from "../../hooks/useInstance";

export default function BayDetailsLogs() {
    const { uuid } = useParams();

    const { instance } = useInstance(uuid);

    const [logs, setLogs] = useState<any[]>(undefined);

    useEffect(() => {
        setLogs(instance?.logs.lines ?? []);
    }, [instance]);

    useEffect(() => {
        if (uuid === undefined) return;

        const sse = registerSSE(route(`/instance/${uuid}/events`));

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

    if (!logs) return null;

    return (
        <Fragment>
            <Title>Logs</Title>
            <Logs lines={logs} />
        </Fragment>
    );
}
