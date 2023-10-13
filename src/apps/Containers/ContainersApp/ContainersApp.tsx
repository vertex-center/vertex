import styles from "./ContainersApp.module.sass";
import Container, { Containers } from "../../../components/Container/Container";
import { BigTitle } from "../../../components/Text/Text";
import { api } from "../../../backend/api/backend";
import { APIError } from "../../../components/Error/APIError";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useServerEvent } from "../../../hooks/useEvent";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import Toolbar from "../../../components/Toolbar/Toolbar";
import Spacer from "../../../components/Spacer/Spacer";
import Button from "../../../components/Button/Button";
import SelectTags from "../components/SelectTags/SelectTags";
import { useState } from "react";

type ToolbarProps = {
    tags?: string[];
    onTagsChange?: (tags: string[]) => void;
};

const ToolbarContainers = (props: ToolbarProps) => {
    const { tags, onTagsChange } = props;

    return (
        <Toolbar className={styles.toolbar}>
            <SelectTags selected={tags} onChange={onTagsChange} />
            <Spacer />
            <Button primary to="/app/vx-containers/add" rightIcon="add">
                Create container
            </Button>
        </Toolbar>
    );
};

export default function ContainersApp() {
    const queryClient = useQueryClient();

    const [tags, setTags] = useState<string[]>([]);

    const queryContainers = useQuery({
        queryKey: ["containers", { tags }],
        queryFn: () => api.vxContainers.containers.search({ tags }),
    });
    const { data: containers, isLoading, isError, error } = queryContainers;

    const mutationPower = useMutation({
        mutationFn: async (uuid: string) => {
            if (
                containers[uuid].status === "off" ||
                containers[uuid].status === "error"
            ) {
                await api.vxContainers.container(uuid).start();
            } else {
                await api.vxContainers.container(uuid).stop();
            }
        },
    });

    let status = "Waiting server response...";
    if (queryContainers.isSuccess) {
        status = "running";
    } else if (queryContainers.isError) {
        status = "off";
    }

    useServerEvent("/app/vx-containers/containers/events", {
        change: () => {
            queryClient.invalidateQueries({
                queryKey: ["containers"],
            });
        },
        open: () => {
            queryClient.invalidateQueries({
                queryKey: ["containers"],
            });
        },
    });

    const onTagsChange = (tags: string[]) => {
        setTags((prev) => {
            queryClient.invalidateQueries({
                queryKey: ["containers", { tags: prev }],
                exact: true,
            });
            return tags;
        });
    };

    return (
        <div className={styles.server}>
            <ProgressOverlay show={isLoading} />
            <div className={styles.title}>
                <BigTitle>All containers</BigTitle>
            </div>

            {error && (
                <div className={styles.containers}>
                    <APIError error={error} />
                </div>
            )}

            <div className={styles.containers}>
                <ToolbarContainers tags={tags} onTagsChange={onTagsChange} />
                {!isLoading && !isError && (
                    <Containers>
                        <Container
                            container={{
                                value: {
                                    display_name: "Vertex",
                                    status,
                                },
                            }}
                        />
                        {Object.values(containers)?.map((inst) => (
                            <Container
                                key={inst.uuid}
                                container={{
                                    value: inst,
                                    to: `/app/vx-containers/${inst.uuid}/`,
                                    onPower: async () =>
                                        mutationPower.mutate(inst.uuid),
                                }}
                            />
                        ))}
                    </Containers>
                )}
            </div>
        </div>
    );
}
