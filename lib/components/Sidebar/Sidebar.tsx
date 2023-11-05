import { HTMLProps } from "react";
import cx from "classnames";
import { SidebarItem } from "./SidebarItem.tsx";
import { SidebarGroup } from "./SidebarGroup.tsx";
import "./Sidebar.sass";

export type SidebarProps = HTMLProps<HTMLDivElement> & {
    rootUrl: string;
    currentUrl: string;
};

export function Sidebar(props: Readonly<SidebarProps>) {
    const { rootUrl, currentUrl, ...others } = props;

    return (
        <nav
            className={cx("sidebar", {
                "sidebar-with-item-selected":
                    !currentUrl.endsWith(rootUrl) &&
                    !currentUrl.endsWith(rootUrl + "/"),
            })}
            {...others}
        />
    );
}

Sidebar.Item = SidebarItem;
Sidebar.Group = SidebarGroup;
