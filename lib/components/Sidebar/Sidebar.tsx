import { Fragment, HTMLProps, useContext } from "react";
import cx from "classnames";
import { SidebarItem } from "./SidebarItem.tsx";
import { SidebarGroup } from "./SidebarGroup.tsx";
import "./Sidebar.sass";
import { PageContext } from "../../contexts/PageContext.tsx";

export type SidebarProps = HTMLProps<HTMLDivElement> & {};

export function Sidebar(props: Readonly<SidebarProps>) {
    const { children, ...others } = props;

    const { showSidebar, setShowSidebar } = useContext(PageContext);

    return (
        <Fragment>
            <div
                className={cx("sidebar-overlay", {
                    "sidebar-overlay-visible": showSidebar,
                })}
                onClick={() => setShowSidebar?.(false)}
            />
            <nav
                className={cx("sidebar", {
                    "sidebar-visible": showSidebar,
                })}
                {...others}
            >
                <div className="sidebar-groups">{children}</div>
            </nav>
        </Fragment>
    );
}

Sidebar.Item = SidebarItem;
Sidebar.Group = SidebarGroup;
