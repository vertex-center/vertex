import Bay from "../../components/Bay/Bay";
import { useEffect, useState } from "react";
import { getService, InstalledService } from "../../backend/backend";
import { useParams } from "react-router-dom";

import styles from "./BayDetails.module.sass";
import Symbol from "../../components/Symbol/Symbol";
import { Horizontal } from "../../components/Layouts/Layouts";

type MenuItemProps = {
    symbol: string;
    name: string;
};

function MenuItem(props: MenuItemProps) {
    const { symbol, name } = props;

    return (
        <div className={styles.menuItem}>
            <Horizontal alignItems="center" gap={12}>
                <Symbol name={symbol} />
                {name}
            </Horizontal>
        </div>
    );
}

export default function BayDetails() {
    const { uuid } = useParams();

    const [instance, setInstance] = useState<InstalledService>();

    useEffect(() => {
        getService(uuid).then((instance: InstalledService) => {
            setInstance(instance);
        });
    }, [uuid]);

    return (
        <div className={styles.details}>
            <div className={styles.bay}>
                <Bay name={instance?.name} status={instance?.status} />
            </div>
            <div className={styles.menu}>
                <MenuItem symbol="terminal" name="Console" />
                {/*<MenuItem symbol="hub" name="Connections" />*/}
                {/*<MenuItem symbol="settings" name="Settings" />*/}
            </div>
            <div className={styles.content}></div>
        </div>
    );
}
