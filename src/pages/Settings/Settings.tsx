import { HeaderHome } from "../../components/Header/Header";
import Sidebar, {
    SidebarItem,
    SidebarSeparator,
    SidebarTitle,
} from "../../components/Sidebar/Sidebar";
import { Horizontal } from "../../components/Layouts/Layouts";

import styles from "./Settings.module.sass";
import { Outlet } from "react-router-dom";

type Props = {};

export default function Settings(props: Props) {
    return (
        <div>
            <HeaderHome />
            <Horizontal className={styles.content}>
                <Sidebar>
                    <SidebarTitle>Settings</SidebarTitle>
                    <SidebarItem
                        to="/settings/theme"
                        symbol="palette"
                        name="Theme"
                    />
                    <SidebarSeparator />
                    <SidebarTitle>Administration</SidebarTitle>
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
                </Sidebar>
                <div className={styles.side}>
                    <Outlet />
                </div>
            </Horizontal>
        </div>
    );
}
