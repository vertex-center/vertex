import { PropsWithChildren } from "react";
import { Link, NavLink } from "react-router-dom";

import styles from "./Header.module.sass";
import Logo from "../Logo/Logo";
import classNames from "classnames";

type ItemProps = PropsWithChildren<{
    to: string;
}>;

function Item({ children, to }: ItemProps) {
    return (
        <NavLink
            to={to}
            className={({ isActive }) =>
                classNames({
                    [styles.item]: true,
                    [styles.itemActive]: isActive,
                })
            }
        >
            <li>{children}</li>
        </NavLink>
    );
}

export default function Header() {
    return (
        <header className={styles.header}>
            <Link to="/" className={styles.logo}>
                <Logo />
            </Link>
            <nav>
                <ul className={styles.items}>
                    <Item to="marketplace">Marketplace</Item>
                    <Item to="installed">Installed</Item>
                </ul>
            </nav>
        </header>
    );
}
