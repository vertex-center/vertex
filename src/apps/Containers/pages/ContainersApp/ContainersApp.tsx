import styles from "./ContainersApp.module.sass";
import Container, {
    Containers,
} from "../../../../components/Container/Container";
import { api } from "../../../../backend/api/backend";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useServerEvent } from "../../../../hooks/useEvent";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import Toolbar from "../../../../components/Toolbar/Toolbar";
import Spacer from "../../../../components/Spacer/Spacer";
import { Button, MaterialIcon } from "@vertex-center/components";
import SelectTags from "../../components/SelectTags/SelectTags";
import { useState } from "react";
import NoItems from "../../../../components/NoItems/NoItems";
import { useContainers } from "../../hooks/useContainers";
import { useNavigate } from "react-router-dom";
import { useTitle } from "../../../../hooks/useTitle";

type ToolbarProps = {
    tags?: string[];
    onTagsChange?: (tags: string[]) => void;
};

const ToolbarContainers = (props: ToolbarProps) => {
    const navigate = useNavigate();

    const { tags, onTagsChange } = props;

    return (
        <Toolbar className={styles.toolbar}>
            <SelectTags selected={tags} onChange={onTagsChange} />
            <Spacer />
            <Button
                variant="colored"
                onClick={() => navigate("/app/vx-containers/add")}
                rightIcon={<MaterialIcon icon="add" />}
            >
                Create container
            </Button>
        </Toolbar>
    );
};

export default function ContainersApp() {
    useTitle("All containers");

    const queryClient = useQueryClient();

    const [tags, setTags] = useState<string[]>([]);
    const {
        data: containers,
        isLoading,
        isError,
        error,
    } = useContainers({ tags });

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
            <div className={styles.containers}>
                <ToolbarContainers tags={tags} onTagsChange={onTagsChange} />

                {error && <APIError error={error} />}

                {!isLoading && !isError && (
                    <Containers>
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
                        {Object.values(containers)?.length === 0 && (
                            <NoItems
                                text="No containers found."
                                icon="deployed_code"
                            />
                        )}
                    </Containers>
                )}
            </div>
        </div>
    );
}
