import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../components/Sidebar/Sidebar";

import styles from "./Settings.module.sass";
import { Outlet } from "react-router-dom";
import { BigTitle } from "../../components/Text/Text";
import { Fragment } from "react";

export default function Settings() {
    return (
        <Fragment>
            <div className={styles.title}>
                <BigTitle>Settings</BigTitle>
            </div>
            <div className={styles.content}>
                <Sidebar>
                    <SidebarGroup title="Settings">
                        <SidebarItem
                            to="/settings/theme"
                            symbol="palette"
                            name="Theme"
                        />
                    </SidebarGroup>
                    <SidebarGroup title="Administration">
                        <SidebarItem
                            to="/settings/notifications"
                            symbol="notifications"
                            name="Notifications"
                        />
                        <SidebarItem
                            to="/settings/updates"
                            symbol="update"
                            name="Updates"
                        />
                        <SidebarItem
                            to="/settings/about"
                            symbol="info"
                            name="About"
                        />
                    </SidebarGroup>
                </Sidebar>
                <div className={styles.side}>
                    <Outlet />
                </div>
            </div>
        </Fragment>
    );
}