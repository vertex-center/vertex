import Bay from "../../components/Bay/Bay";
import { useCallback, useEffect, useState } from "react";
import {
    deleteInstance,
    getInstance,
    Instance,
    route,
    startInstance,
    stopInstance,
} from "../../backend/backend";
import { Outlet, useNavigate, useParams } from "react-router-dom";

import styles from "./BayDetails.module.sass";
import { Horizontal } from "../../components/Layouts/Layouts";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../backend/sse";
import Spacer from "../../components/Spacer/Spacer";
import Header from "../../components/Header/Header";
import Sidebar, { SidebarItem } from "../../components/Sidebar/Sidebar";

export const bayNavItems = [
    {
        label: "Logs",
        to: "/logs",
        symbol: "terminal",
    },
    {
        label: "Environment",
        to: "/environment",
        symbol: "tune",
    },
    {
        label: "Dependencies",
        to: "/dependencies",
        symbol: "widgets",
    },
];

export default function BayDetails() {
    const { uuid } = useParams();
    const navigate = useNavigate();

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
        if (uuid === undefined) return;

        const sse = registerSSE(route(`/instance/${uuid}/events`));

        const onStatusChange = (e) => {
            setInstance((instance) => ({ ...instance, status: e.data }));
        };

        registerSSEEvent(sse, "status_change", onStatusChange);

        return () => {
            unregisterSSEEvent(sse, "status_change", onStatusChange);

            unregisterSSE(sse);
        };
    }, [uuid]);

    const toggleInstance = async (uuid: string) => {
        if (instance.status === "off") {
            await startInstance(uuid);
        } else {
            await stopInstance(uuid);
        }
    };

    const onDeleteInstance = () => {
        deleteInstance(uuid).then(() => {
            navigate("/");
        });
    };

    return (
        <div className={styles.details}>
            <Header />
            <div className={styles.bay}>
                <Bay
                    name={instance?.name}
                    status={instance?.status}
                    onPower={() => toggleInstance(uuid)}
                />
            </div>
            <Horizontal className={styles.content}>
                <Sidebar>
                    <SidebarItem to="/" symbol="arrow_back" name="Back" />
                    <SidebarItem
                        to={`/bay/${uuid}/`}
                        symbol="home"
                        name="Home"
                    />
                    <div className={styles.separator} />
                    {bayNavItems.map((item) => (
                        <SidebarItem
                            to={`/bay/${uuid}${item.to}`}
                            symbol={item.symbol}
                            name={item.label}
                        />
                    ))}
                    <Spacer />
                    <SidebarItem
                        onClick={onDeleteInstance}
                        symbol="delete"
                        name="Delete"
                        red
                    />
                </Sidebar>
                <div className={styles.side}>
                    <Outlet />
                </div>
            </Horizontal>
        </div>
    );
}
