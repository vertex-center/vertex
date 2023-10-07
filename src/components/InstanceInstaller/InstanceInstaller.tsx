import Bay from "../Bay/Bay";
import { APIError } from "../Error/APIError";
import { Fragment, useEffect, useState } from "react";
import { ProgressOverlay } from "../Progress/Progress";
import { useFetch } from "../../hooks/useFetch";
import { Instance, Instances } from "../../models/instance";
import { api } from "../../backend/backend";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../backend/sse";

type Props = {
    name: string;
    tag: string;
    install: () => Promise<any>;
};

type Inst = Partial<
    Instance & {
        onPower: (inst: Instance) => Promise<void>;
    }
>;

export default function InstanceInstaller(props: Readonly<Props>) {
    const { name, tag, install } = props;

    const {
        data: instances,
        reload: reloadInstances,
        loading: loadingInstances,
        error: errorInstances,
    } = useFetch<Instances>(api.instances.get);

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
                reloadInstances().catch(setError);
            });
    };

    const onPower = async (inst: Inst) => {
        if (inst?.status === "off" || inst?.status === "error") {
            await api.instance.start(inst.uuid);
            return;
        }
        await api.instance.stop(inst.uuid);
    };

    useEffect(() => {
        if (instance?.uuid === undefined) return;

        const sse = registerSSE(`/instance/${instance.uuid}/events`);

        const onStatusChange = (e: any) => {
            setInstance((inst) => ({ ...inst, status: e.data }));
        };

        registerSSEEvent(sse, "status_change", onStatusChange);

        return () => {
            unregisterSSEEvent(sse, "status_change", onStatusChange);

            unregisterSSE(sse);
        };
    }, [instance]);

    return (
        <Fragment>
            <ProgressOverlay show={downloading || loadingInstances} />
            <APIError error={error ?? errorInstances} />
            <Bay
                instances={[
                    {
                        value: instance,
                        to: instance?.uuid
                            ? `/instances/${instance?.uuid}`
                            : undefined,
                        onInstall: onInstall,
                        onPower: async () => {
                            await onPower(instance);
                        },
                    },
                ]}
            />
        </Fragment>
    );
}
