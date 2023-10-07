import { BigTitle } from "../../../components/Text/Text";
import { Fragment } from "react";
import styles from "./MonitoringApp.module.sass";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { Horizontal } from "../../../components/Layouts/Layouts";
import { SiGrafana, SiPrometheus } from "@icons-pack/react-simple-icons";
import { Outlet } from "react-router-dom";

export default function MonitoringApp() {
    return (
        <Fragment>
            <BigTitle className={styles.title}>Monitoring</BigTitle>

            <Horizontal className={styles.content}>
                <Sidebar root="/monitoring">
                    <SidebarGroup title="Overview">
                        <SidebarItem
                            symbol="rule"
                            to="/monitoring/metrics"
                            name="Metrics"
                        />
                    </SidebarGroup>
                    <SidebarGroup title="Collectors">
                        <SidebarItem
                            symbol={<SiPrometheus size={20} />}
                            to="/monitoring/prometheus"
                            name="Prometheus"
                        />
                    </SidebarGroup>
                    <SidebarGroup title="Visualizations">
                        <SidebarItem
                            symbol={<SiGrafana size={20} />}
                            to="/monitoring/grafana"
                            name="Grafana"
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
