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
import { NavLink, Outlet, useNavigate, useParams } from "react-router-dom";

import styles from "./BayDetails.module.sass";
import Symbol from "../../components/Symbol/Symbol";
import { Horizontal } from "../../components/Layouts/Layouts";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../backend/sse";
import Spacer from "../../components/Spacer/Spacer";
import classNames from "classnames";

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

type MenuItemProps = {
    to?: string;
    onClick?: () => void;

    symbol: string;
    name: string;

    red?: boolean;
};

function MenuItem(props: MenuItemProps) {
    const { to, symbol, name, onClick, red } = props;

    const content = (
        <Horizontal alignItems="center" gap={12}>
            <Symbol name={symbol} />
            {name}
        </Horizontal>
    );

    const className = classNames({
        [styles.menuItem]: true,
        [styles.menuItemRed]: red,
    });

    if (!to)
        return (
            <div className={className} onClick={onClick}>
                {content}
            </div>
        );

    return (
        <NavLink
            to={to}
            className={({ isActive }) =>
                classNames({
                    [className]: true,
                    [styles.menuItemActive]: isActive,
                })
            }
        >
            {content}
        </NavLink>
    );
}

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
                    <MenuItem to={`/bay/${uuid}/`} symbol="home" name="Home" />
                    <div className={styles.separator} />
                    {bayNavItems.map((item) => (
                        <MenuItem
                            to={`/bay/${uuid}${item.to}`}
                            symbol={item.symbol}
                            name={item.label}
                        />
                    ))}
                    <Spacer />
                    <MenuItem
                        onClick={onDeleteInstance}
                        symbol="delete"
                        name="Delete"
                        red
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
