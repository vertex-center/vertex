import type { Meta, StoryObj } from "@storybook/react";
import { Button } from "./Button";
import { MaterialIcon } from "../../../index.ts";

const meta: Meta<typeof Button> = {
    title: "Components/Button",
    component: Button,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof Button>;

const icon = {
    control: "select",
    options: ["", "deployed_code_update", "arrow_back", "arrow_forward"],
    mapping: {
        "": null,
        deployed_code_update: <MaterialIcon icon="deployed_code_update" />,
        arrow_back: <MaterialIcon icon="arrow_back" />,
        arrow_forward: <MaterialIcon icon="arrow_forward" />,
    },
};

const Default: Story = {
    args: {
        variant: "colored",
        children: "Button",
    },
    argTypes: {
        type: {
            control: "select",
            options: ["colored", "outlined", "danger"],
        },
        disabled: {
            control: "boolean",
            defaultValue: false,
        },
        leftIcon: icon,
        rightIcon: icon,
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
        variant: "colored",
    },
};

export const Outlined: Story = {
    ...Default,
    args: {
        ...Default.args,
        variant: "outlined",
    },
};

export const Danger: Story = {
    ...Default,
    args: {
        ...Default.args,
        variant: "danger",
    },
};

export const WithLeftIcon: Story = {
    ...Default,
    args: {
        ...Default.args,
        leftIcon: <MaterialIcon icon="deployed_code_update" />,
    },
};

export const WithRightIcon: Story = {
    ...Default,
    args: {
        ...Default.args,
        rightIcon: <MaterialIcon icon="deployed_code_update" />,
    },
};

export default meta;
