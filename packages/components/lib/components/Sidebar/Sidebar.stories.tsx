import { Meta, StoryObj } from "@storybook/react";
import { Sidebar } from "./Sidebar.tsx";
import { MaterialIcon } from "../../../index.ts";

const meta: Meta<typeof Sidebar> = {
    title: "Components/Sidebar",
    component: Sidebar,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof Sidebar>;

export const Normal: Story = {
    render: function Render(props) {
        return (
            <Sidebar {...props}>
                <Sidebar.Group title="Group A">
                    <Sidebar.Item
                        label="Element A"
                        icon={<MaterialIcon icon="database" />}
                        link={{ href: "/app/A" }}
                    />
                    <Sidebar.Item
                        label="Element B"
                        icon={<MaterialIcon icon="database" />}
                        link={{ href: "/app/B" }}
                    />
                </Sidebar.Group>
                <Sidebar.Group title="Group B">
                    <Sidebar.Item
                        label="Element C"
                        icon={<MaterialIcon icon="database" />}
                        link={{ href: "/app/C" }}
                    />
                </Sidebar.Group>
            </Sidebar>
        );
    },
};

export default meta;
