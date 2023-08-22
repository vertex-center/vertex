import styles from "./Dock.module.sass";
import Symbol from "../Symbol/Symbol";
import { NavLink } from "react-router-dom";
import classNames from "classnames";

type DockAppProps = {
    to: string;
    symbol: string;
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
            <Symbol name={symbol} />
        </NavLink>
    );
}

type Props = {};

export default function Dock(props: Props) {
    return (
        <div className={styles.dockContainer}>
            <div className={styles.dock}>
                <DockApp to="/infrastructure" symbol="storage" />
                <DockApp to="/settings" symbol="settings" />
            </div>
        </div>
    );
}
