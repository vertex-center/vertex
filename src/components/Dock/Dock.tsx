import styles from "./Dock.module.sass";
import Symbol from "../Symbol/Symbol";
import { NavLink } from "react-router-dom";
import classNames from "classnames";

type DockAppProps = {
    to: string;
    symbol: string;
    color?: string;
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
        </NavLink>
    );
}

export default function Dock() {
    return (
        <div className={styles.dockContainer}>
            <div className={styles.dock}>
                <DockApp to="/instances" symbol="storage" />
                <DockApp to="/proxy" symbol="router" />
                <div className={styles.separator} />
                <DockApp to="/settings" symbol="settings" />
            </div>
        </div>
    );
}
