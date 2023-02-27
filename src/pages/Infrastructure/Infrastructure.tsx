import styles from "./Infrastructure.module.sass";
import Bay from "../../components/Bay/Bay";
import { useEffect, useState } from "react";
import {
    getInstalledServices,
    InstalledService,
    InstalledServices,
    startService,
    stopService,
} from "../../backend/backend";
import Symbol from "../../components/Symbol/Symbol";
import { Link } from "react-router-dom";

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
        const interval = setInterval(fetchServices, 1000);
        return () => clearInterval(interval);
    }, []);

    const toggleService = async (service: InstalledService) => {
        if (installed[service.id].status === "off") {
            await startService(service);
        } else {
            await stopService(service);
        }
    };

    return (
        <div className={styles.server}>
            <div className={styles.bays}>
                <Bay name="Vertex" status={status} />
                {Object.keys(installed)?.map((key) => (
                    <Bay
                        key={installed[key].id}
                        name={installed[key].name}
                        status={installed[key].status}
                        onPower={() => toggleService(installed[key])}
                    />
                ))}
                <Link to="/marketplace" className={styles.addBay}>
                    <Symbol name="add" />
                </Link>
            </div>
        </div>
    );
}
