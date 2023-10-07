import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";

export default function SettingsApp() {
    let sidebar = (
        <Sidebar root="/settings">
            <SidebarGroup title="Settings">
                <SidebarItem to="/settings/theme" icon="palette" name="Theme" />
            </SidebarGroup>
            <SidebarGroup title="Administration">
                <SidebarItem
                    to="/settings/notifications"
                    icon="notifications"
                    name="Notifications"
                />
                <SidebarItem
                    to="/settings/hardware"
                    icon="hard_drive"
                    name="Hardware"
                />
                <SidebarItem
                    to="/settings/security"
                    icon="key"
                    name="Security"
                />
                <SidebarItem
                    to="/settings/updates"
                    icon="update"
                    name="Updates"
                />
                <SidebarItem to="/settings/about" icon="info" name="About" />
            </SidebarGroup>
        </Sidebar>
    );

    return <PageWithSidebar title="Settings" sidebar={sidebar} />;
}
