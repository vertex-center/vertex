import Logs from "../../../../components/Logs/Logs";
import { useParams } from "react-router-dom";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useServerEvent } from "../../../../hooks/useEvent";
import { useQueryClient } from "@tanstack/react-query";
import { produce } from "immer";
import { useContainerLogs } from "../../hooks/useContainer";
import { Title, Vertical } from "@vertex-center/components";

export default function ContainerLogs() {
    const { uuid } = useParams();
    const queryClient = useQueryClient();

    const { data: logs, isLoading, error } = useContainerLogs(uuid);

    const onStdout = (e: MessageEvent) => {
        queryClient.setQueryData(["container_logs", uuid], (logs: any[]) => [
            ...logs,
            {
                kind: "out",
                message: JSON.parse(e.data),
            },
        ]);
    };

    const onStderr = (e: MessageEvent) => {
        queryClient.setQueryData(["container_logs", uuid], (logs: any[]) => [
            ...logs,
            {
                kind: "err",
                message: JSON.parse(e.data),
            },
        ]);
    };

    const onDownload = (e: MessageEvent) => {
        queryClient.setQueryData(["container_logs", uuid], (logs: any[]) => {
            return produce(logs, (draft) => {
                const dl = JSON.parse(e.data);

                let downloads = [];
                if (
                    draft.length > 0 &&
                    draft[draft.length - 1].kind === "downloads"
                ) {
                    downloads = draft[draft.length - 1].message;
                }

                const i = downloads.findIndex((d) => d.id === dl.id);

                if (i === -1) {
                    downloads = [...downloads, dl];
                } else {
                    downloads[i] = dl;
                }

                if (draft.length === 0) return draft;
                if (draft[draft.length - 1].kind === "downloads") {
                    draft[draft.length - 1] = {
                        kind: "downloads",
                        message: downloads,
                    };
                    return draft;
                } else {
                    draft.push({
                        kind: "downloads",
                        message: downloads,
                    });
                    return draft;
                }
            });
        });
    };

    const route = uuid ? `/app/containers/container/${uuid}/events` : "";

    useServerEvent(route, {
        stdout: onStdout,
        stderr: onStderr,
        download: onDownload,
    });

    if (!logs) return null;

    return (
        <Vertical gap={24}>
            <Title variant="h2">Logs</Title>
            <APIError error={error} />
            <ProgressOverlay show={isLoading} />
            <Logs lines={logs} />
        </Vertical>
    );
}
