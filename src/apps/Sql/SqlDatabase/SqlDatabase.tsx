import { Vertical } from "../../../components/Layouts/Layouts";
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
import List from "../../../components/List/List";
import ListItem from "../../../components/List/ListItem";
import ListIcon from "../../../components/List/ListIcon";
import ListInfo from "../../../components/List/ListInfo";
import ListTitle from "../../../components/List/ListTitle";
import { Title } from "../../../components/Text/Text";
import styles from "./SqlDatabase.module.sass";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import NoItems from "../../../components/NoItems/NoItems";
import { MaterialIcon } from "@vertex-center/components";

export default function SqlDatabase() {
    const { uuid } = useParams();
    const queryClient = useQueryClient();

    const {
        data: container,
        isLoading: isLoadingContainer,
        error: errorContainer,
    } = useQuery({
        queryKey: ["containers", uuid],
        queryFn: api.vxContainers.container(uuid).get,
    });

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
            await api.vxContainers.container(inst.uuid).start();
            return;
        }
        await api.vxContainers.container(inst.uuid).stop();
    };

    const route = uuid ? `/app/vx-containers/container/${uuid}/events` : "";

    useServerEvent(route, {
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
        <Vertical gap={30}>
            <ProgressOverlay show={isLoadingContainer ?? isLoadingDatabase} />

            <Vertical gap={20}>
                <APIError error={errorContainer ?? errorDatabase} />
                <Containers>
                    <Container
                        container={{
                            value: container ?? {
                                uuid: uuidv4(),
                            },
                            to: `/app/vx-containers/${container?.uuid}`,
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
            </Vertical>

            <Vertical gap={20}>
                <Title className={styles.title}>Databases</Title>
                {databases}
            </Vertical>
        </Vertical>
    );
}
