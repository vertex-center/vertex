import { Meta, StoryObj } from "@storybook/react";
import { Input } from "./Input.tsx";

const meta: Meta<typeof Input> = {
    title: "Components/Input",
    component: Input,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof Input>;

export const Normal: Story = {
    args: {
        id: "input",
        placeholder: "Placeholder",
        disabled: false,
        required: true,
    },
    argTypes: {
        placeholder: {
            control: "text",
        },
        disabled: {
            control: "boolean",
        },
        required: {
            control: "boolean",
        },
        onChange: { action: "onChange" },
    },
    render: (props) => <Input {...props} />,
};

export default meta;
