import PageWithSidebar from "../../../../components/PageWithSidebar/PageWithSidebar";
import { useSidebar } from "../../../../hooks/useSidebar";
import { MaterialIcon, Sidebar, useTitle } from "@vertex-center/components";
import l from "../../../../components/NavLink/navlink";

export default function Account() {
    useTitle("My Account");

    const sidebar = useSidebar(
        <Sidebar>
            <Sidebar.Group title="General">
                <Sidebar.Item
                    label="Information"
                    icon={<MaterialIcon icon="account_circle" />}
                    link={l("/account/info")}
                />
                <Sidebar.Item
                    label="Security"
                    icon={<MaterialIcon icon="security" />}
                    link={l("/account/security")}
                />
            </Sidebar.Group>
        </Sidebar>
    );

    return <PageWithSidebar sidebar={sidebar}>content</PageWithSidebar>;
}
