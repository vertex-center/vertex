import Logs from "../../../../components/Logs/Logs";
import { Fragment } from "react";
import { useParams } from "react-router-dom";
import { api } from "../../../../backend/backend";
import { Title } from "../../../../components/Text/Text";
import styles from "./InstanceLogs.module.sass";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useServerEvent } from "../../../../hooks/useEvent";
import { useQuery, useQueryClient } from "@tanstack/react-query";

export default function InstanceLogs() {
    const { uuid } = useParams();
    const queryClient = useQueryClient();

    const queryLogs = useQuery({
        queryKey: ["instance_logs", uuid],
        queryFn: api.vxInstances.instance(uuid).logs.get,
    });
    const { data: logs, isLoading, error } = queryLogs;

    const onStdout = (e: MessageEvent) => {
        queryClient.setQueryData(["instance_logs", uuid], (logs: any[]) => [
            ...logs,
            {
                kind: "out",
                message: JSON.parse(e.data),
            },
        ]);
    };

    const onStderr = (e: MessageEvent) => {
        queryClient.setQueryData(["instance_logs", uuid], (logs: any[]) => [
            ...logs,
            {
                kind: "err",
                message: JSON.parse(e.data),
            },
        ]);
    };

    const onDownload = (e: MessageEvent) => {
        queryClient.setQueryData(["instance_logs", uuid], (logs: any[]) => {
            const dl = JSON.parse(e.data);

            let downloads = [];
            if (logs.length > 0 && logs[logs.length - 1].kind === "downloads") {
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

    const route = uuid ? `/app/vx-instances/instance/${uuid}/events` : "";

    useServerEvent(route, {
        stdout: onStdout,
        stderr: onStderr,
        download: onDownload,
    });

    if (!logs) return null;

    return (
        <Fragment>
            <ProgressOverlay show={isLoading} />
            <Title className={styles.title}>Logs</Title>
            <APIError error={error} />
            <Logs lines={logs} />
        </Fragment>
    );
}
