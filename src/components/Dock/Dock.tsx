import styles from "./Dock.module.sass";
import Icon from "../Icon/Icon";
import { NavLink } from "react-router-dom";
import classNames from "classnames";
import DockDrawer from "./DockDrawer";
import { Fragment, useState } from "react";
import { apps } from "../../models/app";

type DockAppProps = {
    to?: string;
    icon: string;
    name: string;
    onClick?: () => void;
};

export function DockApp(props: Readonly<DockAppProps>) {
    const { to, icon, name, onClick } = props;

    const content = (
        <Fragment>
            <Icon className={styles.icon} name={icon} />
            <span className={styles.name}>{name}</span>
        </Fragment>
    );

    if (!to) {
        return (
            <div className={styles.app} onClick={onClick}>
                {content}
            </div>
        );
    }

    return (
        <NavLink
            to={to}
            onClick={onClick}
            className={({ isActive }) =>
                classNames({
                    [styles.app]: true,
                    [styles.appActive]: isActive,
                })
            }
        >
            {content}
        </NavLink>
    );
}

export default function Dock() {
    const [drawerOpen, setDrawerOpen] = useState(false);

    return (
        <Fragment>
            <DockDrawer
                show={drawerOpen}
                onClose={() => setDrawerOpen(false)}
            />
            <div
                className={classNames({
                    [styles.dockContainer]: true,
                    [styles.dockContainerOpen]: drawerOpen,
                })}
            >
                <div className={styles.dock}>
                    {apps.map((app) => (
                        <DockApp
                            key={app.to}
                            {...app}
                            onClick={() => setDrawerOpen(false)}
                        />
                    ))}
                    <div className={styles.separator} />
                    <DockApp
                        to="/settings"
                        icon="settings"
                        name="Settings"
                        onClick={() => setDrawerOpen(false)}
                    />
                    <DockApp
                        icon="apps"
                        name="Apps"
                        onClick={() => {
                            setDrawerOpen((o) => !o);
                        }}
                    />
                </div>
            </div>
        </Fragment>
    );
}
