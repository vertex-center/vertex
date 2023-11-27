import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import { MaterialIcon, Sidebar, useTitle } from "@vertex-center/components";
import l from "../../../components/NavLink/navlink";
import { useSidebar } from "../../../hooks/useSidebar";

export default function SettingsApp() {
    useTitle("Settings");

    let sidebar = useSidebar(
        <Sidebar>
            <Sidebar.Group title="Settings">
                <Sidebar.Item
                    label="Theme"
                    icon={<MaterialIcon icon="palette" />}
                    link={l("/settings/theme")}
                />
            </Sidebar.Group>
            <Sidebar.Group title="Administration">
                <Sidebar.Item
                    label="Notifications"
                    icon={<MaterialIcon icon="notifications" />}
                    link={l("/settings/notifications")}
                />
                <Sidebar.Item
                    label="Hardware"
                    icon={<MaterialIcon icon="hard_drive" />}
                    link={l("/settings/hardware")}
                />
                <Sidebar.Item
                    label="Database"
                    icon={<MaterialIcon icon="database" />}
                    link={l("/settings/database")}
                />
                <Sidebar.Item
                    label="Security"
                    icon={<MaterialIcon icon="key" />}
                    link={l("/settings/security")}
                />
                <Sidebar.Item
                    label="Updates"
                    icon={<MaterialIcon icon="update" />}
                    link={l("/settings/updates")}
                />
                <Sidebar.Item
                    label="About"
                    icon={<MaterialIcon icon="info" />}
                    link={l("/settings/about")}
                />
            </Sidebar.Group>
        </Sidebar>
    );

    return <PageWithSidebar sidebar={sidebar} />;
}
