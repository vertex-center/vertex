import styles from "./Sidebar.module.sass";
import Symbol from "../Symbol/Symbol";
import classNames from "classnames";
import { NavLink, useLocation } from "react-router-dom";
import { Fragment, HTMLProps, PropsWithChildren } from "react";
import { Text } from "../Text/Text";
import Spacer from "../Spacer/Spacer";

function SidebarTitle({ children }: PropsWithChildren) {
    return <Text className={styles.title}>{children}</Text>;
}

type ItemProps = {
    to?: string;
    onClick?: () => void;

    symbol: string | any;
    name: string;
    notifications?: number;

    red?: boolean;
};

export function SidebarItem(props: ItemProps) {
    const { to, name, onClick, red } = props;

    let symbol: any;
    if (typeof props.symbol === "string") {
        symbol = <Symbol name={props.symbol} />;
    } else {
        symbol = props.symbol;
    }

    const content = (
        <Fragment>
            <div className={styles.symbol}>{symbol}</div>
            {name}
            <Spacer />
            {props.notifications && (
                <div className={styles.notifications}>
                    {props.notifications}
                </div>
            )}
        </Fragment>
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

type Props = HTMLProps<HTMLDivElement> & {
    root: string;
};

export default function Sidebar(props: Props) {
    const { children, root } = props;

    const location = useLocation();

    return (
        <nav
            className={classNames({
                [styles.navbar]: true,
                [styles.navbarWithItemSelected]:
                    !location.pathname.endsWith(root) &&
                    !location.pathname.endsWith(root + "/"),
            })}
        >
            {children}
        </nav>
    );
}
