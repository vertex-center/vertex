import Logs, { LogLine } from "../../../../components/Logs/Logs";
import { Fragment, useEffect, useState } from "react";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../../../backend/sse";
import { useParams } from "react-router-dom";
import { api } from "../../../../backend/backend";
import { Title } from "../../../../components/Text/Text";
import styles from "./InstanceLogs.module.sass";
import { APIError } from "../../../../components/Error/Error";

export default function InstanceLogs() {
    const { uuid } = useParams();

    const [logs, setLogs] = useState<LogLine[]>([]);
    const [error, setError] = useState();

    useEffect(() => {
        api.instance.logs
            .get(uuid)
            .then((res) => setLogs(res.data))
            .catch(setError);
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
            <Title className={styles.title}>Logs</Title>
            <APIError error={error} />
            <Logs lines={logs} />
        </Fragment>
    );
}
