import styles from "./Infrastructure.module.sass";
import Bay from "../../components/Bay/Bay";
import { useEffect, useState } from "react";
import { getInstalledServices, Service } from "../../backend/backend";

export default function Infrastructure() {
    const [installed, setInstalled] = useState<Service[]>([]);

    useEffect(() => {
        getInstalledServices().then((installed) => setInstalled(installed));
    }, []);

    return (
        <div className={styles.server}>
            <div className={styles.bays}>
                <Bay name="Vertex" status="running" />
                {installed?.map((service) => (
                    <Bay key={service.id} name={service.name} status="error" />
                ))}
            </div>
        </div>
    );
}
