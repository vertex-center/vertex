import PageWithSidebar from "../../../../components/PageWithSidebar/PageWithSidebar";
import { useSidebar } from "../../../../hooks/useSidebar";
import { Sidebar, useTitle } from "@vertex-center/components";
import l from "../../../../components/NavLink/navlink";
import { Envelope, Key, Palette, UserCircle } from "@phosphor-icons/react";

export default function Account() {
    useTitle("My Account");

    const sidebar = useSidebar(
        <Sidebar>
            <Sidebar.Group title="General">
                <Sidebar.Item
                    label="Information"
                    icon={<UserCircle />}
                    link={l("/account/info")}
                />
                <Sidebar.Item
                    label="Security"
                    icon={<Key />}
                    link={l("/account/security")}
                />
                <Sidebar.Item
                    label="Emails"
                    icon={<Envelope />}
                    link={l("/account/emails")}
                />
            </Sidebar.Group>
            <Sidebar.Group title="Appearance">
                <Sidebar.Item
                    label="Theme"
                    icon={<Palette />}
                    link={l("/account/theme")}
                />
            </Sidebar.Group>
        </Sidebar>
    );

    return <PageWithSidebar sidebar={sidebar}>content</PageWithSidebar>;
}
