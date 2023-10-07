import Sidebar, {
    SidebarGroup,
    SidebarItem,
} from "../../../components/Sidebar/Sidebar";
import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";

export default function SettingsApp() {
    let sidebar = (
        <Sidebar root="/settings">
            <SidebarGroup title="Settings">
                <SidebarItem
                    to="/settings/theme"
                    symbol="palette"
                    name="Theme"
                />
            </SidebarGroup>
            <SidebarGroup title="Administration">
                <SidebarItem
                    to="/settings/notifications"
                    symbol="notifications"
                    name="Notifications"
                />
                <SidebarItem
                    to="/settings/hardware"
                    symbol="hard_drive"
                    name="Hardware"
                />
                <SidebarItem
                    to="/settings/security"
                    symbol="key"
                    name="Security"
                />
                <SidebarItem
                    to="/settings/updates"
                    symbol="update"
                    name="Updates"
                />
                <SidebarItem to="/settings/about" symbol="info" name="About" />
            </SidebarGroup>
        </Sidebar>
    );

    return <PageWithSidebar title="Settings" sidebar={sidebar} />;
}
