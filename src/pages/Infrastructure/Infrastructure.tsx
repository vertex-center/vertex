import styles from "./Infrastructure.module.sass";
import Bay from "../../components/Bay/Bay";
import { useEffect, useState } from "react";
import {
    getInstalledServices,
    InstalledServices,
    startService,
    stopService,
} from "../../backend/backend";
import Symbol from "../../components/Symbol/Symbol";
import { Link } from "react-router-dom";
import SSE from "../../backend/sse";

export default function Infrastructure() {
    const [status, setStatus] = useState("Checking...");
    const [installed, setInstalled] = useState<InstalledServices>({});

    const fetchServices = () => {
        getInstalledServices()
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
        const sse = new SSE("http://localhost:6130/services/events");

        sse.on("open", (e) => {
            console.log(e);
            fetchServices();
        });

        sse.on("change", (e) => {
            console.log(e);
            fetchServices();
        });

        return () => sse.close();
    }, []);

    const toggleService = async (uuid: string) => {
        if (installed[uuid].status === "off") {
            await startService(uuid);
        } else {
            await stopService(uuid);
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
                        onPower={() => toggleService(uuid)}
                    />
                ))}
                <Link to="/marketplace" className={styles.addBay}>
                    <Symbol name="add" />
                </Link>
            </div>
        </div>
    );
}
