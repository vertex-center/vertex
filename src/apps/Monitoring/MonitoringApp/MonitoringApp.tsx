import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { SiGrafana, SiPrometheus } from "@icons-pack/react-simple-icons";
import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import { Fragment } from "react";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useServerEvent } from "../../../hooks/useEvent";
import { useQueryClient } from "@tanstack/react-query";
import { useContainers } from "../../Containers/hooks/useContainers";
import { useTitle } from "../../../hooks/useTitle";

export default function MonitoringApp() {
    useTitle("Monitoring");

    const queryClient = useQueryClient();

    const { containers: prometheusContainers, isLoading: isLoadingPrometheus } =
        useContainers({
            tags: ["Vertex Monitoring - Prometheus Collector"],
        });
    const prometheusContainer = Object.values(prometheusContainers ?? {})[0];

    const { containers: grafanaContainers, isLoading: isLoadingGrafana } =
        useContainers({
            tags: ["Vertex Monitoring - Grafana Visualizer"],
        });
    const grafanaContainer = Object.values(grafanaContainers ?? {})[0];

    useServerEvent("/app/vx-containers/containers/events", {
        change: (e) => {
            queryClient.invalidateQueries({
                queryKey: ["containers"],
            });
        },
    });

    const sidebar = (
        <Sidebar root="/app/vx-monitoring">
            <SidebarGroup title="Overview">
                <SidebarItem
                    icon="rule"
                    to="/app/vx-monitoring/metrics"
                    name="Metrics"
                />
            </SidebarGroup>
            <SidebarGroup title="Collectors">
                <SidebarItem
                    icon={<SiPrometheus size={20} />}
                    to="/app/vx-monitoring/prometheus"
                    name="Prometheus"
                    led={
                        prometheusContainer && {
                            status: prometheusContainer?.status,
                        }
                    }
                />
            </SidebarGroup>
            <SidebarGroup title="Visualizations">
                <SidebarItem
                    icon={<SiGrafana size={20} />}
                    to="/app/vx-monitoring/grafana"
                    name="Grafana"
                    led={
                        grafanaContainer && {
                            status: grafanaContainer?.status,
                        }
                    }
                />
            </SidebarGroup>
        </Sidebar>
    );

    return (
        <Fragment>
            <ProgressOverlay show={isLoadingGrafana || isLoadingPrometheus} />
            <PageWithSidebar sidebar={sidebar} />
        </Fragment>
    );
}
