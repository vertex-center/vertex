import styles from "./ContainersApp.module.sass";
import Container, {
    Containers,
} from "../../../../components/Container/Container";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useServerEvent } from "../../../../hooks/useEvent";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import Toolbar from "../../../../components/Toolbar/Toolbar";
import Spacer from "../../../../components/Spacer/Spacer";
import {
    Button,
    MaterialIcon,
    useTitle,
    Vertical,
} from "@vertex-center/components";
import SelectTags from "../../components/SelectTags/SelectTags";
import { useState } from "react";
import NoItems from "../../../../components/NoItems/NoItems";
import { useContainers } from "../../hooks/useContainers";
import { useNavigate } from "react-router-dom";
import { API } from "../../backend/api";

type ToolbarProps = {
    tags?: string[];
    onTagsChange?: (tags: string[]) => void;
};

const ToolbarContainers = (props: ToolbarProps) => {
    const navigate = useNavigate();

    const { tags, onTagsChange } = props;

    return (
        <Toolbar>
            <SelectTags selected={tags} onChange={onTagsChange} />
            <Spacer />
            <Button
                variant="colored"
                onClick={() => navigate("/app/containers/add")}
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
        mutationFn: async (id: string) => {
            const container = containers?.find((c) => c.id === id);
            if (container.status === "off" || container.status === "error") {
                await API.startContainer(id);
            } else {
                await API.stopContainer(id);
            }
        },
    });

    // @ts-ignore
    useServerEvent(window.api_urls.containers, "/containers/events", {
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
            <Vertical gap={12} className={styles.containers}>
                <ToolbarContainers tags={tags} onTagsChange={onTagsChange} />

                {error && <APIError error={error} />}

                {!isLoading && !isError && (
                    <Containers>
                        {containers?.map((c) => (
                            <Container
                                key={c.id}
                                container={{
                                    value: c,
                                    to: `/app/containers/${c.id}/`,
                                    onPower: async () =>
                                        mutationPower.mutate(c.id),
                                }}
                            />
                        ))}
                        {containers?.length === 0 && (
                            <NoItems
                                text="No containers found."
                                icon="deployed_code"
                            />
                        )}
                    </Containers>
                )}
            </Vertical>
        </div>
    );
}
