import React, { Fragment, ReactNode } from "react";
import cx from "classnames";
import { NavLink, NavLinkProps } from "../NavLink/NavLink.tsx";

export type SidebarItemVariant = "default" | "red";

export type SidebarItemProps<T> = {
    variant?: SidebarItemVariant;
    label: string;
    icon?: React.JSX.Element;
    onClick?: () => void;
    notifications?: number;
    trailing?: ReactNode;
    link?: NavLinkProps<T>;
};

export function SidebarItem<T>(props: Readonly<SidebarItemProps<T>>) {
    const { variant, label, icon, onClick, trailing, link } = props;

    const content = (
        <Fragment>
            {icon && <div className="sidebar-item-icon">{icon}</div>}
            {label}
            {props.notifications !== undefined && (
                <div className="sidebar-item-notification">
                    {props.notifications}
                </div>
            )}
            <div className="sidebar-item-trailing">{trailing}</div>
        </Fragment>
    );

    const className = cx("sidebar-item", {
        "sidebar-item-red": variant === "red",
    });

    if (!link)
        return (
            <div className={className} onClick={onClick}>
                {content}
            </div>
        );

    return (
        <NavLink {...link} className={className}>
            {content}
        </NavLink>
    );
}
