import { useParams } from "react-router-dom";
import { api } from "../../../backend/api/backend";
import { Container as ContainerModel } from "../../../models/container";
import Container, { Containers } from "../../../components/Container/Container";
import { v4 as uuidv4 } from "uuid";
import { useServerEvent } from "../../../hooks/useEvent";
import { APIError } from "../../../components/Error/APIError";
import {
    KeyValueGroup,
    KeyValueInfo,
} from "../../../components/KeyValueInfo/KeyValueInfo";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import {
    List,
    ListIcon,
    ListInfo,
    ListItem,
    ListTitle,
    MaterialIcon,
    Title,
} from "@vertex-center/components";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import NoItems from "../../../components/NoItems/NoItems";
import Content from "../../../components/Content/Content";
import useContainer from "../../Containers/hooks/useContainer";
import { API } from "../../Containers/backend/api";

export default function SqlDatabase() {
    const { uuid } = useParams();
    const queryClient = useQueryClient();

    const {
        container,
        isLoading: isLoadingContainer,
        error: errorContainer,
    } = useContainer(uuid);

    const {
        data: db,
        isLoading: isLoadingDatabase,
        error: errorDatabase,
    } = useQuery({
        queryKey: ["sql_containers", uuid],
        queryFn: api.vxSql.container(uuid).get,
    });

    const onPower = async (inst: ContainerModel) => {
        if (!inst) {
            console.error("Container not found");
            return;
        }
        if (inst?.status === "off" || inst?.status === "error") {
            await API.startContainer(inst.id);
            return;
        }
        await API.stopContainer(inst.id);
    };

    const route = uuid ? `/container/${uuid}/events` : "";

    // @ts-ignore
    useServerEvent(window.api_urls.containers, route, {
        status_change: () => {
            queryClient.invalidateQueries({
                queryKey: ["containers", uuid],
            });
        },
    });

    let databases;
    if (db?.databases?.length === 0) {
        databases = <NoItems text="No databases yet." icon="database" />;
    } else {
        databases = (
            <List>
                {db?.databases?.map((db) => (
                    <ListItem key={db.name}>
                        <ListIcon>
                            <MaterialIcon icon="database" />
                        </ListIcon>
                        <ListInfo>
                            <ListTitle>{db.name}</ListTitle>
                        </ListInfo>
                    </ListItem>
                ))}
            </List>
        );
    }

    return (
        <Content>
            <ProgressOverlay show={isLoadingContainer ?? isLoadingDatabase} />
            <APIError error={errorContainer ?? errorDatabase} />

            <Containers>
                <Container
                    container={{
                        value: container ?? {
                            id: uuidv4(),
                        },
                        to: `/app/containers/${container?.id}`,
                        onPower: () => onPower(container),
                    }}
                />
            </Containers>

            <KeyValueGroup>
                <KeyValueInfo name="Username" loading={isLoadingDatabase}>
                    {db?.username}
                </KeyValueInfo>
                <KeyValueInfo name="Password" loading={isLoadingDatabase}>
                    {db?.password}
                </KeyValueInfo>
            </KeyValueGroup>

            <Title variant="h3">Databases</Title>
            {databases}
        </Content>
    );
}
