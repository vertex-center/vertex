import { Meta, StoryObj } from "@storybook/react";
import Box from "./Box.tsx";
import { Paragraph } from "../Paragraph/Paragraph.tsx";

const meta: Meta = {
    title: "Components/Box",
    component: Box,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof Box>;

export const Normal: Story = {
    args: {
        children: (
            <Paragraph>
                Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
                eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut
                enim ad minim veniam, quis nostrud exercitation ullamco laboris
                nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor
                in reprehenderit in voluptate velit esse cillum dolore eu fugiat
                nulla pariatur.
            </Paragraph>
        ),
    },
    argTypes: {
        type: {
            control: "select",
            options: ["info", "tip", "warning"],
        },
    },
    render: (props) => <Box {...props} />,
};

export default meta;
