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
import { APIError } from "../../../../components/Error/ErrorBox";
import { ProgressOverlay } from "../../../../components/Progress/Progress";

export default function InstanceLogs() {
    const { uuid } = useParams();

    const [logs, setLogs] = useState<LogLine[]>([]);
    const [error, setError] = useState();
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (uuid === undefined) return;
        setLoading(true);
        api.instance.logs
            .get(uuid)
            .then((res) => setLogs(res.data))
            .catch(setError)
            .finally(() => setLoading(false));
    }, [uuid]);

    useEffect(() => {
        if (uuid === undefined) return;

        const sse = registerSSE(`/instance/${uuid}/events`);

        const onStdout = (e) => {
            console.error(e.data);

            setLogs((logs) => [
                ...logs,
                {
                    kind: "out",
                    message: JSON.parse(e.data),
                },
            ]);
        };

        const onStderr = (e) => {
            setLogs((logs) => [
                ...logs,
                {
                    kind: "err",
                    message: JSON.parse(e.data),
                },
            ]);
        };

        const onDownload = (e) => {
            setLogs((logs) => {
                const dl = JSON.parse(e.data);

                let downloads = [];
                if (
                    logs.length > 0 &&
                    logs[logs.length - 1].kind === "downloads"
                ) {
                    downloads = logs[logs.length - 1].message;
                }

                const i = downloads.findIndex((d) => d.id === dl.id);

                if (i === -1) {
                    downloads = [...downloads, dl];
                } else {
                    downloads[i] = dl;
                }

                if (logs.length === 0) return logs;
                if (logs[logs.length - 1].kind === "downloads") {
                    const lgs = [...logs];
                    lgs[logs.length - 1] = {
                        kind: "downloads",
                        message: downloads,
                    };
                    return lgs;
                } else {
                    return [
                        ...logs,
                        {
                            kind: "downloads",
                            message: downloads,
                        },
                    ];
                }
            });
        };

        registerSSEEvent(sse, "stdout", onStdout);
        registerSSEEvent(sse, "stderr", onStderr);
        registerSSEEvent(sse, "download", onDownload);

        return () => {
            unregisterSSEEvent(sse, "stdout", onStdout);
            unregisterSSEEvent(sse, "stderr", onStderr);
            unregisterSSEEvent(sse, "download", onDownload);

            unregisterSSE(sse);
        };
    }, [uuid]);

    if (!logs) return null;

    return (
        <Fragment>
            <ProgressOverlay show={loading} />
            <Title className={styles.title}>Logs</Title>
            <APIError error={error} />
            <Logs lines={logs} />
        </Fragment>
    );
}
