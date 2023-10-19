import { Meta, StoryObj } from "@storybook/react";
import { MaterialIcon } from "./MaterialIcon";

const meta: Meta<typeof MaterialIcon> = {
    title: "Components/Material Icon",
    component: MaterialIcon,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof MaterialIcon>;

export const Icon: Story = {
    name: "Normal",
    args: {
        icon: "deployed_code_update",
    },
    render: (props) => <MaterialIcon {...props} />,
};

export default meta;
