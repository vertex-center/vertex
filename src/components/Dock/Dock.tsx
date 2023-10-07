import styles from "./Dock.module.sass";
import Icon from "../Icon/Icon";
import { NavLink } from "react-router-dom";
import classNames from "classnames";

type DockAppProps = {
    to: string;
    icon: string;
    name: string;
};

function DockApp(props: Readonly<DockAppProps>) {
    const { to, icon } = props;

    return (
        <NavLink
            to={to}
            className={({ isActive }) =>
                classNames({
                    [styles.app]: true,
                    [styles.appActive]: isActive,
                })
            }
        >
            <Icon className={styles.icon} name={icon} />
            <span className={styles.name}>{props.name}</span>
        </NavLink>
    );
}

export default function Dock() {
    return (
        <div className={styles.dockContainer}>
            <div className={styles.dock}>
                <DockApp to="/instances" icon="storage" name="Instances" />
                <DockApp to="/monitoring" icon="monitoring" name="Monitoring" />
                <DockApp to="/tunnels" icon="subway" name="Tunnels" />
                <DockApp
                    to="/reverse-proxy"
                    icon="router"
                    name="Reverse Proxy"
                />
                <div className={styles.separator} />
                <DockApp to="/settings" icon="settings" name="Settings" />
            </div>
        </div>
    );
}
