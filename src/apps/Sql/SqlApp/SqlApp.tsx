import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { api } from "../../../backend/api/backend";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { Fragment } from "react";
import { SiPostgresql } from "@icons-pack/react-simple-icons";
import { useServerEvent } from "../../../hooks/useEvent";
import { useQuery, useQueryClient } from "@tanstack/react-query";

export default function SqlApp() {
    const queryClient = useQueryClient();
    const { data: instances, isLoading } = useQuery({
        queryKey: ["instances"],
        queryFn: api.vxInstances.instances.all,
    });

    const dbs = Object.values(instances ?? {})?.filter((i) =>
        i?.tags?.some((t) => t.includes("vertex-") && t.includes("-sql"))
    );

    const sidebar = (
        <Sidebar root="/app/vx-sql">
            {Object.values(dbs).length > 0 && (
                <SidebarGroup title="DBMS">
                    {dbs?.map((inst) => {
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

    useServerEvent("/app/vx-instances/instances/events", {
        change: () => {
            queryClient.invalidateQueries({
                queryKey: ["instances"],
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
