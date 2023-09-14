import styles from "./Sidebar.module.sass";
import { Horizontal } from "../Layouts/Layouts";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";
import { NavLink } from "react-router-dom";
import { HTMLProps, PropsWithChildren } from "react";
import { Text } from "../Text/Text";

function SidebarTitle({ children }: PropsWithChildren) {
    return <Text className={styles.title}>{children}</Text>;
}

type ItemProps = {
    to?: string;
    onClick?: () => void;

    symbol: string | any;
    name: string;

    red?: boolean;
};

export function SidebarItem(props: ItemProps) {
    const { to, symbol, name, onClick, red } = props;

    const content = (
        <Horizontal alignItems="center" gap={12}>
            <Symbol className={styles.symbol} name={symbol} />
            {name}
        </Horizontal>
    );

    const className = classNames({
        [styles.navbarItem]: true,
        [styles.navbarItemRed]: red,
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
                    [styles.navbarItemActive]: isActive,
                })
            }
        >
            {content}
        </NavLink>
    );
}

type GroupProps = PropsWithChildren<{
    title?: string;
}>;

export function SidebarGroup(props: GroupProps) {
    const { title, children } = props;

    return (
        <div className={styles.group}>
            {title && <SidebarTitle>{title}</SidebarTitle>}
            {children}
        </div>
    );
}

type Props = HTMLProps<HTMLDivElement>;

export default function Sidebar({ children }: Props) {
    return <nav className={styles.navbar}>{children}</nav>;
}
