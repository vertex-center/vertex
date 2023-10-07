import { Fragment } from "react";
import { BigTitle } from "../../../components/Text/Text";
import styles from "./TunnelsApp.module.sass";
import { Horizontal } from "../../../components/Layouts/Layouts";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { SiCloudflare } from "@icons-pack/react-simple-icons";
import { Outlet } from "react-router-dom";

export default function TunnelsApp() {
    return (
        <Fragment>
            <BigTitle className={styles.title}>Tunnels</BigTitle>

            <Horizontal className={styles.content}>
                <Sidebar root="/tunnels">
                    <SidebarGroup title="Providers">
                        <SidebarItem
                            symbol={<SiCloudflare size={20} />}
                            to="/tunnels/cloudflare"
                            name="Cloudflare Tunnel"
                        />
                    </SidebarGroup>
                </Sidebar>
                <div className={styles.side}>
                    <Outlet />
                </div>
            </Horizontal>
        </Fragment>
    );
}
