import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { Fragment } from "react";
import { SiPostgresql } from "@icons-pack/react-simple-icons";
import { useServerEvent } from "../../../hooks/useEvent";
import { useQueryClient } from "@tanstack/react-query";
import { useContainers } from "../../Containers/hooks/useContainers";

export default function SqlApp() {
    const queryClient = useQueryClient();
    const { containers, isLoading } = useContainers({
        tags: ["Vertex SQL"],
    });

    const sidebar = (
        <Sidebar root="/app/vx-sql">
            {Object.values(containers ?? {}).length > 0 && (
                <SidebarGroup title="DBMS">
                    {Object.values(containers ?? {})?.map((inst) => {
                        let icon: string | JSX.Element = "database";
                        const type = inst?.service?.features?.databases?.find(
                            (d) => d.category === "sql"
                        )?.type;
                        if (type === "postgres") {
                            icon = <SiPostgresql />;
                        }

                        return (
                            <SidebarItem
                                key={inst.uuid}
                                icon={icon}
                                name={inst?.display_name ?? inst.service.name}
                                to={`/app/vx-sql/db/${inst.uuid}`}
                                led={{ status: inst?.status }}
                            />
                        );
                    })}
                </SidebarGroup>
            )}
            <SidebarGroup title="Create">
                <SidebarItem
                    icon="download"
                    name="Installer"
                    to="/app/vx-sql/install"
                />
            </SidebarGroup>
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
            <PageWithSidebar title="SQL databases" sidebar={sidebar} />
        </Fragment>
    );
}
