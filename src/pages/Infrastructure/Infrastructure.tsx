import styles from "./Infrastructure.module.sass";
import Bay from "../../components/Bay/Bay";
import { useEffect, useState } from "react";
import {
    getInstances,
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

export default function Infrastructure() {
    const [status, setStatus] = useState("Checking...");
    const [installed, setInstalled] = useState<Instances>({});

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
        if (installed[uuid].status === "off") {
            await startInstance(uuid);
        } else {
            await stopInstance(uuid);
        }
    };

    return (
        <div className={styles.server}>
            <div className={styles.bays}>
                <Bay name="Vertex" status={status} />
                {Object.keys(installed)?.map((uuid) => (
                    <Bay
                        key={uuid}
                        name={installed[uuid].name}
                        status={installed[uuid].status}
                        to={`/bay/${uuid}`}
                        onPower={() => toggleInstance(uuid)}
                    />
                ))}
                <Link to="/marketplace" className={styles.addBay}>
                    <Symbol name="add" />
                </Link>
            </div>
        </div>
    );
}
