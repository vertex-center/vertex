import Logs, { LogLine } from "../../components/Logs/Logs";
import { Fragment, useEffect, useState } from "react";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../backend/sse";
import { useParams } from "react-router-dom";
import { getLatestLogs } from "../../backend/backend";
import { Title } from "../../components/Text/Text";

export default function BayDetailsLogs() {
    const { uuid } = useParams();

    const [logs, setLogs] = useState<LogLine[]>([]);

    useEffect(() => {
        getLatestLogs(uuid)
            .then((res) => setLogs(res.data))
            .catch(console.error);
    }, []);

    useEffect(() => {
        if (uuid === undefined) return;

        const sse = registerSSE(`/instance/${uuid}/events`);

        const onStdout = (e) => {
            setLogs((logs) => [
                ...logs,
                {
                    kind: "out",
                    message: e.data,
                },
            ]);
        };

        const onStderr = (e) => {
            setLogs((logs) => [
                ...logs,
                {
                    kind: "err",
                    message: e.data,
                },
            ]);
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
