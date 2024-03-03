import { Caption } from "../../../../components/Text/Text";
import { useParams } from "react-router-dom";
import useContainer from "../../hooks/useContainer";
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
    Title,
} from "@vertex-center/components";
import Content from "../../../../components/Content/Content";
import { API } from "../../backend/api";
import { DownloadSimple } from "@phosphor-icons/react";

export default function ContainerUpdate() {
    const { uuid } = useParams();
    const queryClient = useQueryClient();

    const { container, isLoading } = useContainer(uuid);

    const [error, setError] = useState();

    const updateVertexIntegration = () => {
        return API.updateService(uuid)
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
                            rightIcon={<DownloadSimple />}
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
