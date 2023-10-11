import { Vertical } from "../../../components/Layouts/Layouts";
import { useParams } from "react-router-dom";
import { api } from "../../../backend/backend";
import { Instance } from "../../../models/instance";
import { useFetch } from "../../../hooks/useFetch";
import Bay from "../../../components/Bay/Bay";
import { v4 as uuidv4 } from "uuid";
import { useServerEvent } from "../../../hooks/useEvent";
import { APIError } from "../../../components/Error/APIError";
import { useEffect } from "react";
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
import Icon from "../../../components/Icon/Icon";
import { Title } from "../../../components/Text/Text";
import styles from "./SqlDatabase.module.sass";

type Props = {};

export default function SqlDatabase(props: Readonly<Props>) {
    const { uuid } = useParams();

    const {
        data: instance,
        loading,
        reload,
        error,
    } = useFetch<Instance>(() => api.vxInstances.instance(uuid).get());

    const {
        data: db,
        loading: dbLoading,
        error: dbError,
    } = useFetch<SqlDBMS>(() => api.vxSql.instance(uuid).get());

    useEffect(() => {
        reload().finally();
    }, [uuid]);

    const onPower = async (inst: Instance) => {
        if (!inst) {
            console.error("Instance not found");
            return;
        }
        if (inst?.status === "off" || inst?.status === "error") {
            await api.vxInstances.instance(inst.uuid).start();
            return;
        }
        await api.vxInstances.instance(inst.uuid).stop();
    };

    const route = uuid ? `/app/vx-instances/instance/${uuid}/events` : "";

    useServerEvent(route, {
        status_change: () => {
            reload().finally();
        },
    });

    return (
        <Vertical gap={30}>
            <ProgressOverlay show={loading ?? dbLoading} />

            <Vertical gap={20}>
                <APIError error={error ?? dbError} />
                <Bay
                    instances={[
                        {
                            value: instance ?? {
                                uuid: uuidv4(),
                            },
                            to: `/app/vx-instances/${instance?.uuid}`,
                            onPower: () => onPower(instance),
                        },
                    ]}
                />

                <KeyValueGroup>
                    <KeyValueInfo name="Username" loading={dbLoading}>
                        {db?.username}
                    </KeyValueInfo>
                    <KeyValueInfo name="Password" loading={dbLoading}>
                        {db?.password}
                    </KeyValueInfo>
                </KeyValueGroup>
            </Vertical>

            <Vertical gap={20}>
                <Title className={styles.title}>Databases</Title>
                <List>
                    {db?.databases?.map((db) => (
                        <ListItem key={db.name}>
                            <ListIcon>
                                <Icon name="database" />
                            </ListIcon>
                            <ListInfo>
                                <ListTitle>{db.name}</ListTitle>
                            </ListInfo>
                        </ListItem>
                    ))}
                </List>
            </Vertical>
        </Vertical>
    );
}
