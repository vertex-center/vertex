import type { Meta, StoryObj } from "@storybook/react";
import { Button } from "./Button";

const meta: Meta<typeof Button> = {
    component: Button,
};

type Story = StoryObj<typeof Button>;

const Default: Story = {
    args: {
        type: "colored",
        children: "Button",
    },
    argTypes: {
        type: {
            control: "select",
            options: ["colored", "outlined", "text"],
        },
        disabled: {
            control: "boolean",
            defaultValue: false,
        },
        onClick: {
            action: "clicked",
        },
    },
    render: (props) => <Button {...props} />,
};

export const Colored: Story = {
    ...Default,
    args: {
        ...Default.args,
        type: "colored",
    },
};

export const Outlined: Story = {
    ...Default,
    args: {
        ...Default.args,
        type: "outlined",
    },
};

export default meta;
