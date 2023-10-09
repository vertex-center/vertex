import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { api } from "../../../backend/backend";
import { Instances } from "../../../models/instance";
import { useFetch } from "../../../hooks/useFetch";
import { ProgressOverlay } from "../../../components/Progress/Progress";
import { Fragment } from "react";
import { SiPostgresql } from "@icons-pack/react-simple-icons";
import { useServerEvent } from "../../../hooks/useEvent";

export default function SqlApp() {
    const {
        data: instances,
        loading,
        reload: reloadInstances,
    } = useFetch<Instances>(api.instances.get);

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

    useServerEvent("/instances/events", {
        change: () => {
            reloadInstances().finally();
        },
    });

    return (
        <Fragment>
            <ProgressOverlay show={loading} />
            <PageWithSidebar title="SQL databases" sidebar={sidebar} />
        </Fragment>
    );
}
