import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { Fragment } from "react";
import { SiPostgresql } from "@icons-pack/react-simple-icons";
import { useServerEvent } from "../../../hooks/useEvent";
import { useQueryClient } from "@tanstack/react-query";
import { useContainers } from "../../Containers/hooks/useContainers";
import { MaterialIcon, Sidebar, useTitle } from "@vertex-center/components";
import l from "../../../components/NavLink/navlink";
import { ContainerLed } from "../../../components/ContainerLed/ContainerLed";
import { useSidebar } from "../../../hooks/useSidebar";

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
                    {Object.values(containers ?? {})?.map((inst) => {
                        let icon = <MaterialIcon icon="database" />;
                        const type = inst?.service?.features?.databases?.find(
                            (d) => d.category === "sql"
                        )?.type;
                        if (type === "postgres") {
                            icon = <SiPostgresql />;
                        }

                        return (
                            <Sidebar.Item
                                key={inst.uuid}
                                label={inst?.display_name ?? inst.service.name}
                                icon={icon}
                                link={l(`/app/vx-sql/db/${inst.uuid}`)}
                                trailing={
                                    inst && (
                                        <ContainerLed
                                            small
                                            status={inst?.status}
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
                    icon={<MaterialIcon icon="download" />}
                    link={l("/app/vx-sql/install")}
                />
            </Sidebar.Group>
        </Sidebar>
    );

    useServerEvent("/app/vx-containers/containers/events", {
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
