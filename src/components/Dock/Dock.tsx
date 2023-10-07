import styles from "./Dock.module.sass";
import Symbol from "../Symbol/Symbol";
import { NavLink } from "react-router-dom";
import classNames from "classnames";

type DockAppProps = {
    to: string;
    symbol: string;
    name: string;
};

function DockApp(props: DockAppProps) {
    const { to, symbol } = props;

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
            <Symbol className={styles.icon} name={symbol} />
            <span className={styles.name}>{props.name}</span>
        </NavLink>
    );
}

export default function Dock() {
    return (
        <div className={styles.dockContainer}>
            <div className={styles.dock}>
                <DockApp to="/instances" symbol="storage" name="Instances" />
                <DockApp
                    to="/monitoring"
                    symbol="monitoring"
                    name="Monitoring"
                />
                <DockApp to="/tunnels" symbol="subway" name="Tunnels" />
                <DockApp
                    to="/reverse-proxy"
                    symbol="router"
                    name="Reverse Proxy"
                />
                <div className={styles.separator} />
                <DockApp to="/settings" symbol="settings" name="Settings" />
            </div>
        </div>
    );
}
