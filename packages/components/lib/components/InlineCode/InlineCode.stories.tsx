import { Meta, StoryObj } from "@storybook/react";
import { InlineCode } from "./InlineCode.tsx";

const meta: Meta = {
    title: "Components/Inline Code",
    component: InlineCode,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof InlineCode>;

export const Default: Story = {
    args: {
        children: "Hello world",
    },
    render: function Render(props) {
        return (
            <>
                Some inline code:
                <InlineCode {...props} />.
            </>
        );
    },
};

export default meta;
