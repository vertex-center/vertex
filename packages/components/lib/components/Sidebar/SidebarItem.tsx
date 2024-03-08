import React, { Fragment, PropsWithChildren, ReactNode, useState } from "react";
import cx from "classnames";
import { NavLink, NavLinkProps } from "../NavLink/NavLink.tsx";
import { CaretDown } from "@phosphor-icons/react";

export type SidebarItemVariant = "default" | "red";

export type SidebarItemProps<T> = PropsWithChildren<{
    variant?: SidebarItemVariant;
    label: string;
    icon?: React.JSX.Element;
    onClick?: () => void;
    notifications?: number;
    trailing?: ReactNode;
    link?: NavLinkProps<T>;
    children?: ReactNode;
}>;

export function SidebarItem<T>(props: Readonly<SidebarItemProps<T>>) {
    const { children, variant, label, icon, trailing, link } = props;

    const hasChildren =
        children !== undefined && children !== null && children !== false;

    const [expanded, setExpanded] = useState(false);

    const content = (
        <Fragment>
            {icon && <div className="sidebar-item-icon">{icon}</div>}
            {label && <div className="sidebar-item-label">{label}</div>}
            {props.notifications !== undefined && (
                <div className="sidebar-item-notification">
                    {props.notifications}
                </div>
            )}
            <div className="sidebar-item-trailing">{trailing}</div>
            {hasChildren && <CaretDown className="sidebar-item-expand" />}
        </Fragment>
    );

    const onClick = () => {
        if (hasChildren) {
            setExpanded(!expanded);
        } else {
            props.onClick?.();
        }
    };

    const className = cx("sidebar-item", {
        "sidebar-item-red": variant === "red",
        "sidebar-item-expanded": expanded,
    });

    let item: React.JSX.Element;
    if (!link) {
        item = (
            <div className={className} onClick={onClick}>
                {content}
            </div>
        );
    } else {
        item = (
            <NavLink {...link} className={className} onClick={onClick}>
                {content}
            </NavLink>
        );
    }

    return (
        <Fragment>
            {item}
            <div
                className={cx("sidebar-item-children", {
                    "sidebar-item-children-expanded": expanded,
                })}
            >
                {children}
            </div>
        </Fragment>
    );
}
