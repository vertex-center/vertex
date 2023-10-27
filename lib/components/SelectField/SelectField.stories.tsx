import { Meta, StoryObj } from "@storybook/react";
import { SelectField, SelectOption } from "./SelectField.tsx";
import { useState } from "react";

const meta: Meta<typeof SelectField> = {
    title: "Components/Fields/Select Field",
    component: SelectField,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof SelectField>;

export const Normal: Story = {
    args: {
        id: "select-field",
        label: "Label",
        required: true,
        description: "A short description",
    },
    argTypes: {
        onChange: { action: "onChange" },
    },
    render: function Render(props) {
        const [value, setValue] = useState<string>();

        const onChange = (value: string) => {
            setValue(value);
            props.onChange?.(value);
        };

        return (
            <SelectField value={value} {...props} onChange={onChange}>
                <SelectOption value="1">One</SelectOption>
                <SelectOption value="2">Two</SelectOption>
                <SelectOption value="3">Three</SelectOption>
            </SelectField>
        );
    },
};

export const Multiple: Story = {
    args: {
        id: "select-field",
        label: "Label",
        required: true,
        description: "A short description",
        multiple: true,
    },
    argTypes: {
        onChange: { action: "onChange" },
    },
    render: function Render(props) {
        const { onChange: _, ...others } = props;

        const [choices, setChoices] = useState<string[]>([]);

        const allChoices = ["Option 1", "Option 2", "Option 3"];

        const onChange = (value: string) => {
            setChoices((choices) => {
                let c = [];
                if (choices.includes(value)) {
                    c = choices.filter((c) => c !== value);
                } else {
                    c = [...choices, value];
                }
                props.onChange?.(c);
                return c;
            });
        };

        return (
            <SelectField onChange={onChange} value={choices.length} {...others}>
                {allChoices.map((choice) => (
                    <SelectOption
                        key={choice}
                        value={choice}
                        selected={choices.includes(choice)}
                    >
                        {choice}
                    </SelectOption>
                ))}
            </SelectField>
        );
    },
};

export default meta;
