import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { SiCloudflare } from "@icons-pack/react-simple-icons";
import { Fragment } from "react";
import { useFetch } from "../../../hooks/useFetch";
import { Instances } from "../../../models/instance";
import { api } from "../../../backend/backend";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useServerEvent } from "../../../hooks/useEvent";

export default function TunnelsApp() {
    const {
        data: instances,
        loading,
        reload: reloadInstances,
    } = useFetch<Instances>(api.instances.get);

    const ledCloudflared = Object.values(instances || {}).find((inst) =>
        inst.tags?.includes("vertex-cloudflare-tunnel")
    );

    useServerEvent("/instances/events", {
        change: (e) => {
            console.log(e);
            reloadInstances().finally();
        },
    });

    const sidebar = (
        <Sidebar root="/tunnels">
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
            <ProgressOverlay show={loading} />
            <PageWithSidebar title="Tunnels" sidebar={sidebar} />
        </Fragment>
    );
}
