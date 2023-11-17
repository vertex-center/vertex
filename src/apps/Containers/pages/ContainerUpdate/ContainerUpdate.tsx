import { Caption } from "../../../../components/Text/Text";
import { useParams } from "react-router-dom";
import useContainer from "../../hooks/useContainer";
import { api } from "../../../../backend/api/backend";
import { useState } from "react";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useQueryClient } from "@tanstack/react-query";
import {
    Button,
    List,
    ListActions,
    ListInfo,
    ListItem,
    ListTitle,
    MaterialIcon,
    Title,
} from "@vertex-center/components";
import Content from "../../../../components/Content/Content";

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
                    <ListInfo>
                        <ListTitle>Vertex integration</ListTitle>
                    </ListInfo>
                    <ListActions>
                        <Button
                            onClick={updateVertexIntegration}
                            rightIcon={<MaterialIcon icon="download" />}
                        >
                            Update
                        </Button>
                    </ListActions>
                </ListItem>
            </List>
        );
    } else {
        content = <Caption>There are no updates available.</Caption>;
    }

    return (
        <Content>
            <Title variant="h2">Update</Title>
            <ProgressOverlay show={isLoading} />
            <APIError error={error} />
            {!error && !isLoading && content}
        </Content>
    );
}
