import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { SiCloudflare } from "@icons-pack/react-simple-icons";
import { Fragment } from "react";
import { api } from "../../../backend/api/backend";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useServerEvent } from "../../../hooks/useEvent";
import { useQuery, useQueryClient } from "@tanstack/react-query";

export default function TunnelsApp() {
    const queryClient = useQueryClient();
    const { data: containers, isLoading } = useQuery({
        queryKey: ["containers"],
        queryFn: api.vxContainers.containers.all,
    });

    const ledCloudflared = Object.values(containers || {}).find((inst) =>
        inst.tags?.includes("Vertex Tunnels - Cloudflare")
    );

    useServerEvent("/app/vx-containers/containers/events", {
        change: (e) => {
            queryClient.invalidateQueries({
                queryKey: ["containers"],
            });
        },
    });

    const sidebar = (
        <Sidebar root="/app/vx-tunnels">
            <SidebarGroup title="Providers">
                <SidebarItem
                    icon={<SiCloudflare size={20} />}
                    to="/app/vx-tunnels/cloudflare"
                    name="Cloudflare Tunnel"
                    led={ledCloudflared && { status: ledCloudflared?.status }}
                />
            </SidebarGroup>
        </Sidebar>
    );

    return (
        <Fragment>
            <ProgressOverlay show={isLoading} />
            <PageWithSidebar title="Tunnels" sidebar={sidebar} />
        </Fragment>
    );
}
