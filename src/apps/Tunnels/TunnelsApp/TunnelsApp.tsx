import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import { SiCloudflare } from "@icons-pack/react-simple-icons";
import { Fragment } from "react";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { useServerEvent } from "../../../hooks/useEvent";
import { useQueryClient } from "@tanstack/react-query";
import { useContainers } from "../../Containers/hooks/useContainers";
import { Sidebar, useTitle } from "@vertex-center/components";
import l from "../../../components/NavLink/navlink";
import { ContainerLed } from "../../../components/ContainerLed/ContainerLed";
import { useSidebar } from "../../../hooks/useSidebar";

export default function TunnelsApp() {
    useTitle("Tunnels");

    const queryClient = useQueryClient();

    const { containers, isLoading } = useContainers({
        tags: ["Vertex Tunnels - Cloudflare"],
    });

    const container = Object.values(containers || {})?.[0];

    // @ts-ignore
    useServerEvent(window.api_urls.containers, "/containers/events", {
        change: (e) => {
            queryClient.invalidateQueries({
                queryKey: ["containers"],
            });
        },
    });

    const sidebar = useSidebar(
        <Sidebar>
            <Sidebar.Group title="Providers">
                <Sidebar.Item
                    label="Cloudflare Tunnel"
                    icon={<SiCloudflare size={20} />}
                    link={l("/app/tunnels/cloudflare")}
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
