import styles from "./Dock.module.sass";
import { NavLink, useLocation } from "react-router-dom";
import classNames from "classnames";
import DockDrawer from "./DockDrawer";
import React, { Fragment, useState } from "react";
import { useApps } from "../../hooks/useApps";
import { IconContext, SquaresFour } from "@phosphor-icons/react";

type DockAppProps = {
    to?: string;
    icon: React.JSX.Element;
    name: string;
    onClick?: () => void;
};

export function DockApp(props: Readonly<DockAppProps>) {
    const { to, icon, name, onClick } = props;

    const content = (
        <Fragment>
            <IconContext.Provider value={{ size: 24 }}>
                {icon}
            </IconContext.Provider>
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
    const { apps } = useApps();

    const location = useLocation();

    let showDevtools = location.pathname.includes("devtools");

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
                    {apps?.length > 0 && (
                        <div className={styles.shortcuts}>
                            {[...(apps ?? [])]
                                ?.filter(
                                    (app) =>
                                        showDevtools ||
                                        app.category !== "devtools"
                                )
                                ?.filter((app) => !app.hidden)
                                ?.sort((a, b) => (a.name > b.name ? 1 : -1))
                                ?.map((app) => (
                                    <DockApp
                                        key={app.id}
                                        to={`/${app.id}`}
                                        {...app}
                                        onClick={() => setDrawerOpen(false)}
                                    />
                                ))}
                            {apps?.length > 0 && (
                                <div className={styles.separator} />
                            )}
                        </div>
                    )}
                    <DockApp
                        icon={<SquaresFour />}
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
