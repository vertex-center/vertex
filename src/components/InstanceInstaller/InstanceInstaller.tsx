import Instance from "../Instance/Instance";
import { APIError } from "../Error/APIError";
import { Fragment, useEffect, useState } from "react";
import { ProgressOverlay } from "../Progress/Progress";
import { Instance as InstanceModel } from "../../models/instance";
import { api } from "../../backend/api/backend";
import { useServerEvent } from "../../hooks/useEvent";
import { useQuery, useQueryClient } from "@tanstack/react-query";

type Props = {
    name: string;
    tag: string;
    install: () => Promise<any>;
};

type Inst = Partial<
    InstanceModel & {
        onPower: (inst: InstanceModel) => Promise<void>;
    }
>;

export default function InstanceInstaller(props: Readonly<Props>) {
    const { name, tag, install } = props;
    const queryClient = useQueryClient();

    const {
        data: instances,
        isLoading: isLoadingInstances,
        error: errorInstances,
    } = useQuery({
        queryKey: ["instances"],
        queryFn: api.vxInstances.instances.all,
    });

    const [downloading, setDownloading] = useState(false);
    const [error, setError] = useState();

    const [instance, setInstance] = useState<Inst>();

    useEffect(() => {
        if (!instances) return;
        const inst = Object.entries(instances).find(
            ([_, instance]) => instance?.tags?.includes(tag) ?? false
        );
        if (!inst) {
            setInstance({
                display_name: name,
                status: "not-installed",
            });
            return;
        }
        setInstance({
            ...inst[1],
            onPower: onPower,
        });
    }, [instances]);

    const onInstall = () => {
        setError(undefined);
        setDownloading(true);

        install()
            .catch(setError)
            .finally(() => {
                setDownloading(false);
                queryClient.invalidateQueries({
                    queryKey: ["instances"],
                });
            });
    };

    const onPower = async (inst: Inst) => {
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

    const route = instance?.uuid
        ? `/app/vx-instances/instance/${instance?.uuid}/events`
        : "";

    useServerEvent(route, {
        status_change: (e) => {
            setInstance((instance) => ({ ...instance, status: e.data }));
        },
    });

    return (
        <Fragment>
            <ProgressOverlay show={downloading || isLoadingInstances} />
            <APIError error={error ?? errorInstances} />
            <Instance
                instance={{
                    value: instance,
                    to: instance?.uuid
                        ? `/app/vx-instances/${instance?.uuid}`
                        : undefined,
                    onInstall: onInstall,
                    onPower: () => onPower(instance),
                }}
            />
        </Fragment>
    );
}
