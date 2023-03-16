import Bay from "../../components/Bay/Bay";
import { useCallback, useEffect, useState } from "react";
import {
    getService,
    InstalledService,
    startService,
    stopService,
} from "../../backend/backend";
import { Link, useParams } from "react-router-dom";

import styles from "./BayDetails.module.sass";
import Symbol from "../../components/Symbol/Symbol";
import { Horizontal } from "../../components/Layouts/Layouts";
import SSE from "../../backend/sse";
import Logs from "../../components/Logs/Logs";

type MenuItemProps = {
    to: string;
    symbol: string;
    name: string;
};

function MenuItem(props: MenuItemProps) {
    const { to, symbol, name } = props;

    return (
        <Link to={to}>
            <div className={styles.menuItem}>
                <Horizontal alignItems="center" gap={12}>
                    <Symbol name={symbol} />
                    {name}
                </Horizontal>
            </div>
        </Link>
    );
}

export default function BayDetails() {
    const { uuid } = useParams();

    const [instance, setInstance] = useState<InstalledService>();

    const [logs, setLogs] = useState<any[]>([]);

    const fetchInstance = useCallback(() => {
        getService(uuid).then((instance: InstalledService) => {
            setInstance(instance);
        });
    }, [uuid]);

    useEffect(() => {
        fetchInstance();
    }, [fetchInstance, uuid]);

    useEffect(() => {
        const sse = new SSE(`http://localhost:6130/service/${uuid}/events`);

        sse.on("stdout", (e) => {
            setLogs((logs) => [
                ...logs,
                {
                    type: "message",
                    message: e.data,
                },
            ]);
        });

        sse.on("stderr", (e) => {
            setLogs((logs) => [
                ...logs,
                {
                    type: "error",
                    message: e.data,
                },
            ]);
        });

        sse.on("status_change", (e) => {
            setInstance((instance) => ({ ...instance, status: e.data }));
        });

        return () => sse.close();
    }, [fetchInstance, uuid]);

    const toggleService = async (uuid: string) => {
        if (instance.status === "off") {
            await startService(uuid);
        } else {
            await stopService(uuid);
        }
    };

    return (
        <div className={styles.details}>
            <div className={styles.bay}>
                <Bay
                    name={instance?.name}
                    status={instance?.status}
                    onPower={() => toggleService(uuid)}
                />
            </div>
            <Horizontal className={styles.content}>
                <div className={styles.menu}>
                    <MenuItem to="/" symbol="arrow_back" name="Back" />
                    <MenuItem
                        to={`/bay/${uuid}/logs`}
                        symbol="terminal"
                        name="Console"
                    />
                    {/*<MenuItem symbol="hub" name="Connections" />*/}
                    {/*<MenuItem symbol="settings" name="Settings" />*/}
                </div>
                <div className={styles.side}>
                    <Logs lines={logs} />
                </div>
            </Horizontal>
        </div>
    );
}
