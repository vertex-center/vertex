import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import { MaterialIcon, Sidebar, useTitle } from "@vertex-center/components";
import { useLocation } from "react-router-dom";
import l from "../../../components/NavLink/navlink";
import React from "react";

export default function ReverseProxyApp() {
    useTitle("Reverse Proxy");

    const { pathname } = useLocation();

    const sidebar = (
        <Sidebar rootUrl="/app/vx-reverse-proxy" currentUrl={pathname}>
            <Sidebar.Group title="Providers">
                <Sidebar.Item
                    label="Vertex Reverse Proxy"
                    icon={<MaterialIcon icon="router" />}
                    link={l("/app/vx-reverse-proxy/vertex")}
                />
            </Sidebar.Group>
        </Sidebar>
    );

    return <PageWithSidebar sidebar={sidebar} />;
}
