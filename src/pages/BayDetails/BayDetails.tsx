import Bay from "../../components/Bay/Bay";
import { useCallback, useEffect, useState } from "react";
import {
    getInstance,
    Instance,
    startInstance,
    stopInstance,
} from "../../backend/backend";
import { Link, Outlet, useParams } from "react-router-dom";

import styles from "./BayDetails.module.sass";
import Symbol from "../../components/Symbol/Symbol";
import { Horizontal } from "../../components/Layouts/Layouts";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../backend/sse";

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

    const [instance, setInstance] = useState<Instance>();

    const fetchInstance = useCallback(() => {
        getInstance(uuid).then((instance: Instance) => {
            setInstance(instance);
        });
    }, [uuid]);

    useEffect(() => {
        fetchInstance();
    }, [fetchInstance, uuid]);

    useEffect(() => {
        const sse = registerSSE(
            `http://localhost:6130/instance/${uuid}/events`
        );

        const onStatusChange = (e) => {
            setInstance((instance) => ({ ...instance, status: e.data }));
        };

        registerSSEEvent(sse, "status_change", onStatusChange);

        return () => {
            unregisterSSEEvent(sse, "status_change", onStatusChange);

            unregisterSSE(sse);
        };
    }, [fetchInstance, uuid]);

    const toggleInstance = async (uuid: string) => {
        if (instance.status === "off") {
            await startInstance(uuid);
        } else {
            await stopInstance(uuid);
        }
    };

    return (
        <div className={styles.details}>
            <div className={styles.bay}>
                <Bay
                    name={instance?.name}
                    status={instance?.status}
                    onPower={() => toggleInstance(uuid)}
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
                    <Outlet />
                </div>
            </Horizontal>
        </div>
    );
}
