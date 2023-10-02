import styles from "./InstancesApp.module.sass";
import Bay from "../../../components/Bay/Bay";
import { useEffect, useState } from "react";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../../backend/sse";
import { Instance, Instances } from "../../../models/instance";
import { BigTitle } from "../../../components/Text/Text";
import Button from "../../../components/Button/Button";
import { Horizontal } from "../../../components/Layouts/Layouts";
import { api } from "../../../backend/backend";
import { APIError } from "../../../components/Error/Error";
import { ProgressOverlay } from "../../../components/Progress/Progress";

export default function InstancesApp() {
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState();
    const [status, setStatus] = useState("Waiting server response...");
    const [installed, setInstalled] = useState<Instances>({});

    const [installedGrouped, setInstalledGrouped] = useState<Instance[][]>([]);

    useEffect(() => {
        const ids = new Set<string>(
            Object.values(installed).map((i) => i.service?.id)
        );
        const final = [];
        ids.forEach((id) => {
            final.push(
                Object.entries(installed)
                    .filter(([_, i]) => i.service?.id === id)
                    .map(([_, i]) => i)
            );
        });
        setInstalledGrouped(final);
    }, [installed]);

    const fetchServices = () => {
        // setError(undefined);

        api.instances
            .get()
            .then((res) => {
                console.log(res.data);
                setInstalled(res.data);
                setStatus("running");
            })
            .catch((error) => {
                setInstalled({});
                setStatus("off");
                setError(error);
            })
            .finally(() => {
                setLoading(false);
            });
    };

    useEffect(() => {
        const sse = registerSSE("/instances/events");

        const onOpen = (e) => {
            console.log(e);
            fetchServices();
        };

        const onChange = (e) => {
            console.log(e);
            fetchServices();
        };

        registerSSEEvent(sse, "open", onOpen);
        registerSSEEvent(sse, "change", onChange);

        return () => {
            unregisterSSEEvent(sse, "open", onOpen);
            unregisterSSEEvent(sse, "change", onChange);

            unregisterSSE(sse);
        };
    }, []);

    const toggleInstance = async (uuid: string) => {
        if (
            installed[uuid].status === "off" ||
            installed[uuid].status === "error"
        ) {
            await api.instance.start(uuid);
        } else {
            await api.instance.stop(uuid);
        }
    };

    const checkForUpdates = async () => {
        api.instances.checkForUpdates().then((res) => {
            setInstalled(res.data);
        });
    };

    return (
        <div className={styles.server}>
            {loading && <ProgressOverlay />}
            <div className={styles.title}>
                <BigTitle>Instances</BigTitle>
            </div>

            {error && (
                <div className={styles.bays}>
                    <APIError error={error} />
                </div>
            )}

            {!loading && !error && (
                <div className={styles.bays}>
                    <Bay instances={[{ name: "Vertex", status }]} />
                    {installedGrouped?.map((instances, i) => (
                        <Bay
                            key={i}
                            instances={instances.map((instance, i) => ({
                                name:
                                    instance?.display_name ??
                                    instance?.service?.name,
                                status: instance.status,
                                count: instances.length > 1 ? i + 1 : undefined,
                                to: `/instances/${instance.uuid}/`,
                                onPower: () => toggleInstance(instance.uuid),
                                method: instance.install_method,
                                update: instance.update,
                            }))}
                        />
                    ))}
                    <Horizontal className={styles.addBay} gap={10}>
                        <Button primary to="/instances/add" leftSymbol="add">
                            Create instance
                        </Button>
                    </Horizontal>
                </div>
            )}
        </div>
    );
}
