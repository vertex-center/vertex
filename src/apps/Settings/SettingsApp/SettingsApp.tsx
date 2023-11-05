import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import { MaterialIcon, Sidebar, useTitle } from "@vertex-center/components";
import { useLocation } from "react-router-dom";
import l from "../../../components/NavLink/navlink";

export default function SettingsApp() {
    useTitle("Settings");

    const { pathname } = useLocation();

    let sidebar = (
        <Sidebar rootUrl="/settings" currentUrl={pathname}>
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
