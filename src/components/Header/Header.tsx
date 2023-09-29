import { HTMLProps, PropsWithChildren } from "react";
import { Link, NavLink } from "react-router-dom";

import styles from "./Header.module.sass";
import Logo from "../Logo/Logo";
import classNames from "classnames";
import Symbol from "../Symbol/Symbol";

type ItemProps = PropsWithChildren<{
    to: string;
    symbol: string;
}>;

function Item({ children, to, symbol }: ItemProps) {
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
            <li className={styles.itemContent}>
                <Symbol className={styles.itemSymbol} name={symbol} />
                <div>{children}</div>
            </li>
        </NavLink>
    );
}

export default function Header({ children }: HTMLProps<HTMLHeadingElement>) {
    return (
        <header className={styles.header}>
            <Link to="/instances" className={styles.logo}>
                <Logo />
            </Link>
            {children}
        </header>
    );
}
