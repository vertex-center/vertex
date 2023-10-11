import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { SiGrafana, SiPrometheus } from "@icons-pack/react-simple-icons";
import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import { api } from "../../../backend/backend";
import { Fragment } from "react";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useServerEvent } from "../../../hooks/useEvent";
import { useQuery, useQueryClient } from "@tanstack/react-query";

export default function MonitoringApp() {
    const queryClient = useQueryClient();
    const { data: instances, isLoading } = useQuery({
        queryKey: ["instances"],
        queryFn: api.vxInstances.instances.all,
    });

    const ledPrometheus = Object.values(instances || {}).find((inst) =>
        inst.tags?.includes("vertex-prometheus-collector")
    );

    const ledGrafana = Object.values(instances || {}).find((inst) =>
        inst.tags?.includes("vertex-grafana-visualizer")
    );

    useServerEvent("/app/vx-instances/instances/events", {
        change: (e) => {
            console.log(e);
            queryClient.invalidateQueries({
                queryKey: ["instances"],
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
                    led={ledPrometheus && { status: ledPrometheus?.status }}
                />
            </SidebarGroup>
            <SidebarGroup title="Visualizations">
                <SidebarItem
                    icon={<SiGrafana size={20} />}
                    to="/app/vx-monitoring/grafana"
                    name="Grafana"
                    led={ledGrafana && { status: ledGrafana?.status }}
                />
            </SidebarGroup>
        </Sidebar>
    );

    return (
        <Fragment>
            <ProgressOverlay show={isLoading} />
            <PageWithSidebar title="Monitoring" sidebar={sidebar} />
        </Fragment>
    );
}
