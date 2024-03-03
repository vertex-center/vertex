import { Meta, StoryObj } from "@storybook/react";
import { Logo } from "./Logo";

const meta: Meta<typeof Logo> = {
    title: "Components/Logo",
    component: Logo,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof Logo>;

export const Default: Story = {
    render: (props) => <Logo {...props} />,
};

export default meta;
