import { Meta, StoryObj } from "@storybook/react";
import { Header } from "./Header.tsx";
import { Title } from "../Title/Title.tsx";
import { LinkProps } from "../Link/Link.tsx";
import { HTMLProps } from "react";
import { MaterialIcon } from "../MaterialIcon/MaterialIcon.tsx";
import { ProfilePicture } from "../ProfilePicture/ProfilePicture.tsx";
import { HeaderItem } from "./HeaderItem.tsx";
import { DropdownItem } from "../Dropdown/Dropdown.tsx";
import { SignOut } from "@phosphor-icons/react";

const meta: Meta<typeof Header> = {
    title: "Components/Header",
    component: Header,
    tags: ["autodocs"],
    parameters: {
        controls: { expanded: true },
        layout: "",
    },
};

type Story = StoryObj<typeof Header>;

const linkLogoProps: LinkProps<HTMLProps<HTMLAnchorElement>> = {
    href: "#",
};

const linkBackProps: LinkProps<HTMLProps<HTMLAnchorElement>> = {
    href: "#",
};

const items = (
    <DropdownItem icon={<SignOut />} red>
        Logout
    </DropdownItem>
);

export const Normal: Story = {
    args: {
        appName: "Vertex App",
        linkBack: linkBackProps,
        linkLogo: linkLogoProps,
        leading: <MaterialIcon icon="arrow_back" />,
        trailing: (
            <HeaderItem items={items}>
                Arra
                <ProfilePicture size={36} />
            </HeaderItem>
        ),
    },
    render: (props) => (
        <Header {...props}>
            <Title variant="h1">Page title</Title>
        </Header>
    ),
};

export default meta;
