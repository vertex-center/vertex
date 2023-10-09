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

type Props = {};

export default function SqlDatabase(props: Readonly<Props>) {
    const { uuid } = useParams();

    const {
        data: instance,
        loading,
        reload,
        error,
    } = useFetch<Instance>(() => api.instance.get(uuid));

    const {
        data: db,
        loading: dbLoading,
        error: dbError,
    } = useFetch<SQLDatabase>(() => api.sql.uuid(uuid).get());

    useEffect(() => {
        reload().finally();
    }, [uuid]);

    const onPower = async (inst: Instance) => {
        if (!inst) {
            console.error("Instance not found");
            return;
        }
        if (inst?.status === "off" || inst?.status === "error") {
            await api.instance.start(inst.uuid);
            return;
        }
        await api.instance.stop(inst.uuid);
    };

    const route = uuid ? `/instance/${uuid}/events` : "";

    useServerEvent(route, {
        status_change: () => {
            reload().finally();
        },
    });

    return (
        <Vertical gap={20}>
            <ProgressOverlay show={loading ?? dbLoading} />
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
    );
}
