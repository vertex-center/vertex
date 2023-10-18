import { Meta, StoryObj } from "@storybook/react";
import { MaterialIcon } from "./MaterialIcon";

const meta: Meta<typeof MaterialIcon> = {
    component: MaterialIcon,
    tags: ["autodocs"],
    title: "Material Icon",
};

type Story = StoryObj<typeof MaterialIcon>;

export const Icon: Story = {
    name: "Normal",
    args: {
        name: "deployed_code_update",
    },
    render: (props) => <MaterialIcon {...props} />,
};

export default meta;
