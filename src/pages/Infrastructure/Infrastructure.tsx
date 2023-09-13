import styles from "./Infrastructure.module.sass";
import Bay from "../../components/Bay/Bay";
import { useEffect, useState } from "react";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../backend/sse";
import Progress from "../../components/Progress";
import { Instance, Instances } from "../../models/instance";
import { BigTitle } from "../../components/Text/Text";
import Button from "../../components/Button/Button";
import { Horizontal } from "../../components/Layouts/Layouts";
import { api } from "../../backend/backend";

export default function Infrastructure() {
    const [loading, setLoading] = useState(true);
    const [status, setStatus] = useState("Waiting server response...");
    const [installed, setInstalled] = useState<Instances>({});

    const [installedGrouped, setInstalledGrouped] = useState<Instance[][]>([]);

    useEffect(() => {
        const ids = new Set<string>(Object.values(installed).map((i) => i.id));
        const final = [];
        ids.forEach((id) => {
            final.push(
                Object.entries(installed)
                    .filter(([_, i]) => i.id === id)
                    .map(([_, i]) => i)
            );
        });
        setInstalledGrouped(final);
    }, [installed]);

    const fetchServices = () => {
        api.instances
            .get()
            .then((res) => {
                console.log(res.data);
                setInstalled(res.data);
                setStatus("running");
            })
            .catch(() => {
                setInstalled({});
                setStatus("off");
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
            {loading && <Progress infinite small />}
            <div className={styles.title}>
                <BigTitle>Infrastructure</BigTitle>
            </div>
            {!loading && (
                <div className={styles.bays}>
                    <Bay instances={[{ name: "Vertex", status }]} />
                    {installedGrouped?.map((instances, i) => (
                        <Bay
                            key={i}
                            instances={instances.map((instance, i) => ({
                                name: instance?.display_name ?? instance.name,
                                status: instance.status,
                                count: instances.length > 1 ? i + 1 : undefined,
                                to: `/infrastructure/${instance.uuid}/`,
                                onPower: () => toggleInstance(instance.uuid),
                                method: instance.install_method,
                                update: instance.update,
                            }))}
                        />
                    ))}
                    <Horizontal className={styles.addBay} gap={10}>
                        <Button primary to="/marketplace" leftSymbol="add">
                            Add services
                        </Button>
                        {/*<Button onClick={checkForUpdates} leftSymbol="update">*/}
                        {/*    Check for updates*/}
                        {/*</Button>*/}
                    </Horizontal>
                </div>
            )}
        </div>
    );
}
