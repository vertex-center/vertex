import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";

export default function ReverseProxyApp() {
    const sidebar = (
        <Sidebar root="/reverse-proxy">
            <SidebarGroup title="Providers">
                <SidebarItem
                    name="Vertex Reverse Proxy"
                    symbol="router"
                    to="/reverse-proxy/vertex"
                />
            </SidebarGroup>
        </Sidebar>
    );

    return <PageWithSidebar title="Reverse Proxy" sidebar={sidebar} />;
}
