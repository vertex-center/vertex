import { Meta, StoryObj } from "@storybook/react";
import { Paragraph } from "./Paragraph.tsx";

const meta: Meta = {
    title: "Typography/Paragraph",
    component: Paragraph,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof Paragraph>;

export const Default: Story = {
    args: {
        children:
            "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
    },
    render: (props) => <Paragraph {...props} />,
};

export default meta;
