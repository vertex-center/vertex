import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { Fragment } from "react";
import { SiPostgresql } from "@icons-pack/react-simple-icons";
import { useServerEvent } from "../../../hooks/useEvent";
import { useQueryClient } from "@tanstack/react-query";
import { useContainers } from "../../Containers/hooks/useContainers";
import { Sidebar, useTitle } from "@vertex-center/components";
import l from "../../../components/NavLink/navlink";
import { ContainerLed } from "../../../components/ContainerLed/ContainerLed";
import { useSidebar } from "../../../hooks/useSidebar";
import { Database, DownloadSimple } from "@phosphor-icons/react";

export default function SqlApp() {
    useTitle("SQL databases");

    const queryClient = useQueryClient();
    const { containers, isLoading } = useContainers({
        tags: ["Vertex SQL"],
    });

    const sidebar = useSidebar(
        <Sidebar>
            {Object.values(containers ?? {}).length > 0 && (
                <Sidebar.Group title="DBMS">
                    {Object.values(containers ?? {})?.map((c) => {
                        let icon = <Database />;
                        const type = c?.service?.features?.databases?.find(
                            (d) => d.category === "sql"
                        )?.type;
                        if (type === "postgres") {
                            icon = <SiPostgresql />;
                        }

                        return (
                            <Sidebar.Item
                                key={c?.id}
                                label={c?.name}
                                icon={icon}
                                link={l(`/sql/db/${c?.id}`)}
                                trailing={
                                    c && (
                                        <ContainerLed
                                            small
                                            status={c?.status}
                                        />
                                    )
                                }
                            />
                        );
                    })}
                </Sidebar.Group>
            )}
            <Sidebar.Group title="Create">
                <Sidebar.Item
                    label="Installer"
                    icon={<DownloadSimple />}
                    link={l("/sql/install")}
                />
            </Sidebar.Group>
        </Sidebar>
    );

    // @ts-ignore
    useServerEvent(window.api_urls.containers, "/containers/events", {
        change: () => {
            queryClient.invalidateQueries({
                queryKey: ["containers"],
            });
        },
    });

    return (
        <Fragment>
            <ProgressOverlay show={isLoading} />
            <PageWithSidebar sidebar={sidebar} />
        </Fragment>
    );
}
