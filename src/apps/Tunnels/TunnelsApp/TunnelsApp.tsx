import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import { SiCloudflare } from "@icons-pack/react-simple-icons";

export default function TunnelsApp() {
    const sidebar = (
        <Sidebar root="/tunnels">
            <SidebarGroup title="Providers">
                <SidebarItem
                    symbol={<SiCloudflare size={20} />}
                    to="/tunnels/cloudflare"
                    name="Cloudflare Tunnel"
                />
            </SidebarGroup>
        </Sidebar>
    );

    return <PageWithSidebar title="Tunnels" sidebar={sidebar} />;
}
