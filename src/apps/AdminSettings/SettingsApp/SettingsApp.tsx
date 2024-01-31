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
                    link={l("/admin/theme")}
                />
            </Sidebar.Group>
            <Sidebar.Group title="Administration">
                <Sidebar.Item
                    label="Notifications"
                    icon={<MaterialIcon icon="notifications" />}
                    link={l("/admin/notifications")}
                />
                <Sidebar.Item
                    label="Updates"
                    icon={<MaterialIcon icon="update" />}
                    link={l("/admin/updates")}
                />
                <Sidebar.Item
                    label="Checks"
                    icon={<MaterialIcon icon="checklist" />}
                    link={l("/admin/checks")}
                />
                <Sidebar.Item
                    label="About"
                    icon={<MaterialIcon icon="info" />}
                    link={l("/admin/about")}
                />
            </Sidebar.Group>
        </Sidebar>
    );

    return <PageWithSidebar sidebar={sidebar} />;
}
