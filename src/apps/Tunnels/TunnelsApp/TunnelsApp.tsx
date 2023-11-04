import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { SiCloudflare } from "@icons-pack/react-simple-icons";
import { Fragment } from "react";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useServerEvent } from "../../../hooks/useEvent";
import { useQueryClient } from "@tanstack/react-query";
import { useContainers } from "../../Containers/hooks/useContainers";
import { useTitle } from "../../../hooks/useTitle";

export default function TunnelsApp() {
    useTitle("Tunnels");

    const queryClient = useQueryClient();

    const { containers, isLoading } = useContainers({
        tags: ["Vertex Tunnels - Cloudflare"],
    });

    const container = Object.values(containers || {})?.[0];

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
                    led={container && { status: container?.status }}
                />
            </SidebarGroup>
        </Sidebar>
    );

    return (
        <Fragment>
            <ProgressOverlay show={isLoading} />
            <PageWithSidebar sidebar={sidebar} />
        </Fragment>
    );
}
