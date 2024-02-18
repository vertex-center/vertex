import { Meta, StoryObj } from "@storybook/react";
import { TextField } from "./TextField.tsx";

const meta: Meta<typeof TextField> = {
    title: "Components/Fields/Text Field",
    component: TextField,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof TextField>;

export const Normal: Story = {
    args: {
        id: "input",
        placeholder: "Placeholder",
        disabled: false,
        label: "Label",
        required: true,
        description: "A short description",
        error: "",
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
    render: (props) => <TextField {...props} />,
};

export default meta;
