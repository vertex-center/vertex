import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { SiGrafana, SiPrometheus } from "@icons-pack/react-simple-icons";
import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";

export default function MonitoringApp() {
    const sidebar = (
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
    );

    return <PageWithSidebar title="Monitoring" sidebar={sidebar} />;
}
