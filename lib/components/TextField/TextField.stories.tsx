import { Meta, StoryObj } from "@storybook/react";
import { TextField } from "./TextField.tsx";

const meta: Meta<typeof TextField> = {
    title: "Components/Text Field",
    component: TextField,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof TextField>;

export const Normal: Story = {
    args: {
        label: "Label",
        required: true,
        description: "A short description",
        error: "",
    },
    argTypes: {
        label: {
            control: "text",
        },
        required: {
            control: "boolean",
        },
    },
    render: (props) => <TextField {...props} />,
};

export default meta;
