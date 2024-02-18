import { Meta, StoryObj } from "@storybook/react";
import { ProfilePicture } from "./ProfilePicture.tsx";

const meta: Meta = {
    title: "Components/Profile Picture",
    component: ProfilePicture,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof ProfilePicture>;

export const Normal: Story = {
    args: {
        src: "https://picsum.photos/200",
        alt: "Profile Picture",
        size: 40,
    },
    render: (props) => <ProfilePicture {...props} />,
};

export default meta;
