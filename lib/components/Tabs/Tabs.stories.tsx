import { Meta, StoryObj } from "@storybook/react";
import { Tabs } from "./Tabs.tsx";
import { TabItem } from "./TabItem.tsx";

const meta: Meta<typeof Tabs> = {
    title: "Components/Tabs",
    component: Tabs,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof Tabs>;

export const Normal: Story = {
    render: function Render(props) {
        return (
            <Tabs {...props}>
                <TabItem label="Element A">Content A</TabItem>
                <TabItem label="Element B">Content B</TabItem>
                <TabItem label="Element C">Content C</TabItem>
            </Tabs>
        );
    },
};

export default meta;
