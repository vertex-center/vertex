import { Meta, StoryObj } from "@storybook/react";
import { TextField } from "./TextField.tsx";
import { Normal as NormalInput } from "../Input/Input.stories.tsx";

const meta: Meta<typeof TextField> = {
    title: "Components/Fields/Text Field",
    component: TextField,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof TextField>;

export const Normal: Story = {
    ...NormalInput,
    render: (props) => <TextField {...props} />,
};

export default meta;
