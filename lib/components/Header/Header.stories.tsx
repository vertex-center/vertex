import { Meta, StoryObj } from "@storybook/react";
import { Header } from "./Header.tsx";
import { Title } from "../Title/Title.tsx";
import { LinkProps } from "../Link/Link.tsx";
import { HTMLProps } from "react";
import { MaterialIcon } from "../MaterialIcon/MaterialIcon.tsx";

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

export const Normal: Story = {
    args: {
        appName: "Vertex App",
        linkBack: linkBackProps,
        linkLogo: linkLogoProps,
        leading: <MaterialIcon icon="arrow_back" />,
    },
    render: (props) => (
        <Header {...props}>
            <Title variant="h1">Page title</Title>
        </Header>
    ),
};

export default meta;
