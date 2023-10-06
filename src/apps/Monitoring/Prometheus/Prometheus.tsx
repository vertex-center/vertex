import { Text, Title } from "../../../components/Text/Text";
import styles from "./Prometheus.module.sass";
import { Vertical } from "../../../components/Layouts/Layouts";
import { useFetch } from "../../../hooks/useFetch";
import { api } from "../../../backend/backend";
import Bay from "../../../components/Bay/Bay";
import { Fragment, useEffect, useState } from "react";
import { Instances } from "../../../models/instance";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { APIError } from "../../../components/Error/Error";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../../backend/sse";

export default function Prometheus() {
    const {
        data: instances,
        reload: reloadInstances,
        loading: loadingInstances,
        error: errorInstances,
    } = useFetch<Instances>(api.instances.get);

    const [downloading, setDownloading] = useState(false);
    const [error, setError] = useState();

    const [instance, setInstance] = useState<any>();

    useEffect(() => {
        if (!instances) return;
        const inst = Object.entries(instances).find(
            ([_, instance]) => instance.service.id === "prometheus"
        );
        if (!inst) {
            setInstance({
                display_name: "Prometheus",
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

        api.metrics
            .collector("prometheus")
            .install()
            .catch(setError)
            .finally(() => {
                setDownloading(false);
                reloadInstances().catch(setError);
            });
    };

    useEffect(() => {
        if (instance?.uuid === undefined) return;

        const sse = registerSSE(`/instance/${instance.uuid}/events`);

        const onStatusChange = (e: any) => {
            setInstance((instance) => ({ ...instance, status: e.data }));
        };

        registerSSEEvent(sse, "status_change", onStatusChange);

        return () => {
            unregisterSSEEvent(sse, "status_change", onStatusChange);

            unregisterSSE(sse);
        };
    }, [instance]);

    const onPower = async (instance: any) => {
        if (instance?.status === "off" || instance?.status === "error") {
            await api.instance.start(instance.uuid);
        }
        {
            await api.instance.stop(instance.uuid);
        }
    };

    return (
        <Vertical gap={30}>
            <ProgressOverlay show={loadingInstances || downloading} />

            <Vertical gap={20}>
                <Title className={styles.title}>Prometheus</Title>
                <Text className={styles.content}>
                    Prometheus allows you to collect metrics gathered by Vertex.
                    {instance?.status === "not-installed" && (
                        <Fragment>
                            {" "}
                            To enable collection, you first need to install
                            Prometheus with Vertex Instances.
                        </Fragment>
                    )}
                </Text>
                <APIError error={error ?? errorInstances} />
                <Bay
                    instances={[
                        {
                            name:
                                instance?.display_name ??
                                instance?.service?.name,
                            status: instance?.status,

                            to: instance?.uuid
                                ? `/instances/${instance?.uuid}`
                                : undefined,

                            onInstall: onInstall,
                            onPower: () => onPower(instance),
                        },
                    ]}
                />
            </Vertical>
        </Vertical>
    );
}
