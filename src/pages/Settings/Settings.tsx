import { HeaderHome } from "../../components/Header/Header";
import Sidebar, { SidebarItem } from "../../components/Sidebar/Sidebar";
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
                    <SidebarItem
                        to="/settings/theme"
                        symbol="palette"
                        name="Theme"
                    />
                </Sidebar>
                <div className={styles.side}>
                    <Outlet />
                </div>
            </Horizontal>
        </div>
    );
}
