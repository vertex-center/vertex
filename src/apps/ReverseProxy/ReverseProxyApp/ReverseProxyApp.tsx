import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import { MaterialIcon, Sidebar, useTitle } from "@vertex-center/components";
import l from "../../../components/NavLink/navlink";
import React from "react";
import { useSidebar } from "../../../hooks/useSidebar";

export default function ReverseProxyApp() {
    useTitle("Reverse Proxy");

    const sidebar = useSidebar(
        <Sidebar>
            <Sidebar.Group title="Providers">
                <Sidebar.Item
                    label="Vertex Reverse Proxy"
                    icon={<MaterialIcon icon="router" />}
                    link={l("/app/reverse-proxy/vertex")}
                />
            </Sidebar.Group>
        </Sidebar>
    );

    return <PageWithSidebar sidebar={sidebar} />;
}
