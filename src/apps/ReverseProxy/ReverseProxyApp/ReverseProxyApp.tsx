import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { useTitle } from "@vertex-center/components";

export default function ReverseProxyApp() {
    useTitle("Reverse Proxy");

    const sidebar = (
        <Sidebar root="/app/vx-reverse-proxy">
            <SidebarGroup title="Providers">
                <SidebarItem
                    name="Vertex Reverse Proxy"
                    icon="router"
                    to="/app/vx-reverse-proxy/vertex"
                />
            </SidebarGroup>
        </Sidebar>
    );

    return <PageWithSidebar sidebar={sidebar} />;
}
