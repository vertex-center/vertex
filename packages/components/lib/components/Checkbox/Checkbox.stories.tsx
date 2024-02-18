import { Meta, StoryObj } from "@storybook/react";
import { Checkbox } from "./Checkbox.tsx";
import { useState } from "react";

const meta: Meta<typeof Checkbox> = {
    title: "Components/Checkbox",
    component: Checkbox,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof Checkbox>;

export const Normal: Story = {
    argTypes: {
        checked: {
            control: "boolean",
        },
        onClick: { action: "onClick" },
    },
    render: function Render(props) {
        const [checked, setChecked] = useState<boolean>(true);
        return (
            <Checkbox
                {...props}
                checked={checked}
                onClick={() => setChecked((c) => !c)}
            />
        );
    },
};

export default meta;
