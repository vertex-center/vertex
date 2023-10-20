import { Caption, Title } from "../../../../components/Text/Text";
import { useParams } from "react-router-dom";
import useContainer from "../../hooks/useContainer";
import styles from "./ContainerUpdate.module.sass";
import { Vertical } from "../../../../components/Layouts/Layouts";
import { api } from "../../../../backend/api/backend";
import { useState } from "react";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useQueryClient } from "@tanstack/react-query";
import List from "../../../../components/List/List";
import ListItem from "../../../../components/List/ListItem";
import ListInfo from "../../../../components/List/ListInfo";
import ListTitle from "../../../../components/List/ListTitle";
import ListActions from "../../../../components/List/ListActions";
import { Button, MaterialIcon } from "@vertex-center/components";

export default function ContainerUpdate() {
    const { uuid } = useParams();
    const queryClient = useQueryClient();

    const { container, isLoading } = useContainer(uuid);

    const [error, setError] = useState();

    const updateVertexIntegration = () => {
        return api.vxContainers
            .container(uuid)
            .update.service()
            .then(() => {
                queryClient.invalidateQueries({
                    queryKey: ["containers", uuid],
                });
            })
            .catch(setError);
    };

    let content;
    if (container?.service_update?.available) {
        content = (
            <List>
                <ListItem>
                    <ListInfo onClick={updateVertexIntegration}>
                        <ListTitle>Vertex integration</ListTitle>
                    </ListInfo>
                    <ListActions>
                        <Button rightIcon={<MaterialIcon icon="download" />}>
                            Update
                        </Button>
                    </ListActions>
                </ListItem>
            </List>
        );
    } else {
        content = (
            <Caption className={styles.content}>
                There are no updates available.
            </Caption>
        );
    }

    return (
        <Vertical gap={20}>
            <ProgressOverlay show={isLoading} />
            <Title className={styles.title}>Update</Title>
            <APIError error={error} />
            {!error && !isLoading && content}
        </Vertical>
    );
}
