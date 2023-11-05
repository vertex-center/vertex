import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import { SiCloudflare } from "@icons-pack/react-simple-icons";
import { Fragment } from "react";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useServerEvent } from "../../../hooks/useEvent";
import { useQueryClient } from "@tanstack/react-query";
import { useContainers } from "../../Containers/hooks/useContainers";
import { Sidebar, useTitle } from "@vertex-center/components";
import { useLocation } from "react-router-dom";
import l from "../../../components/NavLink/navlink";
import { ContainerLed } from "../../../components/ContainerLed/ContainerLed";

export default function TunnelsApp() {
    useTitle("Tunnels");

    const { pathname } = useLocation();
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
        <Sidebar rootUrl="/app/vx-tunnels" currentUrl={pathname}>
            <Sidebar.Group title="Providers">
                <Sidebar.Item
                    label="Cloudflare Tunnel"
                    icon={<SiCloudflare size={20} />}
                    link={l("/app/vx-tunnels/cloudflare")}
                    trailing={
                        container && (
                            <ContainerLed small status={container?.status} />
                        )
                    }
                />
            </Sidebar.Group>
        </Sidebar>
    );

    return (
        <Fragment>
            <ProgressOverlay show={isLoading} />
            <PageWithSidebar sidebar={sidebar} />
        </Fragment>
    );
}
