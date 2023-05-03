import styles from "./Infrastructure.module.sass";
import Bay from "../../components/Bay/Bay";
import { useEffect, useState } from "react";
import {
    getInstances,
    Instance,
    Instances,
    route,
    startInstance,
    stopInstance,
} from "../../backend/backend";
import Symbol from "../../components/Symbol/Symbol";
import { Link } from "react-router-dom";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../backend/sse";
import { HeaderHome } from "../../components/Header/Header";

export default function Infrastructure() {
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
        getInstances()
            .then((installed) => {
                console.log(installed);
                setInstalled(installed);
                setStatus("running");
            })
            .catch(() => {
                setInstalled({});
                setStatus("off");
            });
    };

    useEffect(() => {
        const sse = registerSSE(route("/instances/events"));

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
            await startInstance(uuid);
        } else {
            await stopInstance(uuid);
        }
    };

    return (
        <div className={styles.server}>
            <HeaderHome />
            <div className={styles.bays}>
                <Bay showCables instances={[{ name: "Vertex", status }]} />
                {installedGrouped?.map((instances, i) => (
                    <Bay
                        key={i}
                        showCables
                        instances={instances.map((instance, i) => ({
                            name: instance.name,
                            status: instance.status,
                            count: instances.length > 1 ? i + 1 : undefined,
                            to: `/bay/${instance.uuid}/`,
                            onPower: () => toggleInstance(instance.uuid),
                            use_docker: instance.use_docker,
                        }))}
                    />
                ))}
                <Link to="/marketplace" className={styles.addBay}>
                    <Symbol name="add" />
                </Link>
            </div>
        </div>
    );
}
