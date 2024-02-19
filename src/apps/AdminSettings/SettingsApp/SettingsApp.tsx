import PageWithSidebar from "../../../components/PageWithSidebar/PageWithSidebar";
import { Sidebar, useTitle } from "@vertex-center/components";
import l from "../../../components/NavLink/navlink";
import { useSidebar } from "../../../hooks/useSidebar";
import {
    ClockClockwise,
    Info,
    ListChecks,
    Notification,
} from "@phosphor-icons/react";

export default function SettingsApp() {
    useTitle("Settings");

    let sidebar = useSidebar(
        <Sidebar>
            <Sidebar.Group title="Administration">
                <Sidebar.Item
                    label="Notifications"
                    icon={<Notification />}
                    link={l("/admin/notifications")}
                />
                <Sidebar.Item
                    label="Updates"
                    icon={<ClockClockwise />}
                    link={l("/admin/updates")}
                />
                <Sidebar.Item
                    label="Checks"
                    icon={<ListChecks />}
                    link={l("/admin/checks")}
                />
                <Sidebar.Item
                    label="About"
                    icon={<Info />}
                    link={l("/admin/about")}
                />
            </Sidebar.Group>
        </Sidebar>
    );

    return <PageWithSidebar sidebar={sidebar} />;
}
