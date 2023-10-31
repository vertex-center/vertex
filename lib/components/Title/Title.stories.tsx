import { Meta, StoryObj } from "@storybook/react";
import { Title } from "./Title.tsx";

const meta: Meta<typeof Title> = {
    title: "Components/Title",
    component: Title,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof Title>;

export const Default: Story = {
    args: {
        variant: "h1",
        children: "Title",
    },
    argTypes: {
        variant: {
            control: "select",
            options: ["h1", "h2", "h3", "h4", "h5", "h6"],
        },
    },
    render: (props) => <Title {...props} />,
};

export const All: Story = {
    args: {
        children: "Title",
    },
    render: (props) => (
        <>
            <Title variant="h1" {...props}>
                Title 1
            </Title>
            <Title variant="h2" {...props}>
                Title 2
            </Title>
            <Title variant="h3" {...props}>
                Title 3
            </Title>
            <Title variant="h4" {...props}>
                Title 4
            </Title>
            <Title variant="h5" {...props}>
                Title 5
            </Title>
            <Title variant="h6" {...props}>
                Title 6
            </Title>
        </>
    ),
};

export default meta;
