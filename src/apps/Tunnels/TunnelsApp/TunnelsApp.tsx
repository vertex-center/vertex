import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { SiCloudflare } from "@icons-pack/react-simple-icons";
import { Fragment, useEffect } from "react";
import {
    registerSSE,
    registerSSEEvent,
    unregisterSSE,
    unregisterSSEEvent,
} from "../../../backend/sse";
import { useFetch } from "../../../hooks/useFetch";
import { Instances } from "../../../models/instance";
import { api } from "../../../backend/backend";
import { ProgressOverlay } from "../../../components/Progress/Progress";

export default function TunnelsApp() {
    const {
        data: instances,
        loading,
        reload: reloadInstances,
    } = useFetch<Instances>(api.instances.get);

    const ledCloudflared = Object.values(instances || {}).find((inst) =>
        inst.tags?.includes("vertex-cloudflare-tunnel")
    );

    useEffect(() => {
        const sse = registerSSE("/instances/events");

        const onChange = (e) => {
            console.log(e);
            reloadInstances();
        };

        registerSSEEvent(sse, "change", onChange);

        return () => {
            unregisterSSEEvent(sse, "change", onChange);

            unregisterSSE(sse);
        };
    }, [instances]);

    const sidebar = (
        <Sidebar root="/tunnels">
            <SidebarGroup title="Providers">
                <SidebarItem
                    symbol={<SiCloudflare size={20} />}
                    to="/tunnels/cloudflare"
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
