import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";

export default function SqlApp() {
    const sidebar = (
        <Sidebar root="/app/vx-sql">
            <SidebarGroup title="Create">
                <SidebarItem
                    icon="download"
                    name="Installer"
                    to="/app/vx-sql/install"
                />
            </SidebarGroup>
        </Sidebar>
    );

    return <PageWithSidebar title="SQL databases" sidebar={sidebar} />;
}
